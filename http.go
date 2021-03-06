package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"

	"gopkg.in/mgo.v2"

	"github.com/shurcooL/httpfs/html/vfstemplate"
)

// LoadTemplates 함수는 템플릿을 로딩합니다.
func LoadTemplates() (*template.Template, error) {
	t := template.New("")
	t, err := vfstemplate.ParseGlob(assets, t, "/template/*.html")
	return t, err
}

func webserver() {
	// 템플릿 로딩을 위해서 vfs(가상파일시스템)을 로딩합니다.
	vfsTemplate, err := LoadTemplates()
	if err != nil {
		log.Fatal(err)
	}
	TEMPLATES = vfsTemplate
	// assets 설정
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(assets)))
	// 웹주소 설정
	http.HandleFunc("/", handleIndex)
	http.HandleFunc("/search", handleSearch)
	http.HandleFunc("/add", handleAdd)
	// RestAPI
	http.HandleFunc("/api/schedule", handleAPISchedule)
	// 웹서버 실행
	err = http.ListenAndServe(*flagHTTPPort, nil)
	if err != nil {
		log.Fatal(err)
	}
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	type Today struct {
		Year  int `bson:"year" json:"year"`
		Month int `bson:"month" json:"month"`
		Date  int `bson:"date" json:"date"`
	}
	type recipe struct {
		Theme      string  `bson:"theme" json:"theme"`
		Dates      [42]int `bson:"dates" json:"dates"`
		Today      `bson:"today" json:"today"`
		QueryYear  int `bson:"queryyear" json:"queryyear"`
		QueryMonth int `bson:"querymonth" json:"querymonth"`
		LastYear   int `bson:"lastyear" json:"lastyear"`
		NextYear   int `bson:"nextyear" json:"nextyear"`
		LastMonth  int `bson:"lastmonth" json:"lastmonth"`
		NextMonth  int `bson:"nextmonth" json:"nextmonth"`
	}
	rcp := recipe{
		Theme: "default.css",
	}
	y, m, d := time.Now().Date()
	rcp.Today.Year = y
	rcp.Today.Month = int(m)
	rcp.Today.Date = d
	q := r.URL.Query()
	userID := q.Get("userid")
	month, err := strconv.Atoi(q.Get("month"))
	if err != nil {
		m := rcp.Today.Month // 입력이 제대로 안되면 이번 달을 넣는다
		rcp.QueryMonth = m
		switch m {
		case 1:
			rcp.LastMonth = 12
			rcp.NextMonth = m + 1
		case 12:
			rcp.LastMonth = m - 1
			rcp.NextMonth = 1
		default:
			rcp.LastMonth = m - 1
			rcp.NextMonth = m + 1
		}
		month = m
	}
	rcp.QueryMonth = month
	switch month {
	case 1:
		rcp.LastMonth = 12
		rcp.NextMonth = month + 1
	case 12:
		rcp.LastMonth = month - 1
		rcp.NextMonth = 1
	default:
		rcp.LastMonth = month - 1
		rcp.NextMonth = month + 1
	}
	year, err := strconv.Atoi(q.Get("year"))
	if err != nil {
		y = rcp.Today.Year // 입력이 제대로 안되면 올해 연도를 넣는다.
		rcp.QueryYear = y
		switch month {
		case 1:
			rcp.LastYear = y - 1
			rcp.NextYear = y
		case 12:
			rcp.LastYear = y
			rcp.NextYear = y + 1
		default:
			rcp.LastYear = y
			rcp.NextYear = y
		}
		year = y
	}
	rcp.QueryYear = year
	switch month {
	case 1:
		rcp.LastYear = year - 1
		rcp.NextYear = year
	case 12:
		rcp.LastYear = year
		rcp.NextYear = year + 1
	default:
		rcp.LastYear = year
		rcp.NextYear = year
	}
	// 75mm studio 일때만 css 파일을 변경한다. 이 구조는 개발 초기에만 사용한다.
	if userID == "75mmstudio" {
		rcp.Theme = "75mmstudio.css"
	}
	rcp.Dates, err = genDate(year, month)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = TEMPLATES.ExecuteTemplate(w, "index", rcp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func handleAdd(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("add page"))
}

// handleSearch
func handleSearch(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	userID := q.Get("userid")
	year := q.Get("year")
	month := q.Get("month")
	day := q.Get("day")
	layer := q.Get("layer")
	sortKey := q.Get("sortkey")
	if userID == "" {
		http.Error(w, "URL에 userid를 입력해주세요", http.StatusBadRequest)
		return
	}

	log.Println(year, month, day, layer, sortKey)

	session, err := mgo.Dial(*flagDBIP)
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()
	schedules, err := allSchedules(session, userID)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	err = json.NewEncoder(w).Encode(schedules)
	if err != nil {
		log.Println(err)
	}
}

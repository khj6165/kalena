package main

import (
	"errors"
	"time"
)

// Calendar 자료구조
type Calendar struct {
	Layers []Layer `json:"layers"`
}

//Layer 자료구조
type Layer struct {
	Title     string     `json:"title"`
	Color     string     `json:"color"` //#FF3366
	Greyscale bool       `json:"greyscale"`
	Hidden    bool       `json:"hidden"`
	Schedules []Schedule `json:"schedules"`
}

// Schedule 자료구조
type Schedule struct {
	Title string `json:"title"`
	Start string `json:"start"`
	End   string `json:"end"`
}

// CheckError 매소드는 Layer 자료구조에 에러가 있는지 체크한다.
func (l Layer) CheckError() error {
	if l.Title == "" {
		return errors.New("Layer의 Title이 빈 문자열 입니다")
	}
	if l.Color != "" {
		if !regexWebColor.MatchString(l.Color) {
			return errors.New("#FF0011 형식의 문자열이 아닙니다")
		}
	}
	return nil
}

// CheckError 매소드는 Schedule 자료구조에 에러가 있는지 체크한다.
func (s Schedule) CheckError() error {
	if s.Title == "" {
		return errors.New("Title 이 빈 문자열 입니다")
	}
	if s.Start == "" {
		return errors.New("Start 시간이 빈 문자열 입니다")
	}
	if s.End == "" {
		return errors.New("End 시간이 빈 문자열 입니다")
	}
	if !regexRFC3339Time.MatchString(s.Start) {
		return errors.New("Start 시간이 2019-09-09T14:43:34+09:00 형식의 문자열이 아닙니다")
	}
	if !regexRFC3339Time.MatchString(s.End) {
		return errors.New("End 시간이 2019-09-09T14:43:34+09:00 형식의 문자열이 아닙니다")
	}
	startTime, err := time.Parse("2006-01-02T15:04:05-07:00", s.Start)
	if err != nil {
		return err
	}
	endTime, err := time.Parse("2006-01-02T15:04:05-07:00", s.End)
	if err != nil {
		return err
	}
	// end가 start 시간보다 큰지 체크하는 부분
	if !endTime.After(startTime) {
		return errors.New("끝시간이 시작시간보다 작습니다")
	}
	return nil
}

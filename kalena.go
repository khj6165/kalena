package main

import (
	"flag"
	"log"
	"os"

	"gopkg.in/mgo.v2"
)

var (
	flagAdd            = flag.Bool("add", false, "Add mode")
	flagTitle          = flag.String("title", "", "Title")
	flagLayerTitle     = flag.String("layerTitle", "", "Layer Title")
	flagLayerColor     = flag.String("layerColor", "", "Layer Color")
	flagLayerHidden    = flag.Bool("hidden", false, "Layer hidden")
	flagLayerGreyscale = flag.Bool("greyscale", false, "Layer geyscale")
	flagStart          = flag.String("start", "", "Start time")
	flagEnd            = flag.String("end", "", "End time")
	flagLocation       = flag.String("location", "Asia/Seoul", "location name")
	flagDBIP           = flag.String("dbip", "", "database ip")

	flagHTTPPort = flag.String("http", "", "Web Service Port Number")
)

func main() {
	flag.Parse()
	if *flagAdd {
		c := Calendar{}
		l := Layer{}
		s := Schedule{}

		l.Title = *flagLayerTitle
		l.Color = *flagLayerColor
		l.Hidden = *flagLayerHidden
		l.Greyscale = *flagLayerGreyscale

		s.Title = *flagTitle
		s.Start = *flagStart
		s.End = *flagEnd

		// 체크 구문
		err := l.CheckError()
		if err != nil {
			log.Fatal(err)
		}
		err = s.CheckError()
		if err != nil {
			log.Fatal(err)
		}

		l.Schedules = append(l.Schedules, s)
		c.Layers = append(c.Layers, l)

		session, err := mgo.Dial(*flagDBIP)
		if err != nil {
			log.Fatal(err)
		}
		defer session.Close()
		err = AddCalendar(session, c)
		if err != nil {
			log.Print(err)
		}
	} else if *flagHTTPPort != "" {
		webserver(*flagHTTPPort)
	} else {
		flag.PrintDefaults()
		os.Exit(1)
	}
}

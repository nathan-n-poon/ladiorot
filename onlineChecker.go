package main

import (
	"errors"
	"os"
	"strings"
	"time"
)

var savedLatestDatePath = "lastCommDate.txt"

type onlineChecker struct {
	latestDate time.Time
	MsgChan    chan string
}

func (o *onlineChecker) checkApplicability(subject string) bool {
	if strings.Index(subject, "PING") != -1 {
		return true
	}
	return false
}

func (o *onlineChecker) summary() error {
	if time.Now().Sub(o.latestDate) > 36*time.Hour {
		return errors.New("onion no contact")
	}
	return nil
}

func (o *onlineChecker) run(doneChan chan bool) {
	if _, err := os.Stat(savedLatestDatePath); err == nil {
		dat, err := os.ReadFile(savedLatestDatePath)
		check(err)
		o.latestDate, err = time.Parse("2006-01-02 15:04:05.999999999 -0700 MST", string(dat))
	} else {
		o.latestDate = time.Time{}
	}
	dateDelim := "Date: "

	for msg := range o.MsgChan {
		dateStr := utilGetField(msg, dateDelim)
		date, err := time.Parse(layout, dateStr)
		check(err)

		if date.Sub(o.latestDate) > 0 {
			o.latestDate = date
		}
	}
	file, err := os.Create(savedLatestDatePath)
	check(err)
	defer file.Close()
	file.WriteString(o.latestDate.String())
	doneChan <- true
}

func (o *onlineChecker) getIngressChan() chan string {
	return o.MsgChan
}

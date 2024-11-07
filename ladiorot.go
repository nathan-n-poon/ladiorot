package main

import (
	"os/exec"
	"strconv"
	"strings"
)

var layout string

func check(e error) {
	if e != nil {
		panic(e)
	}
}

type policyChecker interface {
	checkApplicability(string) bool
	summary() error
	run(doneChan chan bool)
	getIngressChan() chan string
}

func main() {
	// some applescript bs
	nnbsp, _ := strconv.Unquote(`"\u202F"`)
	layout = "Monday, January 2, 2006 at 3:04:05" + nnbsp + "PM" //"Wednesday, October 23, 2024 at 7:47:42 PM"

	dbg := &onlineChecker{MsgChan: make(chan string)}
	checkers := []policyChecker{dbg}

	doneChan := make(chan bool)

	go func() {
		for _, checker := range checkers {
			checker.run(doneChan)
		}
	}()

	out, err := exec.Command("osascript", "./toroidal.scptd").Output()
	check(err)

	strOut := strings.TrimSpace(string(out))
	mailDelim := "ENTRY_DELIM"
	mailList := strings.Split(strOut, mailDelim)
	mailList = mailList[:len(mailList)-1]

	for _, mail := range mailList {
		subject := utilGetField(mail, "Subject: ")
		for _, checker := range checkers {
			if checker.checkApplicability(subject) {
				checker.getIngressChan() <- mail
			}
		}
	}
	for _, checker := range checkers {
		close(checker.getIngressChan())
	}

	for _ = range len(checkers) {
		<-doneChan
	}

	for _, checker := range checkers {
		err := checker.summary()
		if err != nil {
			//send email alert
			println(err.Error())
		}
	}
}

func utilGetField(msg, startDelim string) string {
	var fieldDelim = "FIELD_DELIM"

	begIndex := strings.Index(msg, startDelim) + len(startDelim)
	endIndex := strings.Index(msg[begIndex:], fieldDelim) + begIndex
	return strings.TrimSpace(msg[begIndex:endIndex])
}

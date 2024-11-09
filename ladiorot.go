package main

import (
	"github.com/joho/godotenv"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
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

func doCheck(masterDoneChan chan bool) {
	// some applescript bs
	nnbsp, _ := strconv.Unquote(`"\u202F"`)
	layout = "Monday, January 2, 2006 at 3:04:05" + nnbsp + "PM" //"Wednesday, October 23, 2024 at 7:47:42 PM"

	checkers := []policyChecker{
		&onlineChecker{MsgChan: make(chan string)},
		&temperatureChecker{MsgChan: make(chan string)},
	}

	servantDoneChan := make(chan bool)

	for _, checker := range checkers {
		go checker.run(servantDoneChan)
	}

	err := godotenv.Load()
	check(err)
	recAddy := os.Getenv("REC_ADDY")
	cmd := "./bash/readMail.sh" + " " + `"` +
		recAddy + `"`
	out, err := exec.Command("bash", "-c", cmd).Output()
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
		<-servantDoneChan
	}

	for _, checker := range checkers {
		err := checker.summary()
		if err != nil {
			//send email alert
			println(err.Error())
			sendEmail(err.Error())
		}
	}
	masterDoneChan <- true
}

func main() {
	doneChan := make(chan bool)
	timeOutChan := make(chan bool)

	go doCheck(doneChan)
	go func(timeOutChan chan bool) {
		time.Sleep(5 * time.Second)
		timeOutChan <- true
	}(timeOutChan)

	select {
	case <-doneChan:
		return
	case <-timeOutChan:
		panic("Checking emails stalled out :(")
	}
}

func sendEmail(emailMsg string) {
	err := godotenv.Load()
	check(err)
	destAddy := os.Getenv("DEST_ADDY")

	cmd := "./bash/sendMail.sh" + " " + `"` +
		destAddy + "|" +
		emailMsg + `"`
	out, err := exec.Command("bash", "-c", cmd).Output()
	check(err)
	println(string(out))
}

func utilGetField(msg, startDelim string) string {
	var fieldDelim = "FIELD_DELIM"

	begIndex := strings.Index(msg, startDelim) + len(startDelim)
	endIndex := strings.Index(msg[begIndex:], fieldDelim) + begIndex
	return strings.TrimSpace(msg[begIndex:endIndex])
}

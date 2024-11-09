package main

import (
	"errors"
	"strconv"
	"strings"
)

var dangerTemp = 80.0
var stressTemp = dangerTemp - 10.0

type temperatureChecker struct {
	MsgChan        chan string
	triggerWarning bool
}

func (t *temperatureChecker) checkApplicability(subject string) bool {
	if strings.Index(subject, "TEMPERATURE") != -1 {
		return true
	}
	return false
}

func (t *temperatureChecker) summary() error {
	if t.triggerWarning {
		return errors.New("onion too hot")
	}
	return nil
}

func (t *temperatureChecker) run(doneChan chan bool) {
	tempDelim := "Temp: "
	lengthCurRun := 0

	for msg := range t.MsgChan {
		tempStr := utilGetField(msg, tempDelim)
		temp, err := strconv.ParseFloat(tempStr, 64)
		check(err)

		if temp > stressTemp {
			lengthCurRun++
			if temp < dangerTemp && lengthCurRun > 1 {
				t.triggerWarning = true
			} else if temp >= dangerTemp {
				t.triggerWarning = true
			}
		} else {
			lengthCurRun = 0
		}
	}
	doneChan <- true
}

func (t *temperatureChecker) getIngressChan() chan string {
	return t.MsgChan
}

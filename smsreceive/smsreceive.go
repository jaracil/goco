// Package sms-receive is a GopherJS wrapper for cordova plugin sms receive.
//
// Install plugin:
//  cordova plugin add cordova-plugin-sms-receive

package smsreceive

import (
	"errors"

	"github.com/gopherjs/gopherjs/js"
)

type Sms struct {
	*js.Object
	Data SmsData `js:"data"`
}

type SmsData struct {
	*js.Object
	Address       string `js:"address"`
	Body          string `js:"body"`
	Date          int64  `js:"date"`
	DateSent      int64  `js:"date_sent"`
	Read          int    `js:"read"`
	Seen          int    `js:"seen"`
	ServiceCenter string `js:"service_center"`
	Status        int    `js:"status"`
	Type          int    `js:"type"`
}

var instance *js.Object

func mo() *js.Object {
	if instance == nil {
		instance = js.Global.Get("SMSReceive")
	}
	return instance
}

// StartWatch method starts listening for incomming SMS and raises an onSMSArrive event when this happens.
func StartWatch() (err error) {
	ch := make(chan interface{})
	success := func() {
		close(ch)
	}
	failure := func(e string) {
		err = errors.New(e)
		close(ch)
	}
	mo().Call("startWatch", success, failure)
	<-ch
	return
}

// StopWatch method stops listening for incomming SMS. Rember always invoke this method when you have got your SMS to avoid memory leaks.
func StopWatch() (err error) {
	ch := make(chan interface{})
	success := func() {
		close(ch)
	}
	failure := func(e string) {
		err = errors.New(e)
		close(ch)
	}
	mo().Call("stopWatch", success, failure)
	<-ch
	return
}

// SubscribeOnSmsArrive creates the event listener for onSMSArrive event.
func SubscribeOnSmsArrive(cb func(*Sms)) {
	js.Global.Get("document").Call("addEventListener", "onSMSArrive", cb, false)
}

// UnsubscribeOnSmsArrive removes the event listener for onSMSArrive event.
func UnsubscribeOnSmsArrive(cb func(*Sms)) {
	js.Global.Get("document").Call("removeEventListener", "onSMSArrive", cb, false)
}

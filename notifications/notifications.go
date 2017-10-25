// Package push is a GopherJS wrapper for cordova local-notifications plugin.
//
// Install plugin:
//  cordova plugin add https://github.com/katzer/cordova-plugin-local-notifications
package push

import (
	"time"

	"github.com/gopherjs/gopherjs/js"
	"github.com/jaracil/goco"
)

type Notification struct {
	*js.Object
	ID        int       `js:"id"`
	Title     string    `js:"title"`
	Text      string    `js:"text"`
	Every     string    `js:"every"`
	At        time.Time `js:"at"`
	FirstAt   time.Time `js:"firstAt"`
	Badge     int       `js:"badge"`
	Sound     string    `js:"sound"`
	Data      string    `js:"data"`
	Icon      string    `js:"icon"`
	SmallIcon string    `js:"smallIcon"`
	Ongoing   bool      `js:"ongoing"`
	Led       string    `js:"led"`
}

var mo *js.Object

func init() {
	goco.OnDeviceReady(func() {
		mo = js.Global.Get("cordova").Get("plugins").Get("notification").Get("local")
	})
}

// New returns new notification with default values
func New() *Notification {
	return &Notification{Object: js.Global.Get("Object").New()}
}

// HasPermission determines permission to show local notifications
func HasPermission() (res bool) {
	ch := make(chan bool, 1)
	success := func(granted bool) {
		ch <- granted
	}
	mo.Call("hasPermission", success)
	return <-ch
}

// RegisterPermission registers permission to show local notifications
func RegisterPermission() (res bool) {
	ch := make(chan bool, 1)
	success := func(granted bool) {
		ch <- granted
	}
	mo.Call("registerPermission", success)
	return <-ch
}

// Schedule accepts *Notification or []*Notification to schedule.
func Schedule(notif interface{}) {
	ch := make(chan struct{})
	success := func() {
		close(ch)
	}
	mo.Call("schedule", notif, success)
	<-ch
}

// Update accepts *Notification or []*Notification to update.
func Update(notif interface{}) {
	ch := make(chan struct{})
	success := func() {
		close(ch)
	}
	mo.Call("update", notif, success)
	<-ch
}

// Clear clears notification by id or slice of ids.
func Clear(id interface{}) {
	ch := make(chan struct{})
	success := func() {
		close(ch)
	}
	mo.Call("clear", id, success)
	<-ch
}

// ClearAll clears all notifications.
func ClearAll() {
	ch := make(chan struct{})
	success := func() {
		close(ch)
	}
	mo.Call("clearAll", success)
	<-ch
}

// Cancel cancels notification by id or slice of ids.
func Cancel(id interface{}) {
	ch := make(chan struct{})
	success := func() {
		close(ch)
	}
	mo.Call("cancel", id, success)
	<-ch
}

// CancelAll cancels all notifications.
func CancelAll() {
	ch := make(chan struct{})
	success := func() {
		close(ch)
	}
	mo.Call("cancelAll", success)
	<-ch
}

// OnSchedule registers a callback which is invoked when a local notification was scheduled.
func OnSchedule(f func(*Notification, string)) {
	mo.Call("on", "schedule", f)
}

// OnTrigger registers a callback which is invoked when a local notification was triggered.
func OnTrigger(f func(*Notification, string)) {
	mo.Call("on", "trigger", f)
}

// OnUpdate registers a callback which is invoked when a local notification was updated.
func OnUpdate(f func(*Notification, string)) {
	mo.Call("on", "update", f)
}

// OnClick registers a callback which is invoked when a local notification was clicked.
func OnClick(f func(*Notification, string)) {
	mo.Call("on", "click", f)
}

// OnClear registers a callback which is invoked when a local notification was cleared from the notification center.
func OnClear(f func(*Notification, string)) {
	mo.Call("on", "clear", f)
}

// OnCancel registers a callback which is invoked when a local notification was canceled.
func OnCancel(f func(*Notification, string)) {
	mo.Call("on", "cancel", f)
}

// OnClearAll registers a callback which is invoked when all notifications were cleared from the notification center.
func OnClearAll(f func(string)) {
	mo.Call("on", "clearall", f)
}

// OnCancelAll registers a callback which is invoked when all local notification were canceled.
func OnCancelAll(f func(string)) {
	mo.Call("on", "cancelall", f)
}

// IsPresent returns true if notification is still present in the notification center
func IsPresent(id int) (res bool) {
	ch := make(chan bool, 1)
	success := func(present bool) {
		ch <- present
	}
	mo.Call("isPresent", id, success)
	return <-ch
}

// IsScheduled returns true if local notification was scheduled but not yet
// triggered and added to the notification center. Repeating local notifications are always in that state.
func IsScheduled(id int) (res bool) {
	ch := make(chan bool, 1)
	success := func(present bool) {
		ch <- present
	}
	mo.Call("isScheduled", id, success)
	return <-ch
}

// IsTriggered returns true if event has occurred and the local notification is added to the notification center.
func IsTriggered(id int) (res bool) {
	ch := make(chan bool, 1)
	success := func(present bool) {
		ch <- present
	}
	mo.Call("IsTriggered", id, success)
	return <-ch
}

func objToIntSlice(obj *js.Object) []int {
	l := obj.Length()
	ret := make([]int, 0)
	for n := 0; n < l; n++ {
		ret = append(ret, obj.Index(n).Int())
	}
	return ret
}

func GetAllIds() []int {
	ch := make(chan *js.Object, 1)
	success := func(obj *js.Object) {
		ch <- obj
	}
	mo.Call("getAllIds", success)
	return objToIntSlice(<-ch)
}

func GetScheduledIds() []int {
	ch := make(chan *js.Object, 1)
	success := func(obj *js.Object) {
		ch <- obj
	}
	mo.Call("getScheduledIds", success)
	return objToIntSlice(<-ch)
}

func GetTriggeredIds() []int {
	ch := make(chan *js.Object, 1)
	success := func(obj *js.Object) {
		ch <- obj
	}
	mo.Call("getTriggeredIds", success)
	return objToIntSlice(<-ch)
}

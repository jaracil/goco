// Package notifications is a GopherJS wrapper for cordova local-notifications plugin.
//
// Install plugin:
//  cordova plugin add https://github.com/katzer/cordova-plugin-local-notifications
package notifications

import (
	"time"

	"github.com/gopherjs/gopherjs/js"
)

// Notification type contains all information about one notification.
type Notification struct {
	*js.Object
	ID        int       `js:"id"`        // A unique identifier required to clear, cancel, update or retrieve the local notification in the future
	Title     string    `js:"title"`     // First row of the notification - Default: Empty string (iOS) or the app name (Android)
	Text      string    `js:"text"`      // Second row of the notification - Default: Empty string
	Badge     int       `js:"badge"`     // The number currently set as the badge of the app icon in Springboard (iOS) or at the right-hand side of the local notification (Android) - Default: 0 (which means don't show a number)
	Sound     string    `js:"sound"`     // Uri of the file containing the sound to play when an alert is displayed	- Default: res://platform_default
	Data      string    `js:"data"`      // Arbitrary data.
	Icon      string    `js:"icon"`      // Uri of the icon that is shown in the ticker and notification - Default: res://icon
	SmallIcon string    `js:"smallIcon"` // Uri of the resource (only res://) to use in the notification layouts. Different classes of devices may return different sizes - Default: res://ic_popup_reminder
	Ongoing   bool      `js:"ongoing"`   // Ongoing notification
	Led       string    `js:"led"`       // ARGB value that you would like the LED on the device to blink - Default: FFFFFF
	Every     string    `js:"every"`     // The interval at which to reschedule the local notification. That can be a value of second, minute, hour, day, week, month or year.
	At        time.Time // The date and time when the system should deliver the local notification. If the specified value is nil or is a date in the past, the local notification is delivered immediately.
	FirstAt   time.Time // The date and time when the system should first deliver the local notification. If the specified value is nil or is a date in the past, the local notification is delivered immediately.

}

var instance *js.Object

func mo() *js.Object {
	if instance == nil {
		instance = js.Global.Get("cordova").Get("plugins").Get("notification").Get("local")
	}
	return instance
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
	mo().Call("hasPermission", success)
	return <-ch
}

// RegisterPermission registers permission to show local notifications
func RegisterPermission() (res bool) {
	ch := make(chan bool, 1)
	success := func(granted bool) {
		ch <- granted
	}
	mo().Call("registerPermission", success)
	return <-ch
}

// Schedule accepts *Notification or []*Notification to schedule.
func Schedule(notif *Notification) {
	if !notif.At.IsZero() {
		notif.Set("at", notif.At.UnixNano()/1000000)
	}
	if !notif.FirstAt.IsZero() {
		notif.Set("firstAt", notif.FirstAt.UnixNano()/1000000)
	}

	ch := make(chan struct{})
	success := func() {
		close(ch)
	}
	mo().Call("schedule", notif, success)
	<-ch
}

// Update accepts *Notification or []*Notification to update.
func Update(notif *Notification) {
	ch := make(chan struct{})
	success := func() {
		close(ch)
	}
	mo().Call("update", notif, success)
	<-ch
}

// Clear clears notification by ID.
func Clear(id int) {
	ch := make(chan struct{})
	success := func() {
		close(ch)
	}
	mo().Call("clear", id, success)
	<-ch
}

// ClearAll clears all notifications.
func ClearAll() {
	ch := make(chan struct{})
	success := func() {
		close(ch)
	}
	mo().Call("clearAll", success)
	<-ch
}

// Cancel cancels notification by ID.
func Cancel(id int) {
	ch := make(chan struct{})
	success := func() {
		close(ch)
	}
	mo().Call("cancel", id, success)
	<-ch
}

// CancelAll cancels all notifications.
func CancelAll() {
	ch := make(chan struct{})
	success := func() {
		close(ch)
	}
	mo().Call("cancelAll", success)
	<-ch
}

// OnSchedule registers a callback which is invoked when a local notification was scheduled.
func OnSchedule(f func(*Notification, string)) {
	mo().Call("on", "schedule", f)
}

// OnTrigger registers a callback which is invoked when a local notification was triggered.
func OnTrigger(f func(*Notification, string)) {
	mo().Call("on", "trigger", f)
}

// OnUpdate registers a callback which is invoked when a local notification was updated.
func OnUpdate(f func(*Notification, string)) {
	mo().Call("on", "update", f)
}

// OnClick registers a callback which is invoked when a local notification was clicked.
func OnClick(f func(*Notification, string)) {
	mo().Call("on", "click", f)
}

// OnClear registers a callback which is invoked when a local notification was cleared from the notification center.
func OnClear(f func(*Notification, string)) {
	mo().Call("on", "clear", f)
}

// OnCancel registers a callback which is invoked when a local notification was canceled.
func OnCancel(f func(*Notification, string)) {
	mo().Call("on", "cancel", f)
}

// OnClearAll registers a callback which is invoked when all notifications were cleared from the notification center.
func OnClearAll(f func(string)) {
	mo().Call("on", "clearall", f)
}

// OnCancelAll registers a callback which is invoked when all local notification were canceled.
func OnCancelAll(f func(string)) {
	mo().Call("on", "cancelall", f)
}

// IsPresent returns true if notification is still present in the notification center
func IsPresent(id int) (res bool) {
	ch := make(chan bool, 1)
	success := func(present bool) {
		ch <- present
	}
	mo().Call("isPresent", id, success)
	return <-ch
}

// IsScheduled returns true if local notification was scheduled but not yet
// triggered and added to the notification center. Repeating local notifications are always in that state.
func IsScheduled(id int) (res bool) {
	ch := make(chan bool, 1)
	success := func(present bool) {
		ch <- present
	}
	mo().Call("isScheduled", id, success)
	return <-ch
}

// IsTriggered returns true if event has occurred and the local notification is added to the notification center.
func IsTriggered(id int) (res bool) {
	ch := make(chan bool, 1)
	success := func(present bool) {
		ch <- present
	}
	mo().Call("isTriggered", id, success)
	return <-ch
}

// GetAllIds returns all alive ID's.
func GetAllIds() []int {
	ch := make(chan []int, 1)
	success := func(s []int) {
		ch <- s
	}
	mo().Call("getAllIds", success)
	return <-ch
}

// GetScheduledIds returns scheduled ID's.
func GetScheduledIds() []int {
	ch := make(chan []int, 1)
	success := func(s []int) {
		ch <- s
	}
	mo().Call("getScheduledIds", success)
	return <-ch
}

// GetTriggeredIds returns triggered ID's.
func GetTriggeredIds() []int {
	ch := make(chan []int, 1)
	success := func(s []int) {
		ch <- s
	}
	mo().Call("getTriggeredIds", success)
	return <-ch
}

// GetAll returns all alive notifications.
func GetAll() []*Notification {
	ch := make(chan []*Notification, 1)
	success := func(s []*Notification) {
		ch <- s
	}
	mo().Call("getAll", success)
	return <-ch
}

// GetByID returns notification by ID
func GetByID(id int) *Notification {
	ch := make(chan *Notification, 1)
	success := func(n *Notification) {
		ch <- n
	}
	mo().Call("get", id, success)
	return <-ch
}

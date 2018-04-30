// Package dialogs is a GopherJS wrapper for dialogs plugin.
//
// Install plugin:
//  cordova plugin add cordova-plugin-dialogs
package dialogs

import "github.com/gopherjs/gopherjs/js"

var instance *js.Object

func mo() *js.Object {
	if instance == nil {
		instance = js.Global.Get("navigator").Get("notification")
	}
	return instance
}

// Alert shows a custom dialog box, most Cordova implementations use a native dialog box for this feature,
// but some platforms use the browser alert, which is less customizable. It returns true when the dialog is
// dismissed and false in case of error.
func Alert(message, title, button string) bool {
	if mo() == js.Undefined || mo() == nil {
		return false
	}

	ch := make(chan struct{})
	dismissed := func(_ *js.Object) {
		close(ch)
	}
	mo().Call("alert", message, dismissed, title, button)

	<-ch
	return true
}

// Confirm shows a custom dialog with multiple chooses, it returns the index of the button clicked. Given a button
// slice of ["Cancel", "Ok"] it will be 0 if "Cancel" is clicked. It returns -1 in case of error or when the dialog
// is dismissed without a button press.
func Confirm(message, title string, buttons []string) (index int) {
	if mo() == js.Undefined || mo() == nil {
		return -1
	}

	ch := make(chan struct{})
	dismissed := func(obj *js.Object) {
		index = obj.Int() - 1
		close(ch)
	}
	mo().Call("confirm", message, dismissed, title, buttons)

	<-ch
	return index
}

// Prompt shows a native dialog with text-input. It return the index of the button pressed (see Confirm function) and
// also the text-value entered in the prompt.
func Prompt(message, title string, buttons []string, input string) (index int, value string) {
	if mo() == js.Undefined || mo() == nil {
		return -1, input
	}

	ch := make(chan struct{})
	dismissed := func(obj *js.Object) {
		index = obj.Get("buttonIndex").Int() - 1
		value = obj.Get("input1").String()
		close(ch)
	}
	mo().Call("prompt", message, dismissed, title, buttons, input)

	<-ch
	return index, value
}

// Beep will play a beep sound for N times.
func Beep(times uint){
	if mo() == js.Undefined || mo() == nil || times == 0 {
		return
	}

	mo().Call("beep", times)
}

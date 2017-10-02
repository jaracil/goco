package background

import (
	"github.com/gopherjs/gopherjs/js"
)

var bgc *js.Object

func bg() *js.Object {
	if bgc == nil {
		bgc = js.Global.Get("cordova").Get("plugins").Get("backgroundMode")
	}
	return bgc
}

// Enable enables/disables background mode.
func Enable(en bool) {
	bg().Call("setEnabled", en)
}

// InBackground returns true if app is in background
func InBackground() bool {
	return bg().Call("isActive").Bool()
}

// MoveToBackground moves the app to background. (Android only)
func MoveToBackground() {
	bg().Call("moveToBackground")
}

// MoveToForeground moves the app to foreground. (Android only)
func MoveToForeground() {
	bg().Call("moveToForeground")
}

// OverrideBackButton change backbutton behavior. When back button is pressed
// the app is moved to background instead of close it.
func OverrideBackButton() {
	bg().Call("overrideBackButton")
}

// ExcludeFromTaskList exclude the app from the recent task list works on Android 5.0+.
func ExcludeFromTaskList() {
	bg().Call("ExcludeFromTaskList")
}

// IsScreenOff returns false when screen is off.
func IsScreenOff() (ret bool) {
	ch := make(chan struct{})
	success := func(b bool) {
		ret = !b
		close(ch)
	}
	bg().Call("isScreenOff", success)
	<-ch
	return
}

// Wakeup turns on the screen.
func Wakeup() {
	bg().Call("wakeup")
}

// Unlock moves the app to foreground even the device is locked.
func Unlock() {
	bg().Call("unlock")
}

// DisableWebViewOptimizations disable web view optimizations.
// Various APIs like playing media or tracking GPS position in background
// might not work while in background even the background mode is active.
// To fix such issues the plugin provides a method to disable most optimizations done by Android/CrossWalk.
func DisableWebViewOptimizations() {
	bg().Call("disableWebViewOptimizations")
}

// OnEnable sets the function to be called when background mode is enabled.
func OnEnable(f func()) {
	bg().Call("on", "enable", f)
}

// OnDisable sets the function to be called when background mode is disabled.
func OnDisable(f func()) {
	bg().Call("on", "disable", f)
}

// OnActivate sets the function to be called when app enters in background.
func OnActivate(f func()) {
	bg().Call("on", "activate", f)
}

// OnDeactivate sets the function to be called when app enters in foreground.
func OnDeactivate(f func()) {
	bg().Call("on", "deactivate", f)
}

// OnFailure sets the function to be called on failure
func OnFailure(f func()) {
	bg().Call("on", "failure", f)
}

// UnEnable removes OnEnable callback function
func UnEnable(f func()) {
	bg().Call("un", "enable", f)
}

// UnDisable removes OnDisable callback function
func UnDisable(f func()) {
	bg().Call("un", "disable", f)
}

// UnActivate removes OnActivate callback function
func UnActivate(f func()) {
	bg().Call("un", "activate", f)
}

// UnDeactivate removes OnDeactivate callback function
func UnDeactivate(f func()) {
	bg().Call("un", "deactivate", f)
}

// UnFailure removes OnFailure callback function
func UnFailure(f func()) {
	bg().Call("un", "failure", f)
}

// SetDefaults sets notification defaults to indicate that the app is executing tasks in background and being paused would disrupt the user,
// the plug-in has to create a notification while in background - like a download progress bar.
//
// See cordova module documentation for more information.
func SetDefaults(opts interface{}) {
	bg().Call("setDefaults", opts)
}

// Configure modifies the currently displayed notification
func Configure(opts interface{}) {
	bg().Call("configure", opts)
}

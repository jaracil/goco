package background

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/jaracil/goco"
)

var mo *js.Object

func init() {
	goco.OnDeviceReady(func() {
		mo = js.Global.Get("cordova").Get("plugins").Get("backgroundMode")
	})
}

// Enable enables/disables background mode.
func Enable(en bool) {
	mo.Call("setEnabled", en)
}

// InBackground returns true if app is in background
func InBackground() bool {
	return mo.Call("isActive").Bool()
}

// MoveToBackground moves the app to background. (Android only)
func MoveToBackground() {
	mo.Call("moveToBackground")
}

// MoveToForeground moves the app to foreground. (Android only)
func MoveToForeground() {
	mo.Call("moveToForeground")
}

// OverrideBackButton change backbutton behavior. When back button is pressed
// the app is moved to background instead of close it.
func OverrideBackButton() {
	mo.Call("overrideBackButton")
}

// ExcludeFromTaskList exclude the app from the recent task list works on Android 5.0+.
func ExcludeFromTaskList() {
	mo.Call("ExcludeFromTaskList")
}

// IsScreenOff returns false when screen is off.
func IsScreenOff() (ret bool) {
	ch := make(chan struct{})
	success := func(b bool) {
		ret = !b
		close(ch)
	}
	mo.Call("isScreenOff", success)
	<-ch
	return
}

// Wakeup turns on the screen.
func Wakeup() {
	mo.Call("wakeup")
}

// Unlock moves the app to foreground even the device is locked.
func Unlock() {
	mo.Call("unlock")
}

// DisableWebViewOptimizations disable web view optimizations.
// Various APIs like playing media or tracking GPS position in background
// might not work while in background even the background mode is active.
// To fix such issues the plugin provides a method to disable most optimizations done by Android/CrossWalk.
func DisableWebViewOptimizations() {
	mo.Call("disableWebViewOptimizations")
}

// OnEnable sets the function to be called when background mode is enabled.
func OnEnable(f func()) {
	mo.Call("on", "enable", f)
}

// OnDisable sets the function to be called when background mode is disabled.
func OnDisable(f func()) {
	mo.Call("on", "disable", f)
}

// OnActivate sets the function to be called when app enters in background.
func OnActivate(f func()) {
	mo.Call("on", "activate", f)
}

// OnDeactivate sets the function to be called when app enters in foreground.
func OnDeactivate(f func()) {
	mo.Call("on", "deactivate", f)
}

// OnFailure sets the function to be called on failure
func OnFailure(f func()) {
	mo.Call("on", "failure", f)
}

// UnEnable removes OnEnable callback function
func UnEnable(f func()) {
	mo.Call("un", "enable", f)
}

// UnDisable removes OnDisable callback function
func UnDisable(f func()) {
	mo.Call("un", "disable", f)
}

// UnActivate removes OnActivate callback function
func UnActivate(f func()) {
	mo.Call("un", "activate", f)
}

// UnDeactivate removes OnDeactivate callback function
func UnDeactivate(f func()) {
	mo.Call("un", "deactivate", f)
}

// UnFailure removes OnFailure callback function
func UnFailure(f func()) {
	mo.Call("un", "failure", f)
}

// SetDefaults sets notification defaults to indicate that the app is executing tasks in background and being paused would disrupt the user,
// the plug-in has to create a notification while in background - like a download progress bar.
//
// See cordova module documentation for more information.
func SetDefaults(opts interface{}) {
	mo.Call("setDefaults", opts)
}

// Configure modifies the currently displayed notification
func Configure(opts interface{}) {
	mo.Call("configure", opts)
}

// Package push is a GopherJS wrapper for cordova push plugin.
//
// Install plugin:
//  cordova plugin add phonegap-plugin-push
package push

import (
	"errors"

	"github.com/gopherjs/gopherjs/js"
)

// AndroidCfg contains android platform specific configuration
type AndroidCfg struct {
	*js.Object
	Icon               string   `js:"icon"`               // Optional. The name of a drawable resource to use as the small-icon. The name should not include the extension.
	IconColor          string   `js:"iconColor"`          // Optional. Sets the background color of the small icon on Android 5.0 and greater.
	Sound              bool     `js:"sound"`              // Optional. If true it plays the sound specified in the push data or the default system sound.
	Vibrate            bool     `js:"vibrate"`            // Optional. If true the device vibrates on receipt of notification.
	ClearBadge         bool     `js:"clearBadge"`         // Optional. If true the icon badge will be cleared on init and before push messages are processed.
	ClearNotifications bool     `js:"clearNotifications"` // Optional. If true the app clears all pending notifications when it is closed.
	ForceShow          bool     `js:"forceShow"`          // Optional. Controls the behavior of the notification when app is in foreground. If true and app is in foreground, it will show a notification in the notification drawer, the same way as when the app is in background (and on('notification') callback will be called only when the user clicks the notification). When false and app is in foreground, the on('notification') callback will be called immediately.
	Topics             []string `js:"topics"`             // Optional. If the array contains one or more strings each string will be used to subscribe to a FcmPubSub topic.
	MessageKey         string   `js:"messageKey"`         // Optional. The key to search for text of notification.
	TitleKey           string   `js:"titleKey"`           // Optional. The key to search for title of notification.
}

// IOSCfg contains IOS platform specific configuration
type IOSCfg struct {
	*js.Object
	Alert      bool     `js:"alert"`      // Optional. If true the device shows an alert on receipt of notification.
	Badge      bool     `js:"badge"`      // Optional. If true the device sets the badge number on receipt of notification.
	Sound      bool     `js:"sound"`      // Optional. If true the device plays a sound on receipt of notification.
	ClearBadge bool     `js:"clearBadge"` // Optional. If true the badge will be cleared on app startup.
	FCMSandbox bool     `js:"fcmSandbox"` // Whether to use prod or sandbox GCM setting. Defaults to false.
	Topics     []string `js:"topics"`     // Optional. If the array contains one or more strings each string will be used to subscribe to a FcmPubSub topic.
}

// BrowserCfg contains browser platform specific configuration
type BrowserCfg struct {
	*js.Object
	PushServiceURL       string `js:"pushServiceURL"`       // Optional. URL for the push server you want to use.
	ApplicationServerKey string `js:"applicationServerKey"` // Optional. Your GCM API key if you are using VAPID keys.

}

// Config contains all supported platforms configuration
type Config struct {
	*js.Object
	Android *AndroidCfg `js:"android"`
	IOS     *IOSCfg     `js:"ios"`
	Browser *BrowserCfg `js:"browser"`
}

// RegInfo contains registration info (see OnRegistration)
type RegInfo struct {
	*js.Object
	RegistrationID   string `js:"registrationId"`   // The registration ID provided by the 3rd party remote push service.
	RegistrationType string `js:"registrationType"` // The registration type of the 3rd party remote push service. Either FCM or APNS.
}

// NotifExtData contains notification extra data
type NotifExtData struct {
	*js.Object
	Foreground bool `js:"foreground"` // Whether the notification was received while the app was in the foreground.
	Coldstart  bool `js:"coldstart"`  // Will be true if the application is started by clicking on the push notification, false if the app is already started.
	Dismissed  bool `js:"dismissed"`  // Is set to true if the notification was dismissed by the user.
}

// Notification contains notification data
type Notification struct {
	*js.Object
	Message        string        `js:"message"`        // The text of the push message sent from the 3rd party service.
	Title          string        `js:"title"`          // The optional title of the push message sent from the 3rd party service.
	Count          string        `js:"count"`          // The number of messages to be displayed in the badge in iOS/Android or message count in the notification shade in Android. For windows, it represents the value in the badge notification which could be a number or a status glyph.
	Sound          string        `js:"sound"`          // The name of the sound file to be played upon receipt of the notification.
	Image          string        `js:"image"`          // The path of the image file to be displayed in the notification.
	LaunchArgs     string        `js:"launchArgs"`     // The args to be passed to the application on launch from push notification. This works when notification is received in background. (Windows Only)
	AdditionalData *NotifExtData `js:"additionalData"` // See NotifExtData type
}

// NotifError contains error message
type NotifError struct {
	*js.Object
	Mesage string `js:"message"`
}

// Push object
type Push struct {
	*js.Object
}

var instance *js.Object

func mo() *js.Object {
	if instance == nil {
		instance = js.Global.Get("PushNotification")
	}
	return instance
}

// NewConfig returns new Config object with default values
func NewConfig() *Config {
	cfg := &Config{Object: js.Global.Get("Object").New()}
	cfg.Android = &AndroidCfg{Object: js.Global.Get("Object").New()}
	cfg.IOS = &IOSCfg{Object: js.Global.Get("Object").New()}
	cfg.Browser = &BrowserCfg{Object: js.Global.Get("Object").New()}
	return cfg
}

// New returns Push object. The cfg param is a config object returned by NewConfig and filled with valid data.
func New(cfg *Config) *Push {
	obj := mo().Call("init", cfg)
	return &Push{Object: obj}
}

// HasPermission checks whether the push notification permission has been granted.
func HasPermission() (res bool) {
	ch := make(chan struct{})
	success := func(obj *js.Object) {
		res = obj.Get("isEnabled").Bool()
		close(ch)
	}
	mo().Call("hasPermission", success)
	<-ch
	return
}

// OnRegistration registers a function which will be triggered on each successful registration with the 3rd party push service.
func (p *Push) OnRegistration(f func(*RegInfo)) {
	p.Call("on", "registration", f)
}

// OnNotification registers a function which will be triggered each time a push notification is received by a 3rd party push service on the device.
func (p *Push) OnNotification(f func(*Notification)) {
	p.Call("on", "notification", f)
}

// OnError registers a function which will bre triggered when an internal error occurs and the cache is aborted.
func (p *Push) OnError(f func(*NotifError)) {
	p.Call("on", "error", f)
}

// OffRegistration unregisters a function previously registered by OnRegistration function.
func (p *Push) OffRegistration(f func(*RegInfo)) {
	p.Call("off", "registration", f)
}

// OffNotification unregisters a function previously registered by OnNotification function.
func (p *Push) OffNotification(f func(*Notification)) {
	p.Call("off", "notification", f)
}

// OffError unregisters a function previously registered by OnError function.
func (p *Push) OffError(f func(*NotifError)) {
	p.Call("off", "error", f)
}

// UnRegister is used when the application no longer wants to receive push notifications.
// Beware that this cleans up all event handlers previously registered, so you will need to re-register them if you want them to function again without an application reload.
func (p *Push) UnRegister() (err error) {
	ch := make(chan struct{})
	success := func() {
		close(ch)
	}
	fail := func() {
		err = errors.New("Error on unregister")
		close(ch)
	}
	p.Call("unregister", success, fail)
	<-ch
	return
}

// Subscribe is used when the application wants to subscribe a new topic to receive push notifications.
func (p *Push) Subscribe(topic string) (err error) {
	ch := make(chan struct{})
	success := func() {
		close(ch)
	}
	fail := func() {
		err = errors.New("Error on topic subscribe")
		close(ch)
	}
	p.Call("subscribe", topic, success, fail)
	<-ch
	return
}

// UnSubscribe is used when the application no longer wants to receive push notifications from a specific topic but continue to receive other push messages.
func (p *Push) UnSubscribe(topic string) (err error) {
	ch := make(chan struct{})
	success := func() {
		close(ch)
	}
	fail := func() {
		err = errors.New("Error on topic unsubscribe")
		close(ch)
	}
	p.Call("unsubscribe", topic, success, fail)
	<-ch
	return
}

// SetApplicationIconBadgeNumber sets the badge count visible when the app is not running.
func (p *Push) SetApplicationIconBadgeNumber(count int) (err error) {
	ch := make(chan struct{})
	success := func() {
		close(ch)
	}
	fail := func() {
		err = errors.New("Error on set application icon badge number")
		close(ch)
	}
	p.Call("setApplicationIconBadgeNumber", success, fail, count)
	<-ch
	return
}

// GetApplicationIconBadgeNumber gets the current badge count visible when the app is not running.
func (p *Push) GetApplicationIconBadgeNumber() (res int, err error) {
	ch := make(chan struct{})
	success := func(n int) {
		res = n
		close(ch)
	}
	fail := func() {
		err = errors.New("Error on get application icon badge number")
		close(ch)
	}
	p.Call("getApplicationIconBadgeNumber", success, fail)
	<-ch
	return
}

// Finish tells the OS that you are done processing a background push notification.
func (p *Push) Finish(id string) (err error) {
	ch := make(chan struct{})
	success := func(n int) {
		close(ch)
	}
	fail := func() {
		err = errors.New("Error on finish background notification")
		close(ch)
	}
	p.Call("finish", success, fail, id)
	<-ch
	return
}

// ClearAllNotifications tells the OS to clear all notifications from the Notification Center.
func (p *Push) ClearAllNotifications() (err error) {
	ch := make(chan struct{})
	success := func(n int) {
		close(ch)
	}
	fail := func() {
		err = errors.New("Error on clear all notifications")
		close(ch)
	}
	p.Call("clearAllNotifications", success, fail)
	<-ch
	return
}

package push

import (
	"errors"

	"github.com/gopherjs/gopherjs/js"
)

type AndroidCfg struct {
	*js.Object
	Icon               string   `js:"icon"`
	IconColor          string   `js:"iconColor"`
	Sound              bool     `js:"sound"`
	Vibrate            bool     `js:"vibrate"`
	ClearBadge         bool     `js:"clearBadge"`
	ClearNotifications bool     `js:"clearNotifications"`
	ForceShow          bool     `js:"forceShow"`
	Topics             []string `js:"topics"`
	MessageKey         string   `js:"messageKey"`
	TitleKey           string   `js:"titleKey"`
	Channels           []string `js:"channels"`
}

type IOSCfg struct {
	*js.Object
	Alert      bool     `js:"alert"`
	Badge      bool     `js:"badge"`
	Sound      bool     `js:"sound"`
	ClearBadge bool     `js:"clearBadge"`
	FCMSandbox bool     `js:"fcmSandbox"`
	Topics     []string `js:"topics"`
}

type BrowserCfg struct {
	*js.Object
	PushServiceURL       string `js:"pushServiceURL"`
	ApplicationServerKey string `js:"applicationServerKey"`
}

type Config struct {
	*js.Object
	Android *AndroidCfg `js:"android"`
	IOS     *IOSCfg     `js:"ios"`
	Browser *BrowserCfg `js:"browser"`
}

type RegInfo struct {
	*js.Object
	RegistrationId   string `js:"registrationId"`
	RegistrationType string `js:"registrationType"`
}

type NotifExtData struct {
	*js.Object
	Foreground bool `js:"foreground"`
	Coldstart  bool `js:"coldstart"`
	Dismissed  bool `js:"dismissed"`
}

type Notification struct {
	*js.Object
	Message        string        `js:"message"`
	Title          string        `js:"title"`
	Count          string        `js:"count"`
	Sound          string        `js:"sound"`
	Image          string        `js:"image"`
	LaunchArgs     string        `js:"launchArgs"`
	AdditionalData *NotifExtData `js:"additionalData"`
}

type NotifError struct {
	*js.Object
	Mesage string `js:"message"`
}

type Push struct {
	*js.Object
}

func NewConfig() *Config {
	cfg := &Config{Object: js.Global.Get("Object").New()}
	cfg.Android = &AndroidCfg{Object: js.Global.Get("Object").New()}
	cfg.IOS = &IOSCfg{Object: js.Global.Get("Object").New()}
	cfg.Browser = &BrowserCfg{Object: js.Global.Get("Object").New()}
	return cfg
}

func New(cfg *Config) *Push {
	obj := js.Global.Get("PushNotification").Call("init", cfg)
	return &Push{Object: obj}
}

func HasPermission() (res bool) {
	ch := make(chan struct{})
	success := func(obj *js.Object) {
		res = obj.Get("isEnabled").Bool()
		close(ch)
	}
	js.Global.Get("PushNotification").Call("hasPermission", success)
	<-ch
	return
}

func (p *Push) OnRegistration(f func(*RegInfo)) {
	p.Call("on", "registration", f)
}

func (p *Push) OnNotification(f func(*Notification)) {
	p.Call("on", "notification", f)
}

func (p *Push) OnError(f func(*NotifError)) {
	p.Call("on", "error", f)
}

func (p *Push) OffRegistration(f func(*RegInfo)) {
	p.Call("off", "registration", f)
}

func (p *Push) OffNotification(f func(*Notification)) {
	p.Call("off", "notification", f)
}

func (p *Push) OffError(f func(*NotifError)) {
	p.Call("off", "error", f)
}

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

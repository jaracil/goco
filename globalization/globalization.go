// Package globalization is a GopherJS wrapper for cordova cordova-plugin-globalization.
//
// Install plugin:
//  cordova plugin add cordova-plugin-globalization
//
// (Noncomplete implementation)
package globalization

import (
	"errors"

	"github.com/gopherjs/gopherjs/js"
	"github.com/jaracil/goco"
)

var mo *js.Object

func init() {
	goco.OnDeviceReady(func() {
		mo = js.Global.Get("navigator").Get("globalization")
	})
}

// GetPreferredLanguage downloads an update
func GetPreferredLanguage() (lang string, err error) {
	if mo == nil {
		return "", errors.New("Cannot find navigator.globalization object")
	}

	ch := make(chan struct{})
	mo.Call("getPreferredLanguage",
		func(cbLang *js.Object) {
			lang = cbLang.Get("value").String()
			close(ch)
		},
		func() {
			err = errors.New("Error getting preferred language")
			close(ch)
		})
	<-ch

	return lang, err
}

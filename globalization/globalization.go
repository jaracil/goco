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
)

// DatePatternInfo is info returned on GetDatePattern
type DatePatternInfo struct {
	Pattern      string
	Timezone     string
	IANATimezone string
	UTCOffset    int
	DSTOffset    int
}

var instance *js.Object

func mo() *js.Object {
	if instance == nil {
		instance = js.Global.Get("navigator").Get("globalization")
	}
	return instance
}

// GetPreferredLanguage downloads an update
func GetPreferredLanguage() (lang string, err error) {
	if mo() == nil {
		return "", errors.New("Cannot find navigator.globalization object")
	}

	ch := make(chan struct{})
	mo().Call("getPreferredLanguage",
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

// GetDatePattern returns info about date formating
func GetDatePattern() (info *DatePatternInfo, err error) {
	if mo() == nil {
		return nil, errors.New("Cannot find navigator.globalization object")
	}

	ch := make(chan struct{})
	mo().Call("getDatePattern",
		func(jsInfo *js.Object) {
			info = &DatePatternInfo{
				Pattern:      jsInfo.Get("pattern").String(),
				Timezone:     jsInfo.Get("timezone").String(),
				IANATimezone: jsInfo.Get("iana_timezone").String(),
				UTCOffset:    jsInfo.Get("utc_offset").Int(),
				DSTOffset:    jsInfo.Get("dst_offset").Int(),
			}
			close(ch)
		},
		func() {
			err = errors.New("Error getting date pattern")
			close(ch)
		})
	<-ch

	return info, err
}

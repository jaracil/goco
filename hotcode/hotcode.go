// Package hotcode is a GopherJS wrapper for cordova cordova-hot-code-push.
//
// Install plugin:
//  cordova plugin add cordova-hot-code-push
//
// (Noncomplete implementation)
package hotcode

import (
	"errors"
	"fmt"

	"github.com/gopherjs/gopherjs/js"
	"github.com/jaracil/goco"
)

// Options used for FetchUpdate
type Options struct {
	*js.Object
	ConfigFile     string            `js:"configFile"`
	RequestHeaders map[string]string `js:"requestHeaders"`
}

var mo *js.Object

func init() {
	goco.OnDeviceReady(func() {
		mo = js.Global.Get("chcp")
	})
}

// NewOptions creates options object for FetchUpdate
func NewOptions() *Options {
	return &Options{Object: js.Global.Get("Object").New()}
}

// FetchUpdate downloads an update
func FetchUpdate(opts ...*Options) (data interface{}, err error) {
	if mo == nil {
		return nil, errors.New("Cannot find chcp object")
	}

	m := map[string]interface{}{}

	if len(opts) > 0 {
		if opts[0].Get("configFile") != nil {
			m["config-file"] = opts[0].Get("configFile")
		}

		if opts[0].Get("requestHeaders") != nil {
			m["request-headers"] = opts[0].Get("requestHeaders")
		}
	}

	ch := make(chan struct{})
	mo.Call("fetchUpdate", func(cbErr, cbData *js.Object) {
		if cbErr != nil {
			err = fmt.Errorf("Error on fetchUpdate: code=%v, description=%v", cbErr.Get("code"), cbErr.Get("description"))
		}
		data = cbData
		close(ch)
	}, m)
	<-ch
	return data, err
}

// InstallUpdate installs downloaded update
func InstallUpdate() (err error) {
	if mo == nil {
		return errors.New("Cannot find chcp object")
	}

	ch := make(chan struct{})
	mo.Call("installUpdate", func(cbErr *js.Object) {
		if cbErr != nil {
			err = fmt.Errorf("Error on installUpdate: code=%v, description=%v", cbErr.Get("code"), cbErr.Get("description"))
		}
		close(ch)
	})
	<-ch
	return err
}

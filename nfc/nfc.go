// Package nfc is a GopherJS wrapper for cordova nfc plugin.
//
// (INCOMPLETE IMPLEMENTATION!!!)
//
// Install plugin:
//  cordova plugin add phonegap-nfc
package nfc

import (
	"fmt"

	"bitbucket.org/garagemakers/virkey-cloud/frontend/utils"

	"github.com/gopherjs/gopherjs/js"
)

// Tag is an NFC tag
type Tag struct {
	o *js.Object
}

// GetURL returns tag's URL
func (t *Tag) GetURL() string {
	return moNDEF().Get("uriHelper").Call("decodePayload", t.o.Get("ndefMessage").Index(0).Get("payload")).String()
}

var instanceNFC *js.Object
var instanceNDEF *js.Object

func moNFC() *js.Object {
	fmt.Printf("nfc.moNFC")
	if instanceNFC == nil {
		instanceNFC = js.Global.Get("nfc")
	}
	return instanceNFC
}

func moNDEF() *js.Object {
	fmt.Printf("nfc.moNDEF")
	if instanceNDEF == nil {
		instanceNDEF = js.Global.Get("ndef")
	}
	return instanceNDEF
}

// AddMimeTypeListener registers an event listener for NDEF tags matching a specified MIME type.
func AddMimeTypeListener(mimeType string, callback func(*Tag)) (err error) {
	ch := make(chan struct{})
	success := func() {
		close(ch)
	}
	fail := func(obj *js.Object) {
		err = fmt.Errorf("Error registering mimetype listener: %s", utils.Stringify(obj))
		close(ch)
	}
	cb := func(obj *js.Object) {
		tag := Tag{o: obj.Get("tag")}
		callback(&tag)
	}

	moNFC().Call("addMimeTypeListener", mimeType, cb, success, fail)
	<-ch
	return
}

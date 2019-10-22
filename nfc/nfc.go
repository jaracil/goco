// Package nfc is a GopherJS wrapper for cordova nfc plugin.
//
// (INCOMPLETE IMPLEMENTATION!!!)
//
// Install plugin:
//  cordova plugin add phonegap-nfc
package nfc

import (
	"fmt"
	"log"

	"github.com/gopherjs/gopherjs/js"
)

// Status is the status of the NFC support
type Status string

const (
	// NfcEnabled when NFC is enabled
	NfcEnabled Status = "ENABLED"
	// NoNfc the device doesn't support NFC
	NoNfc Status = "NO_NFC"
	// NfcDisabled if the user has disabled NFC
	NfcDisabled Status = "NFC_DISABLED"
	// NoNfcOrNfcDisabled when NFC is not present or disabled (on Windows)
	NoNfcOrNfcDisabled Status = "NO_NFC_OR_NFC_DISABLED"
)

// Tag is an NFC tag
type Tag struct {
	o *js.Object
}

// GetURL returns tag's URL
func (t *Tag) GetURL() string {
	return moNDEF().Get("uriHelper").Call("decodePayload", t.o.Get("ndefMessage").Index(0).Get("payload")).String()
}

func stringify(obj *js.Object) string {
	return js.Global.Get("JSON").Call("stringify", obj).String()
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
		err = fmt.Errorf("Error registering mimetype listener: %s", stringify(obj))
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

// AddNdefListener registers an event listener for any NDEF tag
func AddNdefListener(callback func(*Tag)) (err error) {
	ch := make(chan struct{})
	success := func() {
		close(ch)
	}
	fail := func(obj *js.Object) {
		err = fmt.Errorf("Error registering NDEF listener: %s", stringify(obj))
		close(ch)
	}
	cb := func(obj *js.Object) {
		log.Printf("AddNdefListener cb")
		tag := Tag{o: obj.Get("tag")}
		callback(&tag)
	}

	moNFC().Call("addNdefListener", cb, success, fail)
	<-ch
	return
}

// BeginSession starts scanning for NFC tags (needed on iOS)
func BeginSession() (err error) {
	ch := make(chan struct{})
	success := func() {
		close(ch)
	}
	fail := func(obj *js.Object) {
		err = fmt.Errorf("Error beggining NFC session: %s", stringify(obj))
		close(ch)
	}

	moNFC().Call("beginSession", success, fail)
	<-ch
	return
}

// InvalidateSession stops scanning for NFC tags (needed on iOS)
func InvalidateSession() (err error) {
	ch := make(chan struct{})
	success := func() {
		close(ch)
	}
	fail := func(obj *js.Object) {
		err = fmt.Errorf("Error stopping NFC session: %s", stringify(obj))
		close(ch)
	}

	moNFC().Call("invalidateSession", success, fail)
	<-ch
	return
}

// Enabled returns status of NFC
func Enabled() (status Status) {
	ch := make(chan struct{})
	success := func() {
		status = NfcEnabled
		close(ch)
	}
	fail := func(value string) {
		status = Status(value)
		close(ch)
	}

	moNFC().Call("enabled", success, fail)
	<-ch
	return
}

// ShowSettings shows the NFC settings on the device
func ShowSettings() (err error) {
	ch := make(chan struct{})
	success := func() {
		close(ch)
	}
	fail := func(obj *js.Object) {
		err = fmt.Errorf("Error showing settings: %s", stringify(obj))
		close(ch)
	}

	moNFC().Call("showSettings", success, fail)
	<-ch
	return
}

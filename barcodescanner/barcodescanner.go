// Package barcodescanner is a GopherJS wrapper for barcodescanner plugin.
//
// Install plugin:
//  cordova plugin add cordova-plugin-barcodescanner
//  phonegap plugin add phonegap-plugin-barcodescanner
//
// (Incomplete implementation, missing "Encode" function)
package barcodescanner

import (
	"github.com/gopherjs/gopherjs/js"
	"errors"
	"strings"
)

var (
	ErrCanceled           = errors.New("BarcodeScanner error: User has canceled the scanner")
	ErrUnexpectedError    = errors.New("BarcodeScanner error: Unexpected error")
	ErrPluginNotAvailable = errors.New("BarcodeScanner error: BarcodeScanner not found")
)

type Format string

const (
	QR_CODE      Format = "QR_CODE"
	DATA_MATRIX  Format = "DATA_MATRIX"
	UPC_A        Format = "UPC_A"
	UPC_E        Format = "UPC_E"
	EAN_8        Format = "EAN_8"
	EAN_13       Format = "EAN_13"
	CODE_39      Format = "CODE_39"
	CODE_93      Format = "CODE_93"
	CODE_128     Format = "CODE_128"
	CODABAR      Format = "CODABAR"
	ITF          Format = "ITF"
	RSS_14       Format = "RSS14"
	PDF_417      Format = "PDF417"
	RSS_EXPANDED Format = "RSS_EXPANDED"
	MSI          Format = "MSI"
	AZTEC        Format = "AZTEC"
)

type Orientation string

const (
	Landscape Orientation = "landscape"
	Portrait  Orientation = "portrait"
	Both      Orientation = ""
)

type Scanner struct {
	*js.Object
	PreferFrontCamera     bool        `js:"preferFrontCamera"`     // iOS/Android
	ShowFlipCameraButton  bool        `js:"showFlipCameraButton"`  // iOS/Android
	ShowTorchButton       bool        `js:"showTorchButton"`       // iOS/Android
	TorchOn               bool        `js:"torchOn"`               // Android
	SaveHistory           bool        `js:"saveHistory"`           // Android
	Prompt                string      `js:"prompt"`                // Android
	ResultDisplayDuration int         `js:"resultDisplayDuration"` // Android
	Formats               Format      `js:"formats"`               // See https://www.npmjs.com/package/cordova-plugin-barcodescanner
	Orientation           Orientation `js:"orientation"`           // Android
	DisableAnimations     bool        `js:"disableAnimations"`     // iOS
	DisableSuccessBeep    bool        `js:"disableSuccessBeep"`
}

type Response struct {
	*js.Object
	Text   string `js:"text"`
	Format Format `js:"format"`
}

var instance *js.Object

func mo() *js.Object {
	if instance == nil {
		instance = js.Global.Get("cordova").Get("plugins").Get("barcodeScanner")
	}
	return instance
}

// NewScanner returns new Scanner object with default values
func NewScanner() *Scanner {
	cfg := &Scanner{Object:  js.Global.Get("Object").New()}
	cfg.SetFormats(QR_CODE)

	return cfg
}

// SetFormats defines which format should be decode, it can be multiple formats like `SetFormats(QR_CODE, EAN_13)`.
func (cfg *Scanner) SetFormats(format ...Format) {
	var str strings.Builder

	for _, f := range format {
		str.WriteString(",")
		str.WriteString(string(f))
	}

	cfg.Formats = Format(str.String()[1:])
}

// Scan initialize the scanner, allowing the user to scan the barcode using the camera. It will return a non-nil error
// when is impossible the use the plugin or use the camera or impossible to decode the barcode, if the user cancel the
// scanner a non-nil error are also given.
func (cfg *Scanner) Scan() (resp *Response, err error) {
	if mo() == js.Undefined || mo() == nil {
		return resp, ErrPluginNotAvailable
	}

	ch := make(chan struct{})
	success := func(obj *js.Object) {
		resp = &Response{Object: obj}

		if obj.Get("cancelled").Bool() {
			err = ErrCanceled
		}

		close(ch)
	}
	fail := func(obj *js.Object) {
		err = ErrUnexpectedError
		close(ch)
	}
	mo().Call("scan", success, fail, cfg)
	<-ch
	return
}

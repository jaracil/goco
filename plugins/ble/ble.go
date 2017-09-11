package ble

import (
	"errors"

	"github.com/gopherjs/gopherjs/js"
)

type AdvData struct {
	*js.Object
}

type ScanData struct {
	*js.Object
	Name        string  `js:"Name"`
	ID          string  `js:"id"`
	Advertising AdvData `js:"advertising"`
	Rssi        int     `js:"rssi"`
}

var (
	wantScan = false
	scaning  = false
)

func startScan(srv []string, cbFun func(*ScanData), dups bool) {
	if scaning {
		return
	}
	cb := func(obj *js.Object) {
		cbFun(&ScanData{Object: obj})
	}
	options := map[string]interface{}{"reportDuplicates": dups}
	js.Global.Get("ble").Call("startScanWithOptions", srv, options, cb)
	scaning = true
}

func StartScan(srv []string, cbFun func(*ScanData), dups bool) {
	wantScan = true
	startScan(srv, cbFun, dups)
}

func stopScan() (err error) {
	if !scaning {
		return
	}
	ch := make(chan struct{})
	success := func() {
		scaning = false
		close(ch)
	}
	failure := func() {
		err = errors.New("Error closing BLE scan")
		close(ch)
	}
	js.Global.Get("ble").Call("stopScan", success, failure)
	<-ch
	return
}

func StopScan() (err error) {
	wantScan = false
	return stopScan()
}
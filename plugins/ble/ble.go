package ble

import (
	"errors"

	"github.com/gopherjs/gopherjs/js"
)

var (
	// Status flags
	wantScan   = false
	scaning    = false
	connecting = 0

	// Scan Params
	scanSrv   []string
	scanCbFun func(*Peripheral)
	scanDups  bool
)

func startScan(srv []string, cbFun func(*Peripheral), dups bool) {
	if scaning {
		return
	}
	cb := func(obj *js.Object) {
		if cbFun != nil {
			cbFun(&Peripheral{o: obj})
		}
	}
	options := map[string]interface{}{"reportDuplicates": dups}
	js.Global.Get("ble").Call("startScanWithOptions", srv, options, cb)
	scaning = true
}

func resumeScan() {
	if wantScan && connecting == 0 {
		startScan(scanSrv, scanCbFun, scanDups)
	}
}

func StartScan(srv []string, cbFun func(*Peripheral), dups bool) {
	scanSrv = srv
	scanCbFun = cbFun
	scanDups = dups
	wantScan = true
	resumeScan()
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

func Connect(id string, endConnCb func(per *Peripheral)) (per *Peripheral, err error) {
	ch := make(chan struct{})
	connected := false
	success := func(obj *js.Object) {
		per = &Peripheral{o: obj}
		connected = true
		close(ch)
	}
	failure := func(obj *js.Object) {
		if connected {
			if endConnCb != nil {
				endConnCb(&Peripheral{o: obj})
			}
		} else {
			err = errors.New("Error connecting to BLE peripheral")
			close(ch)
		}
	}
	connecting++
	stopScan()
	js.Global.Get("ble").Call("connect", id, success, failure)
	<-ch
	connecting--
	resumeScan()
	return
}

func Disconnect(id string) (err error) {
	ch := make(chan struct{})
	success := func() {
		close(ch)
	}
	failure := func() {
		err = errors.New("Error closing BLE peripheral: " + id)
		close(ch)
	}
	js.Global.Get("ble").Call("disconnect", id, success, failure)
	<-ch
	return
}

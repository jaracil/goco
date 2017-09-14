package ble

import (
	"errors"
	"time"

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

func stringify(obj *js.Object) string {
	return js.Global.Get("JSON").Call("stringify", obj).String()
}

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
	failure := func(obj *js.Object) {
		err = errors.New("Error closing BLE scan <" + stringify(obj) + ">")
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
	if !IsEnabled() {
		return nil, errors.New("Bluetooth disabled")
	}
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
	failure := func(obj *js.Object) {
		err = errors.New("Error closing BLE peripheral: <" + stringify(obj) + ">")
		close(ch)
	}
	js.Global.Get("ble").Call("disconnect", id, success, failure)
	<-ch
	return
}

func IsConnected(id string) (ret bool) {
	ch := make(chan struct{})
	success := func() {
		ret = true
		close(ch)
	}
	failure := func() {
		ret = false
		close(ch)
	}
	js.Global.Get("ble").Call("isConnected", id, success, failure)
	<-ch
	return
}

func Read(id, srv, char string) (ret []byte, err error) {
	ch := make(chan struct{})
	success := func(obj *js.Object) {
		ret = js.Global.Get("Uint8Array").New(obj).Interface().([]byte)
		close(ch)
	}
	failure := func(obj *js.Object) {
		err = errors.New("BLE read error: <" + stringify(obj) + ">")
		close(ch)
	}
	js.Global.Get("ble").Call("read", id, srv, char, success, failure)
	<-ch
	return
}

func Write(id, srv, char string, data []byte) (err error) {
	ch := make(chan struct{})
	success := func() {
		close(ch)
	}
	failure := func(obj *js.Object) {
		err = errors.New("BLE write error: <" + stringify(obj) + ">")
		close(ch)
	}
	arr := js.NewArrayBuffer(data)
	js.Global.Get("ble").Call("write", id, srv, char, arr, success, failure)
	<-ch
	return
}

func WriteWithoutResponse(id, srv, char string, data []byte) (err error) {
	ch := make(chan struct{})
	success := func() {
		close(ch)
	}
	failure := func(obj *js.Object) {
		err = errors.New("BLE write error: <" + stringify(obj) + ">")
		close(ch)
	}
	arr := js.NewArrayBuffer(data)
	js.Global.Get("ble").Call("writeWithoutResponse", id, srv, char, arr, success, failure)
	<-ch
	return
}

func StartNotification(id, srv, char string, recvCb func([]byte)) (err error) {
	success := func(obj *js.Object) {
		if recvCb != nil {
			recvCb(js.Global.Get("Uint8Array").New(obj).Interface().([]byte))
		}
	}
	failure := func(obj *js.Object) {
		err = errors.New("BLE start notifications error: <" + stringify(obj) + ">")
	}
	js.Global.Get("ble").Call("startNotification", id, srv, char, success, failure)
	time.Sleep(10 * time.Millisecond) // Dirty Hack: Wait for eventual failure callback
	return
}

func StopNotification(id, srv, char string) (err error) {
	ch := make(chan struct{})
	success := func() {
		close(ch)
	}
	failure := func(obj *js.Object) {
		err = errors.New("BLE stop notifications error: <" + stringify(obj) + ">")
		close(ch)
	}
	js.Global.Get("ble").Call("stopNotification", id, srv, char, success, failure)
	<-ch
	return
}

func IsEnabled() (ret bool) {
	ch := make(chan struct{})
	success := func() {
		ret = true
		close(ch)
	}
	failure := func() {
		ret = false
		close(ch)
	}
	js.Global.Get("ble").Call("isEnabled", success, failure)
	<-ch
	return
}

func Enable() (err error) {
	ch := make(chan struct{})
	success := func() {
		close(ch)
	}
	failure := func(obj *js.Object) {
		err = errors.New("Can't enable bluetooth: <" + stringify(obj) + ">")
		close(ch)
	}
	js.Global.Get("ble").Call("enable", success, failure)
	<-ch
	return
}

func StartStateNotifications(recvCb func(string)) (err error) {
	success := func(obj *js.Object) {
		if recvCb != nil {
			recvCb(obj.String())
		}
	}
	js.Global.Get("ble").Call("startStateNotifications", success)
	return
}

func StopStateNotifications() (err error) {
	ch := make(chan struct{})
	success := func() {
		close(ch)
	}
	failure := func(obj *js.Object) {
		err = errors.New("Can't stop bluetooth state notifications: <" + stringify(obj) + ">")
		close(ch)
	}
	js.Global.Get("ble").Call("stopStateNotifications", success, failure)
	<-ch
	return
}

func ShowBluetoothSettings() (err error) {
	ch := make(chan struct{})
	success := func() {
		close(ch)
	}
	failure := func(obj *js.Object) {
		err = errors.New("Can't show bluetooth setings: <" + stringify(obj) + ">")
		close(ch)
	}
	js.Global.Get("ble").Call("showBluetoothSettings", success, failure)
	<-ch
	return
}

func ReadRSSI(id string) (val int, err error) {
	ch := make(chan struct{})
	success := func(obj *js.Object) {
		val = obj.Int()
		close(ch)
	}
	failure := func(obj *js.Object) {
		err = errors.New("Can't get device RSSI: <" + stringify(obj) + ">")
		close(ch)
	}
	js.Global.Get("ble").Call("readRSSI", id, success, failure)
	<-ch
	return
}

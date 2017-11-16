package ble

import (
	"errors"
	"time"

	"github.com/gopherjs/gopherjs/js"
	"github.com/jaracil/goco"
)

var (
	// Status flags
	wantScan = false
	scaning  = false
	paused   = 0

	// Scan Params
	scanSrv   []string
	scanCbFun func(*Peripheral)
	scanDups  bool
)

var mo *js.Object

func init() {
	goco.OnDeviceReady(func() {
		mo = js.Global.Get("ble")
	})
}

func stringify(obj *js.Object) string {
	return js.Global.Get("JSON").Call("stringify", obj).String()
}

func PauseScan() {
	paused++
	stopScan()
}

func ResumeScan() {
	if paused > 0 {
		paused--
	}
	if wantScan && paused == 0 {
		startScan(scanSrv, scanCbFun, scanDups)
	}
}

func Scaning() bool {
	return scaning
}

func startScan(srv []string, cbFun func(*Peripheral), dups bool) {
	if scaning || !wantScan {
		return
	}
	scaning = true
	options := map[string]interface{}{"reportDuplicates": dups}
	mo.Call("startScanWithOptions", srv, options, func(p *Peripheral) {
		p.Parse()
		cbFun(p)
	})
	print("Scan started!!!")
}

func StartScan(srv []string, cbFun func(*Peripheral), dups bool) {
	scanSrv = srv
	scanCbFun = cbFun
	scanDups = dups
	wantScan = true
	if paused == 0 {
		startScan(scanSrv, scanCbFun, scanDups)
	}

}

func stopScan() (err error) {
	if !scaning {
		return
	}
	scaning = false
	ch := make(chan struct{})
	success := func() {
		close(ch)
	}
	failure := func(obj *js.Object) {
		err = errors.New("Error closing BLE scan <" + stringify(obj) + ">")
		close(ch)
	}
	mo.Call("stopScan", success, failure)
	<-ch
	print("Scan stopped!!!")
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
	success := func(p *Peripheral) {
		per = p
		connected = true
		close(ch)
	}
	failure := func(obj *js.Object) {
		if connected {
			if endConnCb != nil {
				endConnCb(&Peripheral{Object: obj})
			}
		} else {
			err = errors.New("Error connecting to BLE peripheral")
			close(ch)
		}
	}
	PauseScan()
	mo.Call("connect", id, success, failure)
	<-ch
	ResumeScan()
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
	mo.Call("disconnect", id, success, failure)
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
	mo.Call("isConnected", id, success, failure)
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
	mo.Call("read", id, srv, char, success, failure)
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
	mo.Call("write", id, srv, char, arr, success, failure)
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
	mo.Call("writeWithoutResponse", id, srv, char, arr, success, failure)
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
	mo.Call("startNotification", id, srv, char, success, failure)
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
	mo.Call("stopNotification", id, srv, char, success, failure)
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
	mo.Call("isEnabled", success, failure)
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
	mo.Call("enable", success, failure)
	<-ch
	return
}

func StartStateNotifications(recvCb func(string)) (err error) {
	mo.Call("startStateNotifications", recvCb)
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
	mo.Call("stopStateNotifications", success, failure)
	<-ch
	return
}

func ShowBluetoothSettings() (err error) {
	ch := make(chan struct{})
	success := func() {
		close(ch)
	}
	failure := func(obj *js.Object) {
		err = errors.New("Can't show bluetooth settings: <" + stringify(obj) + ">")
		close(ch)
	}
	mo.Call("showBluetoothSettings", success, failure)
	<-ch
	return
}

func ReadRSSI(id string) (val int, err error) {
	ch := make(chan struct{})
	success := func(n int) {
		val = n
		close(ch)
	}
	failure := func(obj *js.Object) {
		err = errors.New("Can't get device RSSI: <" + stringify(obj) + ">")
		close(ch)
	}
	mo.Call("readRSSI", id, success, failure)
	<-ch
	return
}

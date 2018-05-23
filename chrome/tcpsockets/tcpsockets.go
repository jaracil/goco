// Package tcpsockets is a GopherJS wrapper for chrome.sockets.tcp plugin
// https://www.npmjs.com/package/cordova-plugin-chrome-apps-sockets-tcp
// https://developer.chrome.com/apps/sockets_tcp
//
// Install plugin:
//  cordova plugin add cordova-plugin-chrome-apps-sockets-tcp
//
// Compatible with "net.Conn" interface.
//
// (Incomplete implementation, missing "setPaused", "setKeepAlive", "setNoDelay", "getSockets)
package tcpsockets

import (
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/gopherjs/gopherjs/js"
)

var (
	ErrPluginNotFound   = errors.New("chrome.sockets.tcp error: Plugin not found")
	ErrConnectionFailed = errors.New("chrome.sockets.tcp error: Cannot connect")
	ErrReadFailed       = errors.New("chrome.sockets.tcp error: Cannot read")
)

type conn struct {
	socketID    int
	ipport      string
	readCh      chan []byte
	readPending chan []byte
	readErrorCh chan int
}

type addr struct {
	ipport string
}

func (a addr) Network() string { return "tcp" }
func (a addr) String() string  { return a.ipport }

var instance *js.Object

func mo() *js.Object {
	if instance == nil {
		instance = js.Global.Get("chrome").Get("sockets").Get("tcp")
	}
	return instance
}

// Create func creates a TCP socket.
func Create() (conn, error) {
	if mo() == js.Undefined || mo() == nil {
		return conn{}, ErrPluginNotFound
	}

	ch := make(chan int)
	infoCallback := func(obj *js.Object) {
		go func() { ch <- obj.Get("socketId").Int() }()
	}

	mo().Call("create", infoCallback)
	socketID := <-ch
	return conn{socketID: socketID}, nil
}

func (c *conn) Connect(peerAddress string, peerPort int) error {
	ch := make(chan int)
	connCallback := func(obj *js.Object) {
		ch <- obj.Int()
	}

	c.ipport = fmt.Sprintf("%s:%d", peerAddress, peerPort)
	mo().Call("connect", c.socketID, peerAddress, peerPort, connCallback)

	readCallback := func(obj *js.Object) {
		var b []byte
		b = js.Global.Get("Uint8Array").New(obj.Get("data")).Interface().([]byte)
		go func() { c.readCh <- b }()
	}

	readErrorCallback := func(obj *js.Object) {
		var resultCode int
		resultCode = obj.Get("data").Int()
		go func() { c.readErrorCh <- resultCode }()
	}

	result := <-ch
	if result >= 0 {
		c.readCh = make(chan []byte)
		c.readPending = make(chan []byte)
		mo().Get("onReceive").Call("addListener", readCallback)
		c.readErrorCh = make(chan int)
		mo().Get("onReceiveError").Call("addListener", readErrorCallback)
		return nil
	}

	return ErrConnectionFailed
}

func (c conn) Close() error {
	mo().Call("disconnect", c.socketID)
	return nil
}

func (c conn) Write(b []byte) (n int, err error) {
	type result struct {
		res, bytes int
	}
	ch := make(chan result)
	connCallback := func(obj *js.Object) {
		go func() { ch <- result{res: obj.Get("resultCode").Int(), bytes: obj.Get("bytesSent").Int()} }()
	}
	mo().Call("send", c.socketID, js.NewArrayBuffer(b), connCallback)

	res := <-ch
	if res.res >= 0 {
		return res.bytes, nil
	}
	return 0, errors.New(fmt.Sprintf("send error %d", res.res))
}

func (c conn) Read(receive []byte) (n int, err error) {
	select {
	case b := <-c.readPending:
		if len(b) > len(receive) {
			extra := b[len(receive):]
			b = b[:len(receive)]
			go func() {
				c.readPending <- extra
			}()
		}
		for i, x := range b {
			receive[i] = x
		}
		return len(b), nil
	default:
	}

	select {
	case b := <-c.readCh:
		if len(b) > len(receive) {
			extra := b[len(receive):]
			b = b[:len(receive)]
			go func() {
				c.readPending <- extra
			}()
		}
		for i, x := range b {
			receive[i] = x
		}
		return len(b), nil
	case resCode := <-c.readErrorCh:
		return 0, errors.New(fmt.Sprintf("recieve error %d", resCode))
	}

	return 0, ErrReadFailed
}

func (c conn) LocalAddr() net.Addr  { return addr{} }
func (c conn) RemoteAddr() net.Addr { return addr{ipport: c.ipport} }

func (c conn) SetDeadline(time.Time) error      { return nil }
func (c conn) SetReadDeadline(time.Time) error  { return nil }
func (c conn) SetWriteDeadline(time.Time) error { return nil }

// Update func updates the socket properties.
func (c conn) Update(socketID int, properties interface{}, cb func()) {
	mo().Call("update", socketID, properties, cb)
}

func (c conn) SetPaused(paused bool) {
	mo().Call("paused", c.socketID, paused)
}

func (c conn) GetInfo() interface{} {
	ch := make(chan interface{})

	infoCallback := func(obj *js.Object) {
		go func() { ch <- obj.Get("socketInfo").Interface() }()
	}
	mo().Call("getInfo", c.socketID, infoCallback)
	res := <-ch
	return res
}

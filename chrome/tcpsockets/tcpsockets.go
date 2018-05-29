// Package tcpsockets is a GopherJS wrapper for chrome.sockets.tcp plugin
// https://www.npmjs.com/package/cordova-plugin-chrome-apps-sockets-tcp
// https://developer.chrome.com/apps/sockets_tcp
//
// Install plugin:
//  cordova plugin add cordova-plugin-chrome-apps-sockets-tcp
//
// Compatible with "net.Conn" interface.
//
// (Incomplete implementation, missing "setNoDelay", "getSockets", "secures")
package tcpsockets

import (
	"fmt"
	"io"
	"net"
	"time"

	"github.com/gopherjs/gopherjs/js"
	"github.com/pkg/errors"
)

var (
	ErrPluginNotFound   = errors.New("chrome.sockets.tcp error: Plugin not found")
	ErrConnectionFailed = errors.New("chrome.sockets.tcp error: Cannot connect")
	ErrReadFailed       = errors.New("chrome.sockets.tcp error: Cannot read")
	ErrKeepAlive        = errors.New("chrome.sockets.tcp error: Keep alive error")
)

type PipeStruct struct {
	w *io.PipeWriter
	r *io.PipeReader
}

type conn struct {
	socketID int
	ipport   string
	readPipe PipeStruct
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

	readCh := make(chan []byte, 100)
	readErrorCh := make(chan int)
	readCallback := func(obj *js.Object) {
		res := js.Global.Get("Uint8Array").New(obj.Get("data")).Interface().([]byte)
		select {
		case readCh <- res:
		default:
		}
	}

	readErrorCallback := func(obj *js.Object) {
		resultCode := obj.Get("data").Int()
		go func() { readErrorCh <- resultCode }()
	}

	result := <-ch
	if result >= 0 {
		c.readPipe.r, c.readPipe.w = io.Pipe()
		mo().Get("onReceive").Call("addListener", readCallback)
		mo().Get("onReceiveError").Call("addListener", readErrorCallback)

		go func() {
		loop:
			for {
				select {
				case receive, ok := <-readCh:
					if _, err := c.readPipe.w.Write(receive); err != nil && ok {
						println("chrome.sockets.tcp error: pipe write error: ", err)
					}
				case resultCode := <-readErrorCh:
					close(readCh)
					close(readErrorCh)
					c.readPipe.w.Close()
					println("chrome.sockets.tcp error: reading error ", resultCode)
					break loop
				}
			}
		}()
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
	return 0, errors.New(fmt.Sprintf("chrome.sockets.tcp error: Send error %d", res.res))
}

func (c conn) Read(receive []byte) (n int, err error) {
	n, err = c.readPipe.r.Read(receive)
	return n, errors.Wrap(err, ErrReadFailed.Error())
}

func (c conn) LocalAddr() net.Addr  { return addr{} }
func (c conn) RemoteAddr() net.Addr { return addr{ipport: c.ipport} }

func (c conn) SetDeadline(time.Time) error      { return nil }
func (c conn) SetReadDeadline(time.Time) error  { return nil }
func (c conn) SetWriteDeadline(time.Time) error { return nil }

func (c conn) Update(socketID int, properties interface{}, cb func()) {
	mo().Call("update", socketID, properties, cb)
}

func (c conn) SetPaused(paused bool) {
	mo().Call("setPaused", c.socketID, paused)
}

func (c conn) SetKeepAlive(enable bool, delaySeconds int) (int, error) {
	ch := make(chan int)
	keepAliveCallback := func(obj *js.Object) {
		go func() { ch <- obj.Get("result").Int() }()
	}

	mo().Call("setKeepAlive", c.socketID, enable, delaySeconds, keepAliveCallback)

	res := <-ch
	if res >= 0 {
		return res, nil
	}
	return -1, ErrKeepAlive
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

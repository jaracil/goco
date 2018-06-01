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
	ErrPluginNotFound       = errors.New("chrome.sockets.tcp error: Plugin not found")
	ErrConnectionTimedOut   = errors.New("chrome.sockets.tcp error: A connection attempt timed out")
	ErrConnectionClosed     = errors.New("chrome.sockets.tcp error: Connection closed")
	ErrConnectionReset      = errors.New("chrome.sockets.tcp error: Connection reset")
	ErrConnectionRefused    = errors.New("chrome.sockets.tcp error: Connection refused")
	ErrConnectionFailed     = errors.New("chrome.sockets.tcp error: Connection failed")
	ErrNameNotResolved      = errors.New("chrome.sockets.tcp error: The host name could not be resolved")
	ErrInternetDisconnected = errors.New("chrome.sockets.tcp error: The Internet connection has been lost")
	ErrGenericFailure       = errors.New("chrome.sockets.tcp error: A generic failure occurred")
	ErrAlreadyConnected     = errors.New("chrome.sockets.tcp error: The socket is already connected")
	ErrInvalidAddress       = errors.New("chrome.sockets.tcp error: The IP address or port number is invalid")
	ErrUnreachableAddress   = errors.New("chrome.sockets.tcp error: The IP address is unreachable")
	ErrConnectionWasClosed  = errors.New("chrome.sockets.tcp error: A connection was closed")
	ErrPipeReadFailed       = errors.New("chrome.sockets.tcp error: Cannot read")
	ErrKeepAlive            = errors.New("chrome.sockets.tcp error: Keep alive error")
	ErrWsProtocolError      = errors.New("chrome.sockets.tcp error: Websocket protocol error")
	ErrAddressInUse         = errors.New("chrome.sockets.tcp error: Address is already in use")
)

const (
	NotFoundErrVal          = -1
	GenericFailureVal       = -2
	AlreadyConnectedVal     = -23
	ConnectionClosedVal     = -100
	ConnectionResetVal      = -101
	ConnectionRefusedVal    = -102
	ConnectionFailedVal     = -104
	NameNotResolvedVal      = -105
	InternetDisconnectedVal = -106
	ConnectionTimedOutVal   = -118
	InvalidAddressVal       = -108
	UnreachableAddressVal   = -109
	WsProtocolErrorVal      = -145
	AddressInUseVal         = -147
)

func (sTcpError *socketTCPError) Error() error {
	switch sTcpError.ResultCode {
	case NotFoundErrVal:
		return ErrPluginNotFound
	case ConnectionClosedVal:
		return ErrConnectionClosed
	case ConnectionResetVal:
		return ErrConnectionReset
	case ConnectionRefusedVal:
		return ErrConnectionRefused
	case ConnectionFailedVal:
		return ErrConnectionFailed
	case NameNotResolvedVal:
		return ErrNameNotResolved
	case InternetDisconnectedVal:
		return ErrInternetDisconnected
	case ConnectionTimedOutVal:
		return ErrConnectionTimedOut
	case InvalidAddressVal:
		return ErrInvalidAddress
	case UnreachableAddressVal:
		return ErrUnreachableAddress
	case GenericFailureVal:
		return ErrGenericFailure
	case AlreadyConnectedVal:
		return ErrAlreadyConnected
	case WsProtocolErrorVal:
		return ErrWsProtocolError
	case AddressInUseVal:
		return ErrAddressInUse
	}
	return fmt.Errorf("Unknown error: %d", sTcpError.ResultCode)
}

type socketTCPError struct {
	ResultCode int
}

type PipeStruct struct {
	w *io.PipeWriter
	r *io.PipeReader
}

type conn struct {
	socketID    int
	ipport      string
	readPipe    PipeStruct
	socketError error
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

func (c *conn) Connect(peerAddress string, peerPort int) (err error) {
	ch := make(chan int)
	connCallback := func(obj *js.Object) {
		go func() { ch <- obj.Int() }()
	}

	c.ipport = fmt.Sprintf("%s:%d", peerAddress, peerPort)
	mo().Call("connect", c.socketID, peerAddress, peerPort, connCallback)

	readCh := make(chan []byte, 100)
	readCallback := func(obj *js.Object) {
		res := js.Global.Get("Uint8Array").New(obj.Get("data")).Interface().([]byte)
		select {
		case readCh <- res:
		default:
			c.Close()
		}
	}

	readErrorCh := make(chan *socketTCPError, 1)
	readErrorCallback := func(obj *js.Object) {
		resultCode := obj.Get("resultCode").Int()
		socketTcpErr := &socketTCPError{ResultCode: resultCode}
		select {
		case readErrorCh <- socketTcpErr:
		default:
		}
	}

	result := <-ch
	if result >= 0 {
		c.readPipe.r, c.readPipe.w = io.Pipe()
		mo().Get("onReceive").Call("addListener", readCallback)
		mo().Get("onReceiveError").Call("addListener", readErrorCallback)

		go func() {
			exit := false
			for {
				select {
				case receive := <-readCh:
					if _, e := c.readPipe.w.Write(receive); e != nil {
						err = e
						exit = true
					}
				case res := <-readErrorCh:
					err = res.Error()
					c.socketError = res.Error()
					exit = true
					break
				}
				if exit {
					readCh = nil
					readErrorCh = nil
					c.readPipe.w.Close()
					break
				}
			}
		}()
		return err
	}
	return ErrConnectionFailed
}

func (c conn) Close() error {
	mo().Call("disconnect", c.socketID)
	return c.socketError
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
	return n, errors.Wrap(err, ErrPipeReadFailed.Error())
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

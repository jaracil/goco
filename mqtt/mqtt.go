// Package mqtt is a GopherJS wrapper for 'CordovaMqTTPlugin'
//
// Install plugin:
// 	cordova plugin add cordova-plugin-mqtt
//  cordova plugin add https://github.com/arcoirislabs/cordova-plugin-mqtt.git
package mqtt

import (
	"errors"

	"github.com/gopherjs/gopherjs/js"
)

type DisconnectObject struct {
	*js.Object
	Success func(obj *js.Object) `js:"success"`
	Error   func(obj *js.Object) `js:"error"`
}

type Server struct {
	*js.Object
	URL      string               `js:"url"`
	Port     int                  `js:"port"`
	ClientID string               `js:"clientId"`
	Success  func(obj *js.Object) `js:"success"`
	Error    func(obj *js.Object) `js:"error"`
}

type PublishOBject struct {
	*js.Object
	Topic   string               `js:"topic"`
	Payload string               `js:"payload"`
	Qos     int                  `js:"qos"`
	Retain  bool                 `js:"retain"`
	Success func(obj *js.Object) `js:"success"`
	Error   func(obj *js.Object) `js:"error"`
}

type UnsubscribeObject struct {
	*js.Object
	Topic   string               `js:"topic"`
	Success func(obj *js.Object) `js:"success"`
	Error   func(obj *js.Object) `js:"error"`
}

type SusbscribeObject struct {
	*js.Object
	Topic   string               `js:"topic"`
	Qos     int                  `js:"qos"`
	Success func(obj *js.Object) `js:"success"`
	Error   func(obj *js.Object) `js:"error"`
}

var (
	instance        *js.Object
	connectedServer *Server
)

func mo() *js.Object {
	if instance == nil {
		instance = js.Global.Get("cordova").Get("plugins").Get("CordovaMqTTPlugin")
	}
	return instance
}

// Connect connects to a MQTT server. Will return an non-nil error if plugin is not installed or there was an error connecting
func Connect(url string, port int, clientID string) (err error) {
	if mo() == nil || mo() == js.Undefined {
		err = errors.New("Couldn't get 'CordovaMqTTPlugin', make sure plugin is installed")
		return err
	}
	if connectedServer != nil {
		Disconnect()
	}
	server := createServer(url, port, clientID)
	ch := make(chan error)
	server.Success = func(obj *js.Object) {
		connectedServer = server
		close(ch)
	}
	server.Error = func(obj *js.Object) {
		err = errors.New(obj.String())
		close(ch)
	}
	mo().Call("connect", server)
	<-ch
	return
}

func createServer(url string, port int, clientID string) *Server {
	server := &Server{Object: js.Global.Get("Object").New()}
	server.URL = url
	server.Port = port
	server.ClientID = clientID
	return server
}

// Disconnect disconnects from the connected server.
func Disconnect() (err error) {
	if connectedServer != nil {
		ch := make(chan error)
		disc := &DisconnectObject{Object: js.Global.Get("Object").New()}
		disc.Success = func(obj *js.Object) {
			connectedServer = nil
			close(ch)
		}
		disc.Error = func(obj *js.Object) {
			err = errors.New(obj.String())
			close(ch)
		}
		mo().Call("disconnect", disc)
		<-ch
	}
	return
}

// SubscribeTopic subscribes to a topic, and calls the function passed as a parameter every time it reads a value from the topic.
func SubscribeTopic(topic string, qos int, resFunc func(string)) error {
	if connectedServer == nil {
		return errors.New("No server connected")
	}
	if err := subscribeToTopic(topic, qos); err != nil {
		return err
	}
	mo().Call("listen", topic, resFunc)
	return nil
}

func subscribeToTopic(topic string, qos int) (err error) {
	ch := make(chan error)
	subs := &SusbscribeObject{Object: js.Global.Get("Object").New()}
	subs.Topic = topic
	subs.Qos = qos
	subs.Success = func(obj *js.Object) {
		close(ch)
	}
	subs.Error = func(obj *js.Object) {
		err = errors.New(obj.String())
		close(ch)
	}
	mo().Call("subscribe", subs)
	<-ch
	return

}

// Publish sends the data from payload to the topic passed as a parameter on the connected server
func Publish(topic, payload string, qos int, retain bool) (err error) {
	if connectedServer == nil {
		return errors.New("No server connected")
	}
	ch := make(chan error)
	pub := &PublishOBject{Object: js.Global.Get("Object").New()}
	pub.Topic = topic
	pub.Payload = payload
	pub.Qos = qos
	pub.Retain = retain
	pub.Success = func(obj *js.Object) {
		close(ch)
	}
	pub.Error = func(obj *js.Object) {
		err = errors.New(obj.String())
		close(ch)
	}
	mo().Call("publish", pub)
	<-ch
	return
}

// UnsubscribeTopic usubscribes and stops reading values from that topic
func UnsubscribeTopic(topic string) (err error) {
	if connectedServer == nil {
		return errors.New("No server connected")
	}
	ch := make(chan error)
	unsub := &UnsubscribeObject{Object: js.Global.Get("Object").New()}
	unsub.Topic = topic
	unsub.Success = func(obj *js.Object) {
		close(ch)
	}
	unsub.Error = func(obj *js.Object) {
		err = errors.New(obj.String())
		close(ch)
	}
	mo().Call("unsubscribe", unsub)
	return
}

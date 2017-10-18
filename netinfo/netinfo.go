// Package netinfo is a GopherJS wrapper for cordova network-information plugin.
//
// This plugin provides two public vars, "Current" and "Kinds"
// Which contains the device's connection kind and available kinds respectively.
//
// Install plugin:
//  cordova plugin add cordova-plugin-network-information
package netinfo

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/jaracil/goco"
)

// ActualKind wraps cordova navigator.connection.type (use netinfo.Current)
type ActualKind struct {
	*js.Object
	Kind string `js:"type"` // Connection kind (See AvailableKinds type)
}

// AvailableKinds wraps cordova global Connection object (use netinfo.Kinds)
type AvailableKinds struct {
	*js.Object
	Unknown  string `js:"UNKNOWN"`
	Ethernet string `js:"ETHERNET"`
	Wifi     string `js:"WIFI"`
	Cell2G   string `js:"CELL_2G"`
	Cell3G   string `js:"CELL_3G"`
	Cell4G   string `js:"CELL_4G"`
	Cell     string `js:"CELL"`
	None     string `js:"NONE"`
}

// Current contains actual connection kind
var Current *ActualKind

// Kinds contains available connection kinds
var Kinds *AvailableKinds

func init() {
	goco.OnDeviceReady(func() {
		Current = &ActualKind{Object: js.Global.Get("navigator").Get("connection")}
		Kinds = &AvailableKinds{Object: js.Global.Get("Connection")}
	})

}

// IsCell returns when actual connection is of kind Cell, Cell2G, Cell3G, Cell4G
func IsCell() bool {
	return Current.Kind == Kinds.Cell || Current.Kind == Kinds.Cell2G || Current.Kind == Kinds.Cell3G || Current.Kind == Kinds.Cell4G
}

// OnOffline registers callback function that runs when device goes offline
func OnOffline(cb func()) {
	js.Global.Get("document").Call("addEventListener", "offline", cb, false)
}

// OnOnline registers callback function that runs when device goes online
func OnOnline(cb func()) {
	js.Global.Get("document").Call("addEventListener", "online", cb, false)
}

// UnOffline clears the previous OnOffline registration
func UnOffline(cb func()) {
	js.Global.Get("document").Call("removeEventListener", "offline", cb, false)
}

// UnOnline clears the previous OnOnline registration
func UnOnline(cb func()) {
	js.Global.Get("document").Call("removeEventListener", "online", cb, false)
}

package netinfo

import (
	"github.com/gopherjs/gopherjs/js"
)

type ActualKind struct {
	*js.Object
	Kind string `js:"type"`
}
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

var Current *ActualKind
var Kinds *AvailableKinds

func Init() {
	Current = &ActualKind{Object: js.Global.Get("navigator").Get("connection")}
	Kinds = &AvailableKinds{Object: js.Global.Get("Connection")}
}

func IsCell() bool {
	return Current.Kind == Kinds.Cell || Current.Kind == Kinds.Cell2G || Current.Kind == Kinds.Cell3G || Current.Kind == Kinds.Cell4G
}

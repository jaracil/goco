package netinfo

import (
	"github.com/gopherjs/gopherjs/js"
)

func Type() string {
	return js.Global.Get("navigator").Get("connection").Get("type").String()
}

func WifiType() string {
	return js.Global.Get("Connection").Get("WIFI").String()
}

func EthernetType() string {
	return js.Global.Get("Connection").Get("ETHERNET").String()
}

func Cell2gType() string {
	return js.Global.Get("Connection").Get("CELL_2G").String()
}

func Cell3gType() string {
	return js.Global.Get("Connection").Get("CELL_3G").String()
}

func Cell4gType() string {
	return js.Global.Get("Connection").Get("CELL_4G").String()
}

func CellType() string {
	return js.Global.Get("Connection").Get("CELL").String()
}

func NoneType() string {
	return js.Global.Get("Connection").Get("NONE").String()
}

func IsCell() bool {
	t := Type()
	return t == Cell2gType() || t == Cell3gType() || t == Cell4gType() || t == CellType()
}

func OnOnline(fn func()) {
	js.Global.Get("document").Call("addEventListener", "online", fn, false)
}

func OnOffline(fn func()) {
	js.Global.Get("document").Call("addEventListener", "offline", fn, false)
}

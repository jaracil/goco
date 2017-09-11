package netinfo

import (
	"github.com/gopherjs/gopherjs/js"
)

const (
	NoneType     = "NONE"
	EthernetType = "ETHERNET"
	WifiType     = "WIFI"
	CellType     = "CELLULAR"
	Cell2gType   = "2G"
	Cell3gType   = "3G"
	Cell4gType   = "4G"
	Cell5gType   = "5G"
)

func Type() string {
	return js.Global.Get("navigator").Get("connection").Get("type").String()
}

func IsCell() bool {
	t := Type()
	return t == Cell2gType || t == Cell3gType || t == Cell4gType || t == Cell5gType || t == CellType
}

func OnOnline(fn func()) {
	js.Global.Get("document").Call("addEventListener", "online", fn, false)
}

func OnOffline(fn func()) {
	js.Global.Get("document").Call("addEventListener", "offline", fn, false)
}

package ble

import (
	"github.com/gopherjs/gopherjs/js"
)

type Peripheral struct {
	o *js.Object
}

func (p *Peripheral) Name() string {
	return p.o.Get("name").String()
}

func (p *Peripheral) ID() string {
	return p.o.Get("id").String()
}

func (p *Peripheral) RSSI() int {
	return p.o.Get("rssi").Int()
}

func (p *Peripheral) RawAdvData() []byte {
	return js.Global.Get("Uint8Array").New(p.o.Get("advertising")).Interface().([]byte)
}

func (p *Peripheral) Services() (ret []string) {
	ret = make([]string, 0)
	servicesJS, ok := p.o.Get("services").Interface().([]interface{})
	if !ok {
		return
	}
	for _, srv := range servicesJS {
		ret = append(ret, srv.(string))
	}
	return
}

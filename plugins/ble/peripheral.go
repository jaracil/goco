package ble

import (
	"encoding/hex"

	"github.com/gopherjs/gopherjs/js"
)

type AdvField struct {
	Type int
	Data []byte
}

type Peripheral struct {
	*js.Object
}

func (p *Peripheral) Name() string {
	return p.Get("name").String()
}

func (p *Peripheral) ID() string {
	return p.Get("id").String()
}

func (p *Peripheral) RSSI() int {
	return p.Get("rssi").Int()
}

func (p *Peripheral) RawAdvData() []byte {
	return js.Global.Get("Uint8Array").New(p.Get("advertising")).Interface().([]byte)
}

func (p *Peripheral) Services() (ret []string) {
	ret = make([]string, 0)
	servicesJS, ok := p.Get("services").Interface().([]interface{})
	if !ok {
		return
	}
	for _, srv := range servicesJS {
		ret = append(ret, srv.(string))
	}
	return
}

func ToUUID(data []byte) (ret string) {
	if data != nil && len(data) == 16 {
		ret = hex.EncodeToString(data[0:4]) + "-" + hex.EncodeToString(data[4:6]) + "-" + hex.EncodeToString(data[6:8]) + "-" + hex.EncodeToString(data[8:10]) + "-" + hex.EncodeToString(data[10:16])
	}
	return
}

func Reverse(data []byte) []byte {
	if data != nil {
		for i, j := 0, len(data)-1; i < j; i, j = i+1, j-1 {
			data[i], data[j] = data[j], data[i]
		}
	}
	return data
}

func GetData(fields []*AdvField, tp int) (ret []byte) {

	for _, field := range fields {
		if field.Type == tp {
			ret = field.Data
			break
		}
	}
	return
}

func ParseAdvRawData(data []byte) (ret []*AdvField) {
	p := 0
	ret = make([]*AdvField, 0)
	for p < len(data)-1 {

		size := int(data[p])
		if size == 0 || (p+size+1) > len(data) {
			break
		}
		ret = append(ret, &AdvField{Type: int(data[p+1]), Data: data[p+2 : p+size+1]})
		p += (size + 1)
	}
	return
}

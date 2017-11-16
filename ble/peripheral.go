package ble

import (
	"encoding/base64"
	"encoding/hex"

	"github.com/gopherjs/gopherjs/js"
)

type AdvField struct {
	Type int
	Data []byte
}

type Peripheral struct {
	*js.Object
	serviceUUID string
	vkID        string
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

func (p *Peripheral) rawAdvData() []byte {
	return js.Global.Get("Uint8Array").New(p.Get("advertising")).Interface().([]byte)
}

func (p *Peripheral) ServiceUUID() string {
	return p.serviceUUID
}

func (p *Peripheral) VkID() string {
	return p.vkID
}

func (p *Peripheral) Parse() {
	p.parseAndroid()
}

func (p *Peripheral) parseAndroid() {
	fields := parseAdvRawData(p.rawAdvData())
	p.serviceUUID = toUUID(reverse(getData(fields, 0x7)))

	srvData := getData(fields, 0x16)
	if srvData != nil && len(srvData) == 8 {
		p.vkID = base64.StdEncoding.EncodeToString(srvData[2:8])
	}
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

func toUUID(data []byte) (ret string) {
	if data != nil && len(data) == 16 {
		ret = hex.EncodeToString(data[0:4]) + "-" + hex.EncodeToString(data[4:6]) + "-" + hex.EncodeToString(data[6:8]) + "-" + hex.EncodeToString(data[8:10]) + "-" + hex.EncodeToString(data[10:16])
	}
	return
}

func reverse(data []byte) []byte {
	if data != nil {
		for i, j := 0, len(data)-1; i < j; i, j = i+1, j-1 {
			data[i], data[j] = data[j], data[i]
		}
	}
	return data
}

func getData(fields []*AdvField, tp int) (ret []byte) {

	for _, field := range fields {
		if field.Type == tp {
			ret = field.Data
			break
		}
	}
	return
}

func parseAdvRawData(data []byte) (ret []*AdvField) {
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

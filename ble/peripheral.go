package ble

import (
	"encoding/hex"
	"fmt"

	"github.com/gopherjs/gopherjs/js"
	"github.com/jaracil/goco/device"
)

type AdvField struct {
	Type int
	Data []byte
}

type Characteristic struct {
	*js.Object
	Service        string      `js:"service"`
	Characteristic string      `js:"characteristic"`
	Properties     []string    `js:"properties"`
	Descriptors    interface{} `js:"descriptors"`
}

type Peripheral struct {
	*js.Object
	name             string
	txPowerLevel     int
	flags            int
	services         []string
	servicesData     map[string][]byte
	manufacturerData map[string][]byte
	unknown          map[int][]byte
	characteristics  []Characteristic `js:"characteristics"`
}

func NewPeripheral(jsObj *js.Object) *Peripheral {
	return &Peripheral{
		Object:           jsObj,
		services:         []string{},
		servicesData:     map[string][]byte{},
		manufacturerData: map[string][]byte{},
		unknown:          map[int][]byte{},
		characteristics:  []Characteristic{},
	}
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

func (p *Peripheral) Characteristics() []Characteristic {
	return p.characteristics
}

func (p *Peripheral) Services() (ret []string) {
	return p.services
}

func (p *Peripheral) ServiceData(key string) []byte {
	return p.servicesData[key]
}

func (p *Peripheral) Parse() {
	if device.DevInfo.Platform == "Android" {
		p.parseAndroid()
	} else {
		p.parseIOS()
	}
}

func (p *Peripheral) parseIOS() {
	advertising := p.Get("advertising")

	p.name = advertising.Get("kCBAdvDataLocalName").String()
	p.txPowerLevel = advertising.Get("kCBAdvDataTxPowerLevel").Int()

	for _, item := range advertising.Get("kCBAdvDataServiceUUIDs").Interface().([]interface{}) {
		p.services = append(p.services, item.(string))
	}

	for _, key := range js.Keys(advertising.Get("kCBAdvDataServiceData")) {
		buffer := advertising.Get("kCBAdvDataServiceData").Get(key)
		data := js.Global.Get("Uint8Array").New(buffer).Interface().([]byte)
		p.servicesData[key] = data
	}

	data := js.Global.Get("Uint8Array").New(advertising.Get("kCBAdvDataManufacturerData")).Interface().([]byte)

	if len(data) >= 2 {
		key := p.formatUUID(reverse(data[0:2]))
		value := data[2:]
		p.manufacturerData[key] = value
	}
}

func (p *Peripheral) parseAndroid() {
	arr := p.rawAdvData()
	i := 0

	for i < len(arr) {
		fieldLength := int(arr[i]) - 1
		i++

		if fieldLength == -1 {
			break
		}

		fieldType := arr[i]
		i++

		switch fieldType {
		case 0x01:
			p.flags = int(arr[i])
			i += int(fieldLength)
		case 0x02:
			fallthrough
		case 0x03:
			i = p.extractUUIDs(arr, i, fieldLength, 2)
		case 0x04:
			fallthrough
		case 0x05:
			i = p.extractUUIDs(arr, i, fieldLength, 4)
		case 0x06:
			fallthrough
		case 0x07:
			i = p.extractUUIDs(arr, i, fieldLength, 16)
		case 0x08:
			fallthrough
		case 0x09:
			fieldData := arr[i : i+fieldLength]
			p.name = string(fieldData)
			i += fieldLength
		case 0x0a:
			p.txPowerLevel = int(arr[i])
			i += fieldLength
		case 0x16:
			key := p.formatUUID(arr[i : i+2])
			value := arr[i+2 : i+2+fieldLength-2]
			p.servicesData[key] = value
			i += fieldLength
		case 0xff:
			key := p.formatUUID(reverse(arr[i : i+2]))
			value := arr[i+2 : i+2+fieldLength-2]
			p.manufacturerData[key] = value
			i += fieldLength
		default:
			p.unknown[int(fieldType)] = arr[i : i+fieldLength]
			i += fieldLength
		}
	}
}

func (p *Peripheral) extractUUIDs(advertising []byte, index int, length int, uuidNumBytes int) int {
	uuids := []string{}
	remaining := length
	i := index

	for remaining > 0 {
		uuids = append(uuids, p.formatUUID(reverse(advertising[i:i+uuidNumBytes])))
		i += uuidNumBytes
		remaining -= uuidNumBytes
	}

	p.services = append(p.services, uuids...)

	return i
}

func (p *Peripheral) formatUUID(data []byte) string {
	result := ""

	for _, val := range data {
		result += fmt.Sprintf("%02X", val)
	}

	if len(result) == 32 {
		result = result[0:8] + "-" + result[8:12] + "-" + result[12:16] + "-" + result[16:20] + "-" + result[20:]
	}

	return result
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

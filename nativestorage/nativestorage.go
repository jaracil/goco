package nativestorage

import (
	"errors"

	"github.com/gopherjs/gopherjs/js"
)

var (
	ErrWriteFailed = errors.New("Storage error: Write failed")
	ErrNotFound    = errors.New("Storage error: Item not found")
	ErrNullRef     = errors.New("Storage error: Null reference")
	ErrUndefined   = errors.New("Storage error: Undefined type")
	ErrJSON        = errors.New("Storage error: JSON error")
	ErrParam       = errors.New("Storage error: Wrong parameter")
	ErrUnknown     = errors.New("Storage error: Unknown")
)

var instance *js.Object

func mo() *js.Object {
	if instance == nil {
		instance = js.Global.Get("NativeStorage")
	}
	return instance
}

func safeClose(ch chan struct{}) {
	select {
		case <- ch:
		default:
			close(ch)
	}
}

func errorByCode(code int) error {
	switch code {
	case 1:
		return ErrWriteFailed
	case 2:
		return ErrNotFound
	case 3:
		return ErrNullRef
	case 4:
		return ErrUndefined
	case 5:
		return ErrJSON
	case 6:
		return ErrParam
	}
	return ErrUnknown
}

func SetItem(key string, val interface{}) (err error) {
	ch := make(chan struct{})
	success := func() {
		safeClose(ch)
	}
	fail := func(obj *js.Object) {
		err = errorByCode(obj.Get("code").Int())
		safeClose(ch)
	}
	mo().Call("setItem", key, val, success, fail)
	<-ch
	return
}

func GetItemJS(key string) (ret *js.Object, err error) {
	ch := make(chan struct{})
	success := func(obj *js.Object) {
		ret = obj
		safeClose(ch)
	}
	fail := func(obj *js.Object) {
		err = errorByCode(obj.Get("code").Int())
		safeClose(ch)
	}
	mo().Call("getItem", key, success, fail)
	<-ch
	return
}

func GetItem(key string) (interface{}, error) {
	r, e := GetItemJS(key)
	if e != nil {
		return nil, e
	}
	return r.Interface(), nil
}

func GetInt(key string) (int, error) {
	r, e := GetItemJS(key)
	if e != nil {
		return 0, e
	}
	return r.Int(), nil
}

func GetInt64(key string) (int64, error) {
	r, e := GetItemJS(key)
	if e != nil {
		return 0, e
	}
	return r.Int64(), nil
}

func GetFloat64(key string) (float64, error) {
	r, e := GetItemJS(key)
	if e != nil {
		return 0, e
	}
	return r.Float(), nil
}

func GetString(key string) (string, error) {
	r, e := GetItemJS(key)
	if e != nil {
		return "", e
	}
	return r.String(), nil
}

func GetBool(key string) (bool, error) {
	r, e := GetItemJS(key)
	if e != nil {
		return false, e
	}
	return r.Bool(), nil
}

func RemoveItem(key string) (err error) {
	ch := make(chan struct{})
	success := func() {
		safeClose(ch)
	}
	fail := func(obj *js.Object) {
		err = errorByCode(obj.Get("code").Int())
		safeClose(ch)
	}
	mo().Call("remove", key, success, fail)
	<-ch
	return
}

func RemoveAll() (err error) {
	ch := make(chan struct{})
	success := func() {
		safeClose(ch)
	}
	fail := func(obj *js.Object) {
		err = errorByCode(obj.Get("code").Int())
		safeClose(ch)
	}
	mo().Call("clear", success, fail)
	<-ch
	return
}

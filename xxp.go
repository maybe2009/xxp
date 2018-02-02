package xxp

import (
	"reflect"
	"encoding/binary"
	"bytes"
	"log"
)

type coder struct {
	Buf *bytes.Buffer
	Order binary.ByteOrder
}

//coder will handle all the error she encounter
func (c *coder) error(e error) {
	panic(e)
}

func makeCoder(size uint32, order binary.ByteOrder) *coder {
	c := &coder{}
	c.Buf = bytes.NewBuffer(make([]byte, size, size))
	c.Buf.Reset()
	c.Order = order

	return c
}

//xxp is a LV (length-value)_format protocol
func (c *coder) encode(i interface{}) {
	v := reflect.ValueOf(i)
	switch v.Type().Kind() {
	case reflect.Struct:
		nf := v.NumField()
		for i := 0; i < nf ; i++ {
			c.marshalField(v.Field(i))
		}
	default:
		c.marshalField(v)
	}
}

func (c *coder) marshalField(v reflect.Value) {
	//encode length
	l := uint32(c.calculateLength(v))  //xxp only support 4 byte length
	log.Println("field length is ", l)
	c.encodeToBuf(reflect.ValueOf(l))

	//encode buffer
	c.encodeToBuf(v)
}

func (c *coder) encodeToBuf(data reflect.Value) {
	var b [8]byte

	k := data.Kind()
	switch k {
	case reflect.Int8:
		b[0] = byte(data.Int())
		c.Buf.Write(b[:1])

	case reflect.Uint8:
		b[0] = byte(data.Uint())
		c.Buf.Write(b[:1])

	case reflect.Int16:
		c.Order.PutUint16(b[:2], uint16(data.Int()))
		c.Buf.Write(b[:2])

	case reflect.Uint16:
		c.Order.PutUint16(b[:2], uint16(data.Uint()))
		c.Buf.Write(b[:2])

	case reflect.Int32:
		c.Order.PutUint32(b[:4], uint32(data.Int()))
		c.Buf.Write(b[:4])

	case reflect.Uint32:
		c.Order.PutUint32(b[:4], uint32(data.Uint()))
		c.Buf.Write(b[:4])

	case reflect.Int64:
		c.Order.PutUint64(b[:8], uint64(data.Int()))
		c.Buf.Write(b[:8])

	case reflect.Uint64:
		c.Order.PutUint64(b[:8], data.Uint())
		c.Buf.Write(b[:8])

	case reflect.String:
		c.Buf.Write([]byte(data.String()))

	case reflect.Slice:
		l := data.Len()
		for i := 0; i < l; i++ {
			c.encodeToBuf(data.Index(i))
		}

	case reflect.Array:
		l := data.Len()
		for i := 0; i < l; i++ {
			c.encodeToBuf(data.Index(i))
		}

	case reflect.Struct:
		nf := data.NumField()
		for f := 0; f < nf; f++ {
			c.encodeToBuf(data.Field(f))
		}

	case reflect.Map:
		keys := data.MapKeys()
		for _, k := range  keys {
			c.encodeToBuf(k)
			c.encodeToBuf(data.MapIndex(k))
		}

	default:
		c.error(&UnsupportEncodeType{data.Type()})
	}
}

func (c *coder) calculateLength(data reflect.Value) uint64 {
	k := data.Kind()
	switch k {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32,
		reflect.Int64, reflect.Uint8, reflect.Uint16, reflect.Uint32,
		reflect.Uint64:
		return uint64(data.Type().Size())

	case reflect.String:
		return uint64(len(data.String()))

	case reflect.Slice:
		var sum uint64 = 0
		l := data.Len()
		for i := 0; i < l; i++ {
			sum += c.calculateLength(data.Index(i))
		}
		return sum

	case reflect.Array:
		var sum uint64 = 0
		l := data.Len()
		for i := 0; i < l; i++ {
			sum += c.calculateLength(data.Index(i))
		}
		return sum

	case reflect.Struct:
		var sum uint64 = 0
		nf := data.NumField()
		for f := 0; f < nf; f++ {
			sum += c.calculateLength(data.Field(f))
		}
		return sum

	case reflect.Map:
		var sum uint64 = 0
		keys := data.MapKeys()
		for _, k := range  keys {
			sum += c.calculateLength(k)
			sum += c.calculateLength(data.MapIndex(k))
		}
		return sum

	default:
		panic("Unsupport type " + data.Type().Name())
	}
}

type UnsupportEncodeType struct {
	t reflect.Type
}

func (e *UnsupportEncodeType) Error() string {
	return "Type " + e.t.Name() + " is not support for encoding"
}
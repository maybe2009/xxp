package xxp
import (
"reflect"
"encoding/binary"
"bytes"
)

type coder struct {
	Buf *bytes.Buffer
	Order binary.ByteOrder
}

//coder will handle all the error she encounter
func (c *coder) error(e error) {
	panic(e)
}

func makeCoder(order binary.ByteOrder) *coder {
	c := &coder{}
	c.Buf = bytes.NewBuffer(make([]byte, 0, 0))
	c.Buf.Reset()
	c.Order = order

	return c
}

func (c *coder) encodeInterface(i interface{}) {
	c.encodeValue(reflect.ValueOf(i))
}

func (c *coder) encodeValue(data reflect.Value) {
	var b [8]byte

	switch k := data.Kind(); k {
	case reflect.Int:
		intSize := data.Type().Size()
		if intSize == 4 {
			c.encodeInterface(int32(data.Int()))
		} else if intSize == 8 {
			c.encodeInterface(data.Int())
		}

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
		c.encodeInterface(data.Len())
		c.Buf.Write([]byte(data.String()))

	case reflect.Slice:
		l := data.Len()
		c.encodeInterface(l)
		for i := 0; i < l; i++ {
			c.encodeValue(data.Index(i))
		}

	case reflect.Array:
		l := data.Len()
		c.encodeInterface(l)
		for i := 0; i < l; i++ {
			c.encodeValue(data.Index(i))
		}

	case reflect.Map:
		keys := data.MapKeys()
		c.encodeInterface(len(keys))
		for _, k := range  keys {
			c.encodeValue(k)
			c.encodeValue(data.MapIndex(k))
		}

	case reflect.Struct:
		nf := data.NumField()
		for f := 0; f < nf; f++ {
			c.encodeValue(data.Field(f))
		}
	default:
		c.error(&UnsupportType{data.Type()})
	}
}

type decoder struct {
	Buf *bytes.Buffer
	Order binary.ByteOrder
}

func makeDecoder(buf *bytes.Buffer, order binary.ByteOrder) *decoder {
	d := &decoder{}
	d.Buf = buf;
	d.Order = order
	return d
}

func (d *decoder) decodeInterface(i interface{}) {
	t := reflect.TypeOf(i)
	if t.Kind() != reflect.Ptr {
		d.error(&UnsupportType{t})
	}
	d.decodeValue(reflect.ValueOf(i))
}

func (d *decoder) decodeValue(data reflect.Value) {
	v := reflect.Indirect(data)
	if !v.CanSet() {
		d.error(&UnsupportType{data.Type()})
	}

	switch v.Type().Kind() {
	case reflect.Int8:
		b, err := d.Buf.ReadByte()
		if err != nil {
			panic(err)
		}
		v.SetInt(int64(b))

	case reflect.Int16:
		b := d.Buf.Next(2)
		v.SetInt(int64(d.Order.Uint16(b)))

	case reflect.Int32:
		b := d.Buf.Next(4)
		v.SetInt(int64(d.Order.Uint32(b)))

	case reflect.Uint32:
		b := d.Buf.Next(4)
		v.SetUint(uint64(d.Order.Uint32(b)))

	case reflect.String:
		l := d.Order.Uint64(d.Buf.Next(8))
		b := d.Buf.Next(int(l))
		v.SetString(string(b))

	case reflect.Map:
		l := d.Order.Uint64(d.Buf.Next(8))

		keyType := v.Type().Key()
		valType := v.Type().Elem()

		for i := 0; i < int(l); i++ {
			keyPtr := reflect.New(keyType)
			valPtr := reflect.New(valType)
			d.decodeValue(keyPtr)
			d.decodeValue(valPtr)

			v.SetMapIndex(reflect.Indirect(keyPtr), reflect.Indirect(valPtr))
		}

	case reflect.Struct:
		nf := v.NumField()
		for f := 0; f < nf; f++ {
			fieldPtr := reflect.New(v.Field(f).Type())
			d.decodeValue(fieldPtr)

			fieldVal := v.Field(f)
			fieldVal.Set(reflect.Indirect(fieldPtr))
		}
	default:
		d.error(&UnsupportType{v.Type()})
	}
}

func (d decoder) error(err error) {
	panic(err)
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

type UnsupportType struct {
	t reflect.Type
}

func (e *UnsupportType) Error() string {
	return "Type " + e.t.Name() + " is not support for encoding"
}
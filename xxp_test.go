package xxp

import (
	"testing"
	"log"
	"encoding/hex"
	"encoding/binary"
	"reflect"
)

func TestInteger(t *testing.T) {
	c := makeCoder(binary.BigEndian)
	var i8 int8 = 1
	var i16 int16 = 123
	var i32 int32 = 3
	var i64 int64 = 4
	var u8 uint8 = 5
	var u16 uint16 = 6
	var u32 uint32 = 7
	var u64 uint64 = 8

	c.encodeInterface(i8)
	c.encodeInterface(i16)
	c.encodeInterface(i32)
	c.encodeInterface(i64)
	c.encodeInterface(u8)
	c.encodeInterface(u16)
	c.encodeInterface(u32)
	c.encodeInterface(u64)
	log.Println("TestEncodeInteger: \n", hex.Dump(c.Buf.Bytes()))

	d := makeDecoder(c.Buf, c.Order)
	var di8 int8 = 0
	var di16 int16 = 0
	//var di32 int32 = 0
	//var di64 int64 = 0
	//var du8 uint8 = 0
	//var du16 uint16 = 0
	//var du32 uint32 = 0
	//var du64 uint64 = 0
	//var di8 int8 = 0

	d.decodeInterface(&di8)
	log.Println("di8: ", di8)
	if i8 != di8 {
		log.Fatal("not equal")
	}

	d.decodeInterface(&di16)
	log.Println("di16", di16)
	if i16 != di16 {
		t.Fatal("not equal")
	}

	log.Println("TestEncodeInteger: \n", hex.Dump(c.Buf.Bytes()))
}

func TestString(t *testing.T) {
	c := makeCoder(binary.BigEndian)
	var str string = "Hello, World"
	c.encodeInterface(str)
	log.Println("TestString: ", hex.Dump(c.Buf.Bytes()))

	d := makeDecoder(c.Buf, c.Order)
	dstr := ""
	d.decodeInterface(&dstr)
	log.Println("TestString dstr ", dstr)
	if str != dstr {
		t.Fatal("not equal")
	}
}

func TestMap(t *testing.T) {
	c := makeCoder(binary.BigEndian)
	m := map[string]uint32{"hello":123}
	m["oh"] = 789
	c.encodeInterface(m)
	log.Println("TestMap: \n", hex.Dump(c.Buf.Bytes()))

	d := makeDecoder(c.Buf, c.Order)
	dm := map[string]uint32{}
	d.decodeInterface(&dm)
	log.Println("TestMap dmap ", dm)
	v, ok := dm["hello"]
	if !ok || v != m["hello"] {
		t.Fatal("key not exist")
	}

	v2, ok := dm["oh"]
	if !ok || v2 != m["oh"] {
		t.Fatal("key not exist")
	}
}

func TestEncodeStruct(t *testing.T) {
	c := makeCoder(binary.BigEndian)
	type ts struct {
		Id uint32
		Name string
	}

	s := ts{7758, "Allen"}
	c.encodeInterface(s)
	log.Println("TestEncodeStruct: \n", hex.Dump(c.Buf.Bytes()))

	d := makeDecoder(c.Buf, c.Order)
	s2 := ts{}
	d.decodeInterface(&s2)
	log.Println("s2: ", s2)
	if s != s2 {
		t.Fatal("struct not equal")
	}
}

func TestEncodeMap(t *testing.T) {
	c := makeCoder(binary.BigEndian)
	m := map[uint32]string{}
	m[1] = "Allen"
	m[2] = "Alice"
	c.encodeInterface(m)
	log.Println("TestEncodeMap: \n", hex.Dump(c.Buf.Bytes()))
}

func TestAll(t *testing.T) {
	c := makeCoder(binary.BigEndian)
	type ts struct {
		Id uint32
		Name string
		m map[uint32]string
	}

	s := ts{}
	s.Id = 123
	s.Name = "Allen"

	m := map[uint32]string{}
	m[1] = "Allen"
	m[2] = "Alice"
	s.m = m
	c.encodeInterface(s)
	log.Println("TestAll: \n", hex.Dump(c.Buf.Bytes()))
}

func change(i interface{}) {
	t := reflect.TypeOf(i)
	switch t.Kind() {
	case reflect.Int32:
		v := reflect.ValueOf(i)
		v.SetInt(7758)
	}
}

func TestValue(t *testing.T) {
	var i uint32 = 0
	change(i)
	log.Println("i is ", i)
}
package xxp

import (
	"testing"
	"log"
	"encoding/hex"
	"encoding/binary"
	"reflect"
)

func TestEncodeInteger(t *testing.T) {
	c := makeCoder(0, binary.BigEndian)
	var i8 int8 = 1
	var i16 int16 = 2
	var i32 int32 = 3
	var i64 int64 = 4
	var u8 uint8 = 5
	var u16 uint16 = 6
	var u32 uint32 = 7
	var u64 uint64 = 8

	c.encode(i8)
	c.encode(i16)
	c.encode(i32)
	c.encode(i64)
	c.encode(u8)
	c.encode(u16)
	c.encode(u32)
	c.encode(u64)

	log.Println("coder: \n", hex.Dump(c.Buf.Bytes()))
}

func TestEncodeStruct(t *testing.T) {
	c := makeCoder(4, binary.BigEndian)
	type ts struct {
		Id uint32
		Name string
	}

	s := ts{7758, "Allen"}
	log.Println("struct size ", len(s.Name))
	c.encode(s)
	log.Println("coder: \n", hex.Dump(c.Buf.Bytes()))
}

func TestCalculateLength(t *testing.T) {
	c := makeCoder(0, binary.BigEndian)
	type ts struct {
		Id uint32
		Ok uint32
		Name string
	}

	s := ts{7758, 123, "Allen!"}
	log.Println("length is ", c.calculateLength(reflect.ValueOf(s)))
	c.encode(s)
	log.Println("coder: \n", hex.Dump(c.Buf.Bytes()))
}

func TestEncodeMap(t *testing.T) {
	c := makeCoder(0, binary.BigEndian)
	m := map[uint32]string{}
	m[1] = "Allen"
	m[2] = "Alice"
	c.encode(m)
	log.Println("coder: \n", hex.Dump(c.Buf.Bytes()))
}

func TestEveryOne(t *testing.T) {
	c := makeCoder(0, binary.BigEndian)
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
	c.encode(s)
	log.Println("coder: \n", hex.Dump(c.Buf.Bytes()))
}
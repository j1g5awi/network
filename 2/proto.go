package main

import (
	"bytes"
)

var ERROR uint8 = 0
var CONNECT uint8 = 1
var AUTH uint8 = 2
var DATA uint8 = 6
var DISCONNECT uint8 = 7

var AUTH_REQUIRE = 'T'
var AUTH_NO_REQUIRE = 'F'

var MAX_DATA_LENGTH_10 = 'A'
var MAX_DATA_LENGTH_20 = 'N'
var MAX_DATA_LENGTH_50 = 'Y'

var ACCEPT = 'T'
var REJECT = 'F'

type Message struct {
	version  uint8
	type_    uint8
	length   uint16
	data     string
	checksum uint16
}

func toByte(uint16 uint16) []byte {
	first := byte(uint16 / 256)
	second := byte(uint16 % 256)
	return []byte{first, second}
}

func (m *Message) Marsal() (data []byte) {
	buffer := bytes.NewBuffer(data)
	buffer.WriteByte(m.version)
	buffer.WriteByte(m.type_)
	buffer.Write(toByte(m.length))
	buffer.WriteString(m.data)
	buffer.Write(toByte(m.checksum))
	return buffer.Bytes()
}

func (m *Message) UnMarsal(data []byte) {
	m.version = uint8(data[0])
	m.type_ = uint8(data[1])
	m.length = uint16(int(data[2])*256 + int(data[3]))
	m.data = string(data[4 : len(data)-2])
	m.checksum = uint16(int(data[len(data)-2])*256 + int(data[len(data)-1]))
}

func NewMessage(type_ uint8, data string) *Message {
	m := &Message{}
	m.version = 19
	m.type_ = type_
	m.data = data
	m.length = uint16(6 + len([]byte(m.data)))
	m.checksum = 0
	return m
}

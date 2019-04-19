package main

import (
	"bytes"
	"encoding/binary"
	"net"
	"time"
)

type ICMP struct {
	Type        uint8
	Code        uint8
	CheckSum    uint16
	Identifier  uint16
	SequenceNum uint16
}

func GetICMP(seq uint16) ICMP {
	icmp := ICMP{
		Type:        8,
		Code:        0,
		CheckSum:    0,
		Identifier:  0,
		SequenceNum: seq,
	}
	var buffer bytes.Buffer
	binary.Write(&buffer, binary.BigEndian, icmp)
	icmp.CheckSum = CheckSum(buffer.Bytes())
	buffer.Reset()
	return icmp
}

func SendICMPRequest(icmp ICMP, destAddr *net.IPAddr) (response_time int64, err error) {
	conn, err := net.DialIP("ip4:icmp", nil, destAddr)
	if err != nil {
		return 0, err
	}
	defer conn.Close()
	var buffer bytes.Buffer
	binary.Write(&buffer, binary.BigEndian, icmp)
	if _, err := conn.Write(buffer.Bytes()); err != nil {
		return 0, err
	}
	tStart := time.Now()
	conn.SetReadDeadline((time.Now().Add(time.Second * 2)))
	recv := make([]byte, 1024)
	_, err = conn.Read(recv)
	if err != nil {
		return 0, err
	}
	tEnd := time.Now()
	response_time = tEnd.Sub(tStart).Nanoseconds() / 1e6

	return response_time, err
}

func CheckSum(data []byte) uint16 {
	var (
		sum    uint32
		length int = len(data)
		index  int
	)
	for length > 1 {
		sum += uint32(data[index])<<8 + uint32(data[index+1])
		index += 2
		length -= 2
	}
	if length > 0 {
		sum += uint32(data[index])
	}
	sum += (sum >> 16)
	return uint16(^sum)
}

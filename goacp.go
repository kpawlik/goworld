package goworld

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
)

var (
	inout *bufio.ReadWriter
)

func init() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	inout = bufio.NewReadWriter(reader, writer)
}

func VerifyConnection(name string) bool {
	var (
		err error
	)
	buffName := []byte(name)
	nameLen := GetShort()
	log.Printf("Bytes to read: %d\n", nameLen)

	buf := make([]byte, 0, nameLen)
	n, err := inout.Read(buf[:cap(buf)])
	log.Printf("Bytes to readed: %d\n", n)
	if err != nil {
		log.Fatalf("Error reading string: %v\n", err)
	}
	buf = buf[:n]
	log.Printf("Data from SW: %s\n", string(buf))
	res := bytes.Compare(buffName, buf)
	log.Printf("Name: %s; From SW: %s; Res: %d\n", string(name), string(buf), res)
	PutBool(true)
	return res == 0
}

func EstablishProtocol(minProtocol, maxProtocol int) bool {
	min, max := minProtocol, maxProtocol
	for {
		in := GetShort()
		log.Printf("Protocol from SW: %d\n", in)
		if in < min || in > max {
			PutBool(false)
			PutUShort(uint16(min))
			PutUShort(uint16(max))
			return false
		} else {
			break
		}
	}
	PutBool(true)
	return true
}

func Connect(processName string, protocolMin, protocolMax int) (err error) {
	log.Printf("ACP started name: %s\n", processName)
	if res := VerifyConnection(processName); !res {
		err = errors.New("Error verify connection")
		return
	}
	log.Println("Connection verified")
	if res := EstablishProtocol(protocolMin, protocolMax); !res {
		err = errors.New("Error establish protocol")
		return
	}
	log.Println("Protocol established")
	return
}

func PutBool(b bool) {
	var ival byte
	if !b {
		ival = 1
	}
	inout.WriteByte(ival)
	inout.Flush()

}

func PutUShort(short uint16) {
	buf := make([]byte, 2, 2)
	binary.LittleEndian.PutUint16(buf, short)
	inout.Write(buf)
	inout.Flush()
}

func PutString(s string) {
	bytes := []byte(s)
	l := len(bytes)
	PutUShort(uint16(l))
	inout.Write(bytes)
	inout.Flush()
}

func readBytes(n int) (buf []byte, err error) {
	buf = make([]byte, 0, n)
	_, err = inout.Read(buf[:cap(buf)])
	if err != nil && err != io.EOF {
		err = errors.New(fmt.Sprintf("Error reading short: %v\n", err))
		return
	}
	buf = buf[:n]
	return
}

func GetShort() int {
	var (
		res uint16
	)
	if err := binary.Read(inout, binary.LittleEndian, &res); err != nil {
		log.Fatal(err)
	}
	return int(res)
}

func GetString() string {
	b := GetShort()
	log.Printf("get string - bytes to read %v\n", b)
	buf, err := readBytes(b)
	if err != nil {
		log.Fatal(err)
	}
	return string(buf)
}

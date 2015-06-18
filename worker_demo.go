package goworld

import (
	"fmt"
	"log"
)

// GetTestResponse returns response object from worker
func (t *Worker) GetTestResponse(request *Request, resp *Response) error {
	var (
		testInt, resInt       int
		ok                    bool
		testFloat, resFloat   float64
		testSFloat, resSFloat float32
		testString, resString string
		testLong, resLong     int64
		testUlong, resUlong   uint64
		testBool, resBool     bool
		bodyElem              BodyElement
	)
	defer func() {
		if err := recover(); err != nil {
			log.Panic("PANIC ", err)
		}
	}()

	body := []BodyElement{}

	path := request.Path
	log.Printf("GetTestResponse - handle path: %s\n", path)

	//test bool
	bodyElem = make(BodyElement)
	testBool = true
	acp.PutBool(testBool)
	resBool = acp.GetBool()
	ok = testBool == !resBool
	bodyElem["Bool1"] = fmt.Sprintf("Sent: %v, Get: %v, OK: %v", testBool, resBool, ok)
	testBool = false
	acp.PutBool(testBool)
	resBool = acp.GetBool()
	ok = testBool == !resBool
	bodyElem["Bool2"] = fmt.Sprintf("Sent: %v, Get: %v, OK: %v", testBool, resBool, ok)
	body = append(body, bodyElem)

	// test ubyte
	bodyElem = make(BodyElement)
	testInt = 12
	acp.PutUbyte(uint8(testInt))
	resInt = acp.GetUbyte()
	ok = testInt+1 == resInt
	bodyElem["UbyteTest"] = fmt.Sprintf("Sent: %d, Get: %d, OK: %v", testInt, resInt, ok)

	// test byte
	testInt = -111
	acp.PutByte(int8(testInt))
	resInt = acp.GetByte()
	ok = testInt+1 == resInt
	bodyElem["ByteTest"] = fmt.Sprintf("Sent: %d, Get: %d, OK: %v", testInt, resInt, ok)
	body = append(body, bodyElem)

	// test ushort
	bodyElem = make(BodyElement)
	testInt = 12
	acp.PutUshort(uint16(testInt))
	resInt = acp.GetUshort()
	ok = testInt+1 == resInt
	bodyElem["UshortTest"] = fmt.Sprintf("Sent: %d, Get: %d, OK: %v", testInt, resInt, ok)
	// test short
	testInt = -111
	acp.PutShort(int16(testInt))
	resInt = acp.GetShort()
	ok = testInt+1 == resInt
	bodyElem["ShortTest"] = fmt.Sprintf("Sent: %d, Get: %d, OK: %v", testInt, resInt, ok)
	body = append(body, bodyElem)

	// test uint
	bodyElem = make(BodyElement)
	testInt = 11112
	acp.PutUint(uint32(testInt))
	resInt = acp.GetUint()
	ok = testInt+1 == resInt
	bodyElem["UintTest"] = fmt.Sprintf("Sent: %d, Get: %d, OK: %v", testInt, resInt, ok)
	// test int
	testInt = -11112
	acp.PutInt(int32(testInt))
	resInt = acp.GetInt()
	ok = testInt+1 == resInt
	bodyElem["IntTest"] = fmt.Sprintf("Sent: %d, Get: %d, OK: %v", testInt, resInt, ok)
	body = append(body, bodyElem)

	// test ulong
	bodyElem = make(BodyElement)
	testUlong = 12345611112
	acp.PutUlong(testUlong)
	resUlong = acp.GetUlong()
	ok = testUlong+1 == resUlong
	bodyElem["UlongTest"] = fmt.Sprintf("Sent: %d, Get: %d, OK: %v", testUlong, resUlong, ok)
	// test long
	testLong = int64(-12345611112)
	acp.PutLong(testLong)
	resLong = acp.GetLong()
	ok = testLong+1 == resLong
	bodyElem["LongTest"] = fmt.Sprintf("Sent: %d, Get: %d, OK: %v", testLong, resLong, ok)
	body = append(body, bodyElem)

	//test short float
	bodyElem = make(BodyElement)
	testSFloat = -11111.44
	acp.PutShortFloat(testSFloat)
	resSFloat = acp.GetShortFloat()
	ok = testSFloat+1 == resSFloat
	bodyElem["ShortFloatTest"] = fmt.Sprintf("Sent: %f, Get: %f, OK: %v", testSFloat, resSFloat, ok)
	body = append(body, bodyElem)

	//test float
	bodyElem = make(BodyElement)
	testFloat = -111112122.44
	acp.PutFloat(float64(testFloat))
	resFloat = acp.GetFloat()
	ok = testFloat+1 == resFloat
	bodyElem["FloatTest"] = fmt.Sprintf("Sent: %f, Get: %f, OK: %v", testFloat, resFloat, ok)
	body = append(body, bodyElem)

	//test string
	bodyElem = make(BodyElement)
	testString = "111112122.44 NON ASCII chars łóśćężąłQążźGFGĘĄÓŁ  źżźć"
	acp.PutString(testString)
	resString = acp.GetString()
	ok = testString == resString
	bodyElem["StringTest"] = fmt.Sprintf("Sent: %s, Get: %s, OK: %v", testString, resString, ok)
	body = append(body, bodyElem)

	resp.Body = body
	return nil
}

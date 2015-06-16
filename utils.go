package goworld

import (
	"encoding/gob"
	"fmt"
	"strconv"
)

func init() {
	gob.Register(NewAcpErr(""))
}

func portNo(port int) string {
	return fmt.Sprintf(":%d", port)
}

type AcpErr struct {
	Err string
}

func (err *AcpErr) Error() string {
	return err.Err
}

func NewAcpErr(msg string) *AcpErr {
	return &AcpErr{msg}
}

func ParseStringParam(value, dataType string) (result interface{}, err error) {

	switch dataType {
	case "boolean":
		return strconv.ParseBool(value)
	case "unsigned_byte":
		if result, err = strconv.ParseUint(value, 10, 8); err != nil {
			break
		}
		result = result.(uint8)
		return
	case "signed_byte":
		if result, err = strconv.ParseInt(value, 10, 8); err != nil {
			break
		}
		result = result.(int8)
		return
	case "unsigned_short":
		if result, err = strconv.ParseUint(value, 10, 16); err != nil {
			break
		}
		result = result.(uint16)
		return
	case "signed_short":
		if result, err = strconv.ParseInt(value, 10, 16); err != nil {
			break
		}
		result = result.(int16)
		return
	case "unsigned_int":
		if result, err = strconv.ParseUint(value, 10, 32); err != nil {
			break
		}
		result = result.(uint32)
		return
	case "signed_int":
		if result, err = strconv.ParseInt(value, 10, 32); err != nil {
			break
		}
		result = result.(uint32)
		return
	case "unsigned_long":
		if result, err = strconv.ParseUint(value, 10, 64); err != nil {
			break
		}
		result = result.(uint64)
		return
	case "signed_long":
		if result, err = strconv.ParseInt(value, 10, 64); err != nil {
			break
		}
		result = result.(uint64)
		return
	case "short_float":
		if result, err = strconv.ParseFloat(value, 32); err != nil {
			break
		}
		result = result.(float32)
		return
	case "float":
		if result, err = strconv.ParseFloat(value, 64); err != nil {
			break
		}
		result = result.(float64)
		return
	case "chars":
		result = value
		return
	default:
		err = NewAcpErr(fmt.Sprintf("Unsuported data type '%s' in ParssStringParam", dataType))
		return
	}
	if err != nil {
		err = NewAcpErr(fmt.Sprintf("Error parsing string parameter '%s' to data type '%s' ", value, dataType))
	}
	return
}

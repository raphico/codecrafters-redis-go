package protocol

import (
	"fmt"
	"strings"
)

type ResponseType int

const (
	SimpleStringType ResponseType = iota
	ErrorType
	NullBulkStringType
	BulkStringType
	IntegerType
	ArrayType
	NullArrayType
)

type Response struct {
	Type  ResponseType
	Value any
}

func (r Response) Serialize() string {
	switch r.Type {
	case ErrorType:
		return fmt.Sprintf("-ERR %s\r\n", r.Value)
	case NullBulkStringType:
		return "$-1\r\n"
	case SimpleStringType:
		return fmt.Sprintf("+%s\r\n", r.Value)
	case IntegerType:
		return fmt.Sprintf(":%d\r\n", r.Value)
	case NullArrayType:
		return "*-1\r\n"
	case ArrayType:
		var builder strings.Builder
		responses := r.Value.([]Response)

		// write array length
		length := len(responses)
		builder.WriteString(fmt.Sprintf("*%d\r\n", length))

		// write bulk strings
		for i := range length {
			builder.WriteString(responses[i].Serialize())
		}

		return builder.String()
	case BulkStringType:
		str := r.Value.(string)
		return fmt.Sprintf("$%d\r\n%s\r\n", len(str), r.Value)
	default:
		panic("unexpected type: invalid response type")
	}
}

func NewSimpleStringResponse(value string) Response {
	return Response{
		Type:  SimpleStringType,
		Value: value,
	}
}

func NewIntegerResponse(value int) Response {
	return Response{
		Type:  IntegerType,
		Value: value,
	}
}

func NewErrorResponse(value string) Response {
	return Response{
		Type:  ErrorType,
		Value: value,
	}
}

func NewBulkStringResponse(value string) Response {
	return Response{
		Type:  BulkStringType,
		Value: value,
	}
}

func NewNullBulkStringResponse() Response {
	return Response{
		Type: NullBulkStringType,
	}
}

func NewArrayResponse(response []Response) Response {
	return Response{
		Type:  ArrayType,
		Value: response,
	}
}

func NewNullArrayResponse() Response {
	return Response{
		Type: NullArrayType,
	}
}

package attribute

import (
	"encoding/binary"
	"fmt"
)

type Attribute interface {
	getType() uint16
	getValue() []byte
}

func Encode(a Attribute) []byte {
	aType := a.getType()
	value := a.getValue()
	length := len(value)

	output := make([]byte, length+4)
	binary.BigEndian.PutUint16(output, aType)
	binary.BigEndian.PutUint16(output, uint16(length))

	output = append(output, value...)

	return output
}

func Decode(attributeBytes []byte) (Attribute, error) {
	attributeType := binary.BigEndian.Uint16(attributeBytes[0:])
	attributeLength := binary.BigEndian.Uint16(attributeBytes[2:])
	attributeValue := attributeBytes[4 : 4+attributeLength]

	switch attributeType {
	case typeXorMappedAddress:
		a, err := decodeXorMappedAddress(attributeValue)
		if err != nil {
			return nil, err
		}
		return a, nil
	default:
		return nil, fmt.Errorf("unsupported attribute type %d", attributeType)
	}
}

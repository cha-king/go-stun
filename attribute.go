package stun

import (
	"encoding/binary"
	"errors"
	"fmt"
	"net"
)

const (
	attributeTypeXorMappedAddress uint16 = 0x0020
)

const (
	familyIpv4 uint8 = 0x01
	familyIpv6 uint8 = 0x02
)

type attribute interface {
	getType() uint16
	getValue() []byte
}

func encodeAttribute(a attribute) []byte {
	aType := a.getType()
	value := a.getValue()
	length := len(value)

	output := make([]byte, length+4)
	binary.BigEndian.PutUint16(output, aType)
	binary.BigEndian.PutUint16(output, uint16(length))

	output = append(output, value...)

	return output
}

func decodeAttribute(attributeBytes []byte) (attribute, error) {
	attributeType := binary.BigEndian.Uint16(attributeBytes[0:])
	attributeLength := binary.BigEndian.Uint16(attributeBytes[2:])
	attributeValue := attributeBytes[4 : 4+attributeLength]

	switch attributeType {
	case attributeTypeXorMappedAddress:
		a, err := decodeXorMappedAddress(attributeValue)
		if err != nil {
			return nil, err
		}
		return a, nil
	default:
		return nil, fmt.Errorf("unsupported attribute type %d", attributeType)
	}
}

// func decodeAttributes(attributeBytes []byte) ([]attribute, error) {

// }

type xorMappedAddress struct {
	family   uint8
	xPort    uint16
	xAddress []byte
}

func (a xorMappedAddress) getType() uint16 {
	return attributeTypeXorMappedAddress
}

func (a xorMappedAddress) getValue() []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint16(b, uint16(a.family))
	binary.BigEndian.PutUint16(b, uint16(a.xPort))
	b = append(b, a.xAddress...)
	return b
}

func (a xorMappedAddress) getAddress() (net.UDPAddr, error) {
	switch a.family {

	case familyIpv4:
		port := int(a.xPort ^ uint16(MessageMagicCookie>>16))
		xAddress := binary.BigEndian.Uint32(a.xAddress)
		addressBytes := xAddress ^ MessageMagicCookie
		addressSlice := make([]byte, 4)
		binary.BigEndian.PutUint32(addressSlice, addressBytes)
		ip := net.IP(addressSlice)
		return net.UDPAddr{IP: ip, Port: port, Zone: ""}, nil

	case familyIpv6:
		return net.UDPAddr{}, errors.New("IPv6 unsupported")

	default:
		return net.UDPAddr{}, fmt.Errorf("unsupported address family: %x", a.family)
	}
}
func decodeXorMappedAddress(addressBytes []byte) (xorMappedAddress, error) {
	family := uint8(binary.BigEndian.Uint16(addressBytes[0:]))
	xPort := binary.BigEndian.Uint16(addressBytes[2:])

	switch family {

	case familyIpv4:
		xAddress := addressBytes[4:8]
		return xorMappedAddress{family, xPort, xAddress}, nil

	case familyIpv6:
		// TODO: Implement me
		return xorMappedAddress{}, errors.New("IPv6 unsupported")

	default:
		return xorMappedAddress{}, fmt.Errorf("unsupported address family: %x", family)
	}
}

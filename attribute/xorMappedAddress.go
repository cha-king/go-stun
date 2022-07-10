package attribute

import (
	"encoding/binary"
	"errors"
	"fmt"
	"net"
)

type XorMappedAddress struct {
	family   uint8
	xPort    uint16
	xAddress []byte
}

func (a XorMappedAddress) getType() uint16 {
	return typeXorMappedAddress
}

func (a XorMappedAddress) getValue() []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint16(b, uint16(a.family))
	binary.BigEndian.PutUint16(b, uint16(a.xPort))
	b = append(b, a.xAddress...)
	return b
}

func (a XorMappedAddress) GetAddress() (net.UDPAddr, error) {
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
func decodeXorMappedAddress(addressBytes []byte) (XorMappedAddress, error) {
	family := uint8(binary.BigEndian.Uint16(addressBytes[0:]))
	xPort := binary.BigEndian.Uint16(addressBytes[2:])

	switch family {

	case familyIpv4:
		xAddress := addressBytes[4:8]
		return XorMappedAddress{family, xPort, xAddress}, nil

	case familyIpv6:
		// TODO: Implement me
		return XorMappedAddress{}, errors.New("IPv6 unsupported")

	default:
		return XorMappedAddress{}, fmt.Errorf("unsupported address family: %x", family)
	}
}

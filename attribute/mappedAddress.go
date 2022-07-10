package attribute

import (
	"encoding/binary"
	"fmt"
	"net"
)

type MappedAddress struct {
	family  uint8
	port    uint16
	address []byte
}

func (a MappedAddress) getType() uint16 {
	return typeMappedAddress
}

func (a MappedAddress) getValue() []byte {
	buf := make([]byte, 16)
	binary.BigEndian.PutUint16(buf[0:], uint16(a.family))
	binary.BigEndian.PutUint16(buf[2:], a.port)
	buf = append(buf, a.address...)
	return buf
}

func (a MappedAddress) GetAddress() net.UDPAddr {
	ip := net.IP(a.address)
	port := int(a.port)
	return net.UDPAddr{IP: ip, Port: port}
}

func decodeMappedAddress(b []byte) (MappedAddress, error) {
	family := uint8(binary.BigEndian.Uint16(b[0:]))
	port := binary.BigEndian.Uint16(b[2:])

	var address []byte
	switch family {
	case familyIpv4:
		address = b[4:8]
	case familyIpv6:
		address = b[4:20]
	default:
		return MappedAddress{}, fmt.Errorf("unsupported address family: %x", family)
	}

	return MappedAddress{family, port, address}, nil
}

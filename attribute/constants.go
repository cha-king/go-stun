package attribute

// TODO: Find a way to share this without circular import
const MessageMagicCookie uint32 = 0x2112A442

const (
	typeMappedAddress    uint16 = 0x0001
	typeXorMappedAddress uint16 = 0x0020
)

const (
	familyIpv4 uint8 = 0x01
	familyIpv6 uint8 = 0x02
)

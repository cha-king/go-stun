package stun

const (
	MessageClassRequest    uint8 = 0b00
	MessageClassIndication uint8 = 0b01
	MessageClassSuccess    uint8 = 0b10
	MessageClassError      uint8 = 0b11
)

const (
	MessageMethodBinding uint8 = 0x0001
)

const MessageMagicCookie uint32 = 0x2112A442

type message struct {
	header
	attributes []attribute
}

func (m message) encode() []byte {
	messageBytes := m.header.encode()
	for _, attribute := range m.attributes {
		attributesBytes := encodeAttribute(attribute)
		messageBytes = append(messageBytes, attributesBytes...)
	}
	return messageBytes
}

func newMessage(class uint8, method uint8, attributes []attribute) message {
	h := newHeader(class, method)
	m := message{h, attributes}
	return m
}

func decodeMessage(messageBytes []byte) (message, error) {
	headerBytes := messageBytes[:20]
	attributeBytes := messageBytes[20:]

	header, err := decodeHeader(headerBytes)
	if err != nil {
		return message{}, err
	}

	a, err := decodeAttribute(attributeBytes)
	if err != nil {
		return message{}, err
	}
	attributes := []attribute{a}

	return message{header, attributes}, nil
}

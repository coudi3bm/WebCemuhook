package protocol

import "fmt"

type Header struct {
	MagicString     [4]byte
	ProtocolVersion uint16
	Length          uint16
	Crc32           uint32
	Id              uint32
	EventType       uint32
}

func (h Header) String() string {
	return fmt.Sprintf(
		"{ MagicString: %s, ProtocolVersion: %d, Length: %d, Crc32: %X, Id: %X, EventType: %X }",
		h.MagicString, h.ProtocolVersion, h.Length, h.Crc32, h.Id, h.EventType,
	)
}

package utils

type MessageData struct {
	TimeStamp    uint64
	RotationRate struct {
		X float32
		Y float32
		Z float32
	}
}

const (
	MessageTypePROTOCOL = 0x100000
	MessageTypeINFO     = 0x100001
	MessageTypeDATA     = 0x100002
)

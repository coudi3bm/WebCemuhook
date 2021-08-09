package protocol

type ControllerDataResponse struct {
	Shared            SharedResponse
	Connected         byte
	PacketId          uint32
	ExtraButtons      byte
	MainButtons       byte
	PsExtraInput      uint16
	LeftStickXY       uint16
	RightStickXY      uint16
	DPadAnalog        uint32
	MainButtonsAnalog uint64

	Touch1 [6]byte
	Touch2 [6]byte

	MotionTimestamp uint64 // In microseconds
	AccelerometerX  float32
	AccelerometerY  float32
	AccelerometerZ  float32
	GyroscopePitch  float32
	GyroscopeYaw    float32
	GyroscopeRoll   float32
}

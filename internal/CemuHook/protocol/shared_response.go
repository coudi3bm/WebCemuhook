package protocol

import "fmt"

type SharedResponse struct {
	Slot           byte
	State          byte
	DeviceModel    byte
	ConnectionType byte
	MacAddress     [6]byte
	BatteryStatus  byte
}

func (c SharedResponse) String() string {
	return fmt.Sprintf(
		"{ Slot: %d, State: %d, DeviceModel: %d, ConnectionType: %d, MacAddress: %v, BatteryStatus: %d }",
		c.Slot, c.State, c.DeviceModel, c.ConnectionType, c.MacAddress, c.BatteryStatus,
	)
}

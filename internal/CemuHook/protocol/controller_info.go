package protocol

import "fmt"

type ControllerInfoRequest struct {
	PortsCount  uint32
	PortIndices [4]byte
}

func (c ControllerInfoRequest) String() string {
	return fmt.Sprintf(
		"{ PortsCount: %d, PortIndices: %v }",
		c.PortsCount, c.PortIndices,
	)
}

type ControllerInfoResponse struct {
	Shared   SharedResponse
	ZeroByte byte
}

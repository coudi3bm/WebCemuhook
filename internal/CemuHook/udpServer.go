package CemuHook

import (
	"bytes"
	"encoding/binary"
	"hash/crc32"
	"math/rand"
	"net"
	"time"
	"web-cemuhook/internal/CemuHook/protocol"
	"web-cemuhook/utils"

	"github.com/rs/zerolog/log"
)

type UdpServer struct {
	protocolVersion uint16
	serverId        uint32
	udpSocket       *net.UDPConn
	clientAddress   net.UDPAddr
	packetCount     uint32
}

func NewUdpServer() *UdpServer {
	u := UdpServer{
		protocolVersion: 1001,
		serverId:        rand.Uint32(),
	}
	u.udpSocket, _ = net.ListenUDP("udp", &net.UDPAddr{
		Port: 26760,
		IP:   net.ParseIP("127.0.0.1"),
	})
	return &u
}

func (u *UdpServer) StartReceive() {
	for {
		buffer := make([]byte, 1024)
		var length, clientAddress, err = u.udpSocket.ReadFromUDP(buffer)
		if err != nil {
			log.Fatal().Err(err)
		}
		u.ProcessIncoming(buffer, length, *clientAddress)
	}
}

func (u *UdpServer) ProcessIncoming(recv_buffer []byte, bufferLength int, clientAddress net.UDPAddr) {
	u.clientAddress = clientAddress

	var packetHeader protocol.Header
	binary.Read(bytes.NewReader(recv_buffer), binary.LittleEndian, &packetHeader)
	log.Debug().Str("Header", packetHeader.String()).Str("from", clientAddress.String()).Int("bufferLength", bufferLength).Msg("UdpServer Received")
	if packetHeader.EventType == utils.MessageTypeINFO {
		var packetControllerInfo protocol.ControllerInfoRequest
		binary.Read(bytes.NewReader(recv_buffer[20:]), binary.LittleEndian, &packetControllerInfo)
		//log.Debug().Str("ControllerInfoRequest", packetControllerInfo.String()).Msg("UdpServer Received - Process Incoming")

		packetBuffer := new(bytes.Buffer)
		u.BeginPacket(packetBuffer, utils.MessageTypeINFO)
		packerControllerInfoResponse := protocol.ControllerInfoResponse{
			Shared: protocol.SharedResponse{
				Slot:           0,
				State:          2,
				DeviceModel:    2,
				ConnectionType: 1,
				MacAddress:     [6]byte{0, 0, 0, 0, 0, 255},
				BatteryStatus:  0xEF,
			},
			ZeroByte: 0,
		}
		binary.Write(packetBuffer, binary.LittleEndian, packerControllerInfoResponse)

		packet := packetBuffer.Bytes()
		u.FinishPacket(packet)

		var packetHeader protocol.Header
		binary.Read(bytes.NewReader(packet), binary.LittleEndian, &packetHeader)
		//log.Debug().Str("packetHeader", packetHeader.String()).Msg("Finish Packet")

		u.udpSocket.WriteTo(packet, &clientAddress)
	}

}

func (u *UdpServer) SendControllerData(data *utils.MessageData) {
	u.packetCount++
	//time.Sleep(10 * time.Microsecond)
	log.Print(data)
	log.Print(u.clientAddress, u.packetCount)

	packetBuffer := new(bytes.Buffer)
	u.BeginPacket(packetBuffer, utils.MessageTypeDATA)
	packerControllerData := protocol.ControllerDataResponse{
		Shared: protocol.SharedResponse{
			Slot:           0,
			State:          2,
			DeviceModel:    2,
			ConnectionType: 1,
			MacAddress:     [6]byte{0, 0, 0, 0, 0, 255},
			BatteryStatus:  0xEF,
		},
		Connected:         1,
		PacketId:          u.packetCount,
		ExtraButtons:      0,
		MainButtons:       0,
		PsExtraInput:      0,
		LeftStickXY:       0,
		RightStickXY:      0,
		DPadAnalog:        0,
		MainButtonsAnalog: 0,

		Touch1: [6]byte{},
		Touch2: [6]byte{},

		MotionTimestamp: uint64(time.Now().UnixNano() / int64(time.Microsecond)),
		AccelerometerX:  0,
		AccelerometerY:  0,
		AccelerometerZ:  0,

		GyroscopePitch: data.RotationRate.X,
		GyroscopeYaw:   data.RotationRate.Y,
		GyroscopeRoll:  data.RotationRate.Z,
	}
	binary.Write(packetBuffer, binary.LittleEndian, packerControllerData)

	packet := packetBuffer.Bytes()
	u.FinishPacket(packet)

	var packetHeader protocol.Header
	binary.Read(bytes.NewReader(packet), binary.LittleEndian, &packetHeader)
	log.Debug().Str("packetHeader", packetHeader.String()).Msg("SendControllerData")

	u.udpSocket.WriteTo(packet, &u.clientAddress)
}

func (u UdpServer) BeginPacket(buffer *bytes.Buffer, messageType uint32) {
	binary.Write(buffer, binary.LittleEndian, protocol.Header{
		MagicString:     [4]byte{'D', 'S', 'U', 'S'},
		ProtocolVersion: u.protocolVersion,
		Length:          0,
		Crc32:           0,
		Id:              u.serverId,
		EventType:       messageType,
	})
}

func (u UdpServer) FinishPacket(buffer []byte) {
	binary.LittleEndian.PutUint32(buffer[6:], uint32(len(buffer)-16))
	binary.LittleEndian.PutUint32(buffer[8:], crc32.ChecksumIEEE(buffer))
}

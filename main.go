package main

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/songgao/water"
	"golang.org/x/net/ipv4"
	"net"
)

func handle(ifce *water.Interface, p []byte) {
	h, err := ipv4.ParseHeader(p)
	check(err)

	log(h)
	if h.Version != 4 || h.Protocol != 17 {
		return
	}

	udpBody := p[h.Len:]
	body := udpBody[8:]
	log("body:", hex.EncodeToString(body))

	h.Dst = net.IP([]byte{10, 1, 0, 1})
	h.Checksum = 0
	hB, err := h.Marshal()
	check(err)

	h.Checksum = int(checkSum(hB))
	hB, err = h.Marshal()
	check(err)

	//rH := &ipv4.Header{
	//	h.Version,
	//	h.Len,
	//	h.TOS,
	//	h.TotalLen,
	//	h.ID,
	//	h.Flags,
	//	h.FragOff,
	//	h.TTL,
	//	h.Protocol,
	//	h.Checksum,
	//	h.Dst,
	//	h.Src,
	//	h.Options,
	//}
	//rHB, err := rH.Marshal()
	//check(err)
	//

	r := bytes.NewBuffer(hB)
	r.Write(udpBody)
	_, err = ifce.Write(r.Bytes())
	log(hex.EncodeToString(r.Bytes()))
	check(err)
}

func main() {
	cfg := water.Config{
		//DeviceType: water.TAP,
		DeviceType: water.TUN,
	}
	cfg.Name = "tun0"

	ifce, err := water.New(cfg)
	check(err)

	fmt.Printf("Interface Name: %s\n", ifce.Name())

	packet := make([]byte, 2000)
	for {
		n, err := ifce.Read(packet)
		if err != nil {
			log("error:", err)
		}
		fmt.Printf("Packet Received: % x\n", packet[:n])
		handle(ifce, packet[:n])
	}
}

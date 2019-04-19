package main

import (
	"fmt"
	"github.com/songgao/water"
	"golang.org/x/net/ipv4"
	"net"
	"railgun"
)

var (
	check = railgun.Check
	log = railgun.Log
)

type RailgunClient struct {
	listenIP net.IP
	ipNet *net.IPNet
	rawConn *ipv4.RawConn
	udpConn *net.UDPConn
}

func (c *RailgunClient) handle(p []byte) {
	h, err := ipv4.ParseHeader(p)
	if err != nil {
		log("parse ip header error:", err)
		return
	}

	//if s.ipNet.Contains(h.Dst) && !h.Dst.Equal(s.listenIP) {
	if h.Dst.Equal(c.listenIP) {
		log(h, p)
		body := p[h.Len:]
		err = c.rawConn.WriteTo(h, body, nil)

		if err != nil {
			log("conn write error:", err)
		}

		//} else if h.Dst.Equal(s.ipNet.IP) || !s.ipNet.Contains(h.Dst) {
		//} else if c.ipNet.Contains(h.Dst) {
	} else {
		log(h, p)
		_, err = c.udpConn.Write(p)
		if err != nil {
			log("udp write error:", err)
		}
	}
}

func (c *RailgunClient) runTUN() {
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
		//fmt.Printf("Packet Received: % x\n", packet[:n])
		c.handle(packet[:n])
	}
}

func (c *RailgunClient) runUDPRecv() {
	buf := make([]byte, 2000)

	for {
		n, err := c.udpConn.Read(buf)
		if err != nil {
			log("udp recv error:", err)
		}

		c.handle(buf[:n])
	}
}

func (c *RailgunClient) Run() {
	go c.runTUN()
	c.runUDPRecv()
}

func NewRailgunClient () *RailgunClient {
	cli := new(RailgunClient)

	cli.listenIP = net.IPv4(10, 1, 0, 2)
	_, cli.ipNet, _ = net.ParseCIDR("10.1.0.1/24")

	var err error
	listener, err := net.ListenPacket("ip4:1", "")
	check(err)

	cli.rawConn, err = ipv4.NewRawConn(listener)
	check(err)

	udpConn, err := net.Dial("udp", "bjong.me:7000")
	check(err)

	cli.udpConn = udpConn.(*net.UDPConn)
	return cli
}

func main() {
	cli := NewRailgunClient()
	cli.Run()
}
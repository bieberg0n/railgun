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

type RailgunServer struct {
	listenIP net.IP
	ipNet *net.IPNet
	rawConn *ipv4.RawConn
	udpServ *net.UDPConn
	ipSrcMap map[string] *net.UDPAddr
}

func (s *RailgunServer) handle(p []byte, from *net.UDPAddr) {
	h, err := ipv4.ParseHeader(p)
	if err != nil {
		log("parse ip header error:", err)
		return
	}

	if from != nil {
		s.ipSrcMap[h.Src.String()] = from
	}

	//if s.ipNet.Contains(h.Dst) && !h.Dst.Equal(s.listenIP) {
	if h.Dst.Equal(s.listenIP) || !s.ipNet.Contains(h.Dst) {
	//if h.Dst.Equal(s.listenIP) {
		log(h, p)

		body := p[h.Len:]
		err = s.rawConn.WriteTo(h, body, nil)
		if err != nil {
			log("conn write error:", err)
		}

	//} else if h.Dst.Equal(s.ipNet.IP) || !s.ipNet.Contains(h.Dst) {
	} else if s.ipNet.Contains(h.Dst) {
		log(h, p)
		udpAddr, exist := s.ipSrcMap[h.Dst.String()]
		if !exist {
			log("udpaddr no exist")
			return
		}

		log("udpaddr:", udpAddr)
		_, err = s.udpServ.WriteToUDP(p, udpAddr)
		if err != nil {
			log("udp write error:", err)
		}
	}
}

func (s *RailgunServer) runTUN() {
	cfg := water.Config{
		//DeviceType: water.TAP,
		DeviceType: water.TUN,
	}
	cfg.Name = "tun1"

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
		s.handle(packet[:n], nil)
	}
}

func (s *RailgunServer) runUDPServ() {
	buf := make([]byte, 2000)
	for {
		n, from, err := s.udpServ.ReadFromUDP(buf)
		if err != nil {
			log("udp recv error:", err)
		}
		log("recv from:", from)

		s.handle(buf[:n], from)
	}
}

func NewRailgunServer() *RailgunServer {
	server := new(RailgunServer)
	server.listenIP = net.IPv4(10, 1, 0, 1)
	_, server.ipNet, _ = net.ParseCIDR("10.1.0.1/24")
	server.ipSrcMap = make(map[string]*net.UDPAddr)

	var err error
	listener, err := net.ListenPacket("ip4:1", "")
	check(err)

	server.rawConn, err = ipv4.NewRawConn(listener)
	check(err)

	udpAddr, _ := net.ResolveUDPAddr("udp", ":7000")
	server.udpServ, err = net.ListenUDP("udp", udpAddr)
	check(err)

	return server
}

func main() {
	serv := NewRailgunServer()
	go serv.runTUN()
	serv.runUDPServ()
}

package main

import (
	"fmt"
	"github.com/songgao/water"
	"golang.org/x/net/ipv4"
	"net"
	"railgun"
	"strconv"
)

var (
	check = railgun.Check
	log = railgun.Log
)

type RailgunServer struct {
	listenIP net.IP
	ipNet *net.IPNet
	ipSrcMap map[string] *net.UDPAddr
}

func (s *RailgunServer) handle(p []byte, from *net.UDPAddr) {
	h, err := ipv4.ParseHeader(p)
	if err != nil {
		log("parse ip header error:", err)
		return
	}

	//log(h)

	if from != nil {
		s.ipSrcMap[h.Src.String()] = from
	}

	//if s.ipNet.Contains(h.Dst) && !h.Dst.Equal(s.listenIP) {
	if h.Dst.Equal(s.listenIP) {
		log(h, p)
		conn, err := net.DialIP("ip4:"+strconv.Itoa(h.Protocol), &net.IPAddr{IP: h.Src}, &net.IPAddr{IP: h.Dst})
		if err != nil {
			log("dial error:", err)
			return
		}
		defer conn.Close()

		body := p[h.Len:]
		_, err = conn.Write(body)
		if err != nil {
			log("conn write error:", err)
		}

	//} else if h.Dst.Equal(s.ipNet.IP) || !s.ipNet.Contains(h.Dst) {
	} else if s.ipNet.Contains(h.Dst) {
		udpAddr, exist := s.ipSrcMap[h.Dst.String()]
		if !exist {
			return
		}

		conn, err := net.Dial("udp", udpAddr.String())
		if err != nil {
			log("dial udp error:", err)
			return
		}

		defer conn.Close()
		_, err = conn.Write(p)
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
		s.handle(packet[:n], nil)
	}
}

func (s *RailgunServer) runUDPServ() {
	udpAddr, _ := net.ResolveUDPAddr("udp", ":7000")
	udpServ, err := net.ListenUDP("udp", udpAddr)
	check(err)
	defer udpServ.Close()

	buf := make([]byte, 2000)

	for {
		n, from, err := udpServ.ReadFromUDP(buf)
		if err != nil {
			log("udp recv error:", err)
		}

		s.handle(buf[:n], from)
	}
}

func NewRailgunServer() *RailgunServer {
	server := new(RailgunServer)
	server.listenIP = net.IPv4(10, 1, 0, 1)
	_, server.ipNet, _ = net.ParseCIDR("10.1.0.1/24")
	server.ipSrcMap = make(map[string]*net.UDPAddr)
	return server
}

func main() {
	serv := NewRailgunServer()
	go serv.runTUN()
	serv.runUDPServ()
}
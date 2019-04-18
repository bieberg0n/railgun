package railgun

//func handle(ifce *water.Interface, p []byte) {
//	h, err := ipv4.ParseHeader(p)
//	check(err)
//
//	log(h)
//	if h.Version != 4 || h.Protocol != 17 {
//		return
//	}
//
//	udpBody := p[h.Len:]
//	body := udpBody[8:]
//	log("body:", hex.EncodeToString(body))
//
//	h.Dst = net.IP([]byte{10, 1, 0, 1})
//	h.Checksum = 0
//	hB, err := h.Marshal()
//	check(err)
//
//	h.Checksum = int(checkSum(hB))
//	hB, err = h.Marshal()
//	check(err)
//
//	conn, err := net.DialIP("ip4:"+strconv.Itoa(h.Protocol), &net.IPAddr{IP: h.Src}, &net.IPAddr{IP: net.ParseIP("10.1.0.1")})
//	check(err)
//
//	udpBody[6], udpBody[7] = 0x00, 0x00
//	_, err = conn.Write(udpBody)
//
//	//conn, err := net.Dial("ip4:udp", "127.0.0.1")
//	//check(err)
//	//_, err = conn.Write(udpBody)
//
//
//	//listener, err := net.ListenPacket("ip4:udp", "10.1.0.1")
//	//check(err)
//	//defer listener.Close()
//	//
//	//conn, err := ipv4.NewRawConn(listener)
//	//_, err = ifce.Write(append(hB, udpBody...))
//
//	check(err)
//}

//func packUDP(data []byte) []byte {
//	//l := len(data)
//	dataBuf := new(bytes.Buffer)
//	//dataBuf.Write(0x27)
//	//dataBuf.Write([]byte{2333, 3333, byte(l), 0})
//	//dataBuf.Write(data)
//
//	p := dataBuf.Bytes()
//	check := checkSum(p)
//	p[18], p[19] = byte(check>>8&255), byte(check&255)
//	return p
//}

//func test() {
//	buff := []byte("test")
//
//	//目的IP
//	dst := net.IPv4(127, 0, 0, 1)
//	//源IP
//	src := net.IPv4(127, 0, 0, 1)
//	//填充ip首部
//	iph := &ipv4.Header{
//		Version:  ipv4.Version,
//		//IP头长一般是20
//		Len:      ipv4.HeaderLen,
//		TOS:      0x00,
//		//buff为数据
//		TotalLen: ipv4.HeaderLen + len(buff),
//		TTL:      64,
//		Flags:    ipv4.DontFragment,
//		FragOff:  0,
//		Protocol: 17,
//		Checksum: 0,
//		Src:      src,
//		Dst:      dst,
//	}
//
//	h, err := iph.Marshal()
//	check(err)
//	//计算IP头部校验值
//	iph.Checksum = int(checkSum(h))
//
//	//填充udp首部
//	//udp伪首部
//	udph := make([]byte, 8)
//	////源ip地址
//	//udph[0], udph[1], udph[2], udph[3] = src[12], src[13], src[14], src[15]
//	////目的ip地址
//	//udph[4], udph[5], udph[6], udph[7] = dst.IP[12], dst.IP[13], dst.IP[14], dst.IP[15]
//	////协议类型
//	//udph[8], udph[9] = 0x00, 0x11
//	////udp头长度
//	//udph[10], udph[11] = 0x00, byte(len(buff)+8)
//	////下面开始就真正的udp头部
//	////源端口号
//	udph[0], udph[1] = 0x09, 0x1D
//	udph[2], udph[3] = 0x0D, 0x05
//	////目的端口号
//	//udph[14], udph[15] = 0x17, 0x70
//	////udp头长度
//	data := []byte("test")
//
//	udph[4], udph[5] = 0x00, byte(len(data)+8)
//	////校验和
//	udph[6], udph[7] = 0x00, 0x00
//	////计算校验值
//
//	d := append(udph, data...)
//
//	c := checkSum(d)
//	d[6], d[7] = byte(c>>8&255), byte(c&255)
//	//
//
//
//	//listener, err := net.ListenPacket("ip4:udp", "192.168.1.140")
//	//check(err)
//	//defer listener.Close()
//	//r, err := ipv4.NewRawConn(listener)
//	//conn, err := net.Dial("ip4:udp", "127.0.0.1")
//	//check(err)
//	//conn, err := net.DialUDP("udp",
//	//	&net.UDPAddr{IP: net.IPv4(192, 168, 1, 140), Port: 0},
//	//	&net.UDPAddr{IP: net.IPv4(127, 0, 0, 1), Port: 3333},
//	//)
//
//	//check(err)
//	//defer conn.Close()
//
//	//err = r.WriteTo(iph, d, nil)
//	//check(err)
//	fd, _ := syscall.Socket(syscall.AF_INET, syscall.SOCK_RAW, syscall.IPPROTO_RAW)
//	addr := syscall.SockaddrInet4{
//		Port: 0,
//		Addr: [4]byte{127, 0, 0, 1},
//	}
//
//	ipHB, _ := iph.Marshal()
//
//	err = syscall.Sendto(fd, append(ipHB, d...), 0, &addr)
//	check(err)
//}

//func test1() {
//	//填充udp首部
//	//udp伪首部
//	udph := make([]byte, 8)
//	////源ip地址
//	//udph[0], udph[1], udph[2], udph[3] = src[12], src[13], src[14], src[15]
//	////目的ip地址
//	//udph[4], udph[5], udph[6], udph[7] = dst.IP[12], dst.IP[13], dst.IP[14], dst.IP[15]
//	////协议类型
//	//udph[8], udph[9] = 0x00, 0x11
//	////udp头长度
//	//udph[10], udph[11] = 0x00, byte(len(buff)+8)
//	////下面开始就真正的udp头部
//	////源端口号
//	udph[0], udph[1] = 0x09, 0x1D
//	udph[2], udph[3] = 0x0D, 0x05
//	////目的端口号
//	//udph[14], udph[15] = 0x17, 0x70
//	////udp头长度
//	data := []byte("test")
//
//	udph[4], udph[5] = 0x00, byte(len(data)+8)
//	////校验和
//	udph[6], udph[7] = 0x00, 0x00
//	////计算校验值
//
//	d := append(udph, data...)
//
//	//c := checkSum1(d)
//	//d[6], d[7] = byte(c>>8&255), byte(c&255)
//	//
//	conn, err := net.Dial("ip4:udp", "127.0.0.1")
//	check(err)
//	_, err = conn.Write(d)
//	check(err)
//
//}

//func main() {
//	//test1()
//
//	cfg := water.Config{
//		//DeviceType: water.TAP,
//		DeviceType: water.TUN,
//	}
//	cfg.Name = "tun0"
//
//	ifce, err := water.New(cfg)
//	check(err)
//
//	fmt.Printf("Interface Name: %s\n", ifce.Name())
//
//	packet := make([]byte, 2000)
//	for {
//		n, err := ifce.Read(packet)
//		if err != nil {
//			log("error:", err)
//		}
//		fmt.Printf("Packet Received: % x\n", packet[:n])
//		handle(ifce, packet[:n])
//	}
//}

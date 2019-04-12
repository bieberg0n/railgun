package main

import (
	"fmt"
	"github.com/songgao/packets/ethernet"
	"github.com/songgao/water"
)

func main() {
	cfg := water.Config{
		DeviceType: water.TAP,
	}
	cfg.Name = "0_0"

	ifce, err := water.New(cfg)
	if err != nil {
		check(err)
	}

	var frame ethernet.Frame

	for {
		frame.Resize(1500)
		n, err := ifce.Read([]byte(frame))
		if err != nil {
			log(err)
		}

		frame := frame[:n]
		log(fmt.Printf("Dst: %s", frame.Destination()))
		log(fmt.Printf("Src: %s", frame.Source()))
		log(fmt.Printf("Ethertype: %x", frame.Ethertype()))
		log(fmt.Printf("Payload: %x", frame.Payload()))
	}
}

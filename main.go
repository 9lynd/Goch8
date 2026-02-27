package main

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/9lynd/Goch8/core"
)

func main() {
	display := &core.Display{}
	keyboard := &core.Keyboard{}
	cpu := core.NewCPU(display, keyboard)

	// Load the ROM
	romData, err := ioutil.ReadFile("roms/Airplane.ch8")
	if err != nil {
		log.Fatalf("Failed to load ROM: %v", err)
	}

	err = cpu.LoadROM(romData)
	if err != nil {
		log.Fatalf("Error loading ROM: %v", err)
	}

	// Check memory at address 0x200
	for i := 0; i < 10; i++ { // Print first 10 bytes for verification
		fmt.Printf("Memory[0x%03X]: 0x%02X\n", 0x200+i, cpu.Memory[0x200+i])
	}
}

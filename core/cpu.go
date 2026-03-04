package core

import (
	"fmt"
	"math/rand"
)

var fontset = [80]byte{
	0xF0, 0x90, 0x90, 0x90, 0xF0, // 0
	0x20, 0x60, 0x20, 0x20, 0x70, // 1
	0xF0, 0x10, 0xF0, 0x80, 0xF0, // 2
	0xF0, 0x10, 0xF0, 0x10, 0xF0, // 3
	0x90, 0x90, 0xF0, 0x10, 0x10, // 4
	0xF0, 0x80, 0xF0, 0x10, 0xF0, // 5
	0xF0, 0x80, 0xF0, 0x90, 0xF0, // 6
	0xF0, 0x10, 0x20, 0x40, 0x40, // 7
	0xF0, 0x90, 0xF0, 0x90, 0xF0, // 8
	0xF0, 0x90, 0xF0, 0x10, 0xF0, // 9
	0xF0, 0x90, 0xF0, 0x90, 0x90, // A
	0xE0, 0x90, 0xE0, 0x90, 0xE0, // B
	0xF0, 0x80, 0x80, 0x80, 0xF0, // C
	0xE0, 0x90, 0x90, 0x90, 0xE0, // D
	0xF0, 0x80, 0xF0, 0x80, 0xF0, // E
	0xF0, 0x80, 0xF0, 0x80, 0x80, // F
}

const fontsetStart = 0x050

type CPU struct {
	V      [16]byte
	I      uint16
	PC     uint16
	SP     uint8
	Stack  [16]uint16
	Memory [4096]byte

	Display  *Display
	Keyboard *Keyboard

	DelayTimer byte
	SoundTimer byte
}

func NewCPU(display *Display, keyboard *Keyboard) *CPU {
	cpu := &CPU{
		PC:       0x200,
		Display:  display,
		Keyboard: keyboard,
	}

	for i, b := range fontset {
		cpu.Memory[fontsetStart+i] = b
	}

	return cpu
}

func (c *CPU) LoadROM(data []byte) error {
	if len(data) > len(c.Memory)-0x200 {
		return fmt.Errorf("ROM too large: %d bytes (max %d)", len(data), len(c.Memory)-0x200)
	}
	copy(c.Memory[0x200:], data)
	return nil
}

func (c *CPU) TickTimers() {
	if c.DelayTimer > 0 {
		c.DelayTimer--
	}
	if c.SoundTimer > 0 {
		c.SoundTimer--
	}
}

func (c *CPU) Cycle() error {
	// 2 bytes, big-endian
	hi := c.Memory[c.PC]
	lo := c.Memory[c.PC+1]
	opcode := uint16(hi)<<8 | uint16(lo)
	c.PC += 2
	return c.execute(opcode)
}

func (c *CPU) execute(opcode uint16) error {
	nnn := opcode & 0x0FFF
	n := byte(opcode & 0x000F)
	x := byte((opcode & 0x0F00) >> 8)
	y := byte((opcode & 0x00F0) >> 4)
	kk := byte(opcode & 0x00FF)

	switch opcode & 0xF000 {
	case 0x0000:
		switch opcode {
		case 0x00E0:
			c.Display.Clear()
		case 0x00EE:
			// RET - return from subroutine
			c.SP--
			c.PC = c.Stack[c.SP]
		default:
			return fmt.Errorf("unknown opcode: 0x%04X", opcode)
		}

	case 0x1000:
		c.PC = nnn

	case 0x2000:
		c.Stack[c.SP] = c.PC
		c.SP++
		c.PC = nnn

	case 0x3000:
		if c.V[x] == kk {
			c.PC += 2
		}

	case 0x4000:
		if c.V[x] != kk {
			c.PC += 2
		}

	case 0x5000:
		// SE Vx, Vy - skip if Vx == Vy
		if c.V[x] == c.V[y] {
			c.PC += 2
		}

	case 0x6000:
		// LD Vx, kk - set Vx = kk
		c.V[x] = kk

	case 0x7000:
		// ADD Vx, kk - set Vx = Vx + kk
		c.V[x] += kk

	case 0x8000:
		switch n {
		case 0x0:
			c.V[x] = c.V[y]

		case 0x1:
			c.V[x] |= c.V[y]

		case 0x2:
			c.V[x] &= c.V[y]

		case 0x3:
			c.V[x] ^= c.V[y]

		case 0x4:
			sum := uint16(c.V[x]) + uint16(c.V[y])
			c.V[0xF] = 0
			if sum > 0xFF {
				c.V[0xF] = 1
			}
			c.V[x] = byte(sum)

		case 0x5:
			// SUB Vx, Vy - with borrow (VF=1 if Vx > Vy, no borrow)
			c.V[0xF] = 0
			if c.V[x] > c.V[y] {
				c.V[0xF] = 1
			}
			c.V[x] -= c.V[y]

		case 0x6:
			c.V[0xF] = c.V[x] & 0x1
			c.V[x] >>= 1

		case 0x7:
			c.V[0xF] = 0
			if c.V[y] > c.V[x] {
				c.V[0xF] = 1
			}
			c.V[x] = c.V[y] - c.V[x]

		case 0xE:
			c.V[0xF] = (c.V[x] & 0x80) >> 7
			c.V[x] <<= 1

		default:
			return fmt.Errorf("unknown opcode: 0x%04X", opcode)
		}

	case 0x9000:
		if c.V[x] != c.V[y] {
			c.PC += 2
		}

	case 0xA000:
		c.I = nnn

	case 0xB000:
		c.PC = nnn + uint16(c.V[0])

	case 0xC000:
		c.V[x] = byte(rand.Intn(256)) & kk

	case 0xD000:
		// DRW Vx, Vy, n - draw n-byte sprite at (Vx, Vy)
		c.V[0xF] = 0
		for row := byte(0); row < n; row++ {
			spriteByte := c.Memory[c.I+uint16(row)]
			for col := byte(0); col < 8; col++ {
				if spriteByte&(0x80>>col) != 0 {
					erased := c.Display.TogglePixel(c.V[x]+col, c.V[y]+row)
					if erased {
						c.V[0xF] = 1
					}
				}
			}
		}

	case 0xE000:
		switch kk {
		case 0x9E: // SKP Vx - skip if key Vx is pressed
			if c.Keyboard.IsPressed(c.V[x]) {
				c.PC += 2
			}
		case 0xA1: // SKNP Vx - skip if key Vx is NOT pressed
			if !c.Keyboard.IsPressed(c.V[x]) {
				c.PC += 2
			}
		default:
			return fmt.Errorf("unknown opcode: 0x%04X", opcode)
		}

	case 0xF000:
		switch kk {
		case 0x07:
			c.V[x] = c.DelayTimer

		case 0x0A:
			pressed, key := c.Keyboard.WaitForKey()
			if pressed {
				c.V[x] = key
			} else {
				c.PC -= 2 // re-execute this opcode next cycle
			}

		case 0x15:
			c.DelayTimer = c.V[x]

		case 0x18:
			c.SoundTimer = c.V[x]

		case 0x1E:
			c.I += uint16(c.V[x])

		case 0x29:
			c.I = uint16(fontsetStart) + uint16(c.V[x])*5

		case 0x33:
			c.Memory[c.I] = c.V[x] / 100
			c.Memory[c.I+1] = (c.V[x] / 10) % 10
			c.Memory[c.I+2] = c.V[x] % 10

		case 0x55:
			for i := byte(0); i <= x; i++ {
				c.Memory[c.I+uint16(i)] = c.V[i]
			}

		case 0x65:
			for i := byte(0); i <= x; i++ {
				c.V[i] = c.Memory[c.I+uint16(i)]
			}

		default:
			return fmt.Errorf("unknown opcode: 0x%04X", opcode)
		}

	default:
		return fmt.Errorf("unknown opcode: 0x%04X", opcode)
	}

	return nil
}

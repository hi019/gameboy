package emu

import (
	"emu/constants"
)

type instructionHandler func(opcode uint16, cpu *CPU)
type instructionMap map[uint16]instructionHandler
type instructionMainMap map[uint16]func(opcode uint16) instructionHandler

var MainMap = instructionMainMap{
	0x0: func(opcode uint16) instructionHandler {
		return OP0[opcode&0xF]
	},
	0x1: func(opcode uint16) instructionHandler {
		return OP[opcode&0xF000>>12]
	},
	0x2: func(opcode uint16) instructionHandler {
		return OP[opcode&0xF000>>12]
	},
	0x3: func(opcode uint16) instructionHandler {
		return OP[opcode&0xF000>>12]
	},
	0x4: func(opcode uint16) instructionHandler {
		return OP[opcode&0xF000>>12]
	},
	0x5: func(opcode uint16) instructionHandler {
		return OP[opcode&0xF000>>12]
	},
	0x6: func(opcode uint16) instructionHandler {
		return OP[opcode&0xF000>>12]
	},
	0x7: func(opcode uint16) instructionHandler {
		return OP[opcode&0xF000>>12]
	},
	0x8: func(opcode uint16) instructionHandler {
		return OP8[opcode&0xF]
	},
	0x9: func(opcode uint16) instructionHandler {
		return OP[opcode&0xF000>>12]
	},
	0xA: func(opcode uint16) instructionHandler {
		return OP[opcode&0xF000>>12]
	},
	0xB: func(opcode uint16) instructionHandler {
		return OP[opcode&0xF000>>12]
	},
	0xC: func(opcode uint16) instructionHandler {
		return OP[opcode&0xF000>>12]
	},
	0xD: func(opcode uint16) instructionHandler {
		return OP[opcode&0xF000>>12]
	},
	0xE: func(opcode uint16) instructionHandler {
		return OPE[opcode&0xFF]
	},
	0xF: func(opcode uint16) instructionHandler {
		return OPE[opcode&0xFF]
	},
}

var OP0 = instructionMap{
	// CLS (0x00E0)
	0x0: func(_ uint16, c *CPU) {
		c.Video = [constants.VideoHeight * constants.VideoWidth]byte{}
	},
	// RET (0x00EE)
	0xE: func(opcode uint16, c *CPU) {
		c.SP--
		c.PC = c.Stack[c.SP]
	},
}

var OP = instructionMap{
	// JP (1nnn)
	0x1: func(opcode uint16, c *CPU) {
		address := opcode & 0x0FFF
		c.PC = address
	},
	// CALL (2nnn)
	0x2: func(opcode uint16, c *CPU) {
		address := opcode & 0x0FFF
		c.Stack[c.SP] = c.PC
		c.SP++
		c.PC = address
	},
	// SE Vx, byte (3xkk)
	0x3: func(opcode uint16, c *CPU) {
		vx := (opcode & 0x0F00) >> 8
		val := uint8(opcode & 0x00FF)

		if c.Registers[vx] == val {
			c.PC += 2
		}
	},
	// SNE Vx, byte (4xkk)
	0x4: func(opcode uint16, c *CPU) {
		vx := (opcode & 0x0F00) >> 8
		val := uint8(opcode & 0x00FF)

		if c.Registers[vx] != val {
			c.PC += 2
		}
	},
	// SE Vx, Vy (5xy0)
	0x5: func(opcode uint16, c *CPU) {
		vx := (opcode & 0x0F00) >> 8
		vy := (opcode & 0x00F0) >> 4

		if c.Registers[vx] == c.Registers[vy] {
			c.PC += 2
		}
	},
	// LD Vx, byte (6xkk)
	0x6: func(opcode uint16, c *CPU) {
		vx := (opcode & 0x0F00) >> 8
		val := uint8(opcode & 0x00FF)

		c.Registers[vx] = val
	},
	// ADD Vx, byte (7xkk)
	0x7: func(opcode uint16, c *CPU) {
		vx := (opcode & 0x0F00) >> 8
		val := uint8(opcode & 0x00FF)

		c.Registers[vx] += val
	},
	// SNE Vx, Vy (9xy0)
	0x9: func(opcode uint16, c *CPU) {
		vx := (opcode & 0x0F00) >> 8
		vy := (opcode & 0x00F0) >> 4

		if c.Registers[vx] != c.Registers[vy] {
			c.PC += 2
		}
	},
	// LD I, addr (Annn)
	0xA: func(opcode uint16, c *CPU) {
		c.I = opcode & 0x0FFF
	},
	// JP V0, addr (Bnnn)
	0xB: func(opcode uint16, c *CPU) {
		address := opcode & 0x0FFF
		c.PC = address + uint16(c.Registers[0])
	},
	// RND Vx, byte (Cxkk)
	0xC: func(opcode uint16, c *CPU) {
		vx := (opcode & 0x0F00) >> 8
		val := uint8(opcode & 0x00FF)

		c.Registers[vx] = c.Random() & val
	},
	// DRW Vx, Vy, nibble (Dxyn)
	// TODO understand this
	0xD: func(opcode uint16, c *CPU) {
		vx := (opcode & 0x0F00) >> 8
		vy := (opcode & 0x00F0) >> 4
		height := uint8(opcode & 0x000F)

		xPos := c.Registers[vx] % constants.VideoWidth
		yPos := c.Registers[vy] % constants.VideoHeight

		c.Registers[0xF] = 0

		for row := 0; row < int(height); row++ {
			spriteByte := c.Memory[c.I+uint16(row)]
			for col := 0; col < 8; col++ {
				spritePixel := spriteByte & (0x80 >> col)
				vidIndex := (int(yPos)+row)*constants.VideoWidth + (int(xPos) + col)

				if vidIndex >= len(c.Video) {
					break
				}

				screenPixel := c.Video[vidIndex]

				if spritePixel > 0 {
					if screenPixel == 0xFF {
						c.Video[vidIndex] = 0
						c.Registers[0xF] = 1
					} else {
						c.Video[vidIndex] = 0xFF
					}
				}
			}
		}
	},
}

var OP8 = instructionMap{
	// LD Vx, Vy (8xy0)
	0x0: func(opcode uint16, c *CPU) {
		vx := (opcode & 0x0F00) >> 8
		vy := (opcode & 0x00F0) >> 4

		c.Registers[vx] = c.Registers[vy]
	},
	// OR Vx, Vy (8xy1)
	0x1: func(opcode uint16, c *CPU) {
		vx := (opcode & 0x0F00) >> 8
		vy := (opcode & 0x00F0) >> 4

		c.Registers[vx] |= c.Registers[vy]
	},
	// AND Vx, Vy (8xy2)
	0x2: func(opcode uint16, c *CPU) {
		vx := (opcode & 0x0F00) >> 8
		vy := (opcode & 0x00F0) >> 4

		c.Registers[vx] &= c.Registers[vy]
	},
	// XOR Vx, Vy (8xy3)
	0x3: func(opcode uint16, c *CPU) {
		vx := (opcode & 0x0F00) >> 8
		vy := (opcode & 0x00F0) >> 4

		c.Registers[vx] ^= c.Registers[vy]
	},
	// ADD Vx, Vy (8xy4)
	0x4: func(opcode uint16, c *CPU) {
		vx := (opcode & 0x0F00) >> 8
		vy := (opcode & 0x00F0) >> 4

		sum := c.Registers[vx] + c.Registers[vy]
		// check for overflow
		if sum > 255 {
			c.Registers[0xF] = 1
		} else {
			c.Registers[0xF] = 0
		}

		c.Registers[vx] = sum & 0xFF
	},
	// SUB Vx, Vy (8xy5)
	0x5: func(opcode uint16, c *CPU) {
		vx := (opcode & 0x0F00) >> 8
		vy := (opcode & 0x00F0) >> 4

		if c.Registers[vx] > c.Registers[vy] {
			c.Registers[0xF] = 1
		} else {
			c.Registers[0xF] = 0
		}

		c.Registers[vx] -= c.Registers[vy]
	},
	// SHR Vx (8xy6)
	0x06: func(opcode uint16, c *CPU) {
		vx := (opcode & 0x0F00) >> 8
		c.Registers[0xF] = c.Registers[vx] & 0x1
		c.Registers[vx] >>= 1
	},
	// SUBN Vx, Vy (8xy7)
	0x07: func(opcode uint16, c *CPU) {
		vx := (opcode & 0x0F00) >> 8
		vy := (opcode & 0x00F0) >> 4

		if c.Registers[vy] > c.Registers[vx] {
			c.Registers[0xF] = 1
		} else {
			c.Registers[0xF] = 0
		}

		c.Registers[vx] = c.Registers[vy] - c.Registers[vx]
	},
	// SHL Vx {, Vy} (8xyE)
	0x0E: func(opcode uint16, c *CPU) {
		vx := (opcode & 0x0F00) >> 8

		c.Registers[0xF] = (c.Registers[vx] & 0x80) >> 7
		c.Registers[vx] <<= 1
	},
}

var OPE = instructionMap{
	// SKP Vx (Ex9E)
	0x9E: func(opcode uint16, c *CPU) {
		vx := (opcode & 0x0F00) >> 8
		key := c.Registers[vx]
		if c.Keypad[key] == 1 {
			c.PC += 2
		}
	},
	// SKNP Vx (ExA1)
	0xA1: func(opcode uint16, c *CPU) {
		vx := (opcode & 0x0F00) >> 8
		key := c.Registers[vx]
		if c.Keypad[key] == 0 {
			c.PC += 2
		}
	},
	// LD Vx, DT (Fx07)
	0x07: func(opcode uint16, c *CPU) {
		vx := (opcode & 0x0F00) >> 8
		c.Registers[vx] = c.DT
	},
	// LD Vx, K (Fx0A)
	0x0A: func(opcode uint16, c *CPU) {
		vx := (opcode & 0x0F00) >> 8
		for i := 0; i < len(c.Keypad); i++ {
			if c.Keypad[i] == 1 {
				c.Registers[vx] = uint8(i)
				return
			}
		}

		// no key pressed, decrement PC by 2 which has the effect of repeating the instruction
		c.PC -= 2
	},
	// LD DT, Vx (Fx15)
	0x15: func(opcode uint16, c *CPU) {
		vx := (opcode & 0x0F00) >> 8
		c.DT = c.Registers[vx]
	},
	// LD ST, Vx (Fx18)
	0x18: func(opcode uint16, c *CPU) {
		vx := (opcode & 0x0F00) >> 8
		c.ST = c.Registers[vx]
	},
	// ADD I, Vx (Fx1E)
	0x1E: func(opcode uint16, c *CPU) {
		vx := (opcode & 0x0F00) >> 8
		c.I += uint16(c.Registers[vx])
	},
	// LD F, Vx (Fx29)
	0x29: func(opcode uint16, c *CPU) {
		vx := (opcode & 0x0F00) >> 8
		digit := c.Registers[vx]

		c.I = FontsetStartAddress + uint16(digit*5)
	},
	// LD B, Vx (Fx33)
	0x33: func(opcode uint16, c *CPU) {
		vx := (opcode & 0x0F00) >> 8
		value := c.Registers[vx]

		// hunderds
		c.Memory[c.I] = value / 100
		// tens
		c.Memory[c.I+1] = (value / 10) % 10
		// ones
		c.Memory[c.I+2] = value % 10
	},
	// LD [I], Vx (Fx55)
	0x55: func(opcode uint16, c *CPU) {
		vx := (opcode & 0x0F00) >> 8
		for i := uint16(0); i <= vx; i++ {
			c.Memory[c.I+i] = c.Registers[i]
		}
	},
	// LD Vx, [I] (Fx65)
	0x65: func(opcode uint16, c *CPU) {
		vx := (opcode & 0x0F00) >> 8
		for i := uint16(0); i <= vx; i++ {
			c.Registers[i] = c.Memory[c.I+i]
		}
	},
}

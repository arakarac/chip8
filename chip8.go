package main

import (
	"fmt"
	"time"
)

const (
	mSize = 4096
	addrSize = 16
	hz = time.Duration(600)
	width = 64 
	height = 32
)

type chip8 struct {
  mem        [mSize]byte     // memory
	v          [16]uint8        // variable registers, V0 through VF
	ir         uint16          // index register, used to point at locations in memory
	pc         uint16          // program counter, points at the current instruction in memory
	keys       [16]uint16
	display    [height][width]uint8  // 64x32
	stack      [addrSize]uint16 // call functions and return from them
	op         uint16          // opcode
	sp         uint16          // stack pointer 
	dTimer     uint16          // delay timer 
	sTimer     uint16          // sound timer
	canDraw    bool
}

func (c *chip8) increment() {
  c.pc += 2
}

func Init() chip8 {
	cpu := chip8{
		canDraw: true,
		pc: 0x200,
	}
  cpu.loadFont()
	return cpu
}

func (c *chip8) cycle() {
  c.op = uint16(c.mem[c.pc])<<8 | uint16(c.mem[c.pc+1])

	x := int((c.op & 0x0F00) >> 8)
	y := int((c.op & 0x00F0) >> 4)
	n := int(c.op & 0x000F)
  nn := uint8(c.op & 0x00FF)
	nnn := c.op &  0x0FFF

	switch c.op & 0xF000 {
	case 0x0000:    //00E0
    switch c.op & 0x000F {
		case 0x0000:
      c.clearDisplay()
      c.increment()
		  c.canDraw = true
    case 0x000E: //0x00EE
			c.sp--
      c.pc = c.stack[c.sp]
      c.increment()
		}
	case 0x1000:    //1NNN
		c.pc = nnn
	case 0x6000:    //6XNN
    c.v[x] = uint8(nn)
		c.increment()
	case 0x7000:    //7XNN
    c.v[x] += uint8(nn)
	  c.increment()
	case 0xA000:    //ANNN
    c.ir = nnn
		c.increment()
	case 0xD000:    //DXYN, n = height, number of rows / p = pixel
    c.v[0xF] = 0
		vy := int(c.v[y] % height)
		vx := int(c.v[x] % width)
		for j := 0; j < n; j++ {
			p := c.mem[c.ir + uint16(j)]
			for i := uint8(0); i < 8; i++ {
				if (p & (0x80 >> i)) != 0 {
					py := (vy + j) % height
					px := (vx + int(i)) % width
					//sp := py * width + px
          if c.display[py][px] == 1 {
				  c.v[0xF] = 1 // collision detection
				  }
				  c.display[py][px] ^= 1
				}
			}
		}
    c.canDraw = true
		c.increment()
	}

  if c.dTimer > 0 {
  		c.dTimer--
  }
  if c.sTimer > 0 {
		if c.sTimer == 1 {
      fmt.Println("BEEP")
		}
	  c.sTimer--
	}

	//fmt.Println(c.op)
}

func (c *chip8) buffer() [32][64]uint8 {
	return c.display
}

func (c *chip8) clearDisplay() {
	for i := 0; i < len(c.display); i++ {
    for j := 0; j < len(c.display[i]); j++ {
			c.display[i][j] = 0x0
		}
	}  
}

func (c *chip8) draw() bool {
	b := c.canDraw
	c.canDraw = false
	return b
}


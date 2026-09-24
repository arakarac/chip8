package main

import (
	"fmt"
	"math/rand"
	"time"
)

const (
	mSize = 4096
	addrSize = 16
	cycles = 10
	seconds = time.Second / 60 
	width = 64 
	height = 32
)

type chip8 struct {
  mem        [mSize]byte     // memory
	v          [16]uint8        // variable registers, V0 through VF
	ir         uint16          // index register, used to point at locations in memory
	pc         uint16          // program counter, points at the current instruction in memory
	key        [16]int
	display    [height][width]uint8  // 64x32
	stack      [addrSize]uint16 // call functions and return from them
	op         uint16          // opcode
	sp         uint16          // stack pointer 
	dTimer     uint16          // delay timer 
	sTimer     uint16          // sound timer
	canDraw    bool
	running    bool
}

func (c *chip8) increment() {
  c.pc += 2
}

func Init() chip8 {
	cpu := chip8{
		canDraw: true,
		pc: 0x200,
		running: true,
	}
  cpu.loadFont()
	return cpu
}

func (c *chip8) cycle() {
  c.op = uint16(c.mem[c.pc])<<8 | uint16(c.mem[c.pc+1])

	x := int((c.op & 0x0F00) >> 8)
	y := int((c.op & 0x00F0) >> 4)
	f := int(0xF)
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
	case 0x2000:
    c.stack[c.sp] = c.pc
		c.sp++
    c.pc = nnn
	case 0x1000:    //1NNN
		c.pc = nnn
	case 0x3000:
		c.increment()
    if c.v[x] == nn {
			c.increment()
	  }
	case 0x4000:   //4XNN
		c.increment()
	  if c.v[x] != nn {
			c.increment()
		}
	case 0x5000:   //5XY0
		c.increment()
    if c.v[x] == c.v[y] {
			c.increment()
		}
	case 0x6000:    //6XNN
    c.v[x] = uint8(nn)
		c.increment()
	case 0x7000:    //7XNN
    c.v[x] += uint8(nn)
	  c.increment()
	case 0x8000:     //math ops
    switch c.op & 0x000F {
		case 0x0000:
      c.v[x] = c.v[y]
			c.increment()
		case 0x0001:
      c.v[x] |= c.v[y]
			c.increment()
		case 0x0002:
      c.v[x] &= c.v[y]
			c.increment()
		case 0x0003:
      c.v[x] ^= c.v[y]
			c.increment()
		case 0x0004:
      c.v[x] += c.v[y]
			c.increment()
		case 0x0005:
      c.v[x] -= c.v[y]
			c.increment()
		case 0x0006:
      c.v[x] = c.v[x] >> 1
			c.v[f] = c.v[x] & 0x1
			c.increment()
		case 0x0007:
			if c.v[y] >= c.v[x] {
				c.v[f] = 1
			} else {
				c.v[f] = 0
			}
      c.v[x] = c.v[x] - c.v[y]
			c.increment()
		case 0x000E:
			c.v[f] = c.v[x] >> 7
      c.v[x] = c.v[x] << 1
      c.increment()
		}
	case 0x9000:    //9XY0
		c.increment()
    if c.v[x] != c.v[y] {
			c.increment()
		}
	case 0xA000:    //ANNN
    c.ir = nnn
		c.increment()
	case 0xB000:    //BNNN
    c.pc = uint16(c.v[0x0]) + nnn
	case 0xC000:    //CXNN
    c.v[x] = uint8(rand.Intn(255)) & nn
		c.increment()
	case 0xD000:    //DXYN, n = height, number of rows, N / p = pixel
    c.v[0xF] = 0
		vy := int(c.v[y] % height)
		vx := int(c.v[x] % width)
		for j := 0; j < n; j++ {
			p := c.mem[c.ir + uint16(j)]
			for i := uint8(0); i < 8; i++ {
				if (p & (0x80 >> i)) != 0 {
					py := (vy + j) % height
					px := (vx + int(i)) % width
          if c.display[py][px] == 1 {
				  c.v[0xF] = 1 // collision detection
				  }
				  c.display[py][px] ^= 1
				}
			}
		}
    c.canDraw = true
		c.increment()
	case 0xE000: 
	switch c.op & 0x00FF {   //input ops
	  case 0x009E:  //EX9E
      c.increment()
		  if c.key[c.v[x]] == 1 {
				c.increment()
			}
		case 0x00A1:  //EXA1
      c.increment()
		  if c.key[c.v[x]] == 0 {
				c.increment()
			}
		}
	case 0xF000:
	switch c.op & 0x00FF {
		case 0x0007:  //FX07
      c.v[x] = uint8(c.dTimer)
		  c.increment()
		case 0x000A:     //FX0A !!! NOT HALTING !!! :(
			for i := 0; i < len(c.key); i++ {
				if c.key[i] == 1 {
					c.v[x] = uint8(i)
				}
				if c.key[i] == 0 {
					fmt.Println("no key pressed")
				}
			}
      c.increment()
		case 0x0015:   //FX15
      c.dTimer = uint16(c.v[x])
		  c.increment()
		case 0x0018:  //FX18
      c.dTimer = uint16(c.v[x])
		  c.increment()
		case 0x001E:  //FX1E
      c.ir += uint16(c.v[x])
		  c.increment()
		case 0x0029:  //FX29
      c.ir = uint16(c.v[x] * 0x5)
		  c.increment()
		case 0x0033:   //FX33
      c.mem[c.ir] = c.v[x] / 100
			c.mem[c.ir + 1] = (c.v[x] / 10) % 10
      c.mem[c.ir + 2] = (c.v[x] % 100) % 10
			c.increment()
		case 0x0055:   //FX55
      for i := 0; i <= x ; i++ {
        c.mem[c.ir + uint16(i)] = c.v[i]
			}
			c.ir = uint16(x) + 1
		  c.increment()
		case 0x0065:   //FX65
			for i := 0; i <= x; i++ {
        c.v[i] = c.mem[c.ir + uint16(i)]
			}
      c.ir = uint16(x) + 1
		  c.increment()
		}
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

func (c *chip8) isRunning() bool {
	b := c.running
	c.running = false
	return b
}

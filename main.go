package main

import (
	"flag"
	"fmt"
	"os"
  //"github.com/veandco/go-sdl2/sdl"
)

var (
	ROM string
)

func (c *chip8) load(path string) {
  file, err := os.OpenFile(path, os.O_RDONLY, 0)
  if err != nil {
		panic("failed loading rom")
	} 
	defer file.Close()

  if path == "" {
		fmt.Println("empty/invalid rom")
		os.Exit(1)
	}

	info, err := file.Stat()
	if err != nil {
		panic(err)
	}

  buffer := make([]byte, info.Size())

  _, err = file.Read(buffer)
  if err != nil {
    panic(err)
  }

	bufferSize := len(buffer)
	//mem := c.mem

	for i := 0; i < bufferSize ; i++ {
		c.mem[i + 512] = buffer[i]
	}
}

func main() {
	flag.StringVar(&ROM, "path","", "path to rom")
	flag.Parse()
	cpu := Init()
	sdlInit()
  cpu.load(ROM)
	//f := os.Args[2]
	//fmt.Println("rom path: " + f)
	for {
	  cpu.cycle()
	  if cpu.draw() {
      cpu.render()
	  }
	}
  //sdl.Delay(1000 / 60)
}


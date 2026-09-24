package main

import (
	"flag"
	"fmt"
	"os"
	"time"
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
	if cpu.isRunning() {
		for {
			cpu.getInput()
			start := time.Now()
			for i := 0; i < cycles; i++ {
        cpu.cycle()
			  if cpu.draw() {
				  cpu.render()
			  }
			}
      past := time.Since(start)
			if past < seconds  {
				time.Sleep(seconds - past)
			}

		}
	} else {
		  os.Exit(0)
	}
}


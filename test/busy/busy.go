package main

import (
	"math/rand"
	"time"
)

func main() {

	sizeMB := 100
	mem := make([]byte, sizeMB*1024*1024)

	for {
		for i := 0; i < len(mem); i += 4096 {
			mem[i] = byte(rand.Intn(256))
		}

		work := rand.Intn(5000000) + 1000000
		x := 0
		for i := range work {
			x += i % 3
		}

		time.Sleep(time.Millisecond * time.Duration(rand.Intn(50)))
	}
}

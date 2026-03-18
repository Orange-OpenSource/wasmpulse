package main

import (
	"fmt"
	"os"
	"time"
)

func main() {
	filename := "output.txt"
	const maxSize = 100 * 1024 * 1024 // 100 MB cap

	if _, err := os.Stat(filename); os.IsNotExist(err) {
		init := []byte("start\n")
		if err := os.WriteFile(filename, init, 0644); err != nil {
			fmt.Println("init write error:", err)
			return
		}
	}

	for {
		data, err := os.ReadFile(filename)
		if err != nil {
			fmt.Println("read error:", err)
			return
		}

		if len(data) >= maxSize {
			fmt.Printf("size capped at %d bytes\n", len(data))
			time.Sleep(3 * time.Second)
			continue
		}

		newLine := fmt.Sprintf("tick %s\n", time.Now().Format(time.RFC3339))
		combined := append(data, data...)
		combined = append(combined, []byte(newLine)...)

		if len(combined) > maxSize {
			combined = combined[:maxSize]
		}

		if err := os.WriteFile(filename, combined, 0644); err != nil {
			fmt.Println("write error:", err)
			return
		}

		fmt.Printf("size now: %d bytes\n", len(combined))

		time.Sleep(3 * time.Second)
	}
}

package main

import (
	p "project02/project02_utils"
	"time"
)

func main() {
	m := p.NewMaps()

	// Create a new instance of the Server struct
	p.Serve(*m)

	for {
		time.Sleep(1 * time.Second)
	}
}

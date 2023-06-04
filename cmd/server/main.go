package main

import (
	"fmt"
	"github.com/bobgromozeka/metrics/internal/server"
)

func main() {
	err := server.Start()

	if err != nil {
		fmt.Println("Error during server start: ", err)
	}
}

package main

import (
	"fmt"

	"github.com/bobgromozeka/metrics/internal/server"
)

func main() {
	setupConfiguration()

	err := server.Start(serverAddr)

	if err != nil {
		fmt.Println("Error during server start: ", err)
	}
}

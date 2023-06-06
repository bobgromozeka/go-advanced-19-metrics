package main

import (
	"fmt"
	"github.com/bobgromozeka/metrics/internal/server"
)

func main() {
	parseFlags()

	err := server.Start(serverAddr)

	if err != nil {
		fmt.Println("Error during server start: ", err)
	}
}

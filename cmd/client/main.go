package main

import (
	"fmt"

	"github.com/avag-sargsyan/word-of-wisdom-pow/internal/client"
	"github.com/avag-sargsyan/word-of-wisdom-pow/internal/pkg/config"
)

func main() {
	fmt.Println("start client")

	// Load necessary configs
	configs, err := config.Load("config/config.json")
	if err != nil {
		fmt.Println("error load config:", err)
		return
	}

	// Generate address based on config for TCP connection
	address := fmt.Sprintf("%s:%d", configs.ServerHost, configs.ServerPort)

	// run client
	err = client.Run(address)
	if err != nil {
		fmt.Println("client error:", err)
	}
}

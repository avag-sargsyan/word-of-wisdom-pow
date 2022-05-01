package main

import (
	"fmt"

	"github.com/avag-sargsyan/word-of-wisdom-pow/internal/pkg/config"
	"github.com/avag-sargsyan/word-of-wisdom-pow/internal/server"
)

func main() {
	fmt.Println("start server")

	// Load necessary configs
	configs, err := config.Load("config/config.json")
	if err != nil {
		fmt.Println("error load config:", err)
		return
	}

	// Generate address based on config for TCP connection
	serverAddress := fmt.Sprintf("%s:%d", configs.ServerHost, configs.ServerPort)

	err = server.Run(serverAddress)
	if err != nil {
		fmt.Println("server error:", err)
	}
}

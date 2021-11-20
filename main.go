package main

import (
	"fmt"
	"os"

	"github.com/t-junjie/fs/internal/router"
)

func main() {
	if err := router.Start(); err != nil {
		fmt.Printf("Error starting server: %w", err)
		os.Exit(1)
	}

}

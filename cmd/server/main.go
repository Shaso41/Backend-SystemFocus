package main

import (
	"flag"
	"log"

	"github.com/yourusername/redis-clone/internal/server"
)

func main() {
	// Parse command-line flags
	address := flag.String("addr", ":6379", "Server address (host:port)")
	flag.Parse()

	// ASCII art banner
	banner := `
╔═══════════════════════════════════════════════════════════╗
║                                                           ║
║   ██████╗ ███████╗██████╗ ██╗███████╗     ██████╗██╗     ║
║   ██╔══██╗██╔════╝██╔══██╗██║██╔════╝    ██╔════╝██║     ║
║   ██████╔╝█████╗  ██║  ██║██║███████╗    ██║     ██║     ║
║   ██╔══██╗██╔══╝  ██║  ██║██║╚════██║    ██║     ██║     ║
║   ██║  ██║███████╗██████╔╝██║███████║    ╚██████╗███████╗║
║   ╚═╝  ╚═╝╚══════╝╚═════╝ ╚═╝╚══════╝     ╚═════╝╚══════╝║
║                                                           ║
║        High-Performance In-Memory Key-Value Store         ║
║                    Version 1.0.0                          ║
║                                                           ║
╚═══════════════════════════════════════════════════════════╝
`
	log.Println(banner)

	// Create and start server
	srv := server.New(*address)
	if err := srv.Start(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

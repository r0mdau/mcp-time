package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	_ "time/tzdata"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/r0mdau/mcp-time/internal/handlers"
	"github.com/r0mdau/mcp-time/internal/timezone"
)

func main() {
	// Define command-line flags matching Python's arguments
	localTimezone := flag.String("local-timezone", "", "Override local timezone (e.g., 'America/New_York')")
	port := flag.Int("port", 8080, "Port to listen on")
	flag.Parse()

	localTZ := timezone.GetLocalTimezone(*localTimezone)
	log.Printf("Using local timezone: %s", localTZ)

	server := mcp.NewServer(&mcp.Implementation{Name: "mcp-time", Version: "v1.0.0"}, nil)
	// Register tools with the determined local timezone
	handlers.RegisterTools(server, localTZ)

	handler := mcp.NewStreamableHTTPHandler(func(*http.Request) *mcp.Server { return server }, nil)
	addr := fmt.Sprintf(":%d", *port)
	log.Printf("MCP Time Server - listening on %s", addr)
	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatal(err)
	}
}

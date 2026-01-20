package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/google/uuid"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// Tool argument types
type CalculateArgs struct {
	Operation string  `json:"operation"`
	A         float64 `json:"a"`
	B         float64 `json:"b"`
}

type ReverseStringArgs struct {
	Text string `json:"text"`
}

type NoArgs struct{}

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds)
	log.Printf("========================================")
	log.Printf("MCP SSE Server Starting")
	log.Printf("========================================")

	// Create MCP server
	server := mcp.NewServer(&mcp.Implementation{
		Name:    "sample-mcp-server-sse",
		Version: "1.0.0",
	}, nil)

	log.Printf("[MAIN] Registering tools...")

	// Create schemas for tool inputs
	calculateSchema, _ := jsonschema.For[CalculateArgs](nil)
	generateUUIDSchema, _ := jsonschema.For[NoArgs](nil)
	reverseStringSchema, _ := jsonschema.For[ReverseStringArgs](nil)

	// Register calculate tool
	mcp.AddTool(server, &mcp.Tool{
		Name:        "calculate",
		Description: "Performs basic arithmetic calculations (add, subtract, multiply, divide)",
		InputSchema: calculateSchema,
	}, calculateHandler)
	log.Printf("[MAIN] - Registered tool: calculate")

	// Register generate_uuid tool
	mcp.AddTool(server, &mcp.Tool{
		Name:        "generate_uuid",
		Description: "Generates a random UUID (v4)",
		InputSchema: generateUUIDSchema,
	}, generateUUIDHandler)
	log.Printf("[MAIN] - Registered tool: generate_uuid")

	// Register reverse_string tool
	mcp.AddTool(server, &mcp.Tool{
		Name:        "reverse_string",
		Description: "Reverses the given string",
		InputSchema: reverseStringSchema,
	}, reverseStringHandler)
	log.Printf("[MAIN] - Registered tool: reverse_string")

	log.Printf("[MAIN] All 3 tools registered successfully")

	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Create SSE handler
	log.Printf("[MAIN] Setting up SSE transport on :%s", port)

	handler := mcp.NewSSEHandler(func(r *http.Request) *mcp.Server {
		log.Printf("[SSE] New connection from %s %s %s", r.RemoteAddr, r.Method, r.URL.Path)
		log.Printf("[SSE] Headers: User-Agent=%s, Origin=%s", r.Header.Get("User-Agent"), r.Header.Get("Origin"))
		return server
	}, nil)

	// Create HTTP server with SSE and health endpoints
	mux := http.NewServeMux()

	// SSE endpoint
	mux.Handle("/sse", handler)

	// Health check endpoint
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[HEALTH] Health check from %s", r.RemoteAddr)
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	// Catch-all for debugging
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/sse" && r.URL.Path != "/health" {
			log.Printf("[WARN] Unknown endpoint requested: %s %s from %s", r.Method, r.URL.Path, r.RemoteAddr)
			http.NotFound(w, r)
		}
	})

	log.Printf("========================================")
	log.Printf("[MAIN] Server ready and listening on :%s", port)
	log.Printf("[MAIN] SSE endpoint: http://localhost:%s/sse", port)
	log.Printf("[MAIN] Health endpoint: http://localhost:%s/health", port)
	log.Printf("========================================")

	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatalf("[MAIN] FATAL: Server failed to start: %v", err)
	}
}

// Tool Handlers

func calculateHandler(ctx context.Context, req *mcp.CallToolRequest, args CalculateArgs) (*mcp.CallToolResult, any, error) {
	log.Printf("[TOOL] Executing calculate with args: %+v", args)

	var result float64

	switch args.Operation {
	case "add":
		result = args.A + args.B
	case "subtract":
		result = args.A - args.B
	case "multiply":
		result = args.A * args.B
	case "divide":
		if args.B == 0 {
			log.Printf("[TOOL] ERROR: Division by zero")
			return &mcp.CallToolResult{
				Content: []mcp.Content{
					&mcp.TextContent{Text: "Error: Division by zero"},
				},
				IsError: true,
			}, nil, nil
		}
		result = args.A / args.B
	default:
		log.Printf("[TOOL] ERROR: Invalid operation: %s", args.Operation)
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("Error: Invalid operation: %s", args.Operation)},
			},
			IsError: true,
		}, nil, nil
	}

	responseText := fmt.Sprintf("Result: %.2f %s %.2f = %.2f", args.A, args.Operation, args.B, result)
	log.Printf("[TOOL] Calculate successful: %s", responseText)

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: responseText},
		},
	}, nil, nil
}

func generateUUIDHandler(ctx context.Context, req *mcp.CallToolRequest, args NoArgs) (*mcp.CallToolResult, any, error) {
	log.Printf("[TOOL] Executing generate_uuid")

	id := uuid.New()
	responseText := fmt.Sprintf("Generated UUID: %s", id.String())

	log.Printf("[TOOL] UUID generated: %s", id.String())

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: responseText},
		},
	}, nil, nil
}

func reverseStringHandler(ctx context.Context, req *mcp.CallToolRequest, args ReverseStringArgs) (*mcp.CallToolResult, any, error) {
	log.Printf("[TOOL] Executing reverse_string with args: %+v", args)

	runes := []rune(args.Text)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}

	reversed := string(runes)
	responseText := fmt.Sprintf("Original: %s\nReversed: %s", args.Text, reversed)

	log.Printf("[TOOL] String reversed successfully")

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: responseText},
		},
	}, nil, nil
}

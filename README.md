# Sample MCP Server with SSE (Golang)

A Model Context Protocol (MCP) server implementation in Go using the official [MCP Go SDK](https://github.com/modelcontextprotocol/go-sdk) with Server-Sent Events (SSE) transport.

## Features

This MCP server provides 3 tools:

1. **calculate** - Performs basic arithmetic operations (add, subtract, multiply, divide)
2. **generate_uuid** - Generates a random UUID (v4)
3. **reverse_string** - Reverses a given string

## Architecture

Built with the official MCP Go SDK (`github.com/modelcontextprotocol/go-sdk`), this server:
- Uses the SDK's built-in SSE transport implementation
- Automatically handles JSON-RPC 2.0 protocol
- Provides comprehensive logging for debugging
- Includes CORS support for browser-based clients

### Endpoints

- `/sse` - SSE endpoint for establishing the event stream (handled by SDK)
- `/health` - Health check endpoint

## Installation

```bash
cd sample-mcp-server-sse
go mod download
go build -o mcp-server
```

## Running the Server

```bash
./mcp-server
```

The server will start on `http://localhost:8080`

## Kubernetes Deployment

### Build Docker Image

```bash
cd sample-mcp-server-sse
docker build -t mcp-server-sse:latest .
```

### Push to Container Registry (Optional)

If using a remote Kubernetes cluster, push the image to your container registry:

```bash
# Tag for your registry
docker tag mcp-server-sse:latest your-registry.com/mcp-server-sse:latest

# Push to registry
docker push your-registry.com/mcp-server-sse:latest
```

Update the image in `k8s/deployment.yaml` to match your registry.

### Deploy to Kubernetes

Using kubectl:

```bash
# Apply manifests
kubectl apply -f k8s/

# Or using kustomize
kubectl apply -k k8s/
```

### Verify Deployment

```bash
# Check deployment status
kubectl get deployments mcp-server-sse

# Check pods
kubectl get pods -l app=mcp-server-sse

# Check service
kubectl get svc mcp-server-sse
```

### Access the Service

#### Port Forward (for testing)

```bash
kubectl port-forward svc/mcp-server-sse 8080:80
```

Then access at `http://localhost:8080`

#### Ingress (for production)

Create an Ingress resource to expose the service:

```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: mcp-server-sse
  annotations:
    nginx.ingress.kubernetes.io/rewrite-target: /
spec:
  rules:
  - host: mcp-server.example.com
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: mcp-server-sse
            port:
              number: 80
```

### Kubernetes Configuration

The deployment includes:
- **Replicas**: 2 pods for high availability
- **Health Checks**: Liveness and readiness probes on `/health`
- **Resources**: CPU and memory limits configured
- **Service**: ClusterIP service on port 80

To scale the deployment:

```bash
kubectl scale deployment mcp-server-sse --replicas=3
```

### Clean Up

```bash
kubectl delete -k k8s/
```

## Testing with MCP Inspector

The [MCP Inspector](https://modelcontextprotocol.io/docs/tools/inspector) is a visual interface for testing and debugging MCP servers. It provides an interactive way to explore your server's tools and test them without writing code.

### Prerequisites

Install the MCP Inspector:

```bash
npx @modelcontextprotocol/inspector
```

### Running the Inspector

1. Start your MCP server:

```bash
# Local development
./mcp-server

# Or with Kubernetes port-forward
kubectl port-forward svc/mcp-server-sse 8080:80
```

2. In a new terminal, launch the inspector with your server's SSE endpoint:

```bash
npx @modelcontextprotocol/inspector http://localhost:8080/sse
```

3. The inspector will open in your browser, showing:
   - **Server Info**: Connection status and server capabilities
   - **Tools Tab**: List of all 3 available tools (calculate, generate_uuid, reverse_string)
   - **Interactive Testing**: Forms to call each tool with parameters

### Using the Inspector

**Test the Calculate Tool:**
1. Click on the "Tools" tab
2. Select "calculate" from the list
3. Fill in the parameters:
   - operation: "add"
   - a: 10
   - b: 5
4. Click "Call Tool" to see the result

**Test the Generate UUID Tool:**
1. Select "generate_uuid"
2. Click "Call Tool" (no parameters needed)
3. View the generated UUID

**Test the Reverse String Tool:**
1. Select "reverse_string"
2. Enter text: "Hello, World!"
3. Click "Call Tool" to see the reversed string

The inspector provides immediate feedback and makes it easy to verify your server is working correctly before integrating it into your application.

## Logging and Debugging

The server includes comprehensive logging to help diagnose connection and protocol issues. All logs are written to stdout with timestamps and categorized by component.

### Log Categories

- `[MAIN]` - Server startup and initialization
- `[SSE]` - SSE connection lifecycle (establish, ping, close)
- `[MESSAGE]` - HTTP message endpoint requests
- `[HANDLER]` - JSON-RPC request processing
- `[INIT]` - Client initialization
- `[TOOL]` - Tool execution
- `[HEALTH]` - Health check requests
- `[WARN]` - Unknown endpoints or warnings

### Example Log Output

```
2026-01-20 10:30:15.123456 main.go:428: ========================================
2026-01-20 10:30:15.123789 main.go:429: MCP SSE Server Starting
2026-01-20 10:30:15.123890 main.go:430: ========================================
2026-01-20 10:30:15.124001 main.go:432: [MAIN] Server initialized with 3 tools
2026-01-20 10:30:15.124112 main.go:434: [MAIN] - Tool: calculate (Performs basic arithmetic calculations (add, subtract, multiply, divide))
2026-01-20 10:30:15.124223 main.go:434: [MAIN] - Tool: generate_uuid (Generates a random UUID (v4))
2026-01-20 10:30:15.124334 main.go:434: [MAIN] - Tool: reverse_string (Reverses the given string)
2026-01-20 10:30:15.124445 main.go:458: [MAIN] Server ready and listening on :8080

# SSE Connection
2026-01-20 10:30:20.567890 main.go:302: [SSE] New connection from 127.0.0.1:54321 GET /sse
2026-01-20 10:30:20.568001 main.go:303: [SSE] Headers: User-Agent=Mozilla/5.0, Origin=http://localhost:3000
2026-01-20 10:30:20.568112 main.go:330: [SSE] Sending endpoint URL: http://localhost:8080/message
2026-01-20 10:30:20.568223 main.go:333: [SSE] Connection established successfully for 127.0.0.1:54321

# JSON-RPC Request
2026-01-20 10:30:21.123456 main.go:361: [MESSAGE] Request from 127.0.0.1:54322 POST /message
2026-01-20 10:30:21.123567 main.go:138: [HANDLER] Processing method: initialize
2026-01-20 10:30:21.123678 main.go:147: [HANDLER] Initialize request from client: test-client v1.0.0
2026-01-20 10:30:21.123789 main.go:170: [INIT] Protocol version requested: 2024-11-05
2026-01-20 10:30:21.123890 main.go:183: [INIT] Server initialized successfully, protocol version: 2024-11-05
2026-01-20 10:30:21.124001 main.go:397: [MESSAGE] Request completed successfully: method=initialize
```

### Troubleshooting Common Issues

**Connection Issues:**

If you see no `[SSE]` logs:
- Check firewall settings
- Verify the server is listening on the correct port
- Check network connectivity

If SSE connection closes immediately:
- Look for `[SSE] Connection closed by client` - indicates client disconnected
- Check for `ERROR: Streaming not supported` - server doesn't support SSE

**Request Failures:**

If you see `ERROR: Method not allowed`:
- The endpoint received wrong HTTP method (e.g., GET instead of POST for /message)
- Check CORS preflight requests with `[SSE] Handling OPTIONS preflight request`

If you see `ERROR: Failed to parse JSON-RPC request`:
- Request body is not valid JSON
- Check the request format matches JSON-RPC 2.0 spec

If you see `ERROR: Unknown method`:
- The requested JSON-RPC method is not supported
- Valid methods: `initialize`, `tools/list`, `tools/call`

**Tool Execution Issues:**

If you see `ERROR: Unknown tool requested`:
- Tool name doesn't match available tools: `calculate`, `generate_uuid`, `reverse_string`

### Viewing Logs in Kubernetes

```bash
# View logs from all pods
kubectl logs -l app=mcp-server-sse --all-containers=true

# Follow logs in real-time
kubectl logs -f deployment/mcp-server-sse

# View logs from specific pod
kubectl logs <pod-name>

# View logs with timestamps
kubectl logs <pod-name> --timestamps
```

## Usage

The easiest way to interact with this MCP server is using the [MCP Inspector](#testing-with-mcp-inspector) or any MCP-compatible client.

### Using the MCP Inspector (Recommended)

```bash
# Start the server
./mcp-server

# In another terminal, run the inspector
npx @modelcontextprotocol/inspector http://localhost:8080/sse
```

The inspector provides a web UI to:
- View all available tools
- Test tools with interactive forms
- See request/response logs

### Using MCP Client Libraries

For programmatic access, use an MCP client library:

**TypeScript/JavaScript:**
```typescript
import { Client } from "@modelcontextprotocol/sdk/client/index.js";
import { SSEClientTransport } from "@modelcontextprotocol/sdk/client/sse.js";

const transport = new SSEClientTransport(
  new URL("http://localhost:8080/sse")
);
const client = new Client({ name: "my-client", version: "1.0.0" }, {});
await client.connect(transport);

// List tools
const tools = await client.listTools();

// Call a tool
const result = await client.callTool({
  name: "calculate",
  arguments: { operation: "add", a: 10, b: 5 }
});
```

**Python:**
```python
from mcp import ClientSession, StdioServerParameters
from mcp.client.sse import sse_client

async with sse_client("http://localhost:8080/sse") as (read, write):
    async with ClientSession(read, write) as session:
        # Initialize
        await session.initialize()

        # List tools
        tools = await session.list_tools()

        # Call tool
        result = await session.call_tool("calculate", {
            "operation": "add",
            "a": 10,
            "b": 5
        })
```

### Direct SSE Connection

You can also connect directly to the SSE endpoint:

```bash
curl -N http://localhost:8080/sse
```

The server will send SSE events with the message endpoint URL and keep the connection alive.

## Tool Specifications

### 1. calculate

Performs basic arithmetic calculations.

**Parameters:**
- `operation` (string, required): One of "add", "subtract", "multiply", "divide"
- `a` (number, required): First operand
- `b` (number, required): Second operand

**Example:**
```json
{
  "operation": "multiply",
  "a": 7,
  "b": 6
}
```

### 2. generate_uuid

Generates a random UUID v4.

**Parameters:** None

### 3. reverse_string

Reverses the provided string.

**Parameters:**
- `text` (string, required): The string to reverse

**Example:**
```json
{
  "text": "Hello, World!"
}
```

## Protocol

This server uses the official MCP Go SDK which implements the MCP protocol version `2024-11-05` with SSE transport. The SDK handles all JSON-RPC 2.0 communication, session management, and transport concerns.

## Dependencies

- `github.com/modelcontextprotocol/go-sdk` - Official MCP Go SDK
- `github.com/google/uuid` - UUID generation for the generate_uuid tool

## Resources

- [MCP Go SDK](https://github.com/modelcontextprotocol/go-sdk)
- [MCP Documentation](https://modelcontextprotocol.io/)
- [MCP Inspector](https://modelcontextprotocol.io/docs/tools/inspector)

## License

MIT

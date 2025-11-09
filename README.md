# MCP Time Server

Simple MCP tools for time and timezone conversions.

Based on the Python MCP Time Server by Anthropic [modelcontextprotocol/servers](https://github.com/modelcontextprotocol/servers/tree/main/src/time).

## Project Structure

This project follows the standard Go project layout:

```
mcp-time/
├── cmd/mcp-time/        # Application entry point
├── internal/
│   ├── types/           # Shared type definitions
│   ├── handlers/        # MCP tool handlers
│   ├── timezone/        # Timezone operations
│   └── timeutil/        # Time utility functions
├── build/               # Compiled binaries
└── docs/                # Documentation
```

## Included Tools

- `get_current_time`: Return current time in a given IANA timezone (default UTC)
- `convert_time`: Convert time between timezones in HH:MM format

Example prompt use in Github Copilot:

- `Get the current time in New York using the MCP Time Server tool.`
- `Convert 14:30 from London time to Tokyo time using the MCP Time Server tool.`

## Development

### Build and Run

1. Build the binary:
   ```bash
   make build
   ```

2. Run the server:
   ```bash
   make run
   ```

3. The MCP server listens on port 8080 by default.

### Command-line Options

- `--local-timezone`: Override local timezone (e.g., 'America/New_York')
- `--port`: Port to listen on (default: 8080)

### Testing

Run the test suite:
```bash
make test
```

Aiming for > 80% test coverage in all internal packages.

### Development Workflow

See [AGENTS.md](AGENTS.md) for detailed development guidelines and testing instructions.

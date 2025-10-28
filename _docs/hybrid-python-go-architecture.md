# Hybrid Python & Go Project Architecture

This document outlines the recommended approach for maintaining a project that uses both Python and Go, including how to organize external packages and manage dependencies.

## Project Structure

```
project-root/
├── .github/                 # GitHub workflows and templates
├── cmd/                     # Main application entry points
│   ├── go-app/              # Go application
│   └── py-app/              # Python application
├── internal/                # Private application code
│   ├── go/                  # Internal Go packages
│   └── py/                  # Internal Python packages
├── pkg/                     # Public Go packages
├── scripts/                 # Build and utility scripts
├── go.mod                   # Go module definition
├── pyproject.toml           # Python project metadata and dependencies
└── README.md                # Project documentation
```

## External Package Management

### Go External Packages

1. **Location**: Place external Go packages under `pkg/external/`
2. **Naming**: Use the original repository name as the directory name
3. **Vendoring**: Use Go modules for dependency management

Example structure:
```
pkg/external/
├── github.com/
│   └── someuser/
│       └── somerepo/        # Cloned external Go package
└── internal/                # Modified external packages
```

### Python External Packages

1. **Location**: Place external Python packages under `_external/`
2. **Virtual Environment**: Always use a virtual environment
3. **Dependencies**: Use `requirements.txt` or `pyproject.toml`

Example structure:
```
_external/
├── some-python-package/     # Cloned external Python package
└── requirements.txt         # Pinned dependencies
```

## Integration Patterns

### 1. REST API Communication

**Go Service (Server)**
```go
// cmd/go-app/main.go
package main

import (
    "net/http"
    "github.com/gin-gonic/gin"
)

func main() {
    r := gin.Default()
    r.GET("/api/data", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{"message": "Hello from Go!"})
    })
    r.Run(":8080")
}
```

**Python Client**
```python
# cmd/py-app/main.py
import requests

def fetch_from_go():
    response = requests.get("http://localhost:8080/api/data")
    return response.json()

if __name__ == "__main__":
    print(fetch_from_go())
```

### 2. gRPC for High-Performance Communication

1. Define your service in a `.proto` file
2. Generate stubs for both Go and Python
3. Implement server in Go and client in Python (or vice versa)

### 3. Shared Configuration

Use a common configuration format like JSON, YAML, or environment variables:

```yaml
# config/config.yaml
database:
  host: localhost
  port: 5432
  name: myapp
  user: admin

server:
  port: 3000
  environment: development
```

## Development Workflow

### Setting Up the Environment

1. **Clone the repository**
   ```bash
   git clone <repo-url>
   cd <project-name>
   ```

2. **Set up Go environment**
   ```bash
   # Install Go dependencies
   go mod download
   ```

3. **Set up Python environment**
   ```bash
   # Create and activate virtual environment
   python -m venv venv
   source venv/bin/activate  # On Windows: .\venv\Scripts\activate
   
   # Install Python dependencies
   pip install -r requirements.txt
   ```

### Building and Running

**Using Makefile**
```makefile
# Build both Go and Python components
build:
    cd cmd/go-app && go build -o ../../bin/go-app
    cd cmd/py-app && pip install -r requirements.txt

# Run services
dev: build
    # Run Go service in background
    ./bin/go-app &
    # Run Python service
    cd cmd/py-app && python main.py

# Clean build artifacts
clean:
    rm -rf bin/*
    rm -rf __pycache__
    rm -rf *.pyc
```

## Dependency Management

### Go Dependencies
- Use Go modules (`go.mod` and `go.sum`)
- To add a dependency: `go get github.com/example/package`
- To update all dependencies: `go get -u ./...`

### Python Dependencies
- Use `requirements.txt` or `pyproject.toml` with `poetry`
- To add a dependency: `pip install package`
- To freeze dependencies: `pip freeze > requirements.txt`

## Testing

### Go Tests
```bash
# Run all Go tests
go test ./...

# Run tests with coverage
go test -coverprofile=coverage.out ./...
```

### Python Tests
```bash
# Install test dependencies
pip install -r requirements-test.txt

# Run tests
pytest tests/

# Run with coverage
pytest --cov=myapp tests/
```

## CI/CD Integration

Example GitHub Actions workflow:

```yaml
# .github/workflows/ci.yml
name: CI

on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    
    services:
      postgres:
        image: postgres:13
        env:
          POSTGRES_PASSWORD: postgres
        ports:
          - 5432:5432
    
    steps:
    - uses: actions/checkout@v2
    
    - name: Set up Go
      uses: actions/setup-go@v2
      with:
        go-version: '1.19'
    
    - name: Set up Python
      uses: actions/setup-python@v2
      with:
        python-version: '3.9'
    
    - name: Install dependencies
      run: |
        go mod download
        python -m pip install --upgrade pip
        pip install -r requirements.txt
        pip install -r requirements-test.txt
    
    - name: Run Go tests
      run: go test -v ./...
    
    - name: Run Python tests
      run: pytest
```

## Best Practices

1. **Separation of Concerns**:
   - Keep Go and Python code in separate directories
   - Define clear interfaces between components

2. **Error Handling**:
   - Handle errors consistently across both languages
   - Implement proper error propagation

3. **Logging**:
   - Use a consistent logging format
   - Consider using a centralized logging solution

4. **Documentation**:
   - Document the interface between Go and Python components
   - Include examples for common operations

5. **Performance Considerations**:
   - Be mindful of the overhead of cross-language calls
   - Use batching for data transfer between languages

## Example: Python Calling Go

1. Build Go as a shared library:
   ```go
   // pkg/example/lib.go
   package main
   
   import "C"
   
   //export Add
   func Add(a, b int) int {
       return a + b
   }
   
   func main() {}
   ```
   Build with: `go build -buildmode=c-shared -o libexample.so`

2. Call from Python:
   ```python
   # python/example.py
   from ctypes import cdll, c_int
   
   lib = cdll.LoadLibrary('./libexample.so')
   result = lib.Add(c_int(2), c_int(3))
   print(result)  # Output: 5
   ```

## Troubleshooting

### Common Issues

1. **Version Conflicts**:
   - Ensure compatible versions of Go and Python
   - Pin dependency versions in both languages

2. **Build Issues**:
   - Check `GOPATH` and `GOROOT` environment variables
   - Verify Python virtual environment activation

3. **Performance Bottlenecks**:
   - Profile both Go and Python components
   - Consider using a message queue for heavy processing

4. **Dependency Management**:
   - Keep `go.mod` and `requirements.txt` up to date
   - Document any manual steps for setting up external dependencies

## Conclusion

This hybrid architecture allows you to leverage the strengths of both Go (performance, concurrency) and Python (ecosystem, data science) while maintaining a clean separation of concerns. Follow these guidelines to ensure a maintainable and scalable codebase.

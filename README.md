# Status Flow Generator

A Go-based tool that generates visual status flow diagrams from JSON definitions using Graphviz.

## Prerequisites

- Go 1.x or higher
- Graphviz (must be installed and accessible in PATH)

## Installation

1. Clone the repository
2. Ensure Graphviz is installed on your system
3. Build the program: `go build status_flow.go`

## Usage

### Command Line Flags

- `-path`: Output directory path (default: current directory)
- `-name`: Base name for output files without extension (default: "status_flow")
- `-static`: Use predefined static data instead of JSON input (default: false)

### Input Format

The program accepts JSON input in two ways:

1. Direct JSON string through stdin
2. JSON file input

#### JSON Structure

The JSON should be an array of status objects, where each object has:
- `Name`: (string) The name of the status
- `NextStatus`: (array of strings) List of possible next statuses

Example JSON structure:
```json
[
    {
        "Name": "Status A",
        "NextStatus": ["Status B", "Status C"]
    },
    {
        "Name": "Status B",
        "NextStatus": ["Status C"]
    },
    {
        "Name": "Status C",
        "NextStatus": ["Status A"]
    }
]
```

### Examples

1. Using JSON file input:
```bash
./status_flow -path=./output -name=my_flow < input.json
```

2. Using JSON string input:
```bash
echo '[{"Name":"Status A","NextStatus":["Status B","Status C"]},{"Name":"Status B","NextStatus":["Status C"]},{"Name":"Status C","NextStatus":["Status A"]}]' | ./status_flow -path=./output -name=my_flow
```

3. Using predefined static data:
```bash
./status_flow -path=./output -name=my_flow -static
```

## Output

The program generates two output files:
1. DOT file: Graphviz DOT file used to generate the PNG
2. PNG file: PNG image of the SVG file

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details
# ChatGPTImport

A command-line tool that processes JSON exports of ChatGPT conversations and converts them into well-formatted markdown files.

## Features

- Converts ChatGPT JSON exports to individual markdown files
- Exports conversations to a specified output directory or current directory
- Supports limiting the number of conversations processed
- Simple and efficient processing

## Installation

### Build from Source

This tool is written in Go and requires compilation before use.

**Requirements:**
- Go 1.16 or higher

**Build:**
```bash
go build -o chatGPTImport
```

**Installation (Linux):**
```bash
sudo install -m 755 chatGPTImport /usr/local/bin/
```

After installation, you can run `chatGPTImport` from any directory.

## Usage
```bash
chatGPTImport [source] [optional output] [flags]
```

### Options

- `-l, --limit <number>` - Maximum number of conversations to process (0 for no limit, Default)
- `-h, --help` - Display help information

### Example

```bash
chatGPTImport ./chatgptExport/conversations.json ./markdown -l 10
```

## License

MIT
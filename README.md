# Huffman Web

A simple web application for encoding and decoding files using Huffman coding.

## Features

- **File Encoding**: Upload files to compress using Huffman coding.
- **File Decoding**: Decode previously Huffman-encoded files.
- **Web Interface**: User-friendly HTML interface with upload and download pages.
- **Structured Logging**: Built-in request and error logging with `slog`.

## Setup

```bash
git clone https://github.com/yourusername/huffman-web.git
cd huffman-web
go mod tidy
```

## Running the Server

```bash
go run cmd/web/main.go -addr :4000
```

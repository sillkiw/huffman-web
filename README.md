# Huffman Web

A simple web application for encoding and decoding files using Huffman coding.

## Screenshot
<img width="925" height="666" alt="image" src="https://github.com/user-attachments/assets/64f963eb-16e2-437c-ac6c-35c9e32d9cb0" />

## Features

- **File Encoding**: Upload files to compress using Huffman coding.
- **File Decoding**: Decode previously Huffman-encoded files.
- **Web Interface**: User-friendly HTML interface with upload and download pages.
- **Structured Logging**: Built-in request and error logging with `slog`.

## Setup

```bash
git clone https://github.com/sillkiw/huffman-web.git
cd huffman-web
go mod tidy
```

## Running the Server

```bash
go run cmd/web/main.go -addr :4000
```

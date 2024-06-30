# 🕵️‍♂️ 1337 Checkbox Pwner



--
1337 Checkbox Pwner is a high-performance, hacker-themed tool designed to dominate the checkbox world of [onemillioncheckboxes.com](https://onemillioncheckboxes.com). This tool continuously checks checkboxes to help you maintain your king-of-the-hill status.

## Features

- 🚀 High performance with multiple parallel workers.
- 📡 Real-time WebSocket connection to continuously check checkboxes.
- ⚙️ Configurable via command-line arguments for ultimate flexibility.
- 🕶️ Cool hacker-themed logging with emojis.

## Installation

### Prerequisites

- Go (1.21.4 or later)
- GNU Make

### Build and Run

Clone the repository and navigate to the project directory:

```sh
git clone https://github.com/yourusername/1337-checkbox-pwner.git
cd 1337-checkbox-pwner
```

Build and run the application:

#### On Linux/OSX

```sh
make build-run
```

#### On Windows

```sh
make build-run-windows
```

## Usage

```sh
./checkbox-pwner [flags]
```

### Flags

- `-w, --workers`: Number of parallel workers (default: 100)
- `-r, --retries`: Maximum number of retries for each batch (default: 3)
- `-s, --sleep`: Sleep duration between requests in milliseconds (default: 200)
- `-b, --batch`: Batch size for processing checkboxes (default: 50)

### Example

```sh
./checkbox-pwner --workers 150 --retries 5 --sleep 100 --batch 100
```

## Makefile Commands

- `make all`: Default target (build)
- `make deps`: Install dependencies
- `make build`: Build for the current OS
- `make build-windows`: Build for Windows
- `make run`: Run the application
- `make clean`: Clean the build artifacts
- `make build-run`: Build and run for the current OS
- `make build-run-windows`: Build and run for Windows

---

Happy pwning! 🕵️‍♂️💻

# PVM - Pulumi Version Manager

PVM is a version manager for Pulumi CLI that allows you to install and switch between different versions of Pulumi.

## Features

- 🚀 Install multiple Pulumi versions
- 🔄 Switch between installed versions
- 📋 List available versions
- 💡 Show current active version
- 🖥️ Cross-platform support (Windows, Linux, macOS)
- 🔒 Secure downloads from official GitHub releases
- 📦 Local version caching

## Installation

### Using Go

```
go install github.com/tomski747/pvm/cmd/pvm@latest
```

Make sure your `GOPATH/bin` is in your PATH. You can check your GOPATH with:
```bash
go env GOPATH
```

### Building from Source

1. Clone the repository
```bash
git clone https://github.com/tomski747/pvm.git
cd pvm
```

2. Build and install
```bash
make install
```

## Usage

```bash
# Install the latest version
pvm install latest

# Install a specific version
pvm install 3.91.1

# Install and use a version
pvm install 3.91.1 --use

# Switch to an installed version
pvm use 3.91.1

# List installed versions
pvm list

# List all available versions
pvm list --all

# Show current version
pvm current

# Remove a version
pvm remove 3.91.1
```

## License

MIT
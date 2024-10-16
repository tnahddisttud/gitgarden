#!/bin/bash

# Function to print error messages
error() {
    echo "Error: $1" >&2
    exit 1
}

# Check if Go is installed
if ! command -v go &> /dev/null; then
    error "Go is not installed. Please install Go and try again."
fi

# Set repository URL
REPO_URL="https://github.com/tnahddisttud/gitgarden.git"

# Clone the repository (or pull if it already exists)
if [ ! -d "./gitgarden" ]; then
    echo "Cloning the repository..."
    git clone $REPO_URL || error "Failed to clone the repository"
else
    echo "Repository already exists. Pulling latest changes..."
    cd gitgarden && git pull || error "Failed to update the repository"
fi

# Compile the Go program from main.go
echo "Building the gitgarden CLI tool..."
go build -o gitgarden ./main.go || error "Go build failed"

mkdir -p "$HOME/.local/bin"

# Move the binary to the user's PATH
mv $TOOL_NAME "$HOME/.local/bin/"

# Make the binary executable
chmod +x "$HOME/.local/bin/$TOOL_NAME"

# Check if ~/.local/bin is in PATH, if not, add it
if [[ ":$PATH:" != *":$HOME/.local/bin:"* ]]; then
    echo 'export PATH="$HOME/.local/bin:$PATH"' >> "$HOME/.bashrc"
    echo "Added $HOME/.local/bin to PATH in ~/.bashrc"
    echo "Please run 'source ~/.bashrc' or start a new terminal session to update your PATH."
fi

echo "$TOOL_NAME has been successfully installed!"
echo "You may need to restart your terminal or run 'source ~/.bashrc' to use the command immediately."

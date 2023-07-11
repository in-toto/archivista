#!/bin/bash

if [ "$#" -ne 1 ]; then
    echo "Usage: $0 <version>"
    exit 1
fi

version=${1#v} # Remove 'v' prefix if present
url="https://github.com/testifysec/witness/releases/download/v${version}/witness_${version}_linux_amd64.tar.gz"

echo "Downloading Witness binary (version: ${version})..."
curl -L -o witness.tar.gz "${url}"

echo "Extracting Witness binary..."
tar -xzf witness.tar.gz

# Check if UPX is installed, and if not, install it
if ! command -v upx &> /dev/null; then
    echo "UPX not found, installing..."

    if [[ "$OSTYPE" == "linux-gnu"* ]]; then
        # Linux
        sudo apt-get update && sudo apt-get install -y upx
    elif [[ "$OSTYPE" == "darwin"* ]]; then
        # macOS
        if ! command -v brew &> /dev/null; then
            echo "Homebrew not found, please install Homebrew to continue."
            exit 1
        fi

        brew install upx
    else
        echo "Unsupported operating system: $OSTYPE"
        exit 1
    fi
fi


echo "Compressing Witness binary using UPX..."
upx --best --ultra-brute witness

echo "Witness binary has been compressed successfully."
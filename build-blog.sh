#!/bin/bash

# Blog build script for req project
# Builds the static site using pyssg

set -e

echo "Building req blog..."

# Check if pyssg is installed
if [ ! -d ".pyssg" ]; then
    echo "Installing pyssg..."
    git clone https://github.com/maniac-en/pyssg.git .pyssg
    cd .pyssg
    git checkout 46d6632bff1259ae045a116ba6ee4332cbae5f26
    cd ..
fi

# Generate the site
echo "Generating static site..."
cd .pyssg
PYTHONPATH=. python3 src/main.py --config ../pyssg.config.json
cd ..

# Remove pyssg example images that we don't need
echo "Cleaning up unnecessary images..."
if [ -d "docs/images" ]; then
    rm -rf docs/images
fi

echo "Blog built successfully! Files are in /docs directory."
echo "The site is ready for GitHub Pages deployment."
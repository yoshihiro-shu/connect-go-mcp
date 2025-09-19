#!/bin/bash

set -e  # Exit on error

echo "Generating proto files for testdata directories..."

# Find all testdata subdirectories and generate proto files
for dir in testdata/*/; do
    if [ -d "$dir" ] && [ -f "$dir/buf.gen.yaml" ]; then
        echo "Processing directory: $dir"
        cd "$dir"
        rm -rf gen
        buf generate
        cd - > /dev/null
        echo "✅ Generated proto files for $dir"
    fi
done

echo "🎉 All proto files generated successfully!"
#! /bin/bash

cd $(dirname ${0})

# We may want to switch to webpack eventually if this doesn't work.
echo "Attempting to pack web UI..."
mkdir -p build

command -v inliner > /dev/null || {
    echo "Dependency 'inliner' not found. Please install with 'npm i -g inliner'."
    echo "More information: https://github.com/remy/inliner"
    echo "Press Enter to continue without building the UI, or Ctrl+C to quit."
    read -n 1
} && inliner -m index.html > build/compiled_webapp.html

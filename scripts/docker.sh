#!/bin/bash

# Check if the correct number of arguments is provided
if [ "$#" -ne 2 ]; then
    echo "Usage: $0 <build|run> <image_name>"
    exit 1
fi

# Determine whether to build or run
if [ "$1" == "build" ]; then
    # Build the Docker image
    docker build -t $2 .
elif [ "$1" == "run" ]; then
    # Run the Docker container
    docker run -p 1531:1531 $2
else
    echo "Usage: $0 <build|run> <image_name>"
    exit 1
fi

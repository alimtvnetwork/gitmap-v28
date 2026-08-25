#!/bin/bash

# Check if the user is provided as an argument
if [ -z "$1" ]; then
  echo "Usage: $0 <username>"
  exit 1
fi

USERNAME=$1

# Kill all processes owned by the given user
pkill -u $USERNAME

echo "All processes owned by $USERNAME have been killed."

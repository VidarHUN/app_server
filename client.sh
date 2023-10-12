#!/bin/bash

# Define usage function to display how to use the script
usage() {
  echo "Usage: $0 -r <#rooms> -c <#clients>"
  echo "Options:"
  echo "  -r <#rooms>   Number of rooms"
  echo "  -c <#clients> Number of room/client"
  exit 1
}

number_of_rooms=""
number_of_clients=""

# Use getopts to parse command-line options and arguments
while getopts "r:c:" opt; do
  case "$opt" in
    r) number_of_rooms="$OPTARG" ;;
    c) number_of_clients="$OPTARG" ;;
    \?) usage ;;
  esac
done

# Create an array to store the pipe names
declare -a pipes
declare -a pids

# Generate pipe names and create named pipes
for ((i = 1; i <= number_of_rooms; i++)); do
    pipe_name="pipe$i"
    # If there is an open pipe from the previous run
    rm "$pipe_name" &> /dev/null
    pipes+=("$pipe_name")
    mkfifo "$pipe_name"
    ./client < "$pipe_name" &
    pids+=($!)

    room_name="room$i"
    echo "createRoom $room_name" > "$pipe_name"

    echo "$number_of_clients"
    if [ "$number_of_clients" -gt 1 ]; then
        for ((j = 1; j <= number_of_clients; j++)); do
          cpipe_name="cpipe$i$j"
          # If there is an open pipe from the previous run
          rm "$cpipe_name" &> /dev/null
          pipes+=("$cpipe_name")
          mkfifo "$cpipe_name"
          ./client < "$cpipe_name" &
          pids+=($!)
          echo "joinRoom room$i" > "$cpipe_name"
        done
    fi
done

sleep 10

# Clean up: close and remove the named pipes
for pipe in "${pipes[@]}"; do
  rm "$pipe" &> /dev/null
done

for pid in "${pids[@]}"; do
  echo "$pid"
  kill -2 "$pid"
done

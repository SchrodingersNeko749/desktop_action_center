#!/bin/bash
level=$1
echo $level
# Check if the input is valid
if ! [[ "$level" =~ ^[0-9]+$ ]] || (( level < 0 || level > 100 )); then
    exit 1
fi

# Find the path to the brightness file
brightness_file="/sys/class/backlight/amdgpu_bl0/brightness"

# Get the maximum brightness level from the max_brightness file
max_brightness=$(cat $(dirname $brightness_file)/max_brightness)

# Calculate the new brightness level as a percentage of the maximum
new_level=$(echo "scale=0; $level * $max_brightness / 100" | bc)

# Write the new brightness level to the brightness file
echo $new_level | sudo tee $brightness_file
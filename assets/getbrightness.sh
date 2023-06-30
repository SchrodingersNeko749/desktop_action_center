#!/bin/bash

# Get the current brightness value
brightness=$(cat /sys/class/backlight/gmux_backlight/brightness)

# Get the maximum brightness value
max_brightness=$(cat /sys/class/backlight/gmux_backlight/max_brightness)

# Convert the brightness value to a percentage
brightness_percentage=$(echo "scale=2; $brightness / $max_brightness * 100" | bc)

# Round the brightness percentage to the nearest integer
brightness_rounded=$(printf "%.0f" "$brightness_percentage")

# Output the current brightness level
echo $brightness_rounded

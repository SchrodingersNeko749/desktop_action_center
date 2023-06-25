#!/bin/bash

for i in {1..100}; do
    notify-send "Hello $i" "This is a message number $i."
done
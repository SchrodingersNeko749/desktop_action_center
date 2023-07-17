#!/bin/bash


mpc volume | cut -d ":" -f 2 | cut -d "%" -f 1	
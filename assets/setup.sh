#!/bin/bash

mkdir -p ~/.config/actionCenter/
cp -r assets/* ~/.config/actionCenter/
bspc rule -a ActionCenter manage=on locked=on sticky=on layer=above border=off

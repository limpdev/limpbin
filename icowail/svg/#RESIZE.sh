#!/bin/zsh

for svg in *.svg; do
  inkscape --export-filename="${svg%.svg}_resized.svg" \
           --export-area-page --export-width=96 --export-height=96 "$svg"
done

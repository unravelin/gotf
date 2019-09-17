#!/bin/sh
mkdir -p /libs
cd /libs
for f in $(ldd /gotf | sed -n 's/.*\s\(\/.*\) .*/\1/p'); do
  cp --parents "$f" .
done
cd ..

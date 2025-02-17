#!/bin/bash

declare -i x=0
echo "checking if directory /mnt/videos has files"
echo "$(date -Ins): starting"

while [ $(ls /mnt/videos | wc -l) -eq 0 -a $x -lt 200  ]
do
    echo "$(date -Ins): still empty"
    ((x=x+1))
done

echo "$(date -Ins): $(ls /mnt/videos | wc -l) file(s) found"

echo "starting server"

dotnet aspnet-webapp.dll

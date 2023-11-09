#!/bin/bash

declare -i x=0
echo "checking if directory /usr/share/aasx has files"
echo "$(date -Ins): starting"

while [ $(ls /usr/share/aasx | wc -l) -eq 0 -a $x -lt 200  ]
do
    echo "$(date -Ins): still empty"
    ((x=x+1))
done

echo "$(date -Ins): $(ls /usr/share/aasx | wc -l) file(s) found"

echo "starting server"

dotnet aspnet-webapp.dll

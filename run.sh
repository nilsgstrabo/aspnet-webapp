#!/bin/bash

echo "checking if directory /usr/share/aasx has files"
echo "$(date -Ins): starting"

while [ $(ls /usr/share/aasx | wc -l) -eq 0 ]
do
    echo "$(date -Ins): still empty"
    
done

echo "$(date -Ins): $(ls /usr/share/aasx | wc -l) file(s) found"

echo "starting server"

dotnet aspnet-webapp.dll

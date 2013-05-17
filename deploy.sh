#!/bin/bash

go build -o go-wrk *.go

while read node
do
    scp go-wrk root@$node:
    scp config.json root@$node:
    ssh root@$node '/root/go-wrk -d s -f /root/config.json </dev/null >/var/log/root-backup.log 2>&1 &'
done < node.txt

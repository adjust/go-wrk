#!/bin/bash

go build -o go-wrk *.go

while read node; do
    echo $node
    scp -i ~/Downloads/roa.pem go-wrk ubuntu@$node:
    scp -i ~/Downloads/roa.pem config.json ubuntu@$node:
    ssh -i ~/Downloads/roa.pem ubuntu@$node 'sudo su root -c "echo \"echo 88.198.77.76 app.adjust.io >> /etc/hosts\" | sudo bash"'
    ssh -i ~/Downloads/roa.pem ubuntu@$node '/home/ubuntu/go-wrk -d s -f /home/ubuntu/config.json </dev/null >/home/ubuntu/err.log 2>&1 &' &
done < node.txt

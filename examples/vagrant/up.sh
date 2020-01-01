#!/bin/bash

status=$(vagrant status | grep "running" | head -1)
if [[ ! -e $status ]]; then
  vagrant destroy -f # start over
fi

echo "Building ISO"
quiso
echo "Waiting for Virtual Machine to start..."
vagrant up

echo "Contents of /var/lib/cloud/instance/user-data.txt"
vagrant ssh -c "sudo cat /var/lib/cloud/instance/user-data.txt"
echo "Contents of /var/log/cloud-init-output.log"
vagrant ssh -c "cat /var/log/cloud-init-output.log"

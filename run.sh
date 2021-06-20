#!/bin/bash

sudo add-apt-repository ppa:longsleep/golang-backports
sudo apt update
sudo apt -y install python3-pip
sudo apt -y install golang-go
cd ~/
git clone https://github.com/bitclout/core.git
cd core
go build
./core run > /dev/null 2>&1 &
cd ~/balajis-script
go mod vendor
go build
mkdir /root/acopy
pip3 install -r requirements.txt
echo "Sleeping..."
sleep 86400
cp -r /root/.config/bitclout/bitclout/MAINNET/badgerdb /root/acopy
rm /root/acopy/badgerdb/*.mem
./balajis /root/acopy/badgerdb > edges.txt
python3 visualize.py

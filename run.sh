#!/bin/bash

sudo add-apt-repository ppa:longsleep/golang-backports
sudo apt update
sudo apt -y install golang-go
git clone https://github.com/bitclout/core.git
cd core
go build
./core run > /dev/null 2>&1 &
sleep 86400

cd ..
git clone https://github.com/andrewarrow/balajis-script.git
cd balajis-script
go mod vendor
go build
mkdir /root/acopy
copy -r /root/.config/bitclout/bitclout/MAINNET/badgerdb /root/acopy
rm /root/acopy/badgerdb/*.mem
./balajis /root/acopy/badgerdb > edges.txt
pip3 install -r requirements.txt
python3 visualize.py

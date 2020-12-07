#!/bin/bash

cd proto
source GenerateGoCode.sh
cd ..

GOOS=linux go build -ldflags="-s -w" -v -o amono_nvidia_agent agent/main.go &&\
GOOS=linux go build -ldflags="-s -w" -v -o amono_nvidia_monitor monitor/main.go &&\

echo "--- compressing ---" &&\
mkdir -p govidia &&\
cp amono_nvidia_agent amono_nvidia_monitor config.yaml govidia/ &&\
tar -zcvf govidia_smi.tar.gz govidia &&\
rm -rf govidia
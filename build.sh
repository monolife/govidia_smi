#!/bin/bash

go build -v -o amono_nvidia_agent agent/main.go &&\
go build -v -o amono_nvidia_monitor monitor/main.go

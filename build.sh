#!/bin/bash

go build -v -o dzyne_nvidia_agent agent/main.go &&\
go build -v -o dzyne_nvidia_monitor monitor/main.go

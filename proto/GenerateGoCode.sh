#!/bin/bash
#
protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    govidia.proto

# The following is a fix for protobuffer omiting fields with
# zero/null values. Incidentally, this means GPUs with index of '0' 
# get their gpu_index field cut off (oh fun!)
ls *.pb.go | xargs -n1 -IX bash -c \
"sed -e '/int32/ s/,omitempty//' X > X.tmp && mv X{.tmp,}"
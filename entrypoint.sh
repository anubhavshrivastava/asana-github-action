#!/bin/bash

NAME=$1
TIME=$(date)

go run /go/src/app/main/main.go $1 $TIME

FROM golang:latest

RUN apt-get update && apt-get install -y upx-ucl libc6-dev-i386

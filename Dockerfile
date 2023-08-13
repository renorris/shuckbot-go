# Build go-llama.cpp & shuckbot-go
FROM golang:1.21-alpine
WORKDIR /shuckbot-go
COPY . .
RUN apk add --update --no-cache cmake alpine-sdk &&\
cd go-llama.cpp &&\
make clean &&\
make libbinding.a &&\
cd /shuckbot-go &&\
go build -o main .

FROM alpine:latest
COPY --from=0 /shuckbot-go/main /
RUN apk add --update --no-cache libgcc libstdc++
CMD ./main
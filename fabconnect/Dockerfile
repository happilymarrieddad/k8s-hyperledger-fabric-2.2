FROM golang:1.19 AS builder
WORKDIR /app

RUN apt install curl git
RUN git clone https://github.com/hyperledger/firefly-fabconnect
RUN cp -R firefly-fabconnect/* ./
RUN go build -o /app/api main.go

FROM alpine:edge
WORKDIR /app

RUN apk add --no-cache libc6-compat

COPY --from=builder /app/api /app/api
COPY ./fabconnect.yaml /app/fabconnect.yaml
COPY ./ccp.yaml /app/ccp.yaml

#CMD ["/app/api","-f","/app/fabconnect.yaml"]
CMD ["sleep", "infinity"]
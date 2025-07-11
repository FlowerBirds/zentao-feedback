FROM  golang:1.24 as build
ENV GO111MODULE=on
ENV CGO_ENABLED=1
WORKDIR /app

COPY vendor vendor
COPY entity.go ./
COPY util.go ./
COPY go.mod ./
COPY main.go ./

RUN go build

FROM ubuntu:22.04 as ship
RUN apt-get update && apt-get install -y curl dirmngr apt-transport-https lsb-release ca-certificates
WORKDIR /app
COPY --from=build /app/zentao-feedback .

CMD ["/app/zentao-feedback"]
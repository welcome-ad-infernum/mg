FROM golang:1.17-alpine AS builder

WORKDIR $GOPATH/src/mg/

COPY . .

RUN go mod tidy
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o /go/bin/mg
RUN chmod 755 /go/bin/mg

FROM scratch

COPY --from=builder /go/bin/mg /go/bin/mg

ENTRYPOINT ["/go/bin/mg"]
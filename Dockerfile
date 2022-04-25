FROM golang:1.17-alpine as builder

RUN apk add --no-cache ca-certificates

COPY . $GOPATH/src/solomon72/baal
WORKDIR $GOPATH/src/solomon72/baal

RUN go get -d -v
RUN CGO_ENABLED=0 \ 
  GOOS=linux \ 
  GOARCH=amd64 \
  go build -a -installsuffix cgo -ldflags "-w -s" -o /go/bin/baal

FROM scratch as runner

ARG PORT=7001
ENV PORT=${PORT}

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /go/bin/baal /bin/baal

EXPOSE ${PORT}
ENTRYPOINT ["bin/baal"]

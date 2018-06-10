# STEP 1 build executable binary
FROM golang:alpine as builder

# Install SSL ca certificates
RUN apk update && apk add git && apk add ca-certificates

# Create appuser
RUN adduser -D -g '' appuser

COPY . $GOPATH/src/github.com/valentine/ccollage/
WORKDIR $GOPATH/src/github.com/valentine/ccollage/cmd/ccollage/

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -a -installsuffix cgo -ldflags="-w -s" \
    -o /go/bin/ccollage

# STEP 2 build the small image

# Start from scratch
FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd

# Copy our static executable
COPY --from=builder /go/bin/ccollage /go/bin/ccollage
USER appuser
EXPOSE 8080
ENTRYPOINT ["/go/bin/ccollage"]
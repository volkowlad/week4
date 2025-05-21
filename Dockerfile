FROM golang:1.24
ENV GOPATH=/

COPY ./ ./

RUN go mod download
RUN go build -o tasks-app ./cmd/main.go

ENTRYPOINT ["./tasks-app"]
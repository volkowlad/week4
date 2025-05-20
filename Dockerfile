FROM golang:1.24
ENV GOPATH=/

COPY ./ ./

RUN go mod download
RUN go build -o post-app ./app/main.go

ENTRYPOINT ["./post-app"]
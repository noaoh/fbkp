FROM golang:1.10
WORKDIR /go/src
CMD ["go", "get", "github.com/spf13/pflag"]
COPY . .


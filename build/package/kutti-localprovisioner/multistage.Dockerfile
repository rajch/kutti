FROM golang:1.12-alpine AS builder
COPY . /go/src/github.com/rajch/kutti
WORKDIR /go/src/github.com/rajch/kutti/cmd/kutti-localprovisioner
RUN apk update && apk add git
RUN go get -v
RUN CGO_ENABLED='0' go build

FROM scratch
WORKDIR /tmp
COPY --from=builder /go/src/github.com/rajch/kutti/cmd/kutti-localprovisioner/kutti-localprovisioner .
CMD ["/tmp/kutti-localprovisioner"] 
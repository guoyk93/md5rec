FROM golang:1.13 AS builder
ENV CGO_ENABLED 0
WORKDIR /go/src/app
ADD . .
RUN go build -o /md5rec

FROM scratch
COPY --from=builder /md5rec /md5rec
WORKDIR /work
CMD ["/md5rec"]
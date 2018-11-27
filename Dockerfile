FROM amd64/golang as builder
RUN echo $GOPATH
RUN go get -v github.com/denisenkom/go-mssqldb
WORKDIR /go/src/app
COPY ./*.go /go/src/app/
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o main

FROM amd64/alpine
RUN mkdir /app
COPY --from=builder /go/src/app/main /app
COPY config.json /app/config.json
COPY Tree.xml /app/Tree.xml
COPY data /app/data
WORKDIR /app
EXPOSE 80
CMD ["/app/main"]
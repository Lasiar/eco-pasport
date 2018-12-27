FROM amd64/golang as builder
RUN go get -v github.com/denisenkom/go-mssqldb
WORKDIR /go/src/app
COPY ./*.go /go/src/app/
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o main

FROM amd64/alpine
RUN mkdir /app
COPY --from=builder /go/src/app/main /app
COPY config.production.json /app/config.json
COPY Tree.xml /app/Tree.xml
WORKDIR /app
EXPOSE 80
CMD ["/app/main"]
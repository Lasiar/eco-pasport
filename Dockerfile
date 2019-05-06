FROM amd64/golang as builder
RUN go get -v github.com/denisenkom/go-mssqldb
WORKDIR /go/src/EcoPasport
COPY ./*.go /go/src/EcoPasport/
COPY ./web/*.go /go/src/EcoPasport/web/
COPY ./web/context/*.go /go/src/EcoPasport/web/context/
COPY ./model/*.go /go/src/EcoPasport/model/
copy ./base/*.go /go/src/EcoPasport/base/
RUN ls
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o main

FROM amd64/alpine
RUN mkdir /app
COPY --from=builder /go/src/EcoPasport/main /app
COPY config.production.json /app/config.json
COPY Tree.xml /app/Tree.xml
WORKDIR /app
EXPOSE 80
CMD ["/app/main"]

FROM golang:1.16-alpine

WORKDIR $GOPATH/src/github.com/harshdonga/product-scraper/dbapi
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

EXPOSE 5001

CMD ["dbapi"]
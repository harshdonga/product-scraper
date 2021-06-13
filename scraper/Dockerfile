FROM golang:1.16-alpine

WORKDIR $GOPATH/src/github.com/harshdonga/product-scraper/scraper
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

EXPOSE 5000

CMD ["scraper"]
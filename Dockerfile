FROM golang:1.16-alpine3.12 as BUILD
WORKDIR /opt/etsy-orders
COPY . .
RUN apk add git 
RUN go get -d -v ./...
RUN go build -o etsy-orders

FROM alpine:3.13 as FINAL
COPY --from=BUILD /opt/etsy-orders/etsy-orders /bin/
EXPOSE 8081
CMD ["etsy-orders"]
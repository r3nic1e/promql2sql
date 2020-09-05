FROM golang:1.15-alpine3.12 as build

WORKDIR /usr/src/promql2sql
COPY . .
RUN go build -o /promql2sql

###
FROM alpine:3.12

COPY --from=build /promql2sql /usr/bin/promql2sql
CMD ["/usr/bin/promql2sql"]

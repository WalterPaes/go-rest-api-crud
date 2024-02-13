FROM golang:1.19 AS BUILDER

WORKDIR /go/src

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on \
 GOOS=linux go build -o api ./cmd/api/main.go

FROM golang:1.19-alpine3.15 as RUNNER

COPY --from=BUILDER /go/src/api .
COPY --from=builder /go/src/.env .  

CMD ["./api"]
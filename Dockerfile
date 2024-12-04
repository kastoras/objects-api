FROM golang:1.23 AS development

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

RUN go install github.com/air-verse/air@latest

COPY . .

EXPOSE 3030

CMD ["air", "-c", ".air.toml"]

FROM golang:1.23 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o ./main

FROM gcr.io/distroless/base-debian11 AS production

WORKDIR /

COPY --from=builder /app/main /main

EXPOSE 8080

ENTRYPOINT ["/main"]
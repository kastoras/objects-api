FROM golang:1.23.2 AS development

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

RUN go install github.com/air-verse/air@latest
COPY . .

EXPOSE 3030

CMD ["air", "-c", ".air.toml"]
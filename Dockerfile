FROM golang:1.22-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /cpf_importer

EXPOSE 8080

CMD ["/cpf_importer"]
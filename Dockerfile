FROM golang:1.22-alpine

WORKDIR /app

COPY . .

EXPOSE 24750

CMD ["go", "run", "main.go"]

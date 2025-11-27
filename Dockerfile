FROM golang:1.24-alpine

WORKDIR /app

#RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o main ./main.go

EXPOSE 7070

CMD ["./main", "start"]
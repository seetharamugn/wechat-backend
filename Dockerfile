FROM golang:1.20-alpine3.18
LABEL authors="seetharam"
WORKDIR /app
COPY . .
RUN go build -o main .
EXPOSE 8081
ENTRYPOINT ["./main"]
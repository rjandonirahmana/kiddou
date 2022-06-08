FROM golang:alpine AS builder
RUN mkdir /app
WORKDIR /app
COPY . /app
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o kiddou main.go

FROM alpine
RUN mkdir /app
WORKDIR /app
RUN apk --no-cache add ca-certificates
RUN apk add --no-cache tzdata
COPY --from=builder /app/kiddou /app/kiddou
EXPOSE 8080
CMD ["/app/kiddou"]
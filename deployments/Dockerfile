# Билд стадия
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Копируем ВСЕ файлы из корня проекта в /app
COPY ../go .

# Билдим приложение
RUN go build -o main .

# Финальная стадия  
FROM alpine:3.18

WORKDIR /app

COPY --from=builder /app/main .

RUN apk --no-cache add ca-certificates

EXPOSE 8080

ENTRYPOINT [ "./main" ]
# Stage 1: Builder
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Устанавливаем минимальные зависимости для сборки
RUN apk add --no-cache git ca-certificates

# Копируем файлы зависимостей для кэширования
COPY go.mod go.sum ./
RUN go mod download

# Копируем весь исходный код
COPY . .

# Собираем бинарник с конкретным именем
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-w -s" \
    -o /app/downloader-converter-pricelists \
    ./cmd/app/main.go

# Stage 2: Runtime с Cron
FROM alpine:3.19

# Устанавливаем Cron, bash и утилиты
RUN apk add --no-cache \
    ca-certificates \
    tzdata \
    bash \
    curl \
    dcron \
    && mkdir -p /app/{logs,configs,data} \
    && ln -sf /usr/share/zoneinfo/Europe/Moscow /etc/localtime \
    && echo "Europe/Moscow" > /etc/timezone

# Копируем бинарник из стадии сборки
COPY --from=builder /app/downloader-converter-pricelists /app/

# Копируем конфигурацию Cron и скрипты
COPY docker/cron/entrypoint.sh /app/entrypoint.sh
COPY docker/cron/healthcheck.sh /app/healthcheck.sh
COPY docker/cron/run-with-retry.sh /app/run-with-retry.sh

# Создаем базовый конфиг если нужно
RUN echo "# Production configuration" > /app/configs/config.yaml \
    && echo "log_dir: /app/logs" >> /app/configs/config.yaml \
    && echo "data_dir: /app/data" >> /app/configs/config.yaml

# Настраиваем Cron расписание (07:00 MSK = 04:00 UTC)
RUN echo "# Run downloader-converter-pricelists daily at 07:00 MSK (04:00 UTC)" > /etc/crontabs/root \
    && echo "0 4 * * * /app/run-with-retry.sh 3 30 >> /app/logs/cron-execution.log 2>&1" >> /etc/crontabs/root \
    && echo "# Clean old logs every Sunday at 03:00" >> /etc/crontabs/root \
    && echo "0 3 * * 0 find /app/logs -name \"*.log\" -mtime +30 -delete" >> /etc/crontabs/root \
    && crontab /etc/crontabs/root

# Делаем скрипты исполняемыми
RUN chmod +x /app/entrypoint.sh \
    /app/healthcheck.sh \
    /app/run-with-retry.sh \
    /app/downloader-converter-pricelists

# Создаем непривилегированного пользователя для безопасности
RUN addgroup -S appgroup && adduser -S appuser -G appgroup \
    && chown -R appuser:appgroup /app

USER appuser
WORKDIR /app

# Health check для мониторинга
HEALTHCHECK --interval=30s --timeout=3s --start-period=10s --retries=3 \
    CMD /app/healthcheck.sh

# Запускаем Cron в foreground
CMD ["/app/entrypoint.sh"]
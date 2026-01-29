#!/bin/bash
set -e

# Проверяем, что Cron работает
if ! pgrep crond > /dev/null; then
    echo "Cron daemon is not running"
    exit 1
fi

# Проверяем, что есть последний лог за последние 24 часа
if [ -f "/app/logs/cron.log" ]; then
    LAST_LOG_TIME=$(stat -c %Y /app/logs/cron.log 2>/dev/null || echo 0)
    CURRENT_TIME=$(date +%s)
    MAX_AGE=$((24 * 3600))  # 24 часа
    
    if [ $((CURRENT_TIME - LAST_LOG_TIME)) -gt $MAX_AGE ]; then
        echo "No recent activity in logs (older than 24h)"
        exit 1
    fi
fi

# Проверяем наличие бинарника
if [ ! -f "/app/downloader-converter-pricelists" ]; then
    echo "Binary not found"
    exit 1
fi

# Проверяем доступность дискового пространства
DISK_USAGE=$(df /app --output=pcent | tail -1 | tr -d '% ')
if [ "$DISK_USAGE" -gt 90 ]; then
    echo "Disk usage is high: ${DISK_USAGE}%"
    exit 1
fi

echo "Container is healthy"
exit 0
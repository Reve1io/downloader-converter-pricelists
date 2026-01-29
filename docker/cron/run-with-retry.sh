#!/bin/bash
set -e

# Параметры:
# $1 - количество попыток (по умолчанию 3)
# $2 - таймаут в секундах (по умолчанию 30)
MAX_RETRIES=${1:-3}
TIMEOUT_SEC=${2:-30}
RETRY_DELAY=10  # Задержка между попытками в секундах

LOG_FILE="/app/logs/execution_$(date +%Y%m%d_%H%M%S).log"
LOCK_FILE="/tmp/downloader.lock"

echo "================================================" >> "$LOG_FILE"
echo "🔄 Starting execution: $(date '+%Y-%m-%d %H:%M:%S')" >> "$LOG_FILE"
echo "Max retries: $MAX_RETRIES" >> "$LOG_FILE"
echo "Timeout: ${TIMEOUT_SEC}s" >> "$LOG_FILE"
echo "================================================" >> "$LOG_FILE"

# Проверяем блокировку
if [ -f "$LOCK_FILE" ]; then
    echo "⚠️  Previous execution is still running. Skipping." >> "$LOG_FILE"
    echo "⚠️  Previous execution is still running. Skipping."
    exit 0
fi

# Создаем блокировку
touch "$LOCK_FILE"
trap 'rm -f "$LOCK_FILE"' EXIT

# Функция для выполнения задачи с таймаутом
execute_task() {
    local attempt=$1
    echo "" >> "$LOG_FILE"
    echo "▶️  Attempt $attempt/$MAX_RETRIES started at $(date '+%H:%M:%S')" >> "$LOG_FILE"
    
    # Запускаем приложение с таймаутом
    timeout $TIMEOUT_SEC /app/downloader-converter-pricelists >> "$LOG_FILE" 2>&1
    return $?
}

# Основной цикл с повторными попытками
for ((attempt=1; attempt<=MAX_RETRIES; attempt++)); do
    echo "🔄 Attempt $attempt/$MAX_RETRIES..." | tee -a "$LOG_FILE"
    
    execute_task $attempt
    EXIT_CODE=$?
    
    if [ $EXIT_CODE -eq 0 ]; then
        echo "✅ Success on attempt $attempt!" >> "$LOG_FILE"
        echo "✅ Successfully completed at $(date '+%H:%M:%S')" >> "$LOG_FILE"
        echo "✅ Execution successful!"
        
        # Копируем успешный лог в основной файл
        tail -50 "$LOG_FILE" >> /app/logs/cron.log
        exit 0
    elif [ $EXIT_CODE -eq 124 ]; then
        echo "⏱️  Timeout exceeded on attempt $attempt" >> "$LOG_FILE"
        echo "⏱️  Timeout exceeded (${TIMEOUT_SEC}s)"
    else
        echo "❌ Failed with code $EXIT_CODE on attempt $attempt" >> "$LOG_FILE"
        echo "❌ Failed with exit code: $EXIT_CODE"
    fi
    
    # Если это не последняя попытка, ждем перед повторением
    if [ $attempt -lt $MAX_RETRIES ]; then
        echo "⏳ Waiting ${RETRY_DELAY}s before next attempt..." >> "$LOG_FILE"
        echo "⏳ Waiting ${RETRY_DELAY}s before retry..."
        sleep $RETRY_DELAY
    fi
done

# Все попытки исчерпаны
echo "================================================" >> "$LOG_FILE"
echo "❌ All $MAX_RETRIES attempts failed!" >> "$LOG_FILE"
echo "Last exit code: $EXIT_CODE" >> "$LOG_FILE"
echo "================================================" >> "$LOG_FILE"

# Копируем лог ошибки в основной файл
cat "$LOG_FILE" >> /app/logs/cron.log

echo "❌ All attempts failed. Check logs: $LOG_FILE"
exit 1
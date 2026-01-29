#!/bin/bash
set -e

echo "================================================"
echo "📅 Downloader-Converter-Pricelists Cron Container"
echo "================================================"
echo "Start Time: $(date '+%Y-%m-%d %H:%M:%S %Z')"
echo "Timezone: $(cat /etc/timezone)"
echo "User: $(whoami)"
echo "Working Directory: $(pwd)"
echo ""

# Создаем необходимые директории
mkdir -p /app/logs /app/data /app/configs

# Инициализируем лог-файлы
touch /app/logs/cron.log
touch /app/logs/cron-execution.log
touch /app/logs/application.log

echo "📁 Directory structure created"
echo "📊 Log files initialized"
echo ""

# Загружаем конфигурацию если есть
if [ -f "/app/configs/config.yaml" ]; then
    echo "⚙️  Configuration loaded: /app/configs/config.yaml"
    export CONFIG_PATH="/app/configs/config.yaml"
fi

# Логируем переменные окружения (без секретов)
echo "🔧 Environment variables:"
env | grep -E "^(TZ|CONFIG_PATH|LOG_|DATA_)" | sort
echo ""

# Показываем расписание Cron
echo "📅 Cron schedule:"
crontab -l
echo ""

echo "🚀 Starting cron daemon in foreground..."
echo "================================================"
echo ""

# Запускаем Cron с логированием в файл и stdout
exec crond -f -l 8
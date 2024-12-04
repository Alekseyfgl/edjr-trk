#!/bin/sh

# Определяем путь к директории, в которой находится скрипт
SCRIPT_DIR=$(dirname "$(realpath "$0")")
PROJECT_DIR=$(dirname "$SCRIPT_DIR") # Поднимаемся на уровень выше
CMD_DIR="$PROJECT_DIR/cmd"
ENV_FILE="$PROJECT_DIR/.env"

# Проверка, существует ли файл .env
if [ ! -f "$ENV_FILE" ]; then
    echo "Ошибка: Файл $ENV_FILE не найден!"
    exit 1
fi

# Экспорт переменных из .env
export $(grep -v '^#' "$ENV_FILE" | xargs)

# Проверка обязательных переменных
if [ -z "$SERV_PORT" ] || [ -z "$MONGO_URI" ]; then
    echo "Ошибка: Не заданы необходимые переменные окружения (SERV_PORT или MONGO_URI)!"
    exit 1
fi

# Переход в папку с main.go
cd "$CMD_DIR" || exit

# Сборка приложения
echo "Сборка приложения..."
go build -o "$PROJECT_DIR/app" ./main.go

if [ $? -ne 0 ]; then
    echo "Ошибка: Сборка завершилась с ошибкой!"
    exit 1
fi

echo "Сборка завершена успешно!"

# Запуск приложения
echo "Запуск приложения на порту $SERV_PORT..."
"$PROJECT_DIR/app"

if [ $? -ne 0 ]; then
    echo "Ошибка: Приложение завершилось с ошибкой!"
    exit 1
fi

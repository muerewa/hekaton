# Hekaton - Система мониторинга с YAML конфигурацией

Hekaton — это легковесная и гибкая система мониторинга, написанная на Go, которая позволяет отслеживать состояние сервера через выполнение bash-команд и отправлять уведомления при срабатывании заданных условий.

## Основные возможности

- Конфигурация через YAML — простая настройка правил мониторинга без перекомпиляции
- Выполнение bash-команд — мониторинг любых системных параметров через shell-скрипты
- Гибкая система сравнения — поддержка операторов `>`, `<`, `>=`, `<=`, `==`, `!=`
- Telegram уведомления — мгновенные оповещения с шаблонизацией сообщений
- Параллельное выполнение — каждое правило работает в отдельной горутине
- Graceful shutdown — корректное завершение всех операций при остановке

## Quick Start

### Требования

- Go 1.19+
- Доступ к bash (Linux/macOS)
- Telegram Bot Token (для уведомлений)

### Установка

git clone https://github.com/muerewa/hekaton.git
cd hekaton
go mod tidy
go build -o hekaton .

### Конфигурация

Создайте файл `config.yaml`:

    name: "cpu_temp_high"
    bash: "sensors | awk '/^Package id 0:/ {print int($4)}'"
    compare:
    operator: ">"
    value: 80
    actions:

        type: "telegram"
        params:
        token: "YOUR_BOT_TOKEN"
        chat_id: "YOUR_CHAT_ID"
        message: "🔥 Температура CPU: {{.Result}}°C!"

    name: "disk_usage_critical"
    bash: "df --output=pcent / | tail -1 | tr -dc '0-9'"
    compare:
    operator: ">="
    value: 95
    actions:

        type: "telegram"
        params:
        token: "YOUR_BOT_TOKEN"
        chat_id: "YOUR_CHAT_ID"
        message: "💾 Диск заполнен на {{.Result}}%!"


### Запуск

./argusyaml -config config.yaml

text

## Структура конфигурации

| Поле               | Описание                         | Тип           | Обязательно |
|--------------------|----------------------------------|---------------|-------------|
| `name`             | Уникальное имя правила           | string        | Да          |
| `bash`             | Команда для выполнения           | string        | Да          |
| `compare.operator` | Оператор сравнения               | string        | Да          |
| `compare.value`    | Пороговое значение               | number/string | Да          |
| `actions`          | Список действий при срабатывании | array         | Да          |

## Поддерживаемые операторы

- `==` — равенство
- `!=` — неравенство  
- `>` — больше
- `<` — меньше
- `>=` — больше или равно
- `<=` — меньше или равно

## Настройка Telegram бота

1. Создайте бота через @BotFather: https://t.me/botfather
2. Получите токен бота
3. Узнайте ваш Chat ID через @userinfobot: https://t.me/userinfobot
4. Добавьте эти данные в конфигурацию


## Сборка для разных платформ

Linux

GOOS=linux GOARCH=amd64 go build -o argusyaml-linux .
macOS

GOOS=darwin GOARCH=amd64 go build -o argusyaml-macos .
Windows

GOOS=windows GOARCH=amd64 go build -o argusyaml.exe .

## Примеры использования

### Мониторинг сервисов

    name: "nginx_status"
    bash: "systemctl is-active nginx"
    compare:
    operator: "!="
    value: "active"
    actions:

        type: "telegram"
        params:
        token: "BOT_TOKEN"
        chat_id: "CHAT_ID"
        message: "🚨 Nginx не запущен! Статус: {{.Result}}"


### Мониторинг памяти

    name: "memory_usage"
    bash: "free | awk '/^Mem:/ {print int($3/$2 * 100)}'"
    compare:
    operator: ">="
    value: 90
    actions:

        type: "telegram"
        params:
        token: "BOT_TOKEN"
        chat_id: "CHAT_ID"
        message: "⚠️ Использование памяти: {{.Result}}%"

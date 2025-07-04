[![CI сборка](https://github.com/muerewa/hekaton/actions/workflows/go.yml/badge.svg)](https://github.com/muerewa/hekaton/actions/workflows/go.yml)

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

- git clone https://github.com/muerewa/hekaton.git
- cd hekaton
- go mod tidy
- go build ./cmd/hekaton

### Конфигурация

Создайте файл `config.yaml`:
```
- name: "cpu_temp_high"
  bash: "sensors | awk '/^Package id 0:/ {print int($4)}'"
  interval: "60s"           # Интервал проверки: 60s, 5m, 1h и т.д.
  compare:
    operator: ">"
    value: 80
  actions:
    - type: "telegram"
      params:
        token: "BOT_TOKEN"
        chat_id: "CHAT_ID"
        message: "Температура CPU: {{.Result}}°C!"
  timeout: 10               # (опционально) Таймаут выполнения bash-команды, секунд
  retries: 2                # (опционально) Количество повторов при ошибке

- name: "disk_usage_critical"
  bash: "df --output=pcent / | tail -1 | tr -dc '0-9'"
  datatype: "int"
  interval: "5m"
  compare:
    operator: ">="
    value: 95
  actions:
    - type: "telegram"
      params:
        token: "BOT_TOKEN"
        chat_id: "CHAT_ID"
        message: "Диск заполнен на {{.Result}}%!"
  timeout: 10
  retries: 1

- name: "nginx_not_running"
  bash: "systemctl is-active nginx"
  datatype: "string"
  interval: "2m"
  compare:
    operator: "!="
    value: "active"
  actions:
      params:
        smtp_host: "smtp.gmail.com"
        smtp_port: "587"
        username: "your-email@gmail.com"
        password: "your-app-password" # Используйте пароль приложения!
        from: "Monitoring System <monitor@example.com>"
        to: "admin@example.com,alert@example.com"
        subject: "Nginx not running"
        body: "Nginx status: {{.Result}}%!"
  timeout: 5
  retries: 3
```

### Запуск

./hekaton -config=config.yaml

## Структура конфигурации

| Поле               | Описание                         | Тип           | Обязательно |
|--------------------|----------------------------------|---------------|-------------|
| `name`             | Уникальное имя правила           | string        | Да          |
| `bash`             | Команда для выполнения           | string        | Да          |
| `compare.operator` | Оператор сравнения               | string        | Да          |
| `compare.value`    | Пороговое значение               | number/string | Да          |
| `actions`          | Список действий при срабатывании | array         | Да          |
| `interval`         | Интервал проверки условия        | number/string | Нет         |
| `timeout`          | Время выполнения команды         | number/string | Нет         |
| `retries`          | Кол-во ретраев команды при ошибке| number        | Нет         |

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

GOOS=linux GOARCH=amd64 go build -o hekaton-linux .
macOS

GOOS=darwin GOARCH=amd64 go build -o hekaton-macos .
Windows

## Примеры использования

### Мониторинг сервисов

```
- name: "nginx_not_running"
  bash: "systemctl is-active nginx"
  datatype: "string"
  interval: "2m"
  compare:
    operator: "!="
    value: "active"
  actions:
    - type: "telegram"
      params:
        token: "BOT_TOKEN"
        chat_id: "CHAT_ID"
        message: "Сервис nginx не запущен, статус: {{.Result}}"
  timeout: 5
  retries: 3
```

### Мониторинг диска
```
- name: "disk_usage_critical"
  bash: "df --output=pcent / | tail -1 | tr -dc '0-9'"
  datatype: "int"
  interval: "5m"
  compare:
    operator: ">="
    value: 95
  actions:
    - type: "telegram"
      params:
        token: "BOT_TOKEN"
        chat_id: "CHAT_ID"
        message: "Диск заполнен на {{.Result}}%!"
  timeout: 10
  retries: 1
```

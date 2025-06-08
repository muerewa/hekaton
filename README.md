🛡️ Hekaton - Система мониторинга с YAML конфигурацией
Hekaton — это легковесная и гибкая система мониторинга, написанная на Go, которая позволяет отслеживать состояние сервера через выполнение bash-команд и отправлять уведомления при срабатывании заданных условий. Система названа в честь Аргуса — многоглазого стража из греческой мифологии, символизируя всевидящее око мониторинга.
✨ Основные возможности

    Конфигурация через YAML — простая настройка правил мониторинга без перекомпиляции
    Выполнение bash-команд — мониторинг любых системных параметров через shell-скрипты
    Гибкая система сравнения — поддержка операторов >, <, >=, <=, ==, !=
    Telegram уведомления — мгновенные оповещения с шаблонизацией сообщений
    Параллельное выполнение — каждое правило работает в отдельной горутине
    Graceful shutdown — корректное завершение всех операций при остановке

🚀 Быстрый запуск
Требования

    Go 1.19+
    Доступ к bash (Linux/macOS)
    Telegram Bot Token (для уведомлений)

Установка

bash
git clone https://github.com/username/argusyaml.git
cd argusyaml
go mod tidy
go build -o argusyaml .

Конфигурация
Создайте файл config.yaml со следующим содержимым:

text
- name: "cpu_temp_high"
  bash: "sensors | awk '/^Package id 0:/ {print int($4)}'"
  compare:
    operator: ">"
    value: 80
  actions:
    - type: "telegram"
      params:
        token: "YOUR_BOT_TOKEN"
        chat_id: "YOUR_CHAT_ID"
        message: "🔥 Температура CPU: {{.Result}}°C!"

- name: "disk_usage_critical"
  bash: "df --output=pcent / | tail -1 | tr -dc '0-9'"
  compare:
    operator: ">="
    value: 95
  actions:
    - type: "telegram"
      params:
        token: "YOUR_BOT_TOKEN"
        chat_id: "YOUR_CHAT_ID"
        message: "💾 Диск заполнен на {{.Result}}%!"

Запуск

bash
./hekaton -config config.yaml

📋 Структура конфигурации
Поле	Описание	Тип	Обязательное
name	Уникальное имя правила	string	✅
bash	Команда для выполнения	string	✅
compare.operator	Оператор сравнения	string	✅
compare.value	Пороговое значение	number/string	✅
actions	Список действий при срабатывании	array	✅
🔧 Поддерживаемые операторы

    == — равенство
    != — неравенство
    > — больше
    < — меньше
    >= — больше или равно
    <= — меньше или равно

📱 Настройка Telegram бота

    Создайте бота через @BotFather
    Получите токен бота
    Узнайте ваш Chat ID через @userinfobot
    Введите эти данные в конфигурацию

Сборка для разных платформ

bash
# Linux
GOOS=linux GOARCH=amd64 go build -o hekaton-linux .

# macOS
GOOS=darwin GOARCH=amd64 go build -o hekaton-macos .

📊 Примеры использования
Мониторинг сервисов

text
- name: "nginx_status"
  bash: "systemctl is-active nginx"
  compare:
    operator: "!="
    value: "active"
  actions:
    - type: "telegram"
      params:
        token: "BOT_TOKEN"
        chat_id: "CHAT_ID"
        message: "🚨 Nginx не запущен! Статус: {{.Result}}"

Мониторинг памяти

text
- name: "memory_usage"
  bash: "free | awk '/^Mem:/ {print int($3/$2 * 100)}'"
  compare:
    operator: ">"
    value: 90
  actions:
    - type: "telegram"
      params:
        token: "BOT_TOKEN"
        chat_id: "CHAT_ID"
        message: "⚠️ Использование памяти: {{.Result}}%"


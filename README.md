Базовая структура проекта (для Go 1.11+ с модулями)

task-manager/
├── cmd/                  # Основные приложения проекта (точки входа)
│   └── server/           # Главный сервер (можно добавить cmd/cli для CLI-версии)
│       └── main.go       # package main — точка входа
├── internal/             # Внутренние пакеты (доступ только внутри проекта)
│   ├── handlers/         # HTTP-обработчики (например, task_handler.go)
│   ├── models/           # Структуры данных (Task, User и т.д.)
│   ├── storage/          # Логика хранения (memory.go, postgres.go)
│   └── middleware/       # Middleware (logging.go, auth.go)
├── pkg/                  # Пакеты, которые можно использовать в других проектах
│   └── utils/            # Утилиты (например, валидация дат)
├── api/                  # OpenAPI/Swagger-спецификации (если нужно)
├── configs/              # Конфиги (config.yaml, .env)
├── migrations/           # SQL-миграции (если используем БД)
├── scripts/              # Вспомогательные скрипты (деплой, генерация кода)
├── go.mod                # Файл модуля Go (создаётся через go mod init)
└── README.md             # Описание проекта
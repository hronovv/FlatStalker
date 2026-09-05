# FlatStalker

Telegram-бот (@flatstalker_bot), который следит за новыми объявлениями аренды на Kufar и присылает уведомления, когда появляется что-то новое по твоему поиску.

Пользователь добавляет ссылку на поиск Kufar (с фильтрами — район, цена, комнаты и т.д.) через Mini App. Воркер периодически опрашивает Kufar API, сравнивает с уже просмотренными объявлениями и шлёт в Telegram только новые.

## Архитектура

```
Telegram Bot  ←→  Go backend (HTTP API + bot + worker)
                      ↓
                 PostgreSQL
                      ↓
                 Kufar API
```

В одном бинарном файле всё сразу:

- **Telegram-бот** — команды `/start`, `/links`, `/status`, `/help`, inline-кнопки для паузы/удаления ссылок
- **HTTP API** — кабинет в Mini App (`/api/me`, `/api/links`, `/api/pay`)
- **Worker** — три горутины по тарифам (free / plus / pro), каждая со своим интервалом

При первом добавлении ссылки все текущие объявления помечаются как просмотренные, уведомления не отправляются. Дальше приходят только новые.

## Тарифы


|                   | Free  | Plus  | Pro    |
| ----------------- | ----- | ----- | ------ |
| Ссылок            | 1     | 3     | 5      |
| Интервал проверки | 5 мин | 2 мин | 30 сек |


Интервалы настраиваются в `.env` (`WORKER_INTERVAL_*`). Оплата — через Telegram Payments (BYN), периоды от 1 до 180 дней.

## Mini App

Статический фронт в `web/` — кабинет внутри Telegram:

- добавление и управление ссылками
- статус тарифа
- покупка Plus/Pro
- русский и белорусский язык

## Стек

- Go 1.26, Gin
- PostgreSQL (pgx)
- [go-telegram/bot](https://github.com/go-telegram/bot)
- Vanilla JS + Telegram WebApp SDK



## Структура

```
cmd/app/           точка входа
internal/
  api/             HTTP handlers, Telegram auth
  bot/             команды и уведомления
  worker/          опрос Kufar по расписанию
  source/kufar/    парсинг URL, запросы к API
  plan/            тарифы, лимиты, цены
  repository/      работа с БД
  tgauth/          валидация initData Mini App
migrations/        SQL-миграции
web/               Mini App (HTML/CSS/JS)
```



## Локальный запуск

```bash
cp .env.example .env   # заполни BOT_TOKEN, DATABASE_URL и остальное
go run ./cmd/app
```

Фронт отдельно (для разработки):

```bash
make web   # http://localhost:5173
```

В `.env` добавь `http://localhost:5173` в `CORS_ORIGINS`.

### Миграции

```bash
psql "$DATABASE_URL" -f migrations/001_listings.sql
psql "$DATABASE_URL" -f migrations/002_seen_ads.sql
# ... по порядку до 006_payments.sql
```



## Переменные окружения

Все обязательные — см. `.env.example`. Основные:


| Переменная               | Назначение                                        |
| ------------------------ | ------------------------------------------------- |
| `BOT_TOKEN`              | токен Telegram-бота                               |
| `DATABASE_URL`           | PostgreSQL                                        |
| `HTTP_ADDR`              | адрес API (`:8080` или `127.0.0.1:8080` за nginx) |
| `WORKER_INTERVAL_*`      | интервалы опроса по тарифам                       |
| `CORS_ORIGINS`           | разрешённые origin для Mini App                   |
| `PAYMENT_PROVIDER_TOKEN` | токен Telegram Payments                           |
| `TELEGRAM_SUPPORT`       | контакт поддержки в `/help`                       |




## Deploy

Сборка:

```bash
GOOS=linux GOARCH=amd64 go build -o flatstalker ./cmd/app
```

На сервере: PostgreSQL, systemd-сервис для бинарника, nginx для статики (`web/`) и прокси `/api/` → backend. Логи с локальной машины:

```bash
make logs   # нужен SERVER=root@IP в .env
```



## API

Все `/api/*` эндпоинты требуют Telegram initData (`Authorization: tma ...` или заголовок `X-Telegram-Init-Data`).


| Метод  | Путь             | Описание                     |
| ------ | ---------------- | ---------------------------- |
| GET    | `/health`        | healthcheck                  |
| GET    | `/api/me`        | профиль, тариф, ссылки, цены |
| POST   | `/api/links`     | добавить ссылку Kufar        |
| PATCH  | `/api/links/:id` | пауза / возобновление        |
| DELETE | `/api/links/:id` | удалить ссылку               |
| POST   | `/api/pay`       | выставить счёт на тариф      |



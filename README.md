# Training planner HTTP API

Сервис слушает HTTP на порту `8080` по умолчанию. Для локального запуска вместе с PostgreSQL:

```bash
docker compose up --build
```

## Создание тренировки

`POST /training/create` принимает `application/json`:

```bash
curl --request POST http://localhost:8080/training/create \
  --header 'Content-Type: application/json' \
  --data '{
    "user_id": 1,
    "training": {
      "date": "2026-08-19T18:30:00+03:00",
      "comment": "Интервальная тренировка",
      "items": [
        {
          "template_id": 2,
          "params": {
            "cross_template_params": {
              "speed": "4:30/km",
              "duration": 1200,
              "max_pulse": 170
            }
          }
        }
      ],
      "result_video_links": ["https://example.com/video/1"]
    }
  }'
```

Успешный ответ имеет статус `201 Created`:

```json
{"id": 42}
```

Для `multipart/form-data` передайте `user_id` и JSON тренировки в поле `training`:

```bash
curl --request POST http://localhost:8080/training/create \
  --form 'user_id=1' \
  --form 'training={"date":"2026-08-19T18:30:00+03:00","comment":"Интервалы","items":[{"template_id":2,"params":{}}]}'
```

Вместо двух полей можно передать поле `payload` с полным JSON из первого примера. Загрузка файлов пока явно отклоняется: для неё нужно определить файловое хранилище и формат возвращаемых ссылок. Максимальный размер запроса — 10 MiB.

## HTTPS

Если задать обе переменные, сервер запускается в TLS-режиме:

```bash
TLS_CERT_FILE=/path/to/cert.pem \
TLS_KEY_FILE=/path/to/key.pem \
DATABASE_URL='postgres://training_planner:training_planner@localhost:5432/training_planner?sslmode=disable' \
go run ./cmd/service
```

Без этих переменных используется обычный HTTP. Проверка готовности доступна по `GET /health`.

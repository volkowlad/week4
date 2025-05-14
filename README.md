# week4

REST API-сервис, написанный на Go с использованием фреймворка Fiber. Сервис предоставляет базовый функционал для управления задачами. Использовалось in memory хранение данных.

Реализовано:

- Создание задач через API
- Валидация входных данных
- Логирование с использованием `zap`
- Хранение данных в in memory

---

## **Настройка проекта**

### **Создание `.env` файла**

Создайте `.env` файл и пропишите параметры:

```
# General application configuration
LOG_LEVEL=info

# REST API configuration
PORT=8080
WRITE_TIMEOUT=15s
SERVER_NAME=SimpleService
TOKEN=123

```

## **Запуск сервиса**

### **Локальный запуск**

```
go run cmd/main.go

```

Сервис будет доступен по адресу `http://localhost:8080`

---

## **Тестирование API**

### **Создание задачи**

**Запрос:**

```
POST http://localhost:8080/v1/tasks
Content-Type: application/json
Authorization: Bearer your_secret_token

```

```
{
  "title": "New Feature",
  "description": "Develop new API endpoint"
}

```

**Ответ:**

```
{
  "status": "success",
  "data": {
    "task_id": "cfbfdf02-ea2e-45fa-ae88-3fd81ab939bf"
  }
}

```

### **Получение задачи по id**

**Запрос:**

```
GET http://localhost:8080/v1/task/cfbfdf02-ea2e-45fa-ae88-3fd81ab939bf
Content-Type: application/json
Authorization: Bearer your_secret_token

```

**Ответ:**

```
{
    "status": "success",
    "data": {
        "task": {
            "id": "fac89782-bfc0-438f-a87c-f2203402b66e",
            "title": "test",
            "description": "test",
            "status": "new",
            "created": "2025-05-14T18:40:57.673258+03:00",
            "updated": "2025-05-14T18:40:57.673259+03:00"
        }
    }
}

```

### **Получение всех задач**

**Запрос:**

```
GET http://localhost:8080/v1/tasks
Content-Type: application/json
Authorization: Bearer your_secret_token

```

**Ответ:**

```
{
    "status": "success",
    "data": {
        "tasks": {
            "id": "fac89782-bfc0-438f-a87c-f2203402b66e",
            "title": "test",
            "description": "test",
            "status": "new",
            "created": "2025-05-14T18:40:57.673258+03:00",
            "updated": "2025-05-14T18:40:57.673259+03:00"
        }
    }
}

```

### **Удаление задачи по id**

**Запрос:**

```
DELETE http://localhost:8080/v1/delete/fac89782-bfc0-438f-a87c-f2203402b66e
Content-Type: application/json
Authorization: Bearer your_secret_token

```

**Ответ:**

```
{
    "status": "success"
}

```

### **Изменение задачи**

**Запрос:**

```
PUT http://localhost:8080/v1/update/fac89782-bfc0-438f-a87c-f2203402b66e
Content-Type: application/json
Authorization: Bearer your_secret_token

```
```
{
    "title": "test1",
    "description": "test1",
    "status": "in_progress"
}
```

**Ответ:**

```
{
    "status": "success",
    "data": {
        "task": {
            "id": "cfbfdf02-ea2e-45fa-ae88-3fd81ab939bf",
            "title": "test1",
            "description": "test1",
            "status": "in_progress",
            "created": "2025-05-14T18:55:54.681415+03:00",
            "updated": "2025-05-14T18:56:05.060933+03:00"
        }
    }
}

```

---

## **Дополнительная информация**

- Логирование ведётся через `zap.Logger`
- Переменные окружения загружаются через `envconfig`

Сервис готов к работе.
# 🎮 ConsoleRent - Аренда игровых консолей

<div align="center">

**Сервис аренды игровых консолей PlayStation и Xbox с доставкой на дом**

## 📋 О проекте

**ConsoleRent** — это веб-приложение для аренды игровых консолей PlayStation 5, PlayStation 4, Xbox Series X|S и Xbox One X. Сервис позволяет пользователям арендовать консоли на нужный срок с доставкой на дом.

### 🎯 Основные возможности

- ✅ **Каталог консолей** - просмотр всех доступных консолей с фильтрацией по типу (PS5, PS4, Xbox)
- ✅ **Детальные карточки** - подробное описание, комплектация, список игр и отзывы
- ✅ **Система аренды** - выбор дат, указание адреса доставки, автоматический расчет стоимости
- ✅ **Личный кабинет** - история аренд, статус, возможность возврата консоли
- ✅ **Профиль пользователя** - просмотр и редактирование данных, смена пароля, статистика
- ✅ **Адаптивный дизайн** - полная поддержка мобильных устройств

### 💰 Цены на аренду

| Консоль                | Цена за сутки |
| ---------------------- | ------------- |
| PlayStation 5 Standard | 800 ₽         |
| PlayStation 5 Digital  | 800 ₽         |
| PlayStation 4 Slim     | 500 ₽         |
| Xbox Series X          | 800 ₽         |
| Xbox Series S          | 800 ₽         |
| Xbox One X             | 500 ₽         |

**🚚 Доставка бесплатная!**

---

## 🛠 Технологический стек

### Frontend

- **Vue 3** - прогрессивный JavaScript фреймворк
- **Vue Router** - маршрутизация
- **Axios** - HTTP клиент
- **Vite** - сборщик проекта
- **CSS** - адаптивная вёрстка

### Backend

- **Go (Golang)** - высокопроизводительный язык
- **Gin** - веб-фреймворк
- **JWT** - аутентификация
- **bcrypt** - хеширование паролей

### База данных

- **PostgreSQL 15** - реляционная база данных

### DevOps

- **Docker** - контейнеризация
- **Docker Compose** - оркестрация
- **Nginx** - веб-сервер

---

## 📥 Скачивание проекта

### Способ 1: Клонирование через Git

```bash
git clone https://github.com/BznikinPlay/GamePlay.git
cd GamePlay
```

### Способ 2: Скачать ZIP архив

1. Перейдите на https://github.com/BznikinPlay/GamePlay
2. Нажмите кнопку **Code** → **Download ZIP**
3. Распакуйте архив в нужную папку

---

## 🚀 Запуск проекта

### Требования для запуска

- **Docker Desktop** (Windows/Mac) или **Docker Engine** (Linux)
- **Git** (опционально, для клонирования)
- **8GB RAM** (рекомендуется)
- **Свободные порты:** 3000, 8080, 5432

### Быстрый запуск (рекомендуемый способ)

```bash
# Перейдите в папку проекта
cd GamePlay

# Запустите все сервисы через Docker Compose
docker-compose up --build

# Или в фоновом режиме
docker-compose up -d --build
```

### Пошаговый запуск

#### 1. Установите Docker Desktop

Скачайте и установите Docker Desktop с официального сайта:

- [Windows](https://www.docker.com/products/docker-desktop/)
- [Mac](https://www.docker.com/products/docker-desktop/)
- [Linux](https://docs.docker.com/engine/install/)

#### 2. Клонируйте репозиторий

```bash
git clone https://github.com/BznikinPlay/GamePlay.git
cd GamePlay
```

#### 3. Запустите Docker Compose

```bash
docker-compose up --build
```

#### 4. Дождитесь запуска всех сервисов

Вы увидите сообщения:

```
postgres_db | database system is ready to accept connections
go_backend   | Listening and serving HTTP on :8080
vue_frontend | Configuration complete; ready for start up
```

#### 5. Откройте приложение

- **Frontend:** http://localhost:3000
- **Backend API:** http://localhost:8080
- **PostgreSQL:** localhost:5432

---

## 🧪 Тестовые данные

### Тестовый пользователь

| Поле   | Значение         |
| ------ | ---------------- |
| Email  | test@example.com |
| Пароль | password123      |

### Или создайте своего пользователя

1. Откройте http://localhost:3000
2. Нажмите **"Войти / Регистрация"**
3. Перейдите на вкладку **"Регистрация"**
4. Заполните форму и нажмите **"Зарегистрироваться"**

---

## 📁 Структура проекта

```
GamePlay/
├── docker-compose.yml          # Docker Compose конфигурация
├── backend/
│   ├── Dockerfile              # Dockerfile для Go бэкенда
│   ├── go.mod                  # Go модуль
│   ├── go.sum                  # Go зависимости
│   ├── main.go                 # Точка входа
│   ├── models/                 # Модели данных
│   ├── handlers/               # Обработчики запросов
│   ├── database/               # Работа с БД
│   └── middleware/             # Middleware (JWT, CORS)
└── frontend/
    ├── Dockerfile              # Dockerfile для Vue фронтенда
    ├── package.json            # NPM зависимости
    ├── index.html              # Главный HTML файл
    ├── src/
    │   ├── main.js             # Точка входа фронтенда
    │   ├── App.vue             # Корневой компонент
    │   ├── router/             # Маршрутизация
    │   └── views/              # Страницы приложения
    │       ├── HomePage.vue    # Главная страница
    │       ├── CheckoutPage.vue # Оформление аренды
    │       ├── MyRentals.vue   # Мои аренды
    │       └── ProfilePage.vue # Профиль пользователя
    └── public/
        └── gamepad.png         # Иконка сайта
```

---

## 🔧 Команды для управления

### Docker Compose команды

```bash
# Запуск всех сервисов
docker-compose up

# Запуск в фоновом режиме
docker-compose up -d

# Остановка всех сервисов
docker-compose down

# Остановка с удалением томов БД
docker-compose down -v

# Пересборка и запуск
docker-compose up --build

# Просмотр логов
docker-compose logs -f

# Просмотр логов конкретного сервиса
docker-compose logs -f backend
docker-compose logs -f frontend
docker-compose logs -f postgres

# Перезапуск сервиса
docker-compose restart backend
```

### Локальный запуск (без Docker)

#### Backend (требуется Go 1.22+)

```bash
cd backend
go mod download
go run main.go
```

#### Frontend (требуется Node.js 20+)

```bash
cd frontend
npm install
npm run dev
```

---

## 📡 API Endpoints

### Публичные маршруты

| Метод | Endpoint            | Описание                  |
| ----- | ------------------- | ------------------------- |
| POST  | `/api/register`     | Регистрация пользователя  |
| POST  | `/api/login`        | Вход в систему            |
| GET   | `/api/consoles`     | Получение списка консолей |
| GET   | `/api/consoles/:id` | Получение консоли по ID   |

### Защищенные маршруты (требуют JWT токен)

| Метод | Endpoint                  | Описание                     |
| ----- | ------------------------- | ---------------------------- |
| GET   | `/api/user/profile`       | Получение профиля            |
| PUT   | `/api/user/profile`       | Обновление профиля           |
| POST  | `/api/rentals`            | Создание аренды              |
| GET   | `/api/my-rentals`         | Получение аренд пользователя |
| PUT   | `/api/rentals/:id/return` | Возврат консоли              |

---

## 🔐 Переменные окружения

### Backend (.env)

```env
DB_HOST=postgres
DB_PORT=5432
DB_USER=rental_user
DB_PASSWORD=rental_pass
DB_NAME=console_rental
JWT_SECRET=your-super-secret-jwt-key-change-in-production
```

### Frontend (.env)

```env
VITE_API_URL=http://localhost:8080/api
```

## 🗄️ Структура базы данных

```sql
-- Пользователи
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(100) NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Консоли
CREATE TABLE consoles (
    id SERIAL PRIMARY KEY,
    type VARCHAR(50) NOT NULL,
    model VARCHAR(100) NOT NULL,
    price_per_day DECIMAL(10,2) NOT NULL,
    is_available BOOLEAN DEFAULT TRUE,
    description TEXT,
    image_url VARCHAR(500)
);

-- Аренды
CREATE TABLE rentals (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id),
    console_id INTEGER REFERENCES consoles(id),
    start_date TIMESTAMP NOT NULL,
    end_date TIMESTAMP NOT NULL,
    total_price DECIMAL(10,2) NOT NULL,
    delivery_address TEXT NOT NULL,
    status VARCHAR(50) DEFAULT 'active',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

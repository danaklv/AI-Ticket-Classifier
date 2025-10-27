# 1. Базовый образ
FROM python:3.12-slim

# 2. Рабочая директория внутри контейнера
WORKDIR /app

# 3. Копирование зависимостей
COPY requirements.txt .

# 4. Установка зависимостей
RUN pip install -r requirements.txt

# 5. Копирование исходников
COPY . .

# 6. Команда запуска
CMD ["python", "consumer.py"]

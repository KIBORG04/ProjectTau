# Dockerfile

# Используем официальный образ Python
FROM python:3.9-slim

# Устанавливаем рабочую директорию внутри контейнера
WORKDIR /app

# Копируем файл с зависимостями
COPY requirements.txt .

# Устанавливаем зависимости
RUN pip install --no-cache-dir -r requirements.txt

# Копируем код нашего приложения
COPY sarima_service.py .

# Сообщаем Docker, что контейнер будет слушать порт 5000
EXPOSE 5000

# Команда для запуска нашего приложения
CMD ["python", "sarima_service.py"]
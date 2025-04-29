FROM python:3.9-slim 

WORKDIR /app

COPY requirements.txt .

RUN pip install --no-cache-dir -r requirements.txt

COPY . .

EXPOSE 5000

ENV REDIS_HOST="localhost"
ENV REDIS_PORT="6379"
ENV REDIS_QUEUE_NAME="garden_actions"

CMD ["python", "app.py"]
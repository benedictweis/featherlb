FROM python:3.9-slim

WORKDIR /app

COPY requirements.txt .

RUN pip install --no-cache-dir -r requirements.txt

COPY backend.py .

ENV FLASK_HOST=0.0.0.0
ENV FLASK_PORT=8080
ENV FLASK_MESSAGE="Hello, World!"

CMD ["python3", "backend.py"]

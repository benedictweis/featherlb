FROM python:3.9-slim

WORKDIR /app

COPY requirements.txt .

RUN pip install --no-cache-dir -r requirements.txt

COPY backend.py .

ENV FLASK_MESSAGE="Hello, World!"
ENV FLASK_PORT=8080

CMD ["python3", "backend.py"]

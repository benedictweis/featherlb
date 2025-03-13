from flask import Flask
import os

DEFAULT_MESSAGE = os.getenv('FLASK_MESSAGE', 'Hello, World!')
DEFAULT_PORT = int(os.getenv('FLASK_PORT', 8080))

app = Flask(__name__)

@app.route('/')
def hello_world():
    return DEFAULT_MESSAGE

if __name__ == '__main__':
    app.run(port=DEFAULT_PORT)

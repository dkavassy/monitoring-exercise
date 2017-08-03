#!/usr/bin/env python

import monitoring
import random
import time
from flask import Flask
app = Flask(__name__)

metrics = monitoring.create_metrics("requests", "errors")
metrics.export_to(app)

@app.route('/')
def hello_world():
  metrics("requests").increment()
  time.sleep(random.random())
  if not random.randint(0,10):
    metrics("errors").increment()
    return 'Internal Server Error', 500
  return 'Hello, World!'

if __name__ == '__main__':
    app.run()


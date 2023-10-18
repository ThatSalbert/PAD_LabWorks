from logging import Logger
import time
from functools import wraps
from flask import jsonify
import requests

def circuit_breaker(f):
    failures = 0
    last_failure_time = 0
    TASK_TIMEOUT = 0.05
    FAILURE_LIMIT = 3
    FAILURE_THRESHOLD = 3.5

    @wraps(f)
    def wrapper(*args, **kwargs):
        nonlocal failures
        nonlocal last_failure_time
        if failures >= FAILURE_LIMIT:
            time_since_last_failure = time.time() - last_failure_time
            if time_since_last_failure < TASK_TIMEOUT * FAILURE_THRESHOLD:
                return jsonify({'message': 'Service is down. Please try again later.'}), 503
        try:
            response = f(*args, **kwargs)
            failures = 0
            return response
        except requests.exceptions.Timeout:
            failures += 1
            last_failure_time = time.time()
            return jsonify({'message': 'Task Timeout'}), 408
    return wrapper
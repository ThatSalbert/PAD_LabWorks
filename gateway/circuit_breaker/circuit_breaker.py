from logging import Logger
import time
from functools import wraps
from flask import jsonify
import requests
import os


# def circuit_breaker(f):
#     failures = 0
#     last_failure_time = 0
#     TASK_TIMEOUT = float(os.getenv('TIMEOUT'))
#     FAILURE_LIMIT = int(os.getenv('FAILURE_LIMIT'))
#     FAILURE_THRESHOLD = float(os.getenv('FAILURE_THRESHOLD'))

#     @wraps(f)
#     def wrapper(*args, **kwargs):
#         nonlocal failures
#         nonlocal last_failure_time
#         if failures >= FAILURE_LIMIT:
#             time_since_last_failure = time.time() - last_failure_time
#             if time_since_last_failure < TASK_TIMEOUT * FAILURE_THRESHOLD:
#                 return jsonify({'message': 'Service is down. Please try again later.'}), 503
#         try:
#             response = f(*args, **kwargs)
#             failures = 0
#             return response
#         except requests.exceptions.Timeout:
#             failures += 1
#             last_failure_time = time.time()
#             return jsonify({'message': 'Task Timeout'}), 408

#     return wrapper

class CircuitBreaker:
    def __init__(self, task_timeout_limit, reroute_limit, threshold):
        self.state = "closed"
        self.failures = 0
        self.reroutes = 0
        self.task_timeout_limit = task_timeout_limit
        self.reroute_limit = reroute_limit
        self.threshold = threshold

    def transition(self, new_state):
        self.state = new_state
        if new_state == "open":
            self.open_time = time.time()

    def __call__(self, f):
        @wraps(f)
        def wrapper(*args, **kwargs):
            if self.state == "open":
                if time.time() - self.open_time > self.task_timeout_limit * self.threshold:
                    self.transition("half_open")
                else:
                    return jsonify({'message': 'Service is down. Please try again later.'}), 503
                    return
            try:
                response = f(*args, **kwargs)
                self.failures = 0
                self.reroutes = 0
                if self.state == "half_open":
                    self.transition("closed")
                return response
            except Exception as e:
                self.failures += 1
                if self.failures >= self.task_timeout_limit or self.reroutes >= self.reroute_limit:
                    self.transition("open")
                raise e
        return wrapper
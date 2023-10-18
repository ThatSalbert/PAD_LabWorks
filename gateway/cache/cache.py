from functools import wraps
from threading import Lock
from expiringdict import ExpiringDict
from flask import request

cache = ExpiringDict(max_len = 100, max_age_seconds = 15)
cache_lock = Lock()


def caching_mechanism(endpoint_function):
    @wraps(endpoint_function)
    def wrapper(*args, **kwargs):
        key = request.url
        with cache_lock:
            if key in cache:
                return cache[key]
        response = endpoint_function(*args, **kwargs)
        with cache_lock:
            cache[key] = response
        return response
    return wrapper

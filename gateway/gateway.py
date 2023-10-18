from functools import wraps
from threading import Lock
from flask import Flask, request, jsonify
from flask_limiter import Limiter
from flask_limiter.util import get_remote_address
from expiringdict import ExpiringDict
import requests

api = Flask(__name__)
cache = ExpiringDict(max_len = 100, max_age_seconds = 60)
cache_lock = Lock()
limiter = Limiter(get_remote_address, app = api, default_limits = ["200 per day", "50 per hour", "10 per minute"])


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


@api.route('/weather/locations', methods = ['GET'])
@limiter.limit("10 per minute", error_message = 'Too Many Requests', override_defaults = False)
@caching_mechanism
def get_locations():
    country = request.args.get('country')
    service_name = str(request.url_rule).split('/')[1]
    try:
        response_from_service_discovery = requests.get('http://localhost:8002/get_service?service_name=' + service_name,
                                                       timeout = 0.05)
    except requests.exceptions.Timeout:
        response = jsonify({'message': 'Service Discovery timed out'})
        return response, 408
    service_address = response_from_service_discovery.json()['service_address']
    if country is None:
        url = service_address + '/weather/locations'
    else:
        url = service_address + '/weather/locations?country=' + country
    try:
        response_from_service = requests.get(url)
    except requests.exceptions.Timeout:
        response = jsonify({'message': service_name + ' timed out'})
        return response, 408
    response = jsonify(response_from_service.json())
    return response, response_from_service.status_code


@api.route('/weather/current', methods = ['GET'])
@limiter.limit("10 per minute", error_message = 'Too Many Requests', override_defaults = False)
@caching_mechanism
def get_current_weather():
    country = request.args.get('country')
    city = request.args.get('city')
    service_name = str(request.url_rule).split('/')[1]
    try:
        response_from_service_discovery = requests.get('http://localhost:8002/get_service?service_name=' + service_name,
                                                       timeout = 0.05)
    except requests.exceptions.Timeout:
        response = jsonify({'message': 'Service Discovery timed out'})
        return response, 408
    service_address = response_from_service_discovery.json()['service_address']
    if country is None or city is None:
        url = service_address + '/weather/current'
    else:
        url = service_address + '/weather/current?country=' + country + '&city=' + city
    try:
        response_from_service = requests.get(url)
    except requests.exceptions.Timeout:
        response = jsonify({'message': service_name + ' timed out'})
        return response, 408
    response = jsonify(response_from_service.json())
    return response, response_from_service.status_code


@api.route('/weather/forecast', methods = ['GET'])
@limiter.limit("10 per minute", error_message = 'Too Many Requests', override_defaults = False)
@caching_mechanism
def get_weather_forecast():
    country = request.args.get('country')
    city = request.args.get('city')
    service_name = str(request.url_rule).split('/')[1]
    try:
        response_from_service_discovery = requests.get('http://localhost:8002/get_service?service_name=' + service_name,
                                                       timeout = 0.05)
    except requests.exceptions.Timeout:
        response = jsonify({'message': 'Service Discovery timed out'})
        return response, 408
    service_address = response_from_service_discovery.json()['service_address']
    if country is None or city is None:
        url = service_address + '/weather/forecast'
    else:
        url = service_address + '/weather/forecast?country=' + country + '&city=' + city
    try:
        response_from_service = requests.get(url, timeout = 0.05)
    except requests.exceptions.Timeout:
        response = jsonify({'message': service_name + ' timed out'})
        return response, 408
    response = jsonify(response_from_service.json())
    return response, response_from_service.status_code


@api.route('/weather/add_data', methods = ['POST'])
@limiter.limit("10 per minute", error_message = 'Too Many Requests', override_defaults = False)
def add_data():
    add_type = request.args.get('type')
    data = request.get_json()
    service_name = str(request.url_rule).split('/')[1]
    try:
        response_from_service_discovery = requests.get('http://localhost:8002/get_service?service_name=' + service_name,
                                                       timeout = 0.05)
    except requests.exceptions.Timeout:
        response = jsonify({'message': 'Service Discovery timed out'})
        return response, 408
    service_address = response_from_service_discovery.json()['service_address']
    url = service_address + '/weather/add_data?type=' + add_type
    try:
        response_from_service = requests.post(url, json = data)
    except requests.exceptions.Timeout:
        response = jsonify({'message': service_name + ' timed out'})
        return response, 408
    response = jsonify(response_from_service.json())
    return response, response_from_service.status_code


@api.route('/weather/update_data', methods = ['PUT'])
@limiter.limit("10 per minute", error_message = 'Too Many Requests', override_defaults = False)
def delete_data():
    type = request.args.get('type')
    data = request.get_json()
    service_name = str(request.url_rule).split('/')[1]
    try:
        response_from_service_discovery = requests.get('http://localhost:8002/get_service?service_name=' + service_name,
                                                       timeout = 0.05)
    except requests.exceptions.Timeout:
        response = jsonify({'message': 'Service Discovery timed out'})
        return response, 408
    service_address = response_from_service_discovery.json()['service_address']
    url = service_address + '/weather/update_data?type=' + type
    try:
        response_from_service = requests.put(url, json = data)
    except requests.exceptions.Timeout:
        response = jsonify({'message': service_name + ' timed out'})
        return response, 408
    response = jsonify(response_from_service.json())
    return response, response_from_service.status_code


@api.route('/disaster', methods = ['GET'])
@limiter.limit("10 per minute", error_message = 'Too Many Requests', override_defaults = False)
@caching_mechanism
def get_disasters():
    service_name = str(request.url_rule).split('/')[1]
    try:
        response_from_service_discovery = requests.get('http://localhost:8002/get_service?service_name=' + service_name,
                                                       timeout = 0.05)
    except requests.exceptions.Timeout:
        response = jsonify({'message': 'Service Discovery timed out'})
        return response, 408
    service_address = response_from_service_discovery.json()['service_address']
    url = service_address + '/disaster'
    try:
        response_from_service = requests.get(url)
    except requests.exceptions.Timeout:
        response = jsonify({'message': service_name + ' timed out'})
        return response, 408
    response = jsonify(response_from_service.json())
    return response, response_from_service.status_code


@api.route('/disaster/list', methods = ['GET'])
@limiter.limit("10 per minute", error_message = 'Too Many Requests', override_defaults = False)
@caching_mechanism
def get_disaster_list():
    country = request.args.get('country')
    city = request.args.get('city')
    active = request.args.get('active')
    service_name = str(request.url_rule).split('/')[1]
    try:
        response_from_service_discovery = requests.get('http://localhost:8002/get_service?service_name=' + service_name,
                                                       timeout = 0.05)
    except requests.exceptions.Timeout:
        response = jsonify({'message': 'Service Discovery timed out'})
        return response, 408
    service_address = response_from_service_discovery.json()['service_address']
    if country is None or city is None or active is None:
        url = service_address + '/disaster/list'
    else:
        url = service_address + '/disaster/list?country=' + country + '&city=' + city + '&active=' + active
    try:
        response_from_service = requests.get(url, timeout = 0.05)
    except requests.exceptions.Timeout:
        response = jsonify({'message': service_name + ' timed out'})
        return response, 408
    response = jsonify(response_from_service.json())
    return response, response_from_service.status_code


@api.route('/disaster/alert', methods = ['POST'])
@limiter.limit("10 per minute", error_message = 'Too Many Requests', override_defaults = False)
def add_alert():
    data = request.get_json()
    service_name = str(request.url_rule).split('/')[1]
    try:
        response_from_service_discovery = requests.get('http://localhost:8002/get_service?service_name=' + service_name,
                                                       timeout = 0.5)
    except requests.exceptions.Timeout:
        response = jsonify({'message': 'Service Discovery timed out'})
        return response, 408
    service_address = response_from_service_discovery.json()['service_address']
    url = service_address + '/disaster/alert'
    try:
        response_from_service = requests.post(url, json = data)
    except requests.exceptions.Timeout:
        response = jsonify({'message': service_name + ' timed out'})
        return response, 408
    response = jsonify(response_from_service.json())
    return response, response_from_service.status_code


@api.route('/disaster/alert', methods = ['PUT'])
@limiter.limit("10 per minute", error_message = 'Too Many Requests', override_defaults = False)
def update_alert():
    alert_id = request.args.get('alert_id')
    data = request.get_json()
    service_name = str(request.url_rule).split('/')[1]
    try:
        response_from_service_discovery = requests.get('http://localhost:8002/get_service?service_name=' + service_name,
                                                       timeout = 0.5)
    except requests.exceptions.Timeout:
        response = jsonify({'message': 'Service Discovery timed out'})
        return response, 408
    service_address = response_from_service_discovery.json()['service_address']
    url = service_address + '/disaster/alert?alert_id=' + alert_id
    try:
        response_from_service = requests.put(url, json = data)
    except requests.exceptions.Timeout:
        response = jsonify({'message': service_name + ' timed out'})
        return response, 408
    response = jsonify(response_from_service.json())
    return response, response_from_service.status_code


if __name__ == '__main__':
    api.run(host = 'localhost', port = 8003)

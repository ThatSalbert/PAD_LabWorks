from flask import Flask, request, jsonify
from flask_limiter import Limiter
from flask_limiter.util import get_remote_address
from cache.cache import caching_mechanism
import requests
import os

from circuit_breaker.circuit_breaker import circuit_breaker

REROUTE_LIMIT = int(os.getenv('REROUTE_LIMIT'))
TASK_TIMEOUT = float(os.getenv('TIMEOUT'))
SERVICEDISC_HOSTNAME = os.getenv('SERVICEDISC_HOSTNAME')
SERVICEDISC_PORT = os.getenv('SERVICEDISC_PORT')
THRESHOLD = float(os.getenv('FAILURE_THRESHOLD'))

api = Flask(__name__)
limiter = Limiter(get_remote_address, app = api, default_limits = ["200 per day", "50 per hour", "20 per minute"])

circuit_breaker = circuit_breaker(TASK_TIMEOUT, REROUTE_LIMIT, THRESHOLD)

def make_service_discovery_request(service_name):
    service_discovery_address = 'http://' + SERVICEDISC_HOSTNAME + ':' + SERVICEDISC_PORT + '/get_service?service_name=' + service_name
    try:
        response = requests.get(service_discovery_address, timeout = TASK_TIMEOUT)
    except requests.exceptions.Timeout:
        raise requests.exceptions.Timeout
    return response


def make_service_request(url):
    try:
        response = requests.get(url, timeout = TASK_TIMEOUT)
    except requests.exceptions.Timeout:
        raise requests.exceptions.Timeout
    return response


def make_service_request_with_data(url, data):
    try:
        response = requests.get(url, data = data, timeout = TASK_TIMEOUT)
    except requests.exceptions.Timeout:
        raise requests.exceptions.Timeout
    return response


@api.route('/weather/locations', methods = ['GET'])
@limiter.limit("20 per minute", error_message = 'Too Many Requests', override_defaults = False)
@caching_mechanism
@circuit_breaker
def get_locations():
    country = request.args.get('country')
    response_from_service_discovery = make_service_discovery_request('weather')
    service_address = response_from_service_discovery.json()['service_address']
    if country is None:
        url = service_address + '/weather/locations'
    else:
        url = service_address + '/weather/locations?country=' + country
    response_from_service = make_service_request(url)
    response = jsonify(response_from_service.json())
    return response, response_from_service.status_code


@api.route('/weather/current', methods = ['GET'])
@limiter.limit("20 per minute", error_message = 'Too Many Requests', override_defaults = False)
@caching_mechanism
@circuit_breaker
def get_current_weather():
    country = request.args.get('country')
    city = request.args.get('city')
    response_from_service_discovery = make_service_discovery_request('weather')
    service_address = response_from_service_discovery.json()['service_address']
    if country is None or city is None:
        url = service_address + '/weather/current'
    else:
        url = service_address + '/weather/current?country=' + country + '&city=' + city
    response_from_service = make_service_request(url)
    response = jsonify(response_from_service.json())
    return response, response_from_service.status_code


@api.route('/weather/forecast', methods = ['GET'])
@limiter.limit("20 per minute", error_message = 'Too Many Requests', override_defaults = False)
@caching_mechanism
@circuit_breaker
def get_weather_forecast():
    country = request.args.get('country')
    city = request.args.get('city')
    response_from_service_discovery = make_service_discovery_request('weather')
    service_address = response_from_service_discovery.json()['service_address']
    if country is None or city is None:
        url = service_address + '/weather/forecast'
    else:
        url = service_address + '/weather/forecast?country=' + country + '&city=' + city
    response_from_service = make_service_request(url)
    response = jsonify(response_from_service.json())
    return response, response_from_service.status_code


@api.route('/weather/add_data', methods = ['POST'])
@limiter.limit("20 per minute", error_message = 'Too Many Requests', override_defaults = False)
@circuit_breaker
def add_data():
    request_type = request.args.get('type')
    data = request.get_json()
    response_from_service_discovery = make_service_discovery_request('weather')
    service_address = response_from_service_discovery.json()['service_address']
    if request_type is None:
        url = service_address + '/weather/add_data'
    else:
        url = service_address + '/weather/add_data?type=' + request_type
    response_from_service = make_service_request_with_data(url, data)
    response = jsonify(response_from_service.json())
    return response, response_from_service.status_code


@api.route('/weather/update_data', methods = ['PUT'])
@limiter.limit("20 per minute", error_message = 'Too Many Requests', override_defaults = False)
@circuit_breaker
def update_data():
    request_type = request.args.get('type')
    data = request.get_json()
    response_from_service_discovery = make_service_discovery_request('weather')
    service_address = response_from_service_discovery.json()['service_address']
    if request_type is None:
        url = service_address + '/weather/update_data'
    else:
        url = service_address + '/weather/update_data?type=' + request_type
    response_from_service = make_service_request_with_data(url, data)
    response = jsonify(response_from_service.json())
    return response, response_from_service.status_code


@api.route('/disaster', methods = ['GET'])
@limiter.limit("20 per minute", error_message = 'Too Many Requests', override_defaults = False)
@caching_mechanism
@circuit_breaker
def get_disasters():
    response_from_service_discovery = make_service_discovery_request('disaster')
    service_address = response_from_service_discovery.json()['service_address']
    url = service_address + '/disaster'
    response_from_service = make_service_request(url)
    response = jsonify(response_from_service.json())
    return response, response_from_service.status_code


@api.route('/disaster/list', methods = ['GET'])
@limiter.limit("20 per minute", error_message = 'Too Many Requests', override_defaults = False)
@circuit_breaker
@caching_mechanism
def get_disaster_list():
    country = request.args.get('country')
    city = request.args.get('city')
    active = request.args.get('active')
    response_from_service_discovery = make_service_discovery_request('disaster')
    service_address = response_from_service_discovery.json()['service_address']
    if country is None or city is None or active is None:
        url = service_address + '/disaster/list'
    else:
        url = service_address + '/disaster/list?country=' + country + '&city=' + city + '&active=' + active
    response_from_service = make_service_request(url)
    response = jsonify(response_from_service.json())
    return response, response_from_service.status_code


@api.route('/disaster/alert', methods = ['POST'])
@limiter.limit("20 per minute", error_message = 'Too Many Requests', override_defaults = False)
@circuit_breaker
def add_alert():
    data = request.get_json()
    response_from_service_discovery = make_service_discovery_request('disaster')
    service_address = response_from_service_discovery.json()['service_address']
    url = service_address + '/disaster/alert'
    response_from_service = make_service_request_with_data(url, data)
    response = jsonify(response_from_service.json())
    return response, response_from_service.status_code


@api.route('/disaster/alert', methods = ['PUT'])
@limiter.limit("20 per minute", error_message = 'Too Many Requests', override_defaults = False)
@circuit_breaker
def update_alert():
    alert_id = request.args.get('alert_id')
    data = request.get_json()
    response_from_service_discovery = make_service_discovery_request('disaster')
    service_address = response_from_service_discovery.json()['service_address']
    url = service_address + '/disaster/alert?alert_id=' + alert_id
    response_from_service = make_service_request_with_data(url, data)
    response = jsonify(response_from_service.json())
    return response, response_from_service.status_code


@api.errorhandler(429)
def ratelimit_handler(e):
    return jsonify({'message': 'Too Many Requests'}), 429

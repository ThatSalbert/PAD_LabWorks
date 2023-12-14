import time
from flask import Flask, request, jsonify
from flask_limiter import Limiter
from flask_limiter.util import get_remote_address
from cache.cache import caching_mechanism
import requests
import os

REROUTE_LIMIT = int(os.getenv('REROUTE_LIMIT'))
FAILURE_LIMIT = int(os.getenv('FAILURE_LIMIT'))
TASK_TIMEOUT = float(os.getenv('TIMEOUT'))
SERVICEDISC_HOSTNAME = os.getenv('SERVICEDISC_HOSTNAME')
SERVICEDISC_PORT = os.getenv('SERVICEDISC_PORT')
THRESHOLD = float(os.getenv('FAILURE_THRESHOLD'))

api = Flask(__name__)
limiter = Limiter(get_remote_address, app = api, default_limits = ["200 per day", "50 per hour", "20 per minute"])

def circuit_breaker(task_timeout_limit, reroute_limit, failure_limit, threshold, service_name, params, data = None):
    for i in range(0, reroute_limit):
        response_from_service_discovery = make_service_discovery_request(service_name)
        service_address = response_from_service_discovery.json()['service_address']
        url = service_address + params
        api.logger.info("Reroute: " + str(i + 1) + " out of " + str(reroute_limit) + " for " + url)
        try:
            if data is None:
                response_from_service = make_service_request(url)
            else:
                response_from_service = make_service_request_with_data(url, data)
            response = jsonify(response_from_service.json())
            return response, response_from_service.status_code
        except (requests.exceptions.Timeout, requests.exceptions.ConnectionError):
            response, code = circuit_breaker_retry(task_timeout_limit, failure_limit, threshold, url, data)
            if response is not None:
                return response, code
            continue
    response = jsonify({'message': 'Service unavailable'})
    return response, 503
                
def circuit_breaker_retry(task_timeout_limit, failure_limit, threshold, url, data = None):
    start_time = time.time()
    timeouts = 1
    api.logger.info("Timeout: " + str(timeouts) + " out of " + str(failure_limit) + " for " + url)
    while True:
        try:
            if data is None:
                response_from_service = make_service_request(url)
            else:
                response_from_service = make_service_request_with_data(url, data)
            response = jsonify(response_from_service.json())
            return response, response_from_service.status_code
        except (requests.exceptions.Timeout, requests.exceptions.ConnectionError):
            timeouts += 1
            api.logger.info("Timeout: " + str(timeouts) + " out of " + str(failure_limit) + " for " + url)
            if timeouts >= failure_limit:
                return None, 503
            continue
                    

def make_service_discovery_request(service_name):
    service_discovery_address = 'http://' + SERVICEDISC_HOSTNAME + ':' + SERVICEDISC_PORT + '/get_service?service_name=' + service_name
    try:
        response = requests.get(service_discovery_address, timeout = TASK_TIMEOUT)
    except requests.exceptions.Timeout:
        raise requests.exceptions.Timeout
    return response

def make_service_request(url):
    api.logger.info("Trying to make request to: " + url)
    try:
        response = requests.get(url, timeout = TASK_TIMEOUT)
    except requests.exceptions.Timeout:
        raise requests.exceptions.Timeout
    return response

def make_service_request_with_data(url, data):
    api.logger.info("Trying to make request to: " + url)
    try:
        response = requests.get(url, data = data, timeout = TASK_TIMEOUT)
    except requests.exceptions.Timeout:
        raise requests.exceptions.Timeout
    return response

@api.route('/weather/locations', methods = ['GET'])
@limiter.limit("20 per minute", error_message = 'Too Many Requests', override_defaults = False)
@caching_mechanism
def get_locations():
    country = request.args.get('country')
    if country is None:
        params = '/weather/locations'
    else:
        params = '/weather/locations?country=' + country
    response, code = circuit_breaker(TASK_TIMEOUT, REROUTE_LIMIT, FAILURE_LIMIT, THRESHOLD, 'weather', params)
    return response, code

@api.route('/weather/current', methods = ['GET'])
@limiter.limit("20 per minute", error_message = 'Too Many Requests', override_defaults = False)
@caching_mechanism
def get_current_weather():
    country = request.args.get('country')
    city = request.args.get('city')
    if country is None or city is None:
        params = '/weather/current'
    else:
        params = '/weather/current?country=' + country + '&city=' + city
    response, code = circuit_breaker(TASK_TIMEOUT, REROUTE_LIMIT, FAILURE_LIMIT, THRESHOLD, 'weather', params)
    return response, code


@api.route('/weather/forecast', methods = ['GET'])
@limiter.limit("20 per minute", error_message = 'Too Many Requests', override_defaults = False)
@caching_mechanism
def get_weather_forecast():
    country = request.args.get('country')
    city = request.args.get('city')
    if country is None or city is None:
        params = '/weather/forecast'
    else:
        params = '/weather/forecast?country=' + country + '&city=' + city
    response, code = circuit_breaker(TASK_TIMEOUT, REROUTE_LIMIT, FAILURE_LIMIT, THRESHOLD, 'weather', params)
    return response, code


@api.route('/weather/add_data', methods = ['POST'])
@limiter.limit("20 per minute", error_message = 'Too Many Requests', override_defaults = False)
def add_data():
    request_type = request.args.get('type')
    data = request.get_json()
    if request_type is None:
        params = '/weather/add_data'
    else:
        params = '/weather/add_data?type=' + request_type
    response, code = circuit_breaker(TASK_TIMEOUT, REROUTE_LIMIT, FAILURE_LIMIT, THRESHOLD, 'weather', params, data)
    return response, code


@api.route('/weather/update_data', methods = ['PUT'])
@limiter.limit("20 per minute", error_message = 'Too Many Requests', override_defaults = False)
def update_data():
    request_type = request.args.get('type')
    data = request.get_json()
    if request_type is None:
        params = '/weather/update_data'
    else:
        params = '/weather/update_data?type=' + request_type
    response, code = circuit_breaker(TASK_TIMEOUT, REROUTE_LIMIT, FAILURE_LIMIT, THRESHOLD, 'weather', params, data)
    return response, code


@api.route('/disaster', methods = ['GET'])
@limiter.limit("20 per minute", error_message = 'Too Many Requests', override_defaults = False)
@caching_mechanism
def get_disasters():
    params = '/disaster'
    response, code = circuit_breaker(TASK_TIMEOUT, REROUTE_LIMIT, FAILURE_LIMIT, THRESHOLD, 'disaster', params)
    return response, code


@api.route('/disaster/list', methods = ['GET'])
@limiter.limit("20 per minute", error_message = 'Too Many Requests', override_defaults = False)
@caching_mechanism
def get_disaster_list():
    country = request.args.get('country')
    city = request.args.get('city')
    active = request.args.get('active')
    if country is None or city is None or active is None:
        params = '/disaster/list'
    else:
        params = '/disaster/list?country=' + country + '&city=' + city + '&active=' + active
    response, code = circuit_breaker(TASK_TIMEOUT, REROUTE_LIMIT, FAILURE_LIMIT, THRESHOLD, 'disaster', params)
    return response, code


@api.route('/disaster/alert', methods = ['POST'])
@limiter.limit("20 per minute", error_message = 'Too Many Requests', override_defaults = False)
def add_alert():
    data = request.get_json()
    params = '/disaster/alert'
    response, code = circuit_breaker(TASK_TIMEOUT, REROUTE_LIMIT, FAILURE_LIMIT, THRESHOLD, 'disaster', params, data)
    return response, code


@api.route('/disaster/alert', methods = ['PUT'])
@limiter.limit("20 per minute", error_message = 'Too Many Requests', override_defaults = False)
def update_alert():
    alert_id = request.args.get('alert_id')
    data = request.get_json()
    params = '/disaster/alert?alert_id=' + alert_id
    response, code = circuit_breaker(TASK_TIMEOUT, REROUTE_LIMIT, FAILURE_LIMIT, THRESHOLD, 'disaster', params, data)
    return response, code


@api.errorhandler(429)
def ratelimit_handler(e):
    return jsonify({'message': 'Too Many Requests'}), 429

from flask import Flask, request, jsonify
from flask_limiter import Limiter
from flask_limiter.util import get_remote_address
from cache.cache import caching_mechanism
import requests
import os

from circuit_breaker.circuit_breaker import circuit_breaker
from load_balancer.load_balancer import load_balancer

#TASK_TIMEOUT from string to float
TASK_TIMEOUT = float(os.getenv('TIMEOUT'))
SERVICEDISC_HOSTNAME = os.getenv('SERVICEDISC_HOSTNAME')
SERVICEDISC_PORT = os.getenv('SERVICEDISC_PORT')

api = Flask(__name__)
limiter = Limiter(get_remote_address, app = api, default_limits = ["200 per day", "50 per hour", "20 per minute"])


@api.route('/weather/locations', methods = ['GET'])
@limiter.limit("20 per minute", error_message = 'Too Many Requests', override_defaults = False)
@caching_mechanism
@circuit_breaker
def get_locations():
    country = request.args.get('country')
    service_name = load_balancer('weather')
    print('calling service:' + service_name + '\n')
    try:
        response_from_service_discovery = requests.get(
            'http://' + SERVICEDISC_HOSTNAME + ':' + SERVICEDISC_PORT + '/get_service?service_name=' + service_name,
            timeout = TASK_TIMEOUT)
    except requests.exceptions.Timeout:
        raise requests.exceptions.Timeout
    service_address = response_from_service_discovery.json()['service_address']
    if country is None:
        url = service_address + '/weather/locations'
    else:
        url = service_address + '/weather/locations?country=' + country
    try:
        response_from_service = requests.get(url, timeout = TASK_TIMEOUT)
    except requests.exceptions.Timeout:
        raise requests.exceptions.Timeout
    response = jsonify(response_from_service.json())
    return response, response_from_service.status_code


@api.route('/weather/current', methods = ['GET'])
@limiter.limit("20 per minute", error_message = 'Too Many Requests', override_defaults = False)
@caching_mechanism
@circuit_breaker
def get_current_weather():
    country = request.args.get('country')
    city = request.args.get('city')
    service_name = load_balancer('weather')
    print('calling service:' + service_name + '\n')
    try:
        response_from_service_discovery = requests.get(
            'http://' + SERVICEDISC_HOSTNAME + ':' + SERVICEDISC_PORT + '/get_service?service_name=' + service_name,
            timeout = TASK_TIMEOUT)
    except requests.exceptions.Timeout:
        raise requests.exceptions.Timeout
    service_address = response_from_service_discovery.json()['service_address']
    if country is None or city is None:
        url = service_address + '/weather/current'
    else:
        url = service_address + '/weather/current?country=' + country + '&city=' + city
    try:
        response_from_service = requests.get(url, timeout = TASK_TIMEOUT)
    except requests.exceptions.Timeout:
        raise requests.exceptions.Timeout
    response = jsonify(response_from_service.json())
    return response, response_from_service.status_code


@api.route('/weather/forecast', methods = ['GET'])
@limiter.limit("20 per minute", error_message = 'Too Many Requests', override_defaults = False)
@caching_mechanism
@circuit_breaker
def get_weather_forecast():
    country = request.args.get('country')
    city = request.args.get('city')
    service_name = load_balancer('weather')
    print('calling service:' + service_name + '\n')
    try:
        response_from_service_discovery = requests.get(
            'http://' + SERVICEDISC_HOSTNAME + ':' + SERVICEDISC_PORT + '/get_service?service_name=' + service_name,
            timeout = TASK_TIMEOUT)
    except requests.exceptions.Timeout:
        raise requests.exceptions.Timeout
    service_address = response_from_service_discovery.json()['service_address']
    if country is None or city is None:
        url = service_address + '/weather/forecast'
    else:
        url = service_address + '/weather/forecast?country=' + country + '&city=' + city
    try:
        response_from_service = requests.get(url, timeout = TASK_TIMEOUT)
    except requests.exceptions.Timeout:
        raise requests.exceptions.Timeout
    response = jsonify(response_from_service.json())
    return response, response_from_service.status_code


@api.route('/weather/add_data', methods = ['POST'])
@limiter.limit("20 per minute", error_message = 'Too Many Requests', override_defaults = False)
@circuit_breaker
def add_data():
    type = request.args.get('type')
    data = request.get_json()
    service_name = load_balancer('weather')
    print('calling service:' + service_name + '\n')
    try:
        response_from_service_discovery = requests.get(
            'http://' + SERVICEDISC_HOSTNAME + ':' + SERVICEDISC_PORT + '/get_service?service_name=' + service_name,
            timeout = TASK_TIMEOUT)
    except requests.exceptions.Timeout:
        raise requests.exceptions.Timeout
    service_address = response_from_service_discovery.json()['service_address']
    if type is None:
        url = service_address + '/weather/add_data'
    else:
        url = service_address + '/weather/add_data?type=' + type
    try:
        response_from_service = requests.post(url, json = data, timeout = TASK_TIMEOUT)
    except requests.exceptions.Timeout:
        raise requests.exceptions.Timeout
    response = jsonify(response_from_service.json())
    return response, response_from_service.status_code


@api.route('/weather/update_data', methods = ['PUT'])
@limiter.limit("20 per minute", error_message = 'Too Many Requests', override_defaults = False)
@circuit_breaker
def delete_data():
    type = request.args.get('type')
    data = request.get_json()
    service_name = load_balancer('weather')
    print('calling service:' + service_name + '\n')
    try:
        response_from_service_discovery = requests.get(
            'http://' + SERVICEDISC_HOSTNAME + ':' + SERVICEDISC_PORT + '/get_service?service_name=' + service_name,
            timeout = TASK_TIMEOUT)
    except requests.exceptions.Timeout:
        raise requests.exceptions.Timeout
    service_address = response_from_service_discovery.json()['service_address']
    if type is None:
        url = service_address + '/weather/update_data'
    else:
        url = service_address + '/weather/update_data?type=' + type
    try:
        response_from_service = requests.put(url, json = data, timeout = TASK_TIMEOUT)
    except requests.exceptions.Timeout:
        raise requests.exceptions.Timeout
    response = jsonify(response_from_service.json())
    return response, response_from_service.status_code


@api.route('/disaster', methods = ['GET'])
@limiter.limit("20 per minute", error_message = 'Too Many Requests', override_defaults = False)
@caching_mechanism
@circuit_breaker
def get_disasters():
    service_name = load_balancer('disaster')
    print('calling service:' + service_name + '\n')
    try:
        response_from_service_discovery = requests.get(
            'http://' + SERVICEDISC_HOSTNAME + ':' + SERVICEDISC_PORT + '/get_service?service_name=' + service_name,
            timeout = TASK_TIMEOUT)
    except requests.exceptions.Timeout:
        raise requests.exceptions.Timeout
    service_address = response_from_service_discovery.json()['service_address']
    url = service_address + '/disaster'
    try:
        response_from_service = requests.get(url, timeout = TASK_TIMEOUT)
    except requests.exceptions.Timeout:
        raise requests.exceptions.Timeout
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
    service_name = load_balancer('disaster')
    print('calling service:' + service_name + '\n')
    try:
        response_from_service_discovery = requests.get(
            'http://' + SERVICEDISC_HOSTNAME + ':' + SERVICEDISC_PORT + '/get_service?service_name=' + service_name,
            timeout = TASK_TIMEOUT)
    except requests.exceptions.Timeout:
        raise requests.exceptions.Timeout
    service_address = response_from_service_discovery.json()['service_address']
    if country is None or city is None or active is None:
        url = service_address + '/disaster/list'
    else:
        url = service_address + '/disaster/list?country=' + country + '&city=' + city + '&active=' + active
    try:
        response_from_service = requests.get(url, timeout = TASK_TIMEOUT)
    except requests.exceptions.Timeout:
        raise requests.exceptions.Timeout
    response = jsonify(response_from_service.json())
    return response, response_from_service.status_code


@api.route('/disaster/alert', methods = ['POST'])
@limiter.limit("20 per minute", error_message = 'Too Many Requests', override_defaults = False)
@circuit_breaker
def add_alert():
    data = request.get_json()
    service_name = load_balancer('disaster')
    print('calling service:' + service_name + '\n')
    try:
        response_from_service_discovery = requests.get(
            'http://' + SERVICEDISC_HOSTNAME + ':' + SERVICEDISC_PORT + '/get_service?service_name=' + service_name,
            timeout = TASK_TIMEOUT)
    except requests.exceptions.Timeout:
        raise requests.exceptions.Timeout
    service_address = response_from_service_discovery.json()['service_address']
    url = service_address + '/disaster/alert'
    try:
        response_from_service = requests.post(url, json = data, timeout = TASK_TIMEOUT)
    except requests.exceptions.Timeout:
        raise requests.exceptions.Timeout
    response = jsonify(response_from_service.json())
    return response, response_from_service.status_code


@api.route('/disaster/alert', methods = ['PUT'])
@limiter.limit("20 per minute", error_message = 'Too Many Requests', override_defaults = False)
@circuit_breaker
def update_alert():
    alert_id = request.args.get('alert_id')
    data = request.get_json()
    service_name = load_balancer('disaster')
    print('calling service:' + service_name + '\n')
    try:
        response_from_service_discovery = requests.get(
            'http://' + SERVICEDISC_HOSTNAME + ':' + SERVICEDISC_PORT + '/get_service?service_name=' + service_name,
            timeout = TASK_TIMEOUT)
    except requests.exceptions.Timeout:
        raise requests.exceptions.Timeout
    service_address = response_from_service_discovery.json()['service_address']
    url = service_address + '/disaster/alert?alert_id=' + alert_id
    try:
        response_from_service = requests.put(url, json = data, timeout = TASK_TIMEOUT)
    except requests.exceptions.Timeout:
        raise requests.exceptions.Timeout
    response = jsonify(response_from_service.json())
    return response, response_from_service.status_code


@api.errorhandler(429)
def ratelimit_handler(e):
    return jsonify({'message': 'Too Many Requests'}), 429

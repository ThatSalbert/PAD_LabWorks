import os
from flask import Flask, request, jsonify
import requests

api = Flask(__name__)

SERVICEDISC_HOSTNAME = os.getenv('SERVICEDISC_HOSTNAME')
SERVICEDISC_PORT = os.getenv('SERVICEDISC_PORT')
TASK_TIMEOUT = float(os.getenv('TIMEOUT'))

def make_service_discovery_request(service_name):
    service_discovery_address = 'http://' + SERVICEDISC_HOSTNAME + ':' + SERVICEDISC_PORT + '/get_service?service_name=' + service_name
    try:
        response = requests.get(service_discovery_address, timeout = TASK_TIMEOUT)
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

@api.route('/prepare', methods = ['POST'])
def prepare():
    request_data = request.get_json()['data']
    params = request.get_json()['params']
    weather_service = make_service_discovery_request('weather')
    disaster_service = make_service_discovery_request('disaster')
    weather_prepare_url = weather_service.json()['service_address'] + params + "/prepare"
    disaster_prepare_url = disaster_service.json()['service_address'] + params + "/prepare"
    api.logger.info("Trying to make request to: " + weather_prepare_url)
    weather_prepare_response = make_service_request_with_data(weather_prepare_url, request_data)
    api.logger.info("Trying to make request to: " + disaster_prepare_url)
    disaster_prepare_response = make_service_request_with_data(disaster_prepare_url, request_data)
    api.logger.info(str(weather_prepare_response.status_code) + ", " + str(disaster_prepare_response.status_code))
    if weather_prepare_response.status_code == 200 and disaster_prepare_response.status_code == 200:
        weather_prepare_url = weather_service.json()['service_address'] + params + "/commit"
        disaster_prepare_url = weather_service.json()['service_address'] + params + "/commit"
        weather_payload_response = weather_prepare_response.json()['payload']
        disaster_payload_response = disaster_prepare_response.json()['payload']
        api.logger.info("Trying to make request to: " + weather_prepare_url)
        weather_prepare_response = make_service_request_with_data(weather_prepare_url, weather_payload_response)
        api.logger.info("Trying to make request to: " + disaster_prepare_url)
        disaster_prepare_response = make_service_request_with_data(disaster_prepare_url, disaster_payload_response)
        api.logger.info(str(weather_prepare_response.status_code) + ", " + str(disaster_prepare_response.status_code))
        if weather_prepare_response.status_code == 200 and disaster_prepare_response.status_code == 200:
            return weather_prepare_response.json(), 200
        else:
            response = jsonify({'message': 'Could not commit'})
            return response, 409
    response = jsonify({'message': 'Could not prepare'})
    return response, 409
    
        
    
        
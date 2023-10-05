from flask import Flask, request, jsonify
import requests

api = Flask(__name__)


@api.route('/weather/locations', methods=['GET'])
def get_locations():
    country = request.args.get('country')
    service_name = str(request.url_rule).split('/')[1]
    try:
        response_from_service_discovery = requests.get('http://localhost:8002/get_service?service_name=' + service_name, timeout=0.05)
    except requests.exceptions.Timeout:
        response = jsonify({'message': 'Service Discovery timed out'})
        return response, 408
    service_address = response_from_service_discovery.json()['service_address']
    url = service_address + '/weather/locations?country=' + country
    try:
        response_from_service = requests.get(url, timeout=0.05)
    except requests.exceptions.Timeout:
        response = jsonify({'message': service_name + ' timed out'})
        return response, 408
    response = jsonify(response_from_service.json())
    return response, response_from_service.status_code


@api.route('/weather/current', methods=['GET'])
def get_current_weather():
    city = request.args.get('city')
    service_name = str(request.url_rule).split('/')[1]
    try:
        response_from_service_discovery = requests.get('http://localhost:8002/get_service?service_name=' + service_name, timeout=0.05)
    except requests.exceptions.Timeout:
        response = jsonify({'message': 'Service Discovery timed out'})
        return response, 408
    service_address = response_from_service_discovery.json()['service_address']
    url = service_address + '/weather/current?city=' + city
    try:
        response_from_service = requests.get(url, timeout=0.05)
    except requests.exceptions.Timeout:
        response = jsonify({'message': service_name + ' timed out'})
        return response, 408
    response = jsonify(response_from_service.json())
    return response, response_from_service.status_code


@api.route('/weather/forecast', methods=['GET'])
def get_weather_forecast():
    city = request.args.get('city')
    service_name = str(request.url_rule).split('/')[1]
    try:
        response_from_service_discovery = requests.get('http://localhost:8002/get_service?service_name=' + service_name, timeout=0.05)
    except requests.exceptions.Timeout:
        response = jsonify({'message': 'Service Discovery timed out'})
        return response, 408
    service_address = response_from_service_discovery.json()['service_address']
    url = service_address + '/weather/forecast?city=' + city
    try:
        response_from_service = requests.get(url, timeout=0.05)
    except requests.exceptions.Timeout:
        response = jsonify({'message': service_name + ' timed out'})
        return response, 408
    response = jsonify(response_from_service.json())
    return response, response_from_service.status_code


if __name__ == '__main__':
    api.run(host='localhost', port=8003)

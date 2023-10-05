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
        response_from_service = requests.get(url)
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
        response_from_service = requests.get(url)
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
        response_from_service = requests.get(url)
    except requests.exceptions.Timeout:
        response = jsonify({'message': service_name + ' timed out'})
        return response, 408
    response = jsonify(response_from_service.json())
    return response, response_from_service.status_code


@api.route('/weather/add_data', methods=['POST'])
def add_data():
    add_type = request.args.get('type')
    data = request.get_json()
    service_name = str(request.url_rule).split('/')[1]
    try:
        response_from_service_discovery = requests.get('http://localhost:8002/get_service?service_name=' + service_name, timeout=0.05)
    except requests.exceptions.Timeout:
        response = jsonify({'message': 'Service Discovery timed out'})
        return response, 408
    service_address = response_from_service_discovery.json()['service_address']
    url = service_address + '/weather/add_data?type=' + add_type
    try:
        response_from_service = requests.post(url, json=data)
    except requests.exceptions.Timeout:
        response = jsonify({'message': service_name + ' timed out'})
        return response, 408
    response = jsonify(response_from_service.json())
    return response, response_from_service.status_code


@api.route('/weather/delete_data', methods=['DELETE'])
def delete_data():
    delete_type = request.args.get('type')
    data = request.get_json()
    service_name = str(request.url_rule).split('/')[1]
    try:
        response_from_service_discovery = requests.get('http://localhost:8002/get_service?service_name=' + service_name, timeout=0.05)
    except requests.exceptions.Timeout:
        response = jsonify({'message': 'Service Discovery timed out'})
        return response, 408
    service_address = response_from_service_discovery.json()['service_address']
    url = service_address + '/weather/delete_data?type=' + delete_type
    try:
        response_from_service = requests.delete(url, json=data)
    except requests.exceptions.Timeout:
        response = jsonify({'message': service_name + ' timed out'})
        return response, 408
    response = jsonify(response_from_service.json())
    return response, response_from_service.status_code


if __name__ == '__main__':
    api.run(host='localhost', port=8003)

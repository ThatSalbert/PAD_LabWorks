from flask import Flask, request, jsonify
import requests
import service_discovery.service_discovery as sd

api = Flask(__name__)


@api.route('/register', methods=['POST'])
def register_service():
    data = request.get_json()
    service_name = data['service_name']
    service_url = data['service_address']
    if service_name not in sd.get_services():
        sd.add_service(service_name, service_url)
        response = jsonify({'message': 'Service registered successfully'})
        return response, 200
    else:
        response = jsonify({'message': 'Service already registered'})
        return response, 409


@api.route('/get_services', methods=['GET'])
def get_services():
    services = sd.get_services()
    if len(services) == 0:
        return jsonify({'message': 'No services registered'}), 209
    else:
        return jsonify(services), 200


@api.route('/weather/locations', methods=['GET'])
def get_locations():
    country = request.args.get('country')
    url = sd.get_service('weather') + '/weather/locations?country=' + country
    response_from_server = requests.get(url)
    response = jsonify(response_from_server.json())
    return response, response_from_server.status_code


@api.route('/weather/current', methods=['GET'])
def get_current_weather():
    city = request.args.get('city')
    url = sd.get_service('weather') + '/weather/current?city=' + city
    response_from_server = requests.get(url)
    response = jsonify(response_from_server.json())
    return response, response_from_server.status_code


if __name__ == '__main__':
    sd.load_services()
    api.run(host='localhost', port=8003)

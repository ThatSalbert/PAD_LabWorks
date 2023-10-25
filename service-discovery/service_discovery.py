from flask import Flask, request, jsonify

api = Flask(__name__)
services = dict()


@api.route('/register', methods = ['POST'])
def add_service():
    data = request.get_json()
    service_name = data['service_name']
    service_url = data['service_address']
    if service_name not in services:
        services[service_name] = service_url
        response = jsonify({'message': 'Service registered successfully'})
        return response, 200
    else:
        response = jsonify({'message': 'Service already registered'})
        return response, 409


@api.route('/get_service', methods = ['GET'])
def get_service():
    service_name = request.args.get('service_name')
    if service_name in services:
        response = jsonify({'service_name': service_name, 'service_address': services[service_name]})
        return response, 200
    else:
        response = jsonify({'message': 'Service not registered'})
        return response, 404
    

@api.route('/get_all_services', methods = ['GET'])
def get_all_services():
    response = jsonify(services)
    return response, 200
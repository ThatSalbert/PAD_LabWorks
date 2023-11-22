from flask import Flask, request, jsonify

api = Flask(__name__)
services = dict()


def next_service(service_name):
    global services
    if service_name in services:
        if services[service_name]['counter'] < len(services[service_name]['services']):
            services[service_name]['counter'] += 1
            return services[service_name]['services'][services[service_name]['counter']]
        else:
            services[service_name]['counter'] = 1
            return services[service_name]['services'][services[service_name]['counter']]
    else:
        return None


@api.route('/register', methods = ['POST'])
def add_service():
    data = request.get_json()
    service_type = data['service_type']
    service_name = data['service_name']
    service_url = data['service_address']
    if service_type not in services:
        services[service_type] = dict()
        services[service_type]['counter'] = 0
        services[service_type]['services'] = []
        service_info = {
            'service_name': service_name,
            'service_address': service_url
        }
        services[service_type]['services'].append(service_info)
        response = jsonify({'message': 'Service registered'})
        return response, 200
    else:
        service_info = {
            'service_name': service_name,
            'service_address': service_url
        }
        if service_info not in services[service_type]['services']:
            services[service_type]['services'].append(service_info)
            response = jsonify({'message': 'Service registered'})
            return response, 200
        else:
            response = jsonify({'message': 'Service already registered'})
            return response, 409


@api.route('/get_service', methods = ['GET'])
def get_service():
    service_type = request.args.get('service_type')
    service = next_service(service_type)
    if service is not None:
        response = jsonify(service)
        return response, 200
    else:
        response = jsonify({'message': 'Service not found'})
        return response, 404


@api.route('/get_all_services', methods = ['GET'])
def get_all_services():
    response = jsonify(services)
    return response, 200

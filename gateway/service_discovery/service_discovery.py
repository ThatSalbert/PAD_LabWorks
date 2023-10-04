import json
import os

services = {}


def add_service(service_name, service_url):
    if service_name not in services:
        services[service_name] = service_url
        save_services()
        return True
    else:
        return False


def update_service(service_name, service_url):
    if service_name in services:
        services[service_name] = service_url
        save_services()
        return True
    else:
        return False


def remove_service(service_name):
    if service_name in services:
        del services[service_name]
        save_services()
        return True
    else:
        return False


def get_service(service_name):
    if service_name in services:
        return services[service_name]
    else:
        return None


def get_services():
    return services


def save_services():
    folder = os.path.dirname(os.path.realpath(__file__))
    with open(folder + '/services.json', 'w') as file:
        file.write(json.dumps(services))


def load_services():
    folder = os.path.dirname(os.path.realpath(__file__))
    global services
    with open(folder + '/services.json', 'r') as file:
        services = json.loads(file.read())

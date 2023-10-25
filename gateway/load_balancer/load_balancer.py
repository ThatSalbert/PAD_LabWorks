weather_minimum_turn = 1
disaster_minimum_turn = 1
weather_current_turn = 1
disaster_current_turn = 1
weather_maximum_turn = 3
disaster_maximum_turn = 3


def load_balancer(service_name):
    global weather_minimum_turn
    global disaster_minimum_turn
    global weather_current_turn
    global disaster_current_turn
    global weather_maximum_turn
    global disaster_maximum_turn
    if service_name == 'weather':
        if weather_current_turn < weather_maximum_turn:
            weather_current_turn += 1
            return 'weather-service-' + str(weather_current_turn)
        else:
            weather_current_turn = weather_minimum_turn
            return 'weather-service-' + str(weather_current_turn)
    elif service_name == 'disaster':
        if disaster_current_turn < disaster_maximum_turn:
            disaster_current_turn += 1
            return 'disaster-service-' + str(disaster_current_turn)
        else:
            disaster_current_turn = disaster_minimum_turn
            return 'disaster-service-' + str(disaster_current_turn)

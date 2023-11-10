class Config:
    class UnicornType:
        STELLAR = 0
        GALACTIC = 1
        COSMIC = 2

    def __init__(self, unicorn_type, wifi_ssid, wifi_psk, wifi_country, listener_port=5000):
        self.unicorn_type = unicorn_type
        self.wifi_ssid = wifi_ssid
        self.wifi_psk = wifi_psk
        self.wifi_country = wifi_country
        self.listener_port = listener_port

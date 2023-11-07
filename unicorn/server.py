import network
import socket
import time

from galactic import GalacticUnicorn


from config import *
from client import *

wlan = network.WLAN(network.STA_IF)
wlan.active(True)
wlan.connect(WIFI_SSID, WIFI_PSK)

max_wait = 10
while max_wait > 0:
    if wlan.status() < 0 or wlan.status() >= 3:
        break
    max_wait -= 1
    print('waiting for connection...')
    time.sleep(1)

if wlan.status() != 3:
    raise RuntimeError('network connection failed')
else:
    print('connected')
    status = wlan.ifconfig()
    print('ip = ' + status[0])

addr = socket.getaddrinfo('0.0.0.0', LISTEN_PORT)[0][-1]

s = socket.socket()
s.bind(addr)
s.listen(1)

print('listening on', addr)

gu = GalacticUnicorn()

width = GalacticUnicorn.WIDTH
height = GalacticUnicorn.HEIGHT

while True:
    try:
        cl, addr = s.accept()

        print('client connected from', addr)
        RUClient(cl, gu).handle()
        cl.close()

    except OSError as e:
        cl.close()
        print('connection closed')

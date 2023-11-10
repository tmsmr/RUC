import machine

from time import sleep

from rulib import *
from config import c

display = Display(c)
info(str(display) + ' initialized')

try:
    ip = Wifi(c).connect()
except RuntimeError as e:
    display.error()
    error(e)
    sleep(60)
    machine.reset()
info('connected to ' + c.wifi_ssid + ', got IP ' + ip)
display.write(ip.split('.')[-1])

server = Server(c)
info('listening for clients on %s:%d' % (server.addr[0], server.addr[1]))

while True:
    sock = None
    try:
        sock, addr = server.accept()
        info('client %s:%d connected' % (addr[0], addr[1]))
        display.clear()
        warn(Client.handle(sock, display.unicorn))
        warn('client %s:%d disconnected' % (addr[0], addr[1]))
        display.clear()
    finally:
        sock.close()

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
    sock, addr = server.accept()
    info('client %s:%d connected' % (addr[0], addr[1]))
    display.clear()
    try:
        Client.handle(sock, display.unicorn)
    except RuntimeError as e:
        warn(e)
    finally:
        sock.close()
    warn('client %s:%d disconnected' % (addr[0], addr[1]))
    display.clear()

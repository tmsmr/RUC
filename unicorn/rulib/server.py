import socket

BIND_ADDR = '0.0.0.0'

class Server:
    def __init__(self, config):
        self.addr = socket.getaddrinfo(BIND_ADDR, config.listener_port)[0][-1]
        self.sock = socket.socket()
        self.sock.bind(self.addr)
        self.sock.listen(1)

    def accept(self):
        return self.sock.accept()

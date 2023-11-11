from .messages import Header, HeaderEncodedLen, MessageType

RECV_CHUNK_SIZE = 256


class Client:
    @staticmethod
    def recv(sock, count):
        buf = bytearray()
        while True:
            miss = count - len(buf)
            if miss == 0:
                return buf
            raw = sock.recv(RECV_CHUNK_SIZE if miss > RECV_CHUNK_SIZE else miss)
            if len(raw) == 0:
                raise RuntimeError('client socket closed')
            buf.extend(raw)

    @staticmethod
    def handle(sock, unicorn):
        while True:
            header = Header(Client.recv(sock, HeaderEncodedLen))
            data = Client.recv(sock, header.data_len)
            if header.type == MessageType.SET_PIXEL:
                for i in range(0, header.data_len, 5):
                    unicorn.set_pixel(data[i], data[i + 1], data[i + 2], data[i + 3], data[i + 4])

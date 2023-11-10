from .messages import Header, HeaderEncodedLen, MessageType

RECV_CHUNK_SIZE = 1


class Client:
    @staticmethod
    def handle(sock, unicorn):
        while True:
            raw = sock.recv(HeaderEncodedLen)
            if len(raw) != HeaderEncodedLen:
                return 'failed to receive header'
            header = Header(raw)
            data = bytearray()
            for _ in range(int(header.data_len / RECV_CHUNK_SIZE)):
                data.extend(sock.recv(RECV_CHUNK_SIZE))
            if header.data_len % RECV_CHUNK_SIZE != 0:
                data.extend(sock.recv(header.data_len % RECV_CHUNK_SIZE))
            if len(data) != header.data_len:
                return 'failed to receive data'
            if header.type == MessageType.SET_PIXEL:
                for i in range(0, header.data_len, 5):
                    unicorn.set_pixel(data[i], data[i + 1], data[i + 2], data[i + 3], data[i + 4])


#             raw = bytearray()
#             for i in range(HeaderEncodedLen):
#                 next_byte = self.cl.recv(1)
#                 if len(next_byte) == 0:
#                     return
#                 raw.append(next_byte[0])
#             header = Header(raw)
#             data = bytearray()
#             for i in range(header.data_len):
#                 next_byte = self.cl.recv(1)
#                 if len(next_byte) == 0:
#                     return
#                 data.append(next_byte[0])

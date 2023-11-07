from messages import *

class RUClient:
    def __init__(self, cl, gu):
        self.cl = cl
        self.gu = gu

    def handle(self):
        while True:
            raw = bytearray()
            for i in range(HeaderEncodedLen):
                next_byte = self.cl.recv(1)
                if len(next_byte) == 0:
                    return
                raw.append(next_byte[0])
            header = Header()
            if not header.parse(raw):
                print("failed to parse header")
                continue
            data = bytearray()
            for i in range(header.data_len):
                next_byte = self.cl.recv(1)
                if len(next_byte) == 0:
                    return
                data.append(next_byte[0])
            if header.type == SetPixelMessage.TYPE:
                msg = SetPixelMessage()
                if not msg.draw(header, data, self.gu):
                    continue



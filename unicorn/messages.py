HeaderEncodedLen = 6


class Header:
    def __init__(self):
        self.data_len = None
        self.type = None
        self.source_id = None

    def parse(self, raw):
        if len(raw) != HeaderEncodedLen:
            return False
        self.source_id = raw[0]
        self.type = raw[2]
        self.data_len = int.from_bytes(raw[4:6], 'big', False)
        return True

    def __str__(self):
        return "Header(type: %d, data_len: %d)" % (self.type, self.data_len)


class Pixel:
    def __init__(self, data):
        self.x = data[0]
        self.y = data[1]
        self.r = data[2]
        self.g = data[3]
        self.b = data[4]


class SetPixelMessage:
    TYPE = 0

    def __init__(self):
        self.pixels = []

    def draw(self, header, data, gu):
        for i in range(0, header.data_len, 5):
            gu.set_pixel(data[i], data[i+1], data[i+2], data[i+3], data[i+4])
        return True

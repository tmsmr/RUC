HeaderEncodedLen = 6


class Header:
    def __init__(self, raw):
        self.source_id = raw[0]
        self.type = raw[2]
        self.data_len = int.from_bytes(raw[4:6], 'big', False)


class MessageType:
    SET_PIXEL = 0

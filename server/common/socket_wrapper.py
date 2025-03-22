
"""
Sockets protocol:
    - 2 bytes que indican el tamanio del mensaje
    - mensaje

    Ejemplo:
        Cuando el server confirma la apuesta, le envia un "OK" al cliente
    escribiendo:
    02OK
"""

MSG_LEN_SIZE = 2


def read_msg(socket):
    # Primero trato de leer los primeros 2 bytes que me dicen el tamanio
    # del mensaje
    msg_len_b = b""
    while len(msg_len_b) < MSG_LEN_SIZE:
        bytes_leidos = socket.recv(MSG_LEN_SIZE - len(msg_len_b))
        if not bytes_leidos:
            return None
        msg_len_b += bytes_leidos

    msg_len = int.from_bytes(msg_len_b, byteorder="big")
    # Ahora leo el mensaje completo
    msg = b""
    while len(msg) < msg_len:
        bytes_leidos = socket.recv(msg_len - len(msg))
        if not bytes_leidos:
            return None
        msg += bytes_leidos
    return msg.decode('utf-8')


def write_msg(socket, msg):
    # Formateo el mensaje
    encoded_msg = msg.encode('utf-8') # con utf8 1 caracter = 1 byte
    msg_length = len(encoded_msg).to_bytes(MSG_LEN_SIZE, byteorder="big")
    full_msg = msg_length + encoded_msg
    # Lo mando envitand short-writes
    bytes_totales = 0
    while bytes_totales < len(full_msg):
        bytes_enviados = socket.send(full_msg[bytes_totales:])
        if bytes_enviados == 0:
            raise RuntimeError("Socket broken")
        bytes_totales += bytes_enviados
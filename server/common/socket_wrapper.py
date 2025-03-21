
# TODO: Modify the receive to avoid short-reads
def read_msg(socket):
    return socket.recv(1024).rstrip().decode('utf-8')

# TODO: Modify the send to avoid short-writes
def write_msg(socket, msg):
    socket.send("{}\n".format(msg).encode('utf-8'))
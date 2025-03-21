import datetime
import os
import socket
import logging

import signal
import sys

from common.utils import Bet
from common import utils

class Server:
    def __init__(self, port, listen_backlog):
        # Initialize server socket
        self._server_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        self._server_socket.bind(('', port))
        self._server_socket.listen(listen_backlog)

        self.clients = []
        self.seguir_conectando = True

        signal.signal(signal.SIGTERM, self.signal_handler)

    def run(self):
        """
        Dummy Server loop

        Server that accept a new connections and establishes a
        communication with a client. After client with communucation
        finishes, servers starts to accept new connections again
        """

        # TODO: Modify this program to handle signal to graceful shutdown
        # the server
        while self.seguir_conectando:
            client_sock = self.__accept_new_connection()
            if not self.handshake(client_sock):
                print("no anduvo")
                continue
            self.clients.append(client_sock)
            self.__handle_client_connection(client_sock)

    def handshake(self, client_sock):
        """
        Realiza el handshake con el cliente
        """
        try:
            client_msg = client_sock.recv(1024).rstrip().decode('utf-8') + "\n"
            addr = client_sock.getpeername()
            client_expected_msg = os.getenv('HANDSHAKE_REQUEST_MESSAGE') + "\n"
            # logging.info(f"action: handshake | result: success | ip: {addr[0]} | recieved: {client_msg}")
            if client_msg == client_expected_msg:
                response = os.getenv('HANDSHAKE_RESPONSE_MESSAGE') + "\n"
                # logging.info(f"action: handshake | result: success |  ip: {addr[0]} | responding: {response}")
                client_sock.send(response.encode('utf-8'))
                return True
            else:
                client_sock.send("Bye Client\n".encode('utf-8'))
                client_sock.close()
                return False
        except OSError as e:
            # logging.error("action: handshake | result: fail | error"
            # ": {e}")
            client_sock.send("Bye Client\n".encode('utf-8'))
            client_sock.close()
            return False
        

    def __handle_client_connection(self, client_sock):
        """
        Read message from a specific client socket and closes the socket

        If a problem arises in the communication with the client, the
        client socket will also be closed
        """
        try:
            # TODO: Modify the receive to avoid short-reads
            msg = client_sock.recv(1024).rstrip().decode('utf-8')
            addr = client_sock.getpeername()
            logging.info(f'action: receive_message | result: success | ip: {addr[0]} | msg: {msg}')
            # TODO: Modify the send to avoid short-writes
            # client_sock.send("{}\n".format(msg).encode('utf-8'))
            print("Mensaje recibido: ", msg)
            action = msg.split(":")[0]
            if action == "BET":
                bet_data = msg.split(":")[1].split(",")
                name = bet_data[0]
                surname = bet_data[1]
                document = bet_data[2]
                birthdate = bet_data[3]
                number = int(bet_data[4])
                bet = Bet(agency=1, first_name=name, last_name=surname, document=document, birthdate=birthdate, number=number)
                print(bet.first_name)
                utils.store_bets([bet])
                print("Bet: ", bet)
                logging.info(f'action: apuesta_almacenada | result: success | dni: {document} | numero: {number}.')
                client_sock.send("bet stored\n".encode('utf-8'))
        except OSError as e:
            logging.error("action: receive_message | result: fail | error: {e}")
        finally:
            pass
            # client_sock.close()

    def __accept_new_connection(self):
        """
        Accept new connections

        Function blocks until a connection to a client is made.
        Then connection created is printed and returned
        """

        # Connection arrived
        logging.info('action: accept_connections | result: in_progress')
        c, addr = self._server_socket.accept()
        logging.info(f'action: accept_connections | result: success | ip: {addr[0]}')
        return c

    def signal_handler(self, sig, frame):
        print(f"Señal recibida: {sig}")
        self.seguir_conectando = False # Por las dudas de justo se cierre cuando se recibe una conexion
        print("Cerrando socket")
        for client in self.clients:
            print(f"Cerrando cliente {client}")
            client.close()
        self._server_socket.close()
        sys.exit(0)
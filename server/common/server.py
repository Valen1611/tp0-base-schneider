import datetime
import os
import socket
import logging

import signal
import sys

from common.utils import Bet
from common import utils
from common import socket_wrapper
from common import bet_protocol

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
            self.clients.append(client_sock)
            self.__handle_client_connection(client_sock)

    def __handle_client_connection(self, client_sock):
        """
        Read message from a specific client socket and store the bet

        If a problem arises in the communication with the client, the
        client socket will also be closed
        """
        try:            
            # Espero a recibir mensaje de cliente
            msg = socket_wrapper.read_msg(client_sock)
            addr = client_sock.getpeername()
            logging.info(f'action: receive_message | result: success | ip: {addr[0]} | msg: {msg}')
            # Me fijo que accion quiere hacer
            action = bet_protocol.get_action(msg)
            if action == "BET":
                # Leo la data de la apuesta
                agency, name, surname, document, birthdate, number = bet_protocol.read_bet_msg(msg)                
                bet = Bet(agency=agency, first_name=name, last_name=surname, document=document, birthdate=birthdate, number=number)
                # Guardo la apuesta
                utils.store_bets([bet])            
                logging.info(f'action: apuesta_almacenada | result: success | dni: {document} | numero: {number}.')
                # Le confirmo al cliente que se guardo la apuesta
                socket_wrapper.write_msg(client_sock, "OK")
        except OSError as e:
            logging.error("action: receive_message | result: fail | error: {e}")
        finally:            
            client_sock.close()

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
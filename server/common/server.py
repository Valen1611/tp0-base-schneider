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

import multiprocessing

TOTAL_CLIENTS = os.getenv("CLIENTS_AMOUNT", 1)
class Server:
    def __init__(self, port, listen_backlog):
        # Initialize server socket
        self._server_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        self._server_socket.bind(('', port))
        self._server_socket.listen(listen_backlog)

        self.clients = {}
        self.seguir_conectando = True

        self.manager = multiprocessing.Manager()
        self.clients_ids = self.manager.dict()
        self.clients_waiting = self.manager.Value('i', 0)
        self.lock = self.manager.Lock()
        self.bets_file_lock = self.manager.Lock()
        self.procesos = []

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
            self.clients[client_sock] = "TALKING"
            proceso = multiprocessing.Process(target=self.__handle_client_connection, args=(client_sock,))
            proceso.start()
            self.procesos.append(proceso)
            
            
    def hacer_sorteo(self):
        logging.info("action: sorteo | result: success")
        bets = utils.load_bets()
        winners = {agency : [] for agency in self.clients_ids.keys()}
        for bet in bets:
            if utils.has_won(bet):
                logging.info(f"action: apuesta_ganadora | result: success | dni: {bet.document} | numero: {bet.number}")
                winners[bet.agency].append(bet)
        self.notify_winners(winners)

    def notify_winners(self, winners):
        for agency, bets in winners.items():
            client_sock = self.clients_ids[agency]
            winners_dnis = [str(bet.document) for bet in bets]
            socket_wrapper.write_msg(client_sock, f"WINNERS:{','.join(winners_dnis)}")

        # client_sock = self.clients_ids[bet.agency]
        # socket_wrapper.write_msg(client_sock, f"WIN:{bet.number}")            

    def __handle_client_connection(self, client_sock):
        """
        Read message from a specific client socket and store the bet

        If a problem arises in the communication with the client, the
        client socket will also be closed
        """
        try:            
            while True:
                # Espero a recibir mensaje de cliente
                msg = socket_wrapper.read_msg(client_sock)
                addr = client_sock.getpeername()
                logging.info(f'action: receive_message | result: success | ip: {addr[0]} | msg: bets')
                # Me fijo que accion quiere hacer
                action = bet_protocol.get_action(msg)
                if action == "BET":
                    # Leo la data de la apuesta
                    agency, name, surname, document, birthdate, number = bet_protocol.read_bet_msg(msg)                
                    bet = Bet(agency=agency, first_name=name, last_name=surname, document=document, birthdate=birthdate, number=number)            
                    # Guardo la apuesta
                    with self.bets_file_lock:
                        utils.store_bets([bet])            
                    logging.info(f'action: apuesta_almacenada | result: success | dni: {document} | numero: {number}.')
                    # Le confirmo al cliente que se guardo la apuesta
                    socket_wrapper.write_msg(client_sock, "OK")
                    continue
                elif action == "BATCH_BET":
                    try:
                        # Leo la data de las apuestas
                        bets = []
                        for bet in msg.split(":")[1].split(";"):
                            if not bet: 
                                continue                      
                            bet_data = bet.split(",")
                            agency = int(bet_data[0])
                            name = bet_data[1]
                            surname = bet_data[2]
                            document = int(bet_data[3])
                            birthdate = bet_data[4]
                            number = int(bet_data[5])

                            bets.append(Bet(agency=agency, first_name=name, last_name=surname, document=document, birthdate=birthdate, number=int(number)))
                        # Guardo las apuestas
                        utils.store_bets(bets)
                        self.clients_ids[bets[0].agency] = client_sock
                        logging.info(f'action: apuesta_recibida | result: success | cantidad: {len(bets)}')
                        # Le confirmo al cliente que se guardaron las apuestas
                        socket_wrapper.write_msg(client_sock, "OK")
                        continue
                    except Exception as e:
                        logging.error(f'action: apuesta_recibida | result: fail | cantidad: {len(bets)}')
                        break
                elif action == "FINISH":
                    logging.info(f'action: finalizar_conexion | result: success | ip: {addr[0]}')
                    self.clients[client_sock] = "WAITING"
                    with self.lock:
                        print("recibi FINISH, aumento contador: ", self.clients_waiting.value)
                        self.clients_waiting.value += 1
                        if self.clients_waiting.value == int(TOTAL_CLIENTS):
                            print("Haciendo sorteo")
                            self.hacer_sorteo()
                    break
                
        except OSError as e:
            logging.error("action: receive_message | result: fail | error: {e}")
        finally:            
            # client_sock.close()
            return True

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
        for process in self.procesos:
            process.join()
        sys.exit(0)
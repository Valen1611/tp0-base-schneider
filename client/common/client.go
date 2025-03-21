package common

import (
	"bufio"
	"fmt"
	"net"
	"time"

	"github.com/op/go-logging"
	
	"os"
	"os/signal"
	"syscall"
)

var log = logging.MustGetLogger("log")

// ClientConfig Configuration used by the client
type ClientConfig struct {
	ID            string
	ServerAddress string
	LoopAmount    int
	LoopPeriod    time.Duration
}

// Client Entity that encapsulates how
type Client struct {
	config ClientConfig
	conn   net.Conn
}

// NewClient Initializes a new client receiving the configuration
// as a parameter
func NewClient(config ClientConfig) *Client {
	client := &Client{
		config: config,
	}
	go signal_handler()
	return client
}

// CreateClientSocket Initializes client socket. In case of
// failure, error is printed in stdout/stderr and exit 1
// is returned
func (c *Client) createClientSocket() error {
	conn, err := net.Dial("tcp", c.config.ServerAddress)
	if err != nil {
		log.Criticalf(
			"action: connect | result: fail | client_id: %v | error: %v",
			c.config.ID,
			err,
		)
		log.Critical("Closing socket")
		conn.Close() // Cierro socket si fallo
	}
	c.conn = conn
	return nil
}

func writer(conn net.Conn, msg string) error {
	totalWritten := 0
    msgBytes := []byte(msg)

    for totalWritten < len(msgBytes) {
        n, err := conn.Write(msgBytes[totalWritten:])
        if err != nil {
            return err
        }
        totalWritten += n
    }
    return nil
}

func reader(conn net.Conn) (string, error) {
	msg, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		return "", err
	}
	return msg, nil
}

func (c *Client) Handshake() bool {
	c.createClientSocket()
	initial_request := os.Getenv("HANDSHAKE_REQUEST_MESSAGE") + "\n"
	// log.Infof(
	// 	"action: handshake | result: success | client_id: %v | sending: %v",
	// 	c.config.ID,
	// 	initial_request,
	// )
	writer(c.conn, initial_request)
	log.Infof("waiting server response")
	response, error := reader(c.conn)
	log.Infof("response from server: %v", response)
	if error != nil {
		// log.Criticalf(
		// 	"action: handshake | result: fail | client_id: %v | error: %v",
		// 	c.config.ID,
		// 	error,
		// )
		log.Critical("Closing socket")
		c.conn.Close()
		return false
	}

	expexted_response := os.Getenv("HANDSHAKE_RESPONSE_MESSAGE") + "\n"
	fmt.Println("response from server:")
	fmt.Println(initial_request) 
	if response != expexted_response {
		// log.Criticalf(
		// 	"action: handshake | result: fail | client_id: %v | error: %v",
		// 	c.config.ID,
		// 	"Server response should be Hello Client but got " + response,
		// )
		log.Critical("Closing socket")
		c.conn.Close()
		return false
	}

	// log.Infof(
	// 	"action: handshake | result: success | client_id: %v | recieved: %v",
	// 	c.config.ID,
	// 	response,
	// )
	return true
}

func (c *Client) SendBet() bool {
	name := os.Getenv("NOMBRE")
	surname := os.Getenv("APELLIDO")
	documento := os.Getenv("DOCUMENTO")
	nacimiento := os.Getenv("NACIMIENTO")
	numero := os.Getenv("NUMERO")

	bet_msg := fmt.Sprintf("BET:%v,%v,%v,%v,%v\n", name, surname, documento, nacimiento, numero)
	log.Infof(
		"action: send_bet | result: success | client_id: %v | sending: %v",
		c.config.ID,
		bet_msg,
	)
	writer(c.conn, bet_msg)
	log.Infof("waiting server response")
	response, error := reader(c.conn)
	if error != nil {
		log.Criticalf(
			"action: send_bet | result: fail | client_id: %v | error: %v",
			c.config.ID,
			error,
		)
		log.Critical("Closing socket")
		c.conn.Close()
		return false
	}
	log.Infof("response from server: %v", response)

	log.Infof("action: apuesta_enviada | result: success | dni: %v | numero: %v",
			documento,
			numero,
		)

	return true
}


// StartClientLoop Send messages to the client until some time threshold is met
func (c *Client) StartClientLoop() {
	// There is an autoincremental msgID to identify every message sent
	// Messages if the message amount threshold has not been surpassed
	

	for true {
		c.createClientSocket()
		name := os.Getenv("NOMBRE")
		fmt.Fprintf(
			c.conn,
			"[CLIENT %v] Nombre: %v\n",
			c.config.ID,
			name,
		)
		msg, err := bufio.NewReader(c.conn).ReadString('\n')
		c.conn.Close()
		if err != nil {
			log.Errorf("action: receive_message | result: fail | client_id: %v | error: %v",
				c.config.ID,
				err,
			)
			log.Critical("Closing socket")
			c.conn.Close()
			return
		}
	
		log.Infof("action: receive_message | result: success | client_id: %v | msg: %v",
			c.config.ID,
			msg,
		)
		break
	}

	for msgID := 1; msgID <= c.config.LoopAmount; msgID++ {
		// Create the connection the server in every loop iteration. Send an
		c.createClientSocket()

		// TODO: Modify the send to avoid short-write
		fmt.Fprintf(
			c.conn,
			"[CLIENT %v] Message N°%v\n",
			c.config.ID,
			msgID,
		)
		msg, err := bufio.NewReader(c.conn).ReadString('\n')
		c.conn.Close()

		if err != nil {
			log.Errorf("action: receive_message | result: fail | client_id: %v | error: %v",
				c.config.ID,
				err,
			)
			log.Critical("Closing socket")
			c.conn.Close()
			return
		}

		log.Infof("action: receive_message | result: success | client_id: %v | msg: %v",
			c.config.ID,
			msg,
		)

		// Wait a time between sending one message and the next one
		time.Sleep(c.config.LoopPeriod)

	}
	log.Infof("action: loop_finished | result: success | client_id: %v", c.config.ID)
}

func signal_handler() {
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    c := <-quit
    fmt.Println("Closing client", c)
	os.Exit(0)
}
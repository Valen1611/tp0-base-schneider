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
	"encoding/csv"
)

var log = logging.MustGetLogger("log")

// ClientConfig Configuration used by the client
type ClientConfig struct {
	ID            string
	ServerAddress string
	LoopAmount    int
	LoopPeriod    time.Duration
	MaxAmount	  int
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


func (c *Client) SendBet() bool {
	// Levanto las variables de entorno
	name := os.Getenv("NOMBRE")
	surname := os.Getenv("APELLIDO")
	documento := os.Getenv("DOCUMENTO")
	nacimiento := os.Getenv("NACIMIENTO")
	numero := os.Getenv("NUMERO")
	agency := c.config.ID
	
	// Creo el mensaje de apuesta
	bet_msg := GenerateBetMessage(agency, name, surname, documento, nacimiento, numero)
	log.Infof("action: send_bet | result: success | client_id: %v | sending: %v", c.config.ID, bet_msg)
	
	// Envio la apuesta al server
	c.createClientSocket()
	SocketWriter(c.conn, bet_msg)

	// Espero confirmacion
	response, error := SocketReader(c.conn)
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
	if response != "OK" {
		log.Criticalf(
			"action: send_bet | result: fail | client_id: %v | response: %v",
			c.config.ID,
			response,
		)
		log.Critical("Closing socket")
		c.conn.Close()
		return false
	}

	log.Infof("action: apuesta_enviada | result: success | dni: %v | numero: %v",
			documento,
			numero,
		)

	return true
}

func (c *Client) ReadBets(agency string) ([][]Bet, error) {

	filePath := fmt.Sprintf(".data/agency-%s.csv", c.config.ID)
	fmt.Println(filePath)
	f, err := os.Open(filePath)
    if err != nil {
        log.Criticalf("Unable to read input file " + filePath, err)
    }
    defer f.Close()

    csvReader := csv.NewReader(f)

	allBets := [][]Bet{}  

	for {
		batch := []Bet{}

		for i := 0; i < c.config.MaxAmount; i++ {
			record, err := csvReader.Read()
			if err != nil {
				break
			}
			
			var name, surname, documento, nacimiento, numero string
	
			if len(record) > 0 {
				name = record[0]
			}
			if len(record) > 1 {
				surname = record[1]
			}
			if len(record) > 2 {
				documento = record[2]
			}
			if len(record) > 3 {
				nacimiento = record[3]
			}
			if len(record) > 4 {
				numero = record[4]
			}

			bet := Bet{
				Agency:    agency,
				Name:      name,
				Surname:   surname,
				Documento: documento,
				Nacimiento: nacimiento,
				Numero:    numero,
			}

			batch = append(batch, bet)
		}

		if len(batch) == 0 {
			break
		}

		allBets = append(allBets, batch)
	}

	return allBets, nil
}

func (c *Client) SendBatchBets() bool {
	// Levanto las bets del csv
	

	betBatches, error := c.ReadBets(c.config.ID)
	if error != nil {
		log.Criticalf(
			"action: read_bets | result: fail | client_id: %v | error: %v",
			c.config.ID,
			error,
		)
		return false
	}
	// Envio las apuestas al server
	c.createClientSocket()

	for _, bets := range betBatches {
		msg := GenerateBatchBetMessage(bets)	
		log.Infof("action: send_batch_bets | result: success | client_id: %v | sending: bets", c.config.ID)
		// Envio la apuesta al server
		SocketWriter(c.conn, msg)

		// Espero confirmacion
		response, error := SocketReader(c.conn)
		if error != nil {
			log.Criticalf(
				"action: read_bets | result: fail | client_id: %v | error: %v",
				c.config.ID,
				error,
			)
			return false
		}

		if response != "OK" {
			log.Criticalf(
				"action: send_batch_bets | result: fail | client_id: %v | response: %v",
				c.config.ID,
				response,
			)
			return false
		}


		log.Infof("action: apuestas_enviadas | result: success | client_id: %v",
				c.config.ID,
			)	
	}

	SocketWriter(c.conn, "FINISH:")

	// Espero ganador
	log.Infof("action: esperando_ganadores | result: success | client_id: %v", c.config.ID)
	response, error := SocketReader(c.conn)
	if error != nil {
		log.Criticalf(
			"action: read_bets | result: fail | client_id: %v | error: %v",
			c.config.ID,
			error,
		)
		return false
	}

	// Consulto ganadores
	ganadores := ParseWinnersMessage(response)
	cant_ganadores := len(ganadores)
	log.Infof("action: consulta_ganadores | result: success | cant_ganadores: %v", cant_ganadores)
	log.Infof("action: exit")
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
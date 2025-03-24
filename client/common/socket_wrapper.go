
package common

import (
	"net"
	"encoding/binary"
	"fmt"
)


// Sockets protocol:
//     - 2 bytes que indican el tamanio del mensaje
//     - mensaje

//     Ejemplo:
//         Cuando el server confirma la apuesta, le envia un "OK" al cliente
//     escribiendo:
//     02OK


func SocketWriter(conn net.Conn, msg string) error {
	MSG_LEN_SIZE := 2
	// Formateo el mensaje
    encodedMsg := []byte(msg) 
	msgLength := make([]byte, MSG_LEN_SIZE)
	binary.BigEndian.PutUint16(msgLength, uint16(len(encodedMsg))) 
	fullMsg :=append(msgLength, encodedMsg...)

	// Lo mando envitando short write
	bytesTotales := 0
	for bytesTotales < len(fullMsg) {
		bytesEnviados, err := conn.Write(fullMsg[bytesTotales:])
		if err != nil {
			return err		
		}
		bytesTotales += bytesEnviados
	}

    return nil
}


func SocketReader(conn net.Conn) (string, error) {
	// Primero leo los primeros 2 bytes para
	// saber el tamanio del mensaje
	MSG_LEN_SIZE := 2

	// Leo los primeros 2 bytes
	msgLenB := make([]byte, MSG_LEN_SIZE)
	bytesLeidos := 0
	for bytesLeidos < MSG_LEN_SIZE {
		n, err := conn.Read(msgLenB[bytesLeidos:])
		if err != nil {
			fmt.Println("Error al leer el tamanio del mensaje")
			fmt.Println(err)
			fmt.Println("Bytes leidos: ", bytesLeidos)
			return "", err
		}
		bytesLeidos += n
		
	}
	msgLen := int(binary.BigEndian.Uint16(msgLenB))

	// Ahora leo el mensaje completo
	msg := make([]byte, msgLen)
	bytesLeidosMsg := 0
	for bytesLeidosMsg < msgLen {
		n, err := conn.Read(msg[bytesLeidosMsg:msgLen])
		if err != nil {
			return "", err
		}
		bytesLeidosMsg += n
	}
	return string(msg), nil
}

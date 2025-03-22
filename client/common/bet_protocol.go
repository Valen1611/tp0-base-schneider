package common

import (
	"fmt"
)


// Bet protocol:
// 	- Accion
// 	- Separador ':'
// 	- Datos

// 	Ejemplo:
// 	BET:1,Santiago Lionel,Lorca,30904465,1999-03-17,7574


func GenerateBetMessage(agency string, name string, surname string, documento string, nacimiento string, numero string) string {
	action := "BET"
	return fmt.Sprintf(
		"%s:%v,%v,%v,%v,%v,%v",
		action,
		agency,
		name,
		surname,
		documento,
		nacimiento,
		numero,
	)
}
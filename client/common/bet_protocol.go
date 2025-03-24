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

// En Batch, cada apuesta separada por ;
// 	Ejemplo:
// 	BET:1,Santiago Lionel,Lorca,30904465,1999-03-17,7574;2,Juan,Perez,111111,2001-09-20,2938



type Bet struct {
    Agency    string
    Name      string
    Surname   string
    Documento string
    Nacimiento string
    Numero    string
}


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

func GenerateBatchBetMessage(bets []Bet) string {
	action := "BATCH_BET"
	message := fmt.Sprintf("%s:", action)
	for _, bet := range bets {
		message += fmt.Sprintf(
			"%v,%v,%v,%v,%v,%v;",
			bet.Agency,
			bet.Name,
			bet.Surname,
			bet.Documento,
			bet.Nacimiento,
			bet.Numero,
		)
	}
	return message
}
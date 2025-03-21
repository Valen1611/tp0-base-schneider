func GenerateBetMessage(name string, surname string, documento string, nacimiento string, numero string) string {
	action := "BET"
	return fmt.Sprintf(
		"%s:%v,%v,%v,%v,%v",
		action,
		name,
		surname,
		documento,
		nacimiento,
		numero,
	)
}
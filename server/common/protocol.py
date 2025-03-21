def read_bet_msg(msg):
    bet_data = msg.split(":")[1].split(",")
    agency = bet_data[0]
    name = bet_data[1]
    surname = bet_data[2]
    document = bet_data[3]
    birthdate = bet_data[4]
    number = int(bet_data[5])

    return agency, name, surname, document, birthdate, number

def get_action(msg):
    return msg.split(":")[0]
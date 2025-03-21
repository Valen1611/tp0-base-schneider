def read_bet_msg(msg):
    bet_data = msg.split(":")[1].split(",")
    name = bet_data[0]
    surname = bet_data[1]
    document = bet_data[2]
    birthdate = bet_data[3]
    number = int(bet_data[4])

    return name, surname, document, birthdate, number

def get_action(msg):
    return msg.split(":")[0]
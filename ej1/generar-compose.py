import sys


RUTA_HEADER = "ej1/header.yaml"
RUTA_SERVER = "ej1/server.yaml"
RUTA_CLIENT = "ej1/client.yaml"
RUTA_NETWORK = "ej1/network.yaml"

SERVER_NAME = "server"
CLIENT_NAME = "client"

def get_file_content():
    with open(RUTA_HEADER, "r") as archivo:
        header = archivo.read()

    with open(RUTA_SERVER, "r") as archivo:
        server = archivo.read()
    server = server.replace("server", SERVER_NAME)

    with open(RUTA_CLIENT, "r") as archivo:
        client = archivo.read()
    client = client.replace("CLIENT_NAME", CLIENT_NAME)
    client = client.replace("SERVER_NAME", SERVER_NAME)

    with open(RUTA_NETWORK, "r") as archivo:
        networks = archivo.read()
    
    return header, server, client, networks

def generate_yaml(header, server, client, networks, ruta_archivo_salida, cant_clientes):
    with open(ruta_archivo_salida, "w") as archivo:
        archivo.write(header)
        archivo.write(f"{server}\n")
        for i in range(1, int(cant_clientes)+1):
            client_i = client.replace("CLIENT_NUM", f"{i}")
            archivo.write(f"{client_i}\n")
        archivo.write(networks)

def main():
    # Argumentos de la terminal
    args = sys.argv[1:]
    if len(args) != 2:
        print("Uso: generar-compose.sh <nombre_archivo_salida> <cantidad_clientes>")
        return
    ruta_archivo_salida = args[0]
    cant_clientes = args[1]

    # Levanto los archivos que tienen la estructura del output
    header, server, client, networks = get_file_content()

    # Guardo el yaml el output
    generate_yaml(header, server, client, networks, ruta_archivo_salida, cant_clientes)
    print(f"Archivo {ruta_archivo_salida} generado con exito")

if __name__ == "__main__":
    main()
import requests
from requests.auth import HTTPBasicAuth
import json
from mediaclasses import Libro, Autor

#endpoint = "http://sibila.website:2480"
endpoint = "http://localhost:2480"
"""Raíz de los endpoint utilizados para las llamadas REST
"""
database = "/portico"
"""Nombre de la base de datos a utilizar
"""
method_query = "/query"+database+"/sql/"
"""Método para consultas (GET)
"""
method_batch = "/batch"+database
"""Método para ejecución (batch) de comandos (POST)
"""

def getDatabaseInfo():
    req_query = endpoint + "/database" + database
    response = requests.get(req_query, auth=HTTPBasicAuth('admin', 'admin'))
    if (response.ok):
        # Convertir la cadena de bytes en un string, con encode utf-8
        parsed = json.loads(response.content.decode("utf-8"))
        # el response tiene un elemento llamado result, que contiene los valores devueltos
        result = parsed['result']
    else:
        result = None
    return result

def getClassInfo(className : str):
    req_query = endpoint + "/class" + database + "/" + className
    return execGETQuery(req_query)

def execGETQuery (query):
    """Ejecuta la consulta pasada por parámetro y devuelve los datos en formato JSON.
    Args:
        query (string): Query a ejecutar
    Returns:
        string: JSON con los datos devueltos
    """    
    req_query = endpoint+method_query
    response = requests.get(req_query, auth=HTTPBasicAuth('admin', 'admin'))
    if (response.ok):
        # Convertir la cadena de bytes en un string, con encode utf-8
        parsed = json.loads(response.content.decode("utf-8"))
        # el response tiene un elemento llamado result, que contiene los valores devueltos
        result = parsed['result']
    else:
        result = None
    return result

def insertLibroAutor (libro: Libro, autor : Autor):
    '''Inserta un libro y su autor correspondiende, obviando los datos que ya existan
    '''
    json_libro = json.dumps(libro.__dict__)
    json_autor = json.dumps(autor.__dict__)
    script = """BEGIN; 
LET libro = SELECT from Libro where titulo.toUpperCase() = '{titulolibro}'.toUpperCase();
if ($libro.size() = 0) {{
    LET libro = CREATE VERTEX Libro SET titulo = '{titulolibro}';
}}
LET autor = SELECT from Autor where nombre.toUpperCase() = '{nombreautor}'.toUpperCase();
if ($autor.size() = 0) {{
    LET autor = CREATE VERTEX Autor SET nombre = '{nombreautor}';
}}
LET autorDe = match
        {{class:Autor, as: a, where: (nombre.toUpperCase() = '{nombreautor}'.toUpperCase())}}.out('autorDe') 
        {{class:Libro, as: l, where: (titulo.toUpperCase() = '{titulolibro}'.toUpperCase())}} return a;
if ($autorDe.size() = 0) {{
    CREATE EDGE autorDe FROM $autor TO $libro RETRY 100;
}}
COMMIT;"""
    script = script.format(titulolibro=libro.titulo,nombreautor=autor.nombre)
    operaciones = [{"type":"script","language":"sql","script":[script]}]
    data = {"transaction":True,"operations":operaciones}
    json_data = json.dumps(data,indent=4)
    # Armar el request y enviarlo
    url = endpoint+method_batch
    response = requests.post(url,json=data,auth=("admin","admin"))
    return response

def main():
    #print (execGETQuery("select from Libro"))
    #print (getDatabaseInfo())
    #print (getClassInfo("Libro"))

    libro = Libro (titulo="2001: Una odisea espacial")
    autor = Autor (nombre="Arthur C. Clarke")
    response = insertLibroAutor(libro,autor)
    if (response.ok):
        print(response)
    else:
        print("Hubo error al insertar el libro y autor (ROLLBACK)",response.reason)

if __name__ == '__main__':
    main()
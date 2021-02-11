import requests
from urllib import parse
from requests.auth import HTTPBasicAuth
import json
from mediaclasses import *

class DBManager:
    host = None
    port = None
    database = None
    user = None
    password = None
    baseURL = None
    getURL = None
    batchURL = None
    
    def __init__ (self, host='http://localhost', port=2480, database='portico', user='admin', password='admin'):
        self.host = host
        self.port = port
        self.database = database
        self.user = user
        self.password = password
        self.baseURL = host+":"+str(port)
        
        self.getURL = parse.urljoin(self.baseURL,"/query/"+self.database+"/sql/") 
        self.batchURL = parse.urljoin(self.baseURL,"/batch/"+self.database)

    def execGETQuery (self,query):
        """Ejecuta la consulta pasada por parámetro y devuelve los datos en formato JSON.
        Args:
            query (string): Query a ejecutar
        Returns:
            string: JSON con los datos devueltos
        """    
        req_query = parse.urljoin(self.getURL,query)
        response = requests.get(req_query, auth=HTTPBasicAuth(self.user, self.password))
        if (response.ok):
            # Convertir la cadena de bytes en un string, con encode utf-8
            parsed = json.loads(response.content.decode("utf-8"))
            # el response tiene un elemento llamado result, que contiene los valores devueltos
            result = parsed['result']
        else:
            result = None
        return result

    def insertLibroAutor (self, libro: Libro, autor : Autor):
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
    CREATE EDGE esGenero from $libro to (select from Genero where genero = 'Sci Fi');
    COMMIT;"""
        script = script.format(titulolibro=libro.titulo,nombreautor=autor.nombre)
        operaciones = [{"type":"script","language":"sql","script":[script]}]
        data = {"transaction":True,"operations":operaciones}
        response = requests.post(self.batchURL,json=data,auth=HTTPBasicAuth(self.user, self.password))
        return response

    def updateLibro (self, libro: Libro):
        ''' Actualiza los datos de un libro (buscando por titulo)
        '''
        json_libro = json.dumps(libro.__dict__)
        
        script = """BEGIN; 
    LET libro = SELECT from Libro where titulo.toUpperCase() = '{titulolibro}'.toUpperCase();
    if ($libro.size() = 1) {{
        UPDATE Libro SET
        paginas = {paginas},
        publicado = {publicado},
        sinopsis = '{sinopsis}',
        urlDownload = '{urlDownload}',
        urlPortada = '{urlPortada}'
        WHERE titulo = '{titulolibro}';
    }}
    COMMIT;"""
        script = script.format(
            titulolibro=libro.titulo,
            paginas=libro.paginas,
            publicado=libro.publicado,
            sinopsis=libro.sinopsis,
            urlDownload=libro.urlDownload,
            urlPortada=libro.urlPortada
            )
        operaciones = [{"type":"script","language":"sql","script":[script]}]
        data = {"transaction":True,"operations":operaciones}
        response = requests.post(self.batchURL,json=data,auth=HTTPBasicAuth(self.user, self.password))
        return response

    def deleteAuthorAndBooks(self,autor: Autor):
        ''' Actualiza los datos de un libro (buscando por titulo)
        '''
        script = """BEGIN; 
    DELETE VERTEX Libro WHERE in('autorDe').nombre = '{autor}';
    DELETE VERTEX Autor WHERE nombre = '{autor}';
    COMMIT;"""
        script = script.format(
            autor=autor.nombre
            )
        operaciones = [{"type":"script","language":"sql","script":[script]}]
        data = {"transaction":True,"operations":operaciones}
        response = requests.post(self.batchURL,json=data,auth=HTTPBasicAuth(self.user, self.password))
        return response 

    def insertTwitt (self, id : str, text : str, author_id : str, conversation_id : str, in_reply_to_user_id : str):
        script = """
    BEGIN; 
    LET twitt = SELECT from Twitt where id = '{id}';
    if ($twitt.size() = 0) {{
        CREATE VERTEX Twitt SET
        id = '{id}',
        text = '{text}',
        author_id = '{author_id}',
        conversation_id = '{conversation_id}',
        in_reply_to_user_id = '{in_reply_to_user_id}';
    }}
    COMMIT;"""
        script = script.format(
            id=id,
            text=text,
            author_id=author_id,
            conversation_id=conversation_id,
            in_reply_to_user_id=in_reply_to_user_id
            )
        operaciones = [{"type":"script","language":"sql","script":[script]}]
        data = {"transaction":True,"operations":operaciones}
        response = requests.post(self.batchURL,json=data,auth=HTTPBasicAuth(self.user, self.password))
        return response 

    def insertTwittRelation (self, id_source : str, id_destination : str, relation_type : str):
        script = """
    BEGIN; 
    CREATE EDGE {tipo_edge} from (select from Twitt where id = '{id_source}') to (select from Twitt where id = '{id_destination}');
    COMMIT;"""
        if (relation_type == 'replied_to'):
            tipo_edge = 'TwittReply'
        elif (relation_type == 'quoted'):
            tipo_edge = 'TwittCite'
        elif (relation_type == 'retweeted'):
            tipo_edge = 'TwittRetweet'
        else:
            tipo_edge = 'E'
        script = script.format(
            tipo_edge=tipo_edge,
            id_source=id_source,
            id_destination=id_destination
            )
        operaciones = [{"type":"script","language":"sql","script":[script]}]
        data = {"transaction":True,"operations":operaciones}
        response = requests.post(self.batchURL,json=data,auth=HTTPBasicAuth(self.user, self.password))
        return response


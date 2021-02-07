
class Libro:
    paginas : int = None
    publicado : int = None
    sinopsis : str = None
    titulo : str = None
    urlDownload : str = None
    urlPortada : str = None
    def __init__ (self,titulo,paginas=None,publicado=None,sinopsis=None,urlDownload=None,urlPortada=None):
        self.titulo = titulo
        self.paginas = paginas
        self.publicado = publicado
        self.sinopsis = sinopsis
        self.titulo = titulo
        self.urlDownload = urlDownload
        self.urlPortada = urlPortada

class Autor:
    nombre : str = None
    urlEpubLibre : str = None
    def __init__ (self, nombre, urlEpubLibre=None):
        self.nombre = nombre
        self.urlEpubLibre = urlEpubLibre

class Pelicula:
    pass

class Actor:
    pass

class Cita:
    pass

class Director:
    pass

class Genero:
    pass

class Personaje:
    pass

class Serie:
    pass


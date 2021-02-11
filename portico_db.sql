# portico_db SQL FILE

# Crear la clase abstracta "Persona" de la cual heredar
create class Persona if not exists extends V abstract;

# Crear subclases de persona
create class Autor if not exists extends Persona;
create class Actor if not exists extends Persona;
create class Director if not exists extends Persona;

# Crear el género (por ahora uno solo, luego puedo agrear Fantasía y Terror)
create class Genero if not exists extends V;

# Crear los medios más tradicionales
create class Pelicula if not exists extends V;
create class Serie if not exists extends V;
create class Libro if not exists extends V;

# Crear algunas clases adicionales
create class Personaje if not exists extends V;
create class Cita if not exists extends V;

# Crear algunas clases para las relaciones
create class esGenero extends E;	# Genero del medio
create class apareceEn extends E;	# Dónde aparece un personaje
create class actuoEn extends E;		# Dónde actuo un actor
create class inspiroA extends E;	# Obra que inspiró a otra
create class autorDe extends E;		# Autor de (libro)
create class directorDe extends E;	# Director de (pelicula)
create class tieneCita extends E;	# Una entidad (libro,pelicula,serie,personaje,etc) tiene una cita asociada

# Comenzar a crear las propiedades de las clases
# PERSONA
create property Persona.nombre string (notnull true);

# AUTOR, ACTOR, DIRECTOR por el momento no tienen propiedades distintas a PERSONA

# AUTOR
create property Autor.urlEpubLibre string;
	
# GENERO
create property Genero.genero string (notnull true);

# PERSONAJE
create property Personaje.nombre string (notnull true);

# PELICULA
create property Pelicula.titulo string (notnull true);
create property Pelicula.anio integer;

# LIBRO
create property Libro.titulo string (notnull true);
create property Libro.publicado integer;
create property Libro.paginas integer;
create property Libro.urlPortada string;
create property Libro.urlDownload string;
create property Libro.sinopsis string;

# CREAR INDICES
create index Libro.titulo on Libro(titulo) unique;
create index Persona.nombre on Persona(nombre) unique;
create index Genero.genero on Genero(genero) unique;
create index Pelicula.titulo on Pelicula(titulo) unique;

BEGIN; 
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
COMMIT;

# Estructuras para la gestión de Tweets
create class Twitt if not exists extends V;
create property Twitt.id string (notnull true);
create property Twitt.text string;
create property Twitt.author_id string (notnull true);
create property Twitt.author_name string;
create property Twitt.conversation_id string;
create property Twitt.in_reply_to_user_id string;

# Crear algunas relaciones utiles
create class TwittReply if not exists extends E;
create class TwittRetweet if not exists extends E;
create class TwittCite if not exists extends E;
	
/*
Un modelo detallado de las entidades puede encontrarse en: 
http://network.graphdemos.com/
El modelo presente aquí es un MVP
*/
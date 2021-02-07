from datetime import datetime, timedelta
import time
import logging
import epublibrescraper
import dbmanager

def main():
    logging.basicConfig(level=logging.WARNING)

    start = time.time()
    db = dbmanager.DBManager()

    # Libro 627: Ready Player One
    # ONLINE
    url = epublibrescraper.getBookUrl(4385)
    # LOCAL
    #url = "test-data/ready-player-one.html"
    epub = epublibrescraper.EPubLibreScraper(url)
    libro,autor = epub.scrap()
    if (libro is not None and autor is not None):    
        response = db.insertLibroAutor(libro,autor)
        if (response.ok):
            print(response)
            response = db.updateLibro(libro)
            if (response.ok):
                print(response)
            else:
                print("Hubo error al actualizar el libro(ROLLBACK)",response.reason)
        else:
            print("Hubo error al insertar el libro y autor (ROLLBACK)",response.reason)
    else:
        print ("No se pudo obtener el libro indicado")
    stop = time.time()

    print("Tiempo de ejecución: {0}".format(timedelta(seconds=stop-start)))
    

if __name__ == '__main__':
    main()
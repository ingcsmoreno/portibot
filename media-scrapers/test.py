from datetime import datetime, timedelta
import time
import logging
import epublibrescraper
import dbmanager
import mediaclasses
import moviescraper
import json
import pandas as pd
import dbmanager

def main():
    logging.basicConfig(level=logging.WARNING)
    start = time.time()
    
    ms = moviescraper.MovieScraper()
    '''
    print("API URL:",ms.apiUrl)
    print("IMAGES URL:", ms.imagesUrl)
    status,result = ms.discoverMovies(True)
    if (status == 200):
        print (result)
    else:
        print ("No se pudieron descubrir peliculas Sci Fi")
    status,result = ms.getMovieById("76341",True)
    if (status == 200):
        print (result)
    else:
        print ("No se pudieron encontrar pelicuas usando Id")
    '''

    #status,cast,crew = ms.getMovieCastCrew("76341")
    #status,backdrop,poster = ms.getMovieImageUrl("76341")
    #status,movie,director,actores = ms.scrap_movie("76341")
    

    stop = time.time()

    print("Tiempo de ejecución: {0}".format(timedelta(seconds=stop-start)))

if __name__ == '__main__':
    main()
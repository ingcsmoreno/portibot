from datetime import datetime, timedelta
import time
import logging
import epublibrescraper
import dbmanager
import mediaclasses


def main():
    logging.basicConfig(level=logging.WARNING)
    start = time.time()
    
    


    
    stop = time.time()

    print("Tiempo de ejecución: {0}".format(timedelta(seconds=stop-start)))

if __name__ == '__main__':
    main()
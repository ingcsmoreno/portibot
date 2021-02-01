import argparse
from datetime import datetime, timedelta
import json
import os
import re
import time
import logging
import requests

import pandas as pd
from bs4 import BeautifulSoup, BeautifulStoneSoup
from mediaclasses import Libro, Autor

def get_title (soup):
    div = soup.find(id='titulo_libro')
    return div.string.strip()

def get_author (soup):
    autor_div = soup.find('div',class_='negrita aut_sec')
    autor_link = autor_div.find('a')
    autor = autor_link.contents[0]
    link = autor_link.attrs['href']
    return (autor,link);

def get_sinopsis (soup):
    pass

def get_portada_url (soup):
    img_portada = soup.find('img',id='portada')
    img_url = img_portada.attrs['src']
    return img_url

def get_book_data (soup):
    # Encontrar primero el panel de detalle
    div_detalle = soup.find('div',class_='cab_detalle')
    # dentro del panel buscar el texto 'Páginas'
    paginas_text = div_detalle.find('td',string='Páginas:')
    # Una vez encontrado buscar el 'sibling' que contiene el número
    paginas_nro = paginas_text.next_sibling.next_sibling
    cant_paginas = paginas_nro.string
    publicado_text = div_detalle.find('td',string='Publicado en:')
    publicado_nro = publicado_text.next_sibling.next_sibling
    publicado = publicado_nro.string 	
    return (cant_paginas,publicado)

def get_sinopsis (soup):
    div_detalle = soup.find('div',class_='detalle')
    div_sinopsis = div_detalle.find('div',class_='negrita',string='Sinopsis')
    div_ali_justi = div_sinopsis.next_sibling
    span = div_ali_justi.find('span')
    sinopsis = span.text.replace("\n","").replace("'","")
    return sinopsis

def get_download_link (soup):
    a_download_link = soup.find('a',id='en_desc')
    download_link = a_download_link.attrs['href']
    return download_link

def scrap_book(source : str):
    sp = BeautifulSoup(source,"lxml")
    titulo = get_title(sp)    
    logging.info("Titulo: "+titulo)
    
    autor, link = get_author(sp)
    logging.info("Autor: "+autor)
    logging.info("Link autor: "+link)

    portada_url = get_portada_url(sp)
    logging.info("Portada URL: "+portada_url)

    cant_paginas,publicado = get_book_data(sp)
    logging.info("Cantidad de páginas: "+cant_paginas+" Publicado en: "+publicado)
    
    sinopsis = get_sinopsis(sp)
    logging.info("Sinopsis: "+sinopsis)

    download_link = get_download_link(sp)
    logging.info("Download link: "+download_link)

    # Crear y devolver un Libro y su correspondiente autor
    libro = Libro(titulo,cant_paginas,publicado,sinopsis,download_link,portada_url)
    autor = Autor(autor,link)
    #time.sleep(2)

    return (libro,autor)


    '''
    return {'book_id_title':        book_id, 
            'book_id':              get_id(book_id), 
            'book_title':                ' '.join(soup.find('h1', {'id': 'bookTitle'}).text.split()), 
            'isbn':                 get_isbn(soup),
            'isbn13':               get_isbn13(soup),
            'year_first_published': get_year_first_published(soup), 
            'author':               ' '.join(soup.find('span', {'itemprop': 'name'}).text.split()), 
            'num_pages':            get_num_pages(soup), 
            'genres':               get_genres(soup), 
            'shelves':              get_shelves(soup), 
            'lists':                get_all_lists(soup), 
            'num_ratings':          soup.find('meta', {'itemprop': 'ratingCount'})['content'].strip(), 
            'num_reviews':          soup.find('meta', {'itemprop': 'reviewCount'})['content'].strip(),
            'average_rating':       soup.find('span', {'itemprop': 'ratingValue'}).text.strip(), 
            'rating_distribution':  get_rating_distribution(soup)}
    '''

def main():
    logging.basicConfig(level=logging.WARNING)

    start = time.time()

    # Libro 627: Ready Player One
    # ONLINE
    url = 'https://www.epublibre.org/libro/detalle/627'
    response = requests.get(url)
    if (response.status_code == 200):
        source = response.text
        # LOCAL
        #url = "./test-data/ready-player-one.html"
        #source = open(url)
        
        libro,autor = scrap_book(source)

        print(libro.__dict__)
        print(autor.__dict__)
    else:
        print (response.reason)
    stop = time.time()

    logging.info("Tiempo de ejecución: {0}".format(timedelta(seconds=stop-start)))
    '''
    start_time = datetime.now()
    script_name = os.path.basename(__file__)

    parser = argparse.ArgumentParser()
    parser.add_argument('--book_ids_path', type=str)
    parser.add_argument('--output_directory_path', type=str)
    parser.add_argument('--format', type=str, action="store", default="json",
                        dest="format", choices=["json", "csv"],
                        help="set file output format")
    args = parser.parse_args()

    book_ids              = [line.strip() for line in open(args.book_ids_path, 'r') if line.strip()]
    books_already_scraped =  [file_name.replace('.json', '') for file_name in os.listdir(args.output_directory_path) if file_name.endswith('.json') and not file_name.startswith('all_books')]
    books_to_scrape       = [book_id for book_id in book_ids if book_id not in books_already_scraped]
    condensed_books_path   = args.output_directory_path + '/all_books'

    for i, book_id in enumerate(books_to_scrape):
        try:
            print(str(datetime.now()) + ' ' + script_name + ': Scraping ' + book_id + '...')
            print(str(datetime.now()) + ' ' + script_name + ': #' + str(i+1+len(books_already_scraped)) + ' out of ' + str(len(book_ids)) + ' books')

            book = scrape_book(book_id)
            json.dump(book, open(args.output_directory_path + '/' + book_id + '.json', 'w'))

            print('=============================')

        except HTTPError as e:
            print(e)
            exit(0)
    '''

if __name__ == '__main__':
    main()
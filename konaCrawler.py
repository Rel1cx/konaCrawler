#!/usr/bin/python3
#coding=utf-8
#Author:Eva1ent
#Date:2016-07-10
import os
import sys
import timeit
import urllib.request
from collections import Counter
from multiprocessing.dummy import Pool as pool
from urllib.parse import unquote
import requests
from bs4 import BeautifulSoup
from retrying import retry

path = os.path.expanduser('./konachan/')
if not os.path.exists(path):
    os.makedirs(path)

page_list = []
img_list = []
headers = {
    'User-Agent': 'Mozilla/5.0 (Windows NT 10.0; Win64; x64)',
    'Referer': 'https://konachan.com/post',
    'Host': 'konachan.com'
}
base_url = 'http://konachan.com/post?page='
n = int(input('输入结束页:')) + 1


def get_page():
    for i in range(1, n):
        url = base_url + str(i) + '&tags=order:score'
        page_list.append(url)
    return page_list


def get_filename(img_url):
    namel = unquote(img_url, encoding="utf8")
    names = namel[80:].encode(
        sys.getfilesystemencoding(), errors='surrogateescape')
    return names.decode('latin-1')
    '''names = namel[80:]
    return names'''


@retry(stop_max_attempt_number=3)
def downloader(img_url):
    # print(img_url)
    filename = get_filename(img_url)
    urllib.request.urlretrieve(img_url, path + filename)
    print('[已下载] ' + filename)


@retry(stop_max_attempt_number=3)
def down(page):
    r = requests.get(page, headers=headers)
    soup = BeautifulSoup(r.text, 'html.parser')
    for url in soup.find_all("a", attrs={"class": "directlink"}):
        img_url = url.get('href')
        downloader(img_url)


def main():
    page = get_page()
    # print(Counter(img_list))
    download_pool = pool(processes=8)
    download_pool.map(down, page)
    download_pool.close()
    download_pool.join()


if __name__ == '__main__':
    main()

# t = timeit.timeit(main, 'from __main__ import main', number=1)
# print(t)

from pytube import YouTube
from tinydb import TinyDB, Query

db = TinyDB('fantano.json')
vidids = db.table('vidids')
descriptions = db.table("descriptions")

count = 0
for data in vidids.all():
    metadata = YouTube(f'http://youtube.com/watch?v={data["videoId"]}')
    descriptions.insert({
        "videoId": data["videoId"],
        "title": metadata.title,
        "thumb": metadata.thumbnail_url,
        "description": metadata.description
    })
    count = count+1
    print(count)
from tinydb import TinyDB
import spotipy
from spotipy.oauth2 import SpotifyClientCredentials
import json

sp = spotipy.Spotify(auth_manager=SpotifyClientCredentials(client_id="***REMOVED***",
                                                           client_secret="***REMOVED***"))

new_db = TinyDB("fantano2.json")
final_info = new_db.table("ratings")

all_data = {}


def chunker(seq, size):
    return (seq[pos:pos + size] for pos in range(0, len(seq), size))


for info in final_info.all():
    all_data[info["spotify_id"]] = info

albums = None
count = 0
for rating_chunk in chunker(final_info.all(), 20):
    rating_ids = [rating["spotify_id"] for rating in rating_chunk]
    albums = sp.albums(rating_ids)
    for album in albums["albums"]:
        all_data[album["id"]]["spotify_obj"] = album
        count = count + 1
    print(count)

with open('output.json', 'w') as fp:
    json.dump(all_data, fp)

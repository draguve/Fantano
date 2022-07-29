from tinydb import TinyDB
import spotipy
from spotipy.oauth2 import SpotifyClientCredentials
import json

sp = spotipy.Spotify(auth_manager=SpotifyClientCredentials(client_id="***REMOVED***",
                                                           client_secret="***REMOVED***"))

new_db = TinyDB("fantano2.json")
final_info = new_db.table("ratings")
final_not_found = new_db.table("not_found")

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

all_data_server = []
all_data_client = []

for key, value in all_data.items():
    value["spotify_avail"] = True
    all_data_server.append(value)
    data = value.copy()
    data["spotify_album"] = data["spotify_obj"]["name"]
    data["spotify_label"] = data["spotify_obj"]["label"]
    data["spotify_art_name"] = data["spotify_obj"]["artists"][0]["name"]
    data["spotify_year"] = 0 #do this
    del data["spotify_obj"]
    del data["least_fav_tracks"]
    del data["fav_tracks"]
    del data["image"]
    all_data_client.append(data)

for not_found in final_not_found.all():
    not_found["spotify_avail"] = False
    all_data_server.append(not_found)
    data = not_found.copy()
    del data["least_fav_tracks"]
    del data["fav_tracks"]
    all_data_client.append(data)

with open('output.client.json', 'w') as fp:
    json.dump(all_data_client, fp)

with open('output.server.json', 'w') as fp:
    json.dump(all_data_server, fp)

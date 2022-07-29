from tinydb import TinyDB
import spotipy
from spotipy.oauth2 import SpotifyClientCredentials
import re

sp = spotipy.Spotify(auth_manager=SpotifyClientCredentials(client_id="***REMOVED***",
                                                           client_secret="***REMOVED***"))

db = TinyDB('fantano.json')
new_db = TinyDB("fantano2.json")
scores = db.table("scores")
final_info = new_db.table("ratings")
final_not_found = new_db.table("not_found")


def add_to_final(score, results):
    if results is None:
        print(f'{score["title"]} --- not found')
        final_not_found.insert(score)
    elif len(results["albums"]["items"]) > 0:
        score["spotify_id"] = results["albums"]["items"][0]["id"]
        score["image"] = results["albums"]["items"][0]["images"][0]
        score["spotify_name"] = results["albums"]["items"][0]["name"]
        score["spotify_artists"] = results["albums"]["items"][0]["artists"]
        print(f'{score["spotify_name"]}')
        final_info.insert(score)
    else:
        print(f'{score["title"]} --- not found')
        final_not_found.insert(score)


results = None
for score in scores.all():
    try:
        if "album" in score:
            results = sp.search(q=f"{score['album']} {score['artist']}", type=["album"], limit=1)
            add_to_final(score, results)

        elif "rating" in score:
            album_name = score["title"].lower().replace("album", "").replace("review", "").replace("ep", "").replace("-",
                                                                                                                     "")  # hopefully
            album_name = re.sub("\\(.*?\\)", "", album_name).strip()
            results = sp.search(q=album_name, type=["album"], limit=1)
            add_to_final(score, results)
    except:
        add_to_final(score,None)

import json
from pprint import pprint

f = open('output.server.json')
data = json.load(f)
x = None
for key, value in data.items():
    if "fav_tracks" in value and "spotify_obj" in value:
        for track in value["fav_tracks"]:
            for spot_track in value["spotify_obj"]["tracks"]["items"]:
                if spot_track["name"].lower() in track.lower():
                    spot_track["is_fav"] = True
        for track in value["least_fav_tracks"]:
            for spot_track in value["spotify_obj"]["tracks"]["items"]:
                if spot_track["name"].lower() in track.lower():
                    spot_track["least_fav"] = True

# Closing file
f.close()
with open('output.server.json', 'w') as fp:
    json.dump(data, fp)

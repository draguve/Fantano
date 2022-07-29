from tinydb import TinyDB, Query
import re

db = TinyDB('fantano.json')
descriptions = db.table("descriptions")
scores = db.table("scores")

for data in descriptions.all():
    newData = {
        "videoId": data["videoId"],
        "title": data["title"],
        "thumb": data["thumb"],
    }
    rating = re.findall(
        '(\\s|^)(\\d+|xx|heyyeahprettycoosoundtrackbrah|NOT GREAT|CLASSIC|watch|\u00a9|\u00ae|[\u2000-\u3300]|\ud83c[\ud000-\udfff]|\ud83d[\ud000-\udfff]|\ud83e[\ud000-\udfff])/10(\\s|$)',
        data["description"], re.IGNORECASE)
    if len(rating) > 0:
        rating = [rate[1] + "/10" for rate in rating]
        rating = ",".join(rating)
        print(f"{data['title']} - {rating} - https://youtube.com/watch?v={data['videoId']}")
        newData["rating"] = rating
    else:
        # if "REVIEW" in data["title"]:
        #     print(f"{data['title']} - no rating - https://youtube.com/watch?v={data['videoId']}")
        rating = "no rating"
        print(f"{data['title']} - no rating - https://youtube.com/watch?v={data['videoId']}")

    fav_tracks = " ".join(re.findall('^FAV TRACK.*$', data["description"], re.IGNORECASE | re.MULTILINE)).split(":")[
        -1].split(",")
    least_fav_tracks = \
        " ".join(re.findall('^LEAST FAV TRACK.*$', data["description"], re.IGNORECASE | re.MULTILINE)).split(":")[
            -1].split(
            ",")
    fav_tracks = [i.strip() for i in fav_tracks]
    least_fav_tracks = [i.strip() for i in least_fav_tracks]
    newData["fav_tracks"] = fav_tracks
    newData["least_fav_tracks"] = least_fav_tracks
    album_info = re.search(
        '^([\\w\\d\\s_@.#,&+-]*)-([\\w\\d\\s_@.,#&+-]*)/([\\d\\s]+)/([\\w\\d\\s_@.#,&+-]*)/([\\w\\d\\s_@.#,&+-]*)$',
        data["description"], re.IGNORECASE | re.MULTILINE)
    if album_info:
        album_info = album_info.groups()
        album_info = [i.strip() for i in album_info]
        newData["artist"] = album_info[0]
        newData["album"] = album_info[1]
        newData["year"] = album_info[2]
        newData["fantano_genre"] = album_info[3]
    scores.insert(newData)

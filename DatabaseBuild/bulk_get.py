from tinydb import TinyDB, Query

db = TinyDB('fantano.json')
table = db.table('vidids')

FANTANO_UPLOADS="UUt7fwAhXDy3oNFTAzF2o8Pwx"

# -*- coding: utf-8 -*-

# Sample Python code for youtube.playlistItems.list
# See instructions for running these code samples locally:
# https://developers.google.com/explorer-help/code-samples#python

import os

import google_auth_oauthlib.flow
import googleapiclient.discovery
import googleapiclient.errors

scopes = ["https://www.googleapis.com/auth/youtube.readonly"]

def main():
    # Disable OAuthlib's HTTPS verification when running locally.
    # *DO NOT* leave this option enabled in production.
    os.environ["OAUTHLIB_INSECURE_TRANSPORT"] = "1"

    api_service_name = "youtube"
    api_version = "v3"
    client_secrets_file = "client_secret_921808097933-ek97ob5celtc7r5vb7r1ul64njd36nd1.apps.googleusercontent.com.json"

    # Get credentials and create an API client
    flow = google_auth_oauthlib.flow.InstalledAppFlow.from_client_secrets_file(
        client_secrets_file, scopes)
    credentials = flow.run_console()
    youtube = googleapiclient.discovery.build(
        api_service_name, api_version, credentials=credentials)

    gotIds = 0
    request = youtube.playlistItems().list(
        part="contentDetails",
        maxResults=50,
        playlistId="UUt7fwAhXDy3oNFTAzF2o8Pw"
    )

    while request:
        response = request.execute()
        got_response(response)
        gotIds = gotIds + len(response["items"])
        print(gotIds)

        request = youtube.playlistItems().list_next(
            request, response)


def got_response(response):
    for item in response["items"]:
        table.insert(item["contentDetails"])



if __name__ == "__main__":
    main()

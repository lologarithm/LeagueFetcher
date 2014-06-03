League Fetcher Appengine Server
==============

Setup
--------------
Get riot API key from developer.riotgames.com

Create a config.json and put in your API key like so:
{
    "ApiKey": "MY-KEY-HERE"
}

go run LeagueFetcher lets you run locally without appengine.

To run in appengine you would use
goapp serve .   (if run from root of source)

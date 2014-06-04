README.md
=======
League Fetcher Appengine Server
==============

Setup
--------------
Get riot API key from developer.riotgames.com

Create a config.json and put in your API key like so:
{
    "ApiKey": "MY-KEY-HERE"
}

Building the Client
--------------
Install globals (if needed):
	npm install -g bower
	npm install -g react-tools

Run the following commands from terminal, while in the client folder
	npm install
	bower install
	jsx src/js build/js

Running the client in test
--------------
If needed to run indpendant of the server, there is mock responses ready.
To begin running on the mock server, navigate to client/src/js/settings.js.
Change the field "useLocal" to "true".

Running
--------------
go run LeagueFetcher lets you run locally without appengine.

To run in appengine you would use
goapp serve .   (if run from root of source)

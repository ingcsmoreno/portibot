# Tweet Bot: Recommendations
A simple Tiwtter bot that listen to a single Hashtag, and respondes to tweets with it.

## What's the use case?

The goal behind this tweeter bot is to give tweeter users the oportunity of asking for random recommendentions about certain topics to be define by the hashtag the bot is filtering. For instance, you could have the bot to respond with a random Simpsons Quote everytime someone tweet with hashtag #WhatWouldHomerDo. ¯\\_(ツ)_/¯

## Building the Image

Simply run ```docker build --build-arg SHA1VER=`git rev-parse HEAD` -t tweet-bot .``` to get the binary compiled and generated.

## Running the bot

Use the image you just built (or the one from the [public repo]()) to get the pod running.

To do it on your local machine you can use the following command:
```
docker run --rm --name tweetbot \
-e CONSUMER_KEY=XXA \
-e CONSUMER_SECRET=XXB \
-e ACCESS_TOKEN=XXC \
-e ACCESS_TOKEN_SECRET=XXD \
tweet-bot
```

## Developing on local

If you're working on some changes to the app, you can run the code on your local with this docker command: 

```
docker run --rm --name tweetbot -it \
-v `pwd`/cmd/tweet-bot-r:/go/src \
--workdir /go/src \
-e CONSUMER_KEY=XXA \
-e CONSUMER_SECRET=XXB \
-e ACCESS_TOKEN=XXC \
-e ACCESS_TOKEN_SECRET=XXD \
golang:1.15-alpine go run ./...
```

## How is this repository organized

It follows an standar GoLang project layout as described [here](https://github.com/golang-standards/project-layout).
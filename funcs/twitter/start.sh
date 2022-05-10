#!/bin/bash

ls -la /usr/local/twitter-fn
echo $TWITTER_API_KEY 
echo $TWITTER_API_SECRET_KEY 
echo $TWITTER_ACCESS_TOKEN 
echo $TWITTER_ACCESS_TOKEN_SECRET

# /usr/local/twitter-fn search NBA -o json -c 10 -p 8080 -S \
#                 --twitter-api-key $TWITTER_API_KEY \
#                 --twitter-api-secret-key $TWITTER_API_SECRET_KEY \
#                 --twitter-access-token $TWITTER_ACCESS_TOKEN \
#                 --twitter-access-token-secret $TWITTER_ACCESS_TOKEN_SECRET
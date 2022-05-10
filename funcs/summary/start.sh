#!/bin/bash

/usr/local/summary-fn NBA -o text -c 10 -p 8080 -S \
	                --twitter-fn-url $TWITTER_FN_URL \
	                --watson-fn-url $WATSON_FN_URL
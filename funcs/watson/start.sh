#!/bin/bash

/usr/local/watson-fn vr classify http://pbs.twimg.com/media/EHb34-KXYAESI46.jpg -o json -p 8080 -S \
                --watson-api-key $WATSON_API_KEY \
                --watson-api-url $WATSON_API_URL \
                --watson-api-version $WATSON_API_VERSION
# Demo 1

## clone and build

git clone git@github.com:maximilien/knfun.git

cd knfun

./hack/build.sh

### create docker images (make sure DOCKER_USERNAME is set)

./hack/build.sh --docker

## local tests

./twitter-fn search NBA -c 20 -o json -S -p 8080

curl localhost:8080?q=NBA&o=json

./watson-fn vr classify https://pbs.twimg.com/media/EYWvnD5VcAEYj_l.jpg -o json -S -p 8081

curl localhost:8081?q=https://pbs.twimg.com/media/EYWvnD5VcAEYj_l.jpg&o=json

export TWITTER_FN=http://localhost:8080
export WATSON_FN=http://localhost:8081

./summary-fn NBA -c 20 -S -p 8082


## deploy to Knative cluser

kn service create twitter-fn --env TWITTER_API_KEY=$TWITTER_API_KEY \
                             --env TWITTER_API_SECRET_KEY=$TWITTER_API_SECRET_KEY \
                             --env TWITTER_ACCESS_TOKEN=$TWITTER_ACCESS_TOKEN \
                             --env TWITTER_ACCESS_TOKEN_SECRET=$TWITTER_ACCESS_TOKEN_SECRET \
                             --image docker.io/drmax/twitter-fn:latest --label knfun=demo

kn service create watson-fn --env WATSON_API_KEY=$WATSON_API_KEY \
                            --env WATSON_API_URL=$WATSON_API_URL \
                            --env WATSON_API_VERSION=$WATSON_API_VERSION \
                            --image docker.io/drmax/watson-fn:latest --label knfun=demo

export TWITTER_FN=URL_TO_TWITTER_FN
export WATSON_FN=URL_TO_WATSON_FN

kn service create summary-fn --env TWITTER_FN_URL=$TWITTER_FN_URL \
                             --env WATSON_FN_URL=$WATSON_FN_URL \
                             --image docker.io/drmax/summary-fn:latest

## scaling

kn service update watson-fn --concurrency-limit 1 --min-scale 1 --max-scale 5

## traffic splitting

kn service update summary-fn --env ASYNC=true

kn service update summary-fn --tag summary-fn-lpdgl-1=sync

kn service update summary-fn --tag summary-fn-swhzd-2=async

kn service update summary-fn --traffic async=50,sync=50
# Setup

This repository is self-contained, except for the following dependencies which you should resolve before getting started:

1. Go version go1.12.x or later. Go to [Golang download](https://golang.org/dl/) page to get the version for your OS.
2. [Docker engine](https://docs.docker.com/engine/installation/) for your machine. Specifically ensure you can execute the `docker` executable at the command line.
3. [Knative](https://knative.dev/) cluster. Depending on your Kubernetes cluster provider you can find different means for [getting Knative installed](https://knative.dev/docs/install/) into your cluster.
4. [Knative's `kn` client](https://github.com/knative/client). Once you have your Knative cluster, follow the steps to build `kn` or get one from the latest [released builds](https://github.com/knative/client/releases).

Once these four dependencies are met, continue below with the steps to get and set your APIs credentials, and then build, test, deploy, and run the functions for the demo.

## Credentials

As mentioned above, you need to get credentials for both the [Twitter API](https://developer.twitter.com/en/docs) and the [IBM Watson API](https://cloud.ibm.com/apidocs/visual-recognition/visual-recognition-v3). Once you do, following the links, then you should be able to continue with the steps below. 

To facilitate using and passing these secrets at the command line, and to avoid accidentally divulging secrets during demos, I recommend setting environment variables for each secret string in your shell. Feel free to use other means, but in the description below I am assuming that the secret values are set to the environment variables correspondingly named.

### Twitter

Once you have access to the [Twitter API](https://developer.twitter.com/en/docs), you will have four different keys (as strings). They are:

1. Twitter API Key - TWITTER_API_KEY
2. Twitter API Secret Key - TWITTER_API_SECRET_KEY
3. Twitter Access Token - TWITTER_ACCESS_TOKEN
4. Twitter Access Token Secret - TWITTER_ACCESS_TOKEN_SECRET

To set these in your shell, do the following replacing the content in `"..."` (without the `"`) with your key's value

```bash
export TWITTER_API_KEY="your Twitter API key value here"
export TWITTER_API_SECRET_KEY="your Twitter API secret key value here"
export TWITTER_ACCESS_TOKEN="your Twitter access token value here"
export TWITTER_ACCESS_TOKEN_SECRET="your Twitter access token secret value here"
```

### Watson

Access to the [IBM Watson API](https://cloud.ibm.com/apidocs/visual-recognition/visual-recognition-v3) requires one secret and two constant values. They are:

1. Watson API Key - WATSON_API_KEY
2. Watson API URL - WATSON_API_URL
3. Watson APU Version - WATSON_API_VERSION

```bash
export WATSON_API_KEY="your Watson API key value here"
export WATSON_API_URL=https://gateway.watsonplatform.net/visual-recognition/api
export WATSON_API_VERSION=2018-03-19
```

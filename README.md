# knfun

Knfun is a set of Knative micro-services / micro-functions intended to be used
for live demos. The main goal is to illustrate the end-to-end developer
experience of using [Knative](https://knative.dev) and its `kn` CLI.

The first demo uses three functions in a live demo setting (ideally with an
audience participating). The following diagram is a sketch of this demo and
these [slides](docs/demo1-slides.pdf) can be used for introduction.

![Demo sketch](docs/demo1-sketch.png)

## Functions

1. _TwitterFn_ search function (via
   [Twitter's API](https://developer.twitter.com/en/docs))

   - _in_: hashtags or string to search
   - _in_: count (max number of tweets to return)
   - _out_: recent tweets with images

2. _WatsonFn_ visual recognition image classifier function (via IBM's
   [Watson APIs](https://cloud.ibm.com/apidocs/visual-recognition/visual-recognition-v3))

   - _in_: image URL
   - _out_: image features (class) with confidence (score)

3. _SummaryFn_ function
   - _in_: search string and _TwitterFn_ and _WatsonFn_ URLs
   - _in_: count (max number of tweets)
   - _out_: HTML page displaying summary

Both the _TwitterFn_ and _WatsonFn_ require credentials to execute. This means
that the demoer's secret keys for both the _TwitterFn_ and _WatsonFn_ are also
required as input, however, to simplify the discussion, the keys are sometimes
ommited in diagrams and other places.

# Table of Contents

- [Introduction](#knfun)
- [Setup](docs/setup.md)
  - [Credentials](docs/setup.md/#credentials) 
    - [Twitter](docs/setup.md/#twitter) 
    - [Watson](docs/setup.md/#watson)
- [Build](docs/build.md)
- [Test](docs/test.md)
  - [twitter-fn](docs/test.md/#twitter-fn)
  - [watson-fn](docs/test.md/#watson-fn)
  - [summary-fn](docs/test.md/#summary-fn)
  - [Credentials config](docs/test.md/#credentials-config)
  - [e2e](docs/test.md/#e2e)
- [Deploy](docs/deploy.md)
- [Debug](docs/debug.md)
- [Run](docs/run.md)
  - [TwitterFn](docs/run.md/#TwitterFn)
  - [WatsonFn](docs/run.md/#WatsonFn)
  - [SummaryFn](docs/run.md/#SummaryFn)
  - [Scaling](docs/run.md/#scaling)
  - [A/B Testing or Blue/Green Deployment](docs/run.md/#ab-testing-or-bluegreen-deployment)
    - [Tagging Stable Revisions](docs/run.md/#tagging-stable-revisions) 
    - [Deploy New Async Revision](docs/run.md/#deploy-new-async-revision)
    - [Splitting Traffic](docs/run.md/#splitting-traffic)
- [Future](docs/future.md)
  - [Next steps](docs/future.md/#next-step)
  - [Participate](docs/future.md/#participate)

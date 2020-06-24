# Run

Deploying into Knative means you have a Knative cluster ready and your
`$KUBECONFIG` is set or that the `~/.kube/config` file is pointing to the
cluster's Kubernetes configuration file.

## TwitterFn

```bash
kn service create twitter-fn \
		   --env TWITTER_API_KEY=$TWITTER_API_KEY \
		   --env TWITTER_API_SECRET_KEY=$TWITTER_API_SECRET_KEY \
		   --env TWITTER_ACCESS_TOKEN=$TWITTER_ACCESS_TOKEN \
		   --env TWITTER_ACCESS_TOKEN_SECRET=$TWITTER_ACCESS_TOKEN_SECRET \
		   --image docker.io/drmax/twitter-fn:latest
Creating service 'twitter-fn' in namespace 'default':

  0.496s Configuration "twitter-fn" is waiting for a Revision to become ready.
  0.521s ...
  6.986s ...
  7.148s ...
  7.268s Ready to serve.

Service 'twitter-fn' created with latest revision 'twitter-fn-njnks-1' and URL:
http://twitter-fn.knative-cluster.us-south.containers.cloud.ibm.com
```

## WatsonFn

```bash
kn service create watson-fn \
		   --env WATSON_API_KEY=$WATSON_API_KEY \
		   --env WATSON_API_URL=$WATSON_API_URL \
		   --env WATSON_API_VERSION=$WATSON_API_VERSION \
		   --image docker.io/drmax/watson-fn:latest
Creating service 'watson-fn' in namespace 'default':

  0.197s Configuration "watson-fn" is waiting for a Revision to become ready.
  7.643s ...
  7.752s ...
  7.870s Ready to serve.

Service 'watson-fn' created with latest revision 'watson-fn-dnfxc-1' and URL:
http://watson-fn.default.knative-cluster.us-south.containers.cloud.ibm.com
```

## SummaryFn

You need to then set the environment variables `TWITTER_FN_URL` and
`WATSON_FN_URL` to the URLs that `kn service create ...` showed for the
respective function creation. For instance:

```bash
export TWITTER_FN_URL=twitter-fn.knative-cluster.us-south.containers.cloud.ibm.com
export WATSON_FN_URL=watson-fn.knative-cluster.us-south.containers.cloud.ibm.com
```

Deploy the `summary-fn` service with:

```bash
kn service create summary-fn \
		   --env TWITTER_FN_URL=$TWITTER_FN_URL \
		   --env WATSON_FN_URL=$WATSON_FN_URL \
		   --image docker.io/drmax/summary-fn:latest
...
```

You can test the `summary-fn` function by going to the deployed function URL
with your browser or by using `curl`.

## Scaling

While by default, all services deployed to Knative will autoscale on demand,
that is, scaling down to 0 if there are no requests to their endpoint, and
automatically scaling up when a request arrives again, there are times when you
need to control the scaling characteristics for each service.

In the [current scenario](#knfun), the requests to the _WatsonFn_ service will
be proportional to the number of images found from the _TwitterFn_ requests.
This could be orders of magnitude larger than _TwitterFn_ requests and therefore
we may want to parallelize the requests to the _WatsonFn_ function.

We can
[control the scaling characteristics](https://knative.dev/docs/serving/configuring-the-autoscaler/)
of any Knative service using the `kn` client `service update` command using the
following flags:

- `--concurrency-target` - recommendation for when to scale up based on the
  concurrent number of incoming requests. Defaults to `--concurrency-limit` when
  given.
- `-concurrency-limit` - hard limit of concurrent requests to be processed by a
  single replica (single service instance).
- `--max-scale` - maximum number of replicas.
- `--min-scale` - minimum number of replicas.

There are also means to control the amount of `memory` and `cpu` to request for
and to limit any function. However, for this example, let's use the concurrency
flags above to modify the current deployment of the `WatsonFn` service such that
we increase the concurency target and limit as well as the maximum and minimum
scale values. The following command will achieve this:

```bash
kn service update watson-fn --concurrency-limit 1 \
                            --concurrency-target 1 \
                            --max-scale 10

Updating Service 'watson-fn' in namespace 'default':

  0.124s The Configuration is still working to reflect the latest desired specification.
  6.161s Traffic is not yet migrated to the latest revision.
  6.193s Ingress has not yet been reconciled.
  6.282s Waiting for VirtualService to be ready
  7.597s Ready to serve.

Service 'watson-fn' updated with latest revision 'watson-fn-gmdxr-7' and URL:
http://watson-fn.default.knative-cluster.us-south.containers.cloud.ibm.com
```

By limiting and targeting the concurrency to 1 and setting the max scale to 10,
the _WatsonFn_ will start with 1 replica and autoscale up to 10 replicas as
requests comes in while each replica will be allowed to only process one request
at a time.

## A/B Testing or Blue/Green Deployment

In situation when you want to experiment with different working versions of your
functions, Knative supports the ability to setup your service deployment with
traffic splitting. This allows some percentage of requests to a service to be
routed to a specific version and some other percentage to another. This in
effect allows you to deploy your service to implement some  
[A/B test](https://en.wikipedia.org/wiki/A/B_testing) or do a
[Blue/Green deploy](https://martinfowler.com/bliki/BlueGreenDeployment.html)
between revisions.

The Knative client `kn` has commands to help you configure traffic splitting
easily on any of your services and their revisions. Let's explore one example
with the `summary-fn` service. In the default implementation, the requests to
`twitter-fn` and `watson-fn` are serialized. This results in slow responses in
rendering the UI and therefore a suboptimal experience for end users. A better
one would be to make the calls to `watson-fn` services asynchronous and populate
the image classifications and the tag cloud dynamically as they become
available.

### Tagging Stable Revisions

However, let's first tag the current revision as `stable`. We do this by listing
the revisions for `summary-fn` and tagging the latest one. We first need to find
the revision name for the current (latest) and stable one.

```bash
kn revision list -s summary-fn
NAME                 SERVICE      GENERATION   AGE   CONDITIONS   READY   REASON
summary-fn-tcthz-1   summary-fn   8            10h   4 OK / 4     True
summary-fn-lljhz-1   summary-fn   7            20h   3 OK / 4     True
```

We can notice that the revision listed at the top is the most current revision
with generation number 8. It is also ready and has name `summary-fn-tcthz-1`. We
can use the name to tag it as `latest`, `sync`, and `stable`.

```bash
kn service update summary-fn --tag summary-fn-tcthz-1=latest \
							 --tag summary-fn-tcthz-1=stable \
							 --tag summary-fn-tcthz-1=sync
Updating Service 'summary-fn' in namespace 'default':

  0.310s Waiting for VirtualService to be ready
  1.114s Ready to serve.

Service 'summary-fn' updated with latest revision 'summary-fn-tcthz-1' (unchanged) and URL:
http://summary-fn-default.kndemo-267179.sjc03.containers.appdomain.cloud
```

### Deploy New Async Revision

Next step is to deploy a better, faster, async version of the `summary-fn`
service. The branch
[`summary-fn-async`](https://github.com/maximilien/knfun/tree/summary-fn-async)
has one such implemenation. The primary change in that branch is in the
`summary-fn` service root `/` now points to the new async handler so that when
users hit the `summary-fn` endpoint they will be using the async implementation
of the UI that make use [jQuery](https://jquery.com/) to call the `watson-fn`
function versus doing so in a synchronous fashion, as in the stable (or sync)
revision. A brief summary of the changes are in the following files.

```golang
cat ./funcs/summary/cmds.go
...
func (summaryFn *SummaryFn) summary(cmd *cobra.Command, args []string) error {
	if summaryFn.StartServer {
		http.HandleFunc("/", summaryFn.SummaryAsyncHandler)
...
```

The other change is in the file `./funcs/summary/async_layout.html` which is now
used to render the `summary-fn` UI. In there the primary change is to call the
`watson-fn` function via jQuery.

```javascript
...
	$(document).ready(function() {
	    $.get( "{{ $WatsonFnURL }}?q={{ $imageURL }}", function( data ) {
	      data.classifiers.forEach(function(c){
	        c.classes.forEach(function(b) {
	            $("#tw{{$i}}_{{$j}}").append("<div>"+b["class"]+" "+b.score+"</div>");
	            words.push({text: b["class"],
	                        weight: b.score*1000,
	                        link: "#tw{{$i}}_{{$j}}"})
	        });
	      });
	      $('#cloud').jQCloud(words);
	    });
	});
...
```

Next you will need checkout that branch and build the `summary-fn` function with
the code in that branch; and then create an updated Docker image for the
`SummaryFn` service; and push it to `docker.io`.

```bash
git co summary-fn-async
Switched to branch 'summary-fn-async'
./hack/build.sh
🕸️  Update
🧽  Format
⚖️  License
📖 Docs
🚧 Compile
────────────────────────────────────────────
success

./hack/build.sh --docker
🚧 🐳 build images
   🚧 🐳 twitter-fn
Sending build context to Docker daemon  58.05MB
...
🚧 🐳 summary-fn
Sending build context to Docker daemon  58.05MB
...
   🐳 summary-fn
The push refers to repository [docker.io/drmax/summary-fn]
...
```

Next, let's deploy the new `async` revision and tag it as such.

```bash
kn service update summary-fn \
		   --image docker.io/drmax/summary-fn:latest
...

Service 'watson-fn' updated with latest revision 'summary-fn-jtpyt-1' and URL:
http://watson-fn.default.knative-cluster.us-south.containers.cloud.ibm.com

kn service update summary-fn --tag summary-fn-jtpyt-1=async
Updating Service 'summary-fn' in namespace 'default':

  0.310s Waiting for VirtualService to be ready
  1.114s Ready to serve.

Service 'summary-fn' updated with latest revision 'summary-fn-jtpyt-1' (unchanged) and URL:
http://summary-fn-default.kndemo-267179.sjc03.containers.appdomain.cloud
```

### Splitting Traffic

We can show the details of the `summary-fn` function so we can see the various
revisions and tags in order split traffic 50/50 between the revisions tagged
`sync` and `async`.

```bash
kn service describe summary-fn
Name:       summary-fn
Namespace:  default
Labels:     knfun=demo
Age:        10d
URL:        http://summary-fn-default.kndemo-267179.sjc03.containers.appdomain.cloud
Cluster:    http://summary-fn.default.svc.cluster.local

Revisions:
  100%  @latest (summary-fn-tcthz-1) [8] (24m)
        Image:  docker.io/drmax/summary-fn:latest (pinned to f43b43)
     +  summary-fn-tcthz-1 (current @latest) #latest [8] (24m)
        Image:  docker.io/drmax/summary-fn:latest (pinned to f43b43)
     +  summary-fn-tcthz-1 (current @latest) #stable [8] (24m)
        Image:  docker.io/drmax/summary-fn:latest (pinned to f43b43)
     +  summary-fn-tcthz-1 (current @latest) #sync [8] (24m)
        Image:  docker.io/drmax/summary-fn:latest (pinned to f43b43)
     +  summary-fn-jtpyt-1 (current @latest) #async [9] (28m)
        Image:  docker.io/drmax/summary-fn:latest (pinned to f43b43)

Conditions:
  OK TYPE                   AGE REASON
  ++ Ready                   5s
  ++ ConfigurationsReady    24m
  ++ RoutesReady             5s
```

Splitting traffic is a `service update` command that uses the `--traffic` flag
to set the different revisions and the traffic percentage values.

```bash
kn service update summary-fn --traffic async=50,sync=50
Updating Service 'summary-fn' in namespace 'default':

  0.310s Waiting for VirtualService to be ready
  1.114s Ready to serve.

Service 'summary-fn' updated with latest revision 'summary-fn-jtpyt-1' (unchanged) and URL:
http://summary-fn-default.kndemo-267179.sjc03.containers.appdomain.cloud
```

Now when users access the URL for the `summary-fn` function they will (on
average) see 50% the sync page (which takes a few seconds to reder) and 50% the
async version which is instantaneous and asynchronously shows image features and
the tag cloud of all the images features.

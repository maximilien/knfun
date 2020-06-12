# Debug

There are many ways to debug Knative (and Kubernetes) services. And this section is not meant to describe all approaches, nor be comprehensive. Instead, I want to discuss how to access logs from deployed Knative services, since that is a generally useful first feature when debugging.

## Access logs using [kapp](https://get-kapp.io/)

The `kapp` tool is part of the comprehensive [k14s](https://k14s.io/) Kubernetes toolset. Kapp is a generic tool for deploying and managing Kubernetes services. One of its important features is that it only requires you to just label your Kubernetes objects to enable operating on them.

### Installing `kapp`

You can install kapp by following the instructions in [get-kapp.io](https://get-kapp.io/) or [downloading the latest `kapp` binary](https://github.com/k14s/kapp/releases) for your platform.

### Labeling services

We will use `kapp` to view and tail logs from our Knative services by simply labeling them with a common label. For instance the label: `knfun=demo`. You can do this by passing this label (or whatever label you choose) when creating the service with `kn`.

```bash
kn service create watson-fn --env WATSON_API_KEY=$WATSON_API_KEY \
							--env WATSON_API_URL=$WATSON_API_URL \
							--env WATSON_API_VERSION=$WATSON_API_VERSION \
							--image docker.io/drmax/watson-fn:latest \
							--label knfun=demo
```

Or you can apply this label after creation by updating the service.

```bash
kn service update watson-fn --label knfun=demo
```
### Using `kapp` to tail logs

Once all three services `twitter-fn`, `watson-fn`, and `summary-fn` have the same label, you can then use `kapp` to access the logs for all services at once, and even tail the logs in realtime.

```bash
kapp logs -f -a label:knfun=demo -n default --color
# waiting for 'twitter-fn-qmtyh-1-deployment-57b7d9fc6-pwc8j > user-container' logs to become available...
# waiting for 'twitter-fn-qmtyh-1-deployment-57b7d9fc6-pwc8j > queue-proxy' logs to become available...
# starting tailing 'twitter-fn-qmtyh-1-deployment-57b7d9fc6-pwc8j > user-container' logs
# starting tailing 'twitter-fn-qmtyh-1-deployment-57b7d9fc6-pwc8j > queue-proxy' logs
twitter-fn-qmtyh-1-deployment-57b7d9fc6-pwc8j > queue-proxy | {"level":"info","ts":"
...
```

Now, as the services are used and accessed, the window running the `kapp` command will continously display the output from all services. Of course, you can access the log of only one of the services at a time by labeling the service with its own label, e.g., `--label knfun=twitter`, and running `kapp` to filter only on that label.

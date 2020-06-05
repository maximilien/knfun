# Deploy

In order to deploy into a Knative cluster, you must first create images and publish them into a repository. We will be using Docker for that purpose. 

The `./hack/build.sh --docker-images` and `./hack/build.sh --docker-push` will respectively create Docker images and push them in your [Docker Hub](https://docker.io) account for you. Or use `--docker` to both build and push the images at once.

You need to make sure the `docker` executable is visible to your shell and that the environment variable `DOCKER_USERNAME` is set to your Docker Hub user ID.

```bash
./hack/build.sh --docker-images
🚧 🐳 build images
   🚧 🐳 twitter-fn
Sending build context to Docker daemon  55.42MB
...
...
Successfully tagged drmax/twitter-fn:latest
   🚧 🐳 watson-fn
...
Successfully tagged drmax/watson-fn:latest
   🚧 🐳 summary-fn
...
Successfully tagged drmax/summary-fn:latest
```

And when publishing the images.

```bash
./hack/build.sh --docker-push
🐳 push images
   🐳 twitter-fn
The push refers to repository [docker.io/drmax/twitter-fn]
...
   🐳 watson-fn
The push refers to repository [docker.io/drmax/watson-fn]
...
   🐳 summary-fn
...
```

Your images should not be avaible at:

1. `docker.io/${DOCKER_USERNAME}/twitter-fn`
2. `docker.io/${DOCKER_USERNAME}/watson-fn`
3. `docker.io/${DOCKER_USERNAME}/summary-fn`

You can use these images to deploy the functions into your Knative cluster with the `kn` CLI. 

You are welcome to use my images at:

1. ``docker.io/drmax/twitter-fn:latest`
2. `docker.io/drmax/watson-fn:latest`
3. `docker.io/drmax/summary-fn:latest`

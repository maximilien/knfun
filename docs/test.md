# Test

Let's first explore how to test the functions locally. Later we will show how you can run [e2e](#e2e) tests automatically for a quick sanity check.

## twitter-fn

For each function you can run them locally as a CLI to get immediate response. You can also run them as a local server and use your browser to see the responses and changing input. For example, the following will display recent tweets (up to 20) that have the word `NBA`

```bash
./twitter-fn search NBA -c 20 -o text \
			 --twitter-api-key $TWITTER_API_KEY \
			 --twitter-api-secret-key $TWITTER_API_SECRET_KEY \
			 --twitter-access-token $TWITTER_ACCESS_TOKEN \
			 --twitter-access-token-secret $TWITTER_ACCESS_TOKEN_SECRET
```

To run this function as a server and see JSON output on your browser or with curl, do the following (not the `-S` flag):

```bash
./twitter-fn search NBA -c 20 -o json -S -p 8080 \
			 --twitter-api-key $TWITTER_API_KEY \
			 --twitter-api-secret-key $TWITTER_API_SECRET_KEY \
			 --twitter-access-token $TWITTER_ACCESS_TOKEN \
			 --twitter-access-token-secret $TWITTER_ACCESS_TOKEN_SECRET
```

Then open your browser at `http://localhost:8080` or do `curl http://localhost:8080`. You can then `CTRL-C` to stop the server. 

You can change the input at the browser by passing the query parameter `q` and make it equal the value to search for. For example: `http://localhost:8080?q=NFL&o=json`. If you change the `-o` value to `text` then the Twitter search results will display as formatted text.

To see what other options are available for the `twitter-fn` `search` function get the CLI help with: `./twitter-fn search --help` or `./twitter-fn search -h`.

## watson-fn

Similarly for the `watson-fn` function you can test it locally with any image for which you have a public URL with the following.

```bash
./watson-fn vr classify https://upload.wikimedia.org/wikipedia/commons/c/c3/Jordan_by_Lipofsky_16577.jpg -o text \
			   --watson-api-key $WATSON_API_KEY \
			   --watson-api-url $WATSON_API_URL \
			   --watson-api-version $WATSON_API_VERSION
```

To run this as a local server and see JSON output on your browser or with `curl`, do the following:

```bash
./watson-fn vr classify https://upload.wikimedia.org/wikipedia/commons/c/c3/Jordan_by_Lipofsky_16577.jpg -o json -S -p 8081 \
			   --watson-api-key $WATSON_API_KEY \
			   --watson-api-url $WATSON_API_URL \
			   --watson-api-version $WATSON_API_VERSION &
...
# curl the local server to classify an image
curl http://localhost:8081?q=https://upload.wikimedia.org/wikipedia/commons/c/c3/Jordan_by_Lipofsky_16577.jpg&o=json
```

You can change the input at the browser by passing the URL with the `q` or `query` URL parameter. For example: `http://localhost:8081?q=http://pbs.twimg.com/media/EHpWVAvWoAEfVzO.jpg&o=json`. If you change the `-o` value to `text` then the image classification will display as formatted text.

## summary-fn

Finally, you can test the `summary-fn` function locally after running the `twitter-fn` and `watson-fn` as servers. For instance, if they are running respectively at ports `8080` and `8081`, use the following to run the `summary-fn`.

```bash
./summary-fn NBA -o text -c 10 -o text \
 	     --twitter-fn-url http://localhost:8080 \
	     --watson-fn-url http://localhost:8081
```

To run `summary-fn` as a local server and see output on your browser or with curl, do the following:

```bash
./summary-fn NBA -o text -c 10 -o text -S -p 8082 \
             --twitter-fn-url http://localhost:8080 \
             --watson-fn-url http://localhost:8081
```

Open your browser at `http://localhost:8082` or `curl http://localhost:8082` to see output at the terminal.

## Credentials config

You can avoid passing all the credentials everytime as flags by creating a file named `.knfun.yaml` in your home directory and adding the credential values in it. The following command will create that file and set the values from the environment variables we discussed [above](#credentials).

```bash
touch ~/.knfun.yaml
cat <<EOF >> ~/.knfun.yaml
twitter-api-key: $TWITTER_API_KEY
twitter-api-secret-key: $TWITTER_API_SECRET_KEY
twitter-access-token: $TWITTER_ACCESS_TOKEN
twitter-access-token-secret: $TWITTER_ACCESS_TOKEN_SECRET

# watson-fn
watson-api-key: $WATSON_API_KEY
watson-api-url: https://gateway.watsonplatform.net/visual-recognition/api
watson-api-version: 2018-03-19
EOF
```

Once this config `.knfun.yaml` is present in your home directory, execution of the `twitter-fn` and `watson-fn` will pick up the keys values automatically. So instead of invoking:

```bash
./watson-fn vr classify http://pbs.twimg.com/media/EHpWVAvWoAEfVzO.jpg -o json \
						--watson-api-key $WATSON_API_KEY \
						--watson-api-url $WATSON_API_URL \
						--watson-api-version $WATSON_API_VERSION
```

You can simply do:

```bash
./watson-fn vr classify http://pbs.twimg.com/media/EHpWVAvWoAEfVzO.jpg -o json
```

## e2e

You can easily run end-to-end (e2e) tests once you have configured your credentials in a `~/.knfun.yaml` by invoking the `./test/e2e-tests-local.sh`. 

```bash
./test/e2e-tests-local.sh
📋 Formatting
🧪  Testing
=== RUN   TestSmoke
=== PAUSE TestSmoke
=== CONT  TestSmoke
=== RUN   TestSmoke/verifies_twitter-fn_search
Running 'twitter-fn search NBA -c 10 -o json'...
=== RUN   TestSmoke/verifies_watson-fn_vr_classify
Running 'watson-fn vr classify http://pbs.twimg.com/media/EHb34-KXYAESI46.jpg -o json'...
--- PASS: TestSmoke (3.88s)
    --- PASS: TestSmoke/verifies_twitter-fn_search (0.53s)
    --- PASS: TestSmoke/verifies_watson-fn_vr_classify (3.35s)
...
PASS
ok  	github.com/maximilien/knfun/test/e2e	4.052s
```

To run a quick "smoke" test that verifies each function, run the `./hack/build.sh --test`. It will display any errors and place output content in a file named `/tmp/knfun-test-output.XXXXXX`.

```bash
./hack/build.sh --test
🧪  Tests
  🧪 e2e
```

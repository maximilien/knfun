# Build

After cloning [this repository](https://github.com/maximilien/knfun). You can
build all the functions with the following command:

```bash
./hack/build.sh
🕸️  Update
🧽  Format
⚖️  License
📖 Docs
🚧 Compile
────────────────────────────────────────────
success
```

The result is that you should have three executables: `twitter-fn`, `watson-fn`,
and `summary-fn`.

```bash
ls
LICENSE    docs       go.mod     hack       twitter-fn watson-fn
README.md  funcs      go.sum     summary-fn vendor
```

These executables are designed as both CLIs and server functions that you can
test locally as well as deploy and run on Knative.

# Future

This initial demo was [presented](docs/kubecon-2019-sandiego.pdf) at the IBM mini-theater at KubeCon San Diego on Thursday November 21st, 2019.

## Next steps

There are two immediate next steps I would like to see for v2 of this demo:

1. Improve the UI of the `SummaryFn`. In particular making it more dynamic and usable for live demos. For example having a table view in addition to the list view it currently shows.
2. Use Knative's eventing to refresh the `SummaryFn` page automatically when new tweets are available. This assumes a Knative Twitter API event importer and broker is available to use.

## Participate

I welcome your feedback as [issues](https://github.com/maximilien/knfun/issues) and [pull requests](https://github.com/maximilien/knfun/pulls). Feel free to also reuse this in your own demos. [Contact me with links](mailto:maxim@us.ibm.com?subject=[KnFun]demo%20links) if you do so I can list them.
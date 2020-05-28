# Future

This demo was presented at the following dates and venues:

* IBM mini-theater at KubeCon San Diego on Thursday November 21st, 2019 [slides](docs/kubecon-2019-sandiego.pdf)
* IBM User Group Days online on May 201th, 2020 ([free registration](https://ibm-ugd-platform.bemyapp.com/#/event) required) [slides](docs/ibm-user-group-days-2020-online.pdf)[video](https://ibm-ugd-platform.bemyapp.com/#/conference/5eb1d06bfe3f0f001be7e3c4)

## Next steps

There are two immediate next steps I would like to see for v2 of this demo:

1. Improve the UI of the `SummaryFn`. In particular making it more dynamic and usable for live demos. For example having a table view in addition to the list view it currently shows.
2. Use Knative's eventing to refresh the `SummaryFn` page automatically when new tweets are available. This assumes a Knative Twitter API Knative event source is available to use.

## Participate

I welcome your feedback as [issues](https://github.com/maximilien/knfun/issues) and [pull requests](https://github.com/maximilien/knfun/pulls). Feel free to also reuse this in your own demos. [Contact me with links](mailto:maxim@us.ibm.com?subject=[KnFun]demo%20links) if you do so I can list them.
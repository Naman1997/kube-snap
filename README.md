# kube-snap

## ConfigMap

### event_based_snaps

event_based_snaps defines when a snapshot should be taken. 
New commits are made when a new change is detected.
Possible values: 
true - Take a snapshot whenever an event with type != "Normal" is detected. [kubectl get events --field-selector type!=Normal ]
false - Take a snapshot every time the watcher resyncs irrespective of event type. [kubectl get events -w]

### resync_duration

resync_duration defines the retry dutation for the kubernetes watcher.
Default is 0 seconds, but can be increased to save resources. Increasing this might result in missed events.

## References
- [Kubernetes shared informer](https://gianarb.it/blog/kubernetes-shared-informer)
- [Golang standard folder structure](https://github.com/golang-standards/project-layout)
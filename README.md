# kubesnap

kubesnap is a tool for observing events in kuberenetes and reacting to them in an automated manner.

Currently, kubesnap can react to an event in kubernetes by commiting changes found in the cluster resources to a git repo. This helps in obeservability in terms of what resources were changed inside the cluster whilst the event occoured.

In the future releases, kubesnap will also support [reacting to specific types of events using an ansible palybook.](https://github.com/Naman1997/kubesnap/issues/11)

kubesnap is a [small (<50MB)](https://hub.docker.com/layers/237378881/namanarora/kubesnap/latest/images/sha256-216355533e9c08b7794f476e7a337466fdaaf948513995243cdf6f8d1dd25369?context=repo) image and runs natively inside kubernetes while requiring minimal amounts of memory and cpu.

## Installation

```sh
<!--Clone the repo-->
git clone https://github.com/Naman1997/kubesnap.git
cd kubesnap
<!--Update values for installation-->
vim values.yaml
<!--Begin installation-->
helm install kubesnap ./deploy -n kubesnap -f ./deploy/values.yaml
```

## References
- [Kubernetes shared informer](https://gianarb.it/blog/kubernetes-shared-informer)
- [Golang standard folder structure](https://github.com/golang-standards/project-layout)
# kubesnap
[![Docker Image CI](https://github.com/Naman1997/kubesnap/actions/workflows/docker-image.yml/badge.svg)](https://github.com/Naman1997/kubesnap/actions/workflows/docker-image.yml)
[![Go](https://github.com/Naman1997/kubesnap/actions/workflows/go.yml/badge.svg)](https://github.com/Naman1997/kubesnap/actions/workflows/go.yml)

kubesnap is a tool for observing events in kuberenetes and reacting to them in an automated manner.

Currently, kubesnap can react to an event in kubernetes by commiting changes found in the cluster resources to a git repo. This helps in obeservability in terms of what resources were changed inside the cluster whilst the event occoured.

kubesnap is a [small](https://hub.docker.com/r/namanarora/kubesnap/tags) image and runs natively inside kubernetes while requiring minimal amounts of memory and cpu.

## Installation

```sh
<!--Clone the repo-->
git clone https://github.com/Naman1997/kubesnap.git
cd kubesnap/deploy
<!--Update values for installation-->
vim values.yaml
<!--Update values for ansible[Create a yaml file for each event you want to auto-remediate]-->
cd ansible
vim ansible.cfg
<!--You can remove/update the default playbook provided in that dir-->
vim Failed.yaml
<!--Go back to the dir with the values file-->
cd ..
<!--Begin installation-->
helm install kubesnap ./deploy -n kubesnap -f ./deploy/values.yaml
```

## References
- [Kubernetes shared informer](https://gianarb.it/blog/kubernetes-shared-informer)
- [Golang standard folder structure](https://github.com/golang-standards/project-layout)

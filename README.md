# kubesnap
[![Docker Image CI](https://github.com/Naman1997/kubesnap/actions/workflows/docker-image.yml/badge.svg)](https://github.com/Naman1997/kubesnap/actions/workflows/docker-image.yml)
[![Go](https://github.com/Naman1997/kubesnap/actions/workflows/go.yml/badge.svg)](https://github.com/Naman1997/kubesnap/actions/workflows/go.yml)

kubesnap is a tool for observing events in kuberenetes and reacting to them in an automated manner.

Kubesnap can react to an event in kubernetes by commiting changes found in the cluster resources to a git repo. This helps in obeservability in terms of what resources were changed inside the cluster whilst the event occoured.

Kubesnap remediatiates issues by running an ansible playbook specified by the user for a particular event. Yuo can find the default playbook for a Failed event in deploy/ansible dir.

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

## Notes

`kubectl` is available inside the pod. However the default access provided to the pod is "get", "list" and "watch". Therefore, if you want to run commands like `kubectl create` or `kubectl delete` - you'll need to manually edit the cluster role access in [clusterrole.yaml](https://github.com/Naman1997/kubesnap/blob/main/deploy/templates/clusterrole.yaml).

## Todo List

- [x] [Take a cluster snapshot in a git repo on a failure event](https://github.com/Naman1997/kubesnap/issues/6)
- [x] [Enable auto-remediation using ansible playbooks](https://github.com/Naman1997/kubesnap/issues/11)
- [x] [Add a helm chart](https://github.com/Naman1997/kubesnap/issues/3)
- [ ] [Add unit tests for methods in /pkg](https://github.com/Naman1997/kubesnap/issues/16)
- [ ] [Create a UI with auth - user can see logs from here initially](https://github.com/Naman1997/kubesnap/issues/21)
- [ ] Save state of events and remediations using a database
- [ ] Show each auto-remediation in a table with the commitId for it's snapshot in the UI
- [ ] Send an update event to the UI after the database has been updated (Some sort of auto-refresh)
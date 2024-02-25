# k8s-image-watcher
**NB.** Still WIP.

Application for watching and notifying image changes from daemonsets, deployments and statefulsets.

# Usage

```sh
helm repo add terotuomala https://terotuomala.github.io/k8s-image-watcher/
```

```sh
helm upgrade --install k8s-image-watcher terotuomala/k8s-image-watcher --create-namespace -n k8s-image-watcher
```

# Local development

### Install depedencies
```sh
task deps
```

### Run the application in dev mode using skaffold
```sh
task dev
```

### Deploy the application using skaffold 
```sh
task deploy
```

### Generate helm manifest from templates
```sh
task generate-manifests
```
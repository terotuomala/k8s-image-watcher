# k8s-image-watcher

Application for watching and notifying image changes from daemonsets, deployments and statefulsets.

# Install

### Helm

```sh
helm repo add terotuomala https://terotuomala.github.io/k8s-image-watcher/
```
If notifications to Slack is enabled (disabled by default) add your Slack-token:
```sh
kubectl create ns k8s-image-watcher
kubectl -n k8s-image-watcher create secret generic slack-secret --from-literal=token=<your-slack-token>
````

```sh
helm upgrade --install k8s-image-watcher \
-n k8s-image-watcher \
--create-namespace \
terotuomala/k8s-image-watcher
```

# Configuration
The following tables lists the configurable parameters of the k8s-image-watcher chart and their default values.

| Parameter                         | Default                | Description                                                                                                            |
| --------------------------------- | ---------------------- | ---------------------------------------------------------------------------------------------------------------------- |
| `replicaCount`                    | `1`                    | Desired number of pods                                                                                                 |
| `logging.level`                   | `info`                 | Log level: `debug`, `info`, `warn`, `error`                                                                            |
| `image.repository`                | `ghcr.io/terotuomala/k8s-image-watcher` | Image repository                                                                                      |
| `image.tag`                       | `""`                     | Image tag                                                                                                            |
| `image.pullPolicy`                | `IfNotPresent`         | Image pull policy                                                                                                      |
| `serviceAccount.create`           | `true`                 | Whether a service account should be created                                                                            |
| `securityContext`                 | `{}`                   | The security context to be set on the k8s-image-watcher container                                                      |
| `resources.requests.cpu`          | `1m`                   | Pod CPU request                                                                                                        |
| `resources.requests.memory`       | `5Mi`                  | Pod memory request                                                                                                     |
| `resources.limits.cpu`            | `5m`                   | Pod CPU limit                                                                                                          |
| `resources.limits.memory`         | `15Mi`                 | Pod memory limit                                                                                                       |
| `nodeSelector`                    | `{}`                   | Node labels for pod assignment                                                                                         |
| `tolerations`                     | `[]`                   | List of node taints to tolerate                                                                                        |
| `affinity`                        | `{}`                   | Node/pod affinities                                                                                                    |
| `podAnnotations`                  | `{}`                   | Pod annotations                                                                                                        |
| `watch.namespace`                 | `""`                   | Which namespace(s) will be watched. Single namespace or all namespaces ("" = all namespaces)                           |
| `watch.deployment`                | `true`                 | Watch Deployments from `watch.namespace` namespace(s)                                                                  |
| `watch.daemonset`                 | `false`                | Watch DaemonSets from `watch.namespace` namespace(s)                                                                   |
| `watch.statefulset`               | `false`                | Watch StatefulSets from `watch.namespace` namespace(s)                                                                 |
| `slack.enabled`                   | `false`                | Enable notifications to Slack                                                                                          |
| `slack.channel`                   | `""`                   | Slack channel name where to send notifications                                                                         |
| `slack.messageTitle`              | `""`                   | Slack message tittle                                                                                                   | 

# Local development
> I like to use [go-task](https://taskfile.dev/installation/) in order to make the usage a bit easier:

### Install depedencies
```sh
task deps
```

### Install tools (k3d, skaffold) for local development
```sh
task tools
```

### Create local k3s cluster using k3d
```sh
task create-k3s-cluster
```

### Run the application in dev mode using skaffold (uses gcr.io/distroless/base-debian11:debug image)
```sh
task dev
```

### Deploy the application using skaffold (uses cgr.dev/chainguard/static image)
```sh
task deploy
```

### Generate helm manifest from templates
```sh
task generate-manifests
```

### Delete local k3s cluster
```sh
task create-k3s-delete
```
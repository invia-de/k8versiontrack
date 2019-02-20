## K8VersionTrack Management
Installs K8VersionTrack in our Cluster.


| Parameter                       | Description                                                          | Default                                   |
| ------------------------------- | -------------------------------------------------------------------- | ----------------------------------------- |
| `fullnameOverride`              | Override the full resource names                                     | `{release-name}-k8versiontrack (or k8versiontrack if release-name is k8versiontrack`|
|`image`			| Used image							| ??  |
|`feedsConfig` 			| Custom feeds.yaml |  |
|`ingress.enabled` | If true, K8VersionTrack Ingress will be created | `false`|
|`ingress.annotations` | K8VersionTrack Ingress annotations | `{}`|
|`ingress.extraLabels` | K8VersionTrack Ingress additional labels | `{}`|
|`ingress.hosts` | K8VersionTrack Ingress hostnames | `[]`|
|`ingress.tls` | K8VersionTrack Ingress TLS configuration (YAML) | `[]`|
| resources | Deployment Resources | ```resources: 
  limits:
   cpu: 100m
   memory: 128Mi
  requests:
   cpu: 10m
   memory: 128Mi
``` |

Example yaml:

```


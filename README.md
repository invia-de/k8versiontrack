## K8s Application Version Checks
  * show what Version of different Tools are Installed in a Cluster, at what Version and whats the newest Version on Github
  * Latest Releases are parsed from Github Atom Feeds

## Development

1. Git Clone

2. Create and fill Configuration ./config/feeds.yaml

3. Run make build_container

4. Run make run_container

5. Once in Container run
    ```
    kubetoken for k8s login
    kubectl port-forward $(kubectl get pod --selector app=helm,name=tiller -o jsonpath={.items..metadata.name} -n kube-system) 44134:44134 -n kube-system
    go run main.go
    ```

6. Point your Browser to localhost:8888

## Application and Cluster Configuration

* Fill config/feeds.yaml following Example:
```
TillerConnectionURI: "127.0.0.1:44134"
StaticFeeds:
  -
    Link: https://github.com/kubernetes/kops/releases.atom
    Name: Kops
    Installed: 1.11.0
FeedMap:
  -
    Link: https://github.com/prometheus/prometheus/releases.atom
    Name: prometheus

```
* TillerConnectionURI should be the full URI to Tiller, for example: tiller-deploy.kube-system:44134#
* StaticFeeds are Feeds that would be parsed for latest Versions. The Installed Version is static
* FeedMap maps `Name` to Helms deployment `Name` and uses the configured Feed URL for finding the latest Version

## Environment Variables for Configuration
* **KVT_CONFIG_PATH:** Path to feeds.yaml config file. Default: `"./config"`

* **KVT_HTTP_ADDR:** The host and port. Default: `":8888"`

* **KVT_HTTP_CERT_FILE:** Path to cert file. Default: `""`

* **KVT_HTTP_KEY_FILE:** Path to key file. Default: `""`

* **KVT_HTTP_DRAIN_INTERVAL:** How long application will wait to drain old requests before restarting. Default: `"1s"`

## Installation
* Use provided Helm Chart from `charts/` Directory

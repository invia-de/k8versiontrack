## K8s Application Version Checks
  * Offers the possibility to monitor versions of installed tools in a kubernetes cluster
  * Automagically finds tools that were installed with helm
  * Offers prometheus metrics endpoint for integration in grafana
  * Lists the most recent versions from github
  * Latest releases are parsed from github atom feeds

## Development

1. Git clone

2. Create/copy and fill configuration from template ./config/feeds.yaml

3. Run ```make build_container```

4. Run ```make run_container```

5. Once in container run

    ```kubetoken``` for k8s login  
    ```kubectl port-forward $(kubectl get pod --selector app=helm,name=tiller -o jsonpath={.items..metadata.name} -n kube-system) 44134:44134 -n kube-system```  
    ```go run main.go```


6. Point your browser to localhost:8888

## Application and Cluster Configuration

* Fill config/feeds.yaml following this example:
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
* TillerConnectionURI should be the full URI to tiller, for example: tiller-deploy.kube-system:44134#
* StaticFeeds are feeds that should be parsed for latest versions. The installed version is static and has to be filled manually for components that don't allow automatic parsing from feeds
* FeedMap maps `Name` to helm's deployment `Name` and uses the configured feed URL to get the latest version

## Environment Variables for Configuration
* **KVT_CONFIG_PATH:** Path to feeds.yaml config file. Default: `"./config"`

* **KVT_CACHE_ENABLE:** Enable File Cache for Results. Default: true

* **KVT_CACHE_TIME:** Cache TTL in Seconds. Default: 600

* **KVT_HTTP_ADDR:** The host and port. Default: `":8888"`

* **KVT_HTTP_CERT_FILE:** Path to cert file. Default: `""`

* **KVT_HTTP_KEY_FILE:** Path to key file. Default: `""`

* **KVT_HTTP_DRAIN_INTERVAL:** How long application will wait to drain old requests before restarting. Default: `"1s"`

## Installation
* Use provided helm chart from `charts/` directory

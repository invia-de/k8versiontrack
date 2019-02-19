## K8s Application Version Checks
  * show what Version of different Tools are Installed in a Cluster, at what Version and whats the newest Version on Github
  * Latest Releases are parsed from Github Atom Feeds

## Installation

1. Git Clone

2. Run make build_container

3. Run make run_container

4. Once in Container run 
    ```
    kubetoken for k8s login
    kubectl port-forward $(kubectl get pod --selector app=helm,name=tiller -o jsonpath={.items..metadata.name} -n kube-system) 44134:44134 -n kube-system 
    go run main.go
    ```

5. Point your Browser to localhost:8888

## Application and Cluster Configuration
* Fill config/feeds.yaml

## Environment Variables for Configuration

* **HTTP_ADDR:** The host and port. Default: `":8888"`

* **HTTP_CERT_FILE:** Path to cert file. Default: `""`

* **HTTP_KEY_FILE:** Path to key file. Default: `""`

* **HTTP_DRAIN_INTERVAL:** How long application will wait to drain old requests before restarting. Default: `"1s"`


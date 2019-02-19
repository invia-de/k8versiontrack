package controller

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/mmcdole/gofeed"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	//"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/helm/pkg/helm"
)

type VersionList struct {
	LatestVersion string
	Name          string
	Link          string
	Installed     []install
}
type Result struct {
	Objects []VersionList
}

type config struct {
	Feeds []feed
}
type install struct {
	Cluster string
	Version string
}
type feed struct {
	Name      string
	Link      string
	Installed []install
}

func readConfig(filename string) (*viper.Viper, error) {
	v := viper.New()
	v.SetConfigName(filename)
	v.AddConfigPath("./config/")
	err := v.ReadInConfig()
	return v, err
}
func connectK8s() (*kubernetes.Clientset, error) {
	kubeconfig := filepath.Join(os.Getenv("HOME"), ".kube", "config")
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err.Error())
	}
	clientset, err := kubernetes.NewForConfig(config)
	return clientset, err
}

/*
func searchInstall(api *v1.CoreV1Interface, search string) (pvc *v1.PersistentVolumeClaimList, err error) {
	listOptions := metav1.ListOptions{LabelSelector: "", FieldSelector: ""}
	pvc, err = api.PersistentVolumeClaims("jenkins").List(listOptions)
	return pvc, err
}
**/
func printPVCs(pvcs *v1.PersistentVolumeClaimList) {
	if len(pvcs.Items) == 0 {
		log.Println("No claims found")
		return
	}
	template := "%-32s%-8s%-8s\n"
	fmt.Println("--- PVCs ----")
	fmt.Printf(template, "NAME", "STATUS", "CAPACITY")
	var cap resource.Quantity
	for _, pvc := range pvcs.Items {
		quant := pvc.Spec.Resources.Requests[v1.ResourceStorage]
		cap.Add(quant)
		fmt.Printf(template, pvc.Name, string(pvc.Status.Phase), quant.String())
	}

	fmt.Println("-----------------------------")
	fmt.Printf("Total capacity claimed: %s\n", cap.String())
	fmt.Println("-----------------------------")
}
func IndexAction(w http.ResponseWriter, r *http.Request) {
	//@// TODO: Env Var for Local / Cluster Usage
	helmclient := helm.NewClient(helm.Host("127.0.0.1:44134"))
	releases, err := helmclient.ListReleases()
	fp := gofeed.NewParser()
	ret := releases.GetReleases()
	for _, element := range ret {

		md := element.GetChart().GetMetadata()
		sources := md.GetSources()
		//relStatus, _ := helmclient.ReleaseContent(element.Name)
		logrus.Error(element.Name, "installed Version: ", md.GetAppVersion())
		if len(sources) > 0 {
			for _, link := range sources {
				if strings.Contains(link, element.Name) && strings.Contains(link, "github") {
					parseUrl := link + "/releases.atom"
					feed, feedParseError := fp.ParseURL(parseUrl)
					if feedParseError != nil {
						//logrus.Error(feedParseError)
						//http.Error(w, "Service Unavailable", http.StatusInternalServerError)
						logrus.Error("Feed Parse Error of ", parseUrl, " : ", feedParseError)
						continue
					}
					if feed.FeedType == "atom" {
						for _, ele := range feed.Items {
							logrus.Error("Latest Version of ", element.Name, ": ", ele.Title)
							// Only the First Elemnt would be latest
							break
						}
					}

					//logrus.Error("Link %s", el)
				}
			}
			//link := sc[0:1]
			//logrus.Error(sc)
			//logrus.Error(link)
		}
		//logrus.Error(element.Name)
		//logrus.Error(element.AppVersion)
		// index is the index where we are
		// element is the element from someSlice for where we are
	}

	/*
	       clientset, err := connectK8s()
	   	if err != nil {
	   		log.Fatal(err)
	   	}
	   	api := clientset.CoreV1()
	   	listOptions := metav1.ListOptions{LabelSelector: "", FieldSelector: ""}
	   	pvcs, err := api.PersistentVolumeClaims("jenkins").List(listOptions)
	   	//pvcs, err := searchInstall(api, "bla")
	   	if err != nil {
	   		log.Fatal(err)
	   	}

	   	printPVCs(pvcs)
	*/
	var res Result
	res.Objects = make([]VersionList, 0)
	v1, err := readConfig("feeds")
	if err != nil {
		logrus.Error(err)
		http.Error(w, "Service Unavailable", http.StatusInternalServerError)
	}
	var C config
	errUnmarshal := v1.Unmarshal(&C)
	if errUnmarshal != nil {
		logrus.Error(errUnmarshal)
		http.Error(w, "Service Unavailable", http.StatusInternalServerError)
	}

	///toParse := v1.GetStringSlice("feeds")
	//fp := gofeed.NewParser()
	//fmt.Println(C)
	for _, feedToParse := range C.Feeds {
		//fmt.Println(feedToParse.Name)
		feed, feedParseError := fp.ParseURL(feedToParse.Link)
		if feedParseError != nil {
			//logrus.Error(feedParseError)
			//http.Error(w, "Service Unavailable", http.StatusInternalServerError)
			res.Objects = append(res.Objects,
				VersionList{
					Name:          feedToParse.Name,
					LatestVersion: "No Feed Provided",
					Link:          feedToParse.Link,
					Installed:     feedToParse.Installed,
				})
			continue
		}
		if feed.FeedType == "atom" {
			for _, element := range feed.Items {
				ms := VersionList{
					Name:          feedToParse.Name,
					LatestVersion: element.Title,
					Link:          feed.Link,
					Installed:     feedToParse.Installed,
				}
				res.Objects = append(res.Objects, ms)
				// Only the First Elemnt would be latest
				break
			}
		}
	}

	w.Header().Set("Content-Type", "text/html")

	tmpl, err := template.ParseFiles("templates/dashboard.html.tmpl", "templates/index/index.html.tmpl")
	if err != nil {
		logrus.Error(err)
		http.Error(w, "Service Unavailable", http.StatusInternalServerError)
		//libhttp.HandleErrorJson(w, err)
		return
	}
	//fmt.Println(res)
	err = tmpl.Execute(w, res)
	if err != nil {
		//libhttp.HandleErrorJson(w, err)
		logrus.Error(err)
		http.Error(w, "Service Unavailable", http.StatusInternalServerError)
		return
		//fmt.Println("executing template:", err)
	}
}

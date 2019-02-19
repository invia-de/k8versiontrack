package controller

import (
	//"fmt"
	//"log"

	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/mmcdole/gofeed"
	//"k8s.io/api/core/v1"
	//"k8s.io/apimachinery/pkg/api/resource"
	//"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/helm/pkg/helm"
)

type PrometheusList struct {
	LatestVersion    string
	Name             string
	Link             string
	InstalledVersion string
}
type PrometheusResult struct {
	Objects []PrometheusList
}

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
	FeedMap     []feed
	StaticFeeds []feed
}
type install struct {
	Cluster string
	Version string
}
type feed struct {
	Name      string
	Link      string
	Installed string
}

func readConfig(filename string) (*viper.Viper, error) {
	v := viper.New()
	v.SetConfigName(filename)
	v.AddConfigPath("./config/")
	err := v.ReadInConfig()
	return v, err
}
func cleanVersion(version string) string {
	re := regexp.MustCompile("(unknown|[0-9.]+)")
	return re.FindString(version)
}
func getLatestVersionByFeedUrl(parseUrl string) string {
	fp := gofeed.NewParser()
	lv := "unknown"
	feed, feedParseError := fp.ParseURL(parseUrl)
	if feedParseError != nil {
		logrus.Error("Feed Parse Error of ", parseUrl, " : ", feedParseError)
		return lv
	}
	if feed.FeedType == "atom" {
		for _, ele := range feed.Items {
			//logrus.Error("Latest Version of ", element.Name, ": ", ele.Title)
			lv = ele.Title
			// Only the First Elemnt would be latest
			break
		}
	}
	return lv
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
/*
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
*/
func IndexAction(w http.ResponseWriter, r *http.Request) {
	//@// TODO: Env Var for Local / Cluster Usage
	helmclient := helm.NewClient(helm.Host("127.0.0.1:44134"))
	releases, _ := helmclient.ListReleases()

	v1, err := readConfig("feeds")
	if err != nil {
		logrus.Error(err)
	}
	var C config
	errUnmarshal := v1.Unmarshal(&C)
	if errUnmarshal != nil {
		logrus.Error(errUnmarshal)
	}
	var res PrometheusResult
	res.Objects = make([]PrometheusList, 0)

	ret := releases.GetReleases()
	for _, element := range ret {

		parseUrl := "unknown"
		latestVersion := "unknown"

		md := element.GetChart().GetMetadata()

		//Map Overridden feeds
		for _, feedMap := range C.FeedMap {
			if strings.Contains(feedMap.Name, element.Name) {
				parseUrl = feedMap.Link
			}
		}
		//Determine Chart Release Link by Name
		sources := md.GetSources()
		if parseUrl == "unknown" && len(sources) > 0 {
			for _, link := range sources {
				if strings.Contains(link, element.Name) && strings.Contains(link, "github") {
					parseUrl = link + "/releases.atom"
					//logrus.Error("Link %s", el)
				}
			}
		}

		logrus.Error("Linkparse", parseUrl, " ", element.Name)
		if parseUrl != "unknown" {
			latestVersion = getLatestVersionByFeedUrl(parseUrl)

		}
		ms := PrometheusList{
			Name:             element.Name,
			LatestVersion:    cleanVersion(latestVersion),
			InstalledVersion: cleanVersion(md.GetAppVersion()),
		}
		res.Objects = append(res.Objects, ms)
	}

	for _, feedToParse := range C.StaticFeeds {
		latestVersion := "unknown"
		latestVersion = getLatestVersionByFeedUrl(feedToParse.Link)
		ms := PrometheusList{
			Name:             feedToParse.Name,
			LatestVersion:    cleanVersion(latestVersion),
			InstalledVersion: cleanVersion(feedToParse.Installed),
		}
		res.Objects = append(res.Objects, ms)
	}

	logrus.Error(res)

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

	/*
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
	*/
}

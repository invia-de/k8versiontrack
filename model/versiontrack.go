package model

import (
	"github.com/Sirupsen/logrus"
	"github.com/spf13/viper"

	"regexp"
	"strings"

	"github.com/mmcdole/gofeed"
	//"k8s.io/api/core/v1"
	//"k8s.io/apimachinery/pkg/api/resource"
	//"k8s.io/apimachinery/pkg/watch"
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
type config struct {
	FeedMap             []feed
	StaticFeeds         []feed
	TillerConnectionURI string
}
type feed struct {
	Name      string
	Link      string
	Installed string
}

func GetVersions() PrometheusResult {
	v1, err := readConfig("feeds")
	if err != nil {
		logrus.Error(err)
	}
	var C config
	errUnmarshal := v1.Unmarshal(&C)
	if errUnmarshal != nil {
		logrus.Error(errUnmarshal)
	}

	helmclient := helm.NewClient(helm.Host(C.TillerConnectionURI))
	releases, err := helmclient.ListReleases()
	if err != nil {
		logrus.Error(err)
	}

	var res PrometheusResult
	res.Objects = make([]PrometheusList, 0)

	ret := releases.GetReleases()
	for _, element := range ret {

		parseUrl := "unknown"
		latestVersion := "unknown"

		md := element.GetChart().GetMetadata()
		//logrus.Error(md)
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

		//logrus.Error("Linkparse", parseUrl, " ", element.Name)
		if parseUrl != "unknown" {
			latestVersion = getLatestVersionByFeedUrl(parseUrl)
		}
		ms := PrometheusList{
			Name:             element.Name,
			Link:             strings.Replace(parseUrl, ".atom", "", -1),
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
			Link:             strings.Replace(feedToParse.Link, ".atom", "", -1),
			LatestVersion:    cleanVersion(latestVersion),
			InstalledVersion: cleanVersion(feedToParse.Installed),
		}
		res.Objects = append(res.Objects, ms)
	}
	return res
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

func readConfig(filename string) (*viper.Viper, error) {
	v := viper.New()
	v.SetConfigName(filename)
	v.SetDefault("config_path", "./config")
	v.SetEnvPrefix("kvt")
	v.BindEnv("config_path")
	v.AddConfigPath(v.Get("config_path").(string))
	err := v.ReadInConfig()

	return v, err
}

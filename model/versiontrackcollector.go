package model

import (
	"github.com/prometheus/client_golang/prometheus"
)

//Define a struct for you collector that contains pointers
//to prometheus descriptors for each metric you wish to expose.
//Note you can also include fields of other types if they provide utility
//but we just won't be exposing them as metrics.
type versionCollector struct {
	versionMetric *prometheus.Desc
}

//You must create a constructor for you collector that
//initializes every descriptor and returns a pointer to the collector
func NewVersionCollector() *versionCollector {
	return &versionCollector{
		versionMetric: prometheus.NewDesc("invia_k8s_version_tracks",
			"Installed and Recent Versions for Applications",
			[]string{"latestVersion", "name", "installedVersion"}, nil,
		),
	}
}

//Each and every collector must implement the Describe function.
//It essentially writes all descriptors to the prometheus desc channel.
func (collector *versionCollector) Describe(ch chan<- *prometheus.Desc) {

	//Update this section with the each metric you create for a given collector
	ch <- collector.versionMetric
}

//Collect implements required collect function for all promehteus collectors
func (collector *versionCollector) Collect(ch chan<- prometheus.Metric) {

	//Implement logic here to determine proper metric value to return to prometheus
	//for each descriptor or call other functions that do so.
	var metricValue float64
	metricValue = 1
	var res PrometheusResult
	res = GetVersions()
	for _, r := range res.Objects {
		//helmInfo.With(prometheus.Labels{"latestVersion": r.LatestVersion, "name": r.Name, "installedVersion": r.InstalledVersion}).Add(1)
		ch <- prometheus.MustNewConstMetric(collector.versionMetric, prometheus.CounterValue, metricValue, r.LatestVersion, r.Name, r.InstalledVersion)
	}

	//Write latest value for each metric in the prometheus metric channel.
	//Note that you can pass CounterValue, GaugeValue, or UntypedValue types here.
	//ch <- prometheus.MustNewConstMetric(collector.versionMetric, prometheus.CounterValue, metricValue, "v2", "helmchart", "v1")

}

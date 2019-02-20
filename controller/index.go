package controller

import (
	//"fmt"
	//"log"

	"html/template"
	"net/http"

	"github.com/Sirupsen/logrus"
	//"k8s.io/api/core/v1"
	//"k8s.io/apimachinery/pkg/api/resource"
	//"k8s.io/apimachinery/pkg/watch"
	//"github.com/invia-de/K8VersionTrack/main"
	"github.com/invia-de/K8VersionTrack/model"
)

/*
func connectK8s() (*kubernetes.Clientset, error) {
	kubeconfig := filepath.Join(os.Getenv("HOME"), ".kube", "config")
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err.Error())
	}
	clientset, err := kubernetes.NewForConfig(config)
	return clientset, err
}


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

	var res model.PrometheusResult
	res = model.GetVersions()

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

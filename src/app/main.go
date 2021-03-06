package main

import (
	"go-cron-schedules/src/types"
	"html/template"
	"net/http"
	"os"

	discovery "github.com/gkarthiks/k8s-discovery"

	"strings"

	log "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	k8s                               *discovery.K8s
	namespace, scope, templateFilePath string
	avail, scopeAvail                 bool
	err                               error
)

func init() {
	log.SetFormatter(&log.JSONFormatter{TimestampFormat: "2006-01-02 15:04:05"})

	scope, scopeAvail = os.LookupEnv("SCOPE")
	if !scopeAvail {
		log.Infof("The listing will be on the namespace scope.")
	} else {
		log.Infof("The listing will be on the cluster scope.")
	}

	templateFilePath, avail = os.LookupEnv("TMPL_FILE_PATH")
	if !avail {
		log.Panicf("Template file is not available to serve")
	}

}

func main() {
	k8s, _ := discovery.NewK8s()
	if !scopeAvail {
		namespace, err = k8s.GetNamespace()
		if err != nil {
			log.Fatalf("Couldn't get the namespace")
		}
	}

	tmpl := template.Must(template.ParseFiles(templateFilePath))
	http.HandleFunc("/schedule", func(w http.ResponseWriter, r *http.Request) {
		log.Infof("Request URL %s ", r.URL)
		dataToServe := getCronJobsInTypesOnDemand(k8s)
		tmpl.Execute(w, dataToServe)
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Infof("Request URL %s ", r.URL)
		w.Write([]byte("OK"))
	})

	http.ListenAndServe(":8080", nil)
}

func getCronTabLinked(schedule string) string {
	return "https://crontab.guru/#" + strings.Replace(schedule, " ", "_", -1)
}

func getCronJobsInTypesOnDemand(k8s *discovery.K8s) interface{} {
	log.Infof("Collecting the cronJob list on %s namespace on demand", getScope(namespace))
	cronJobs, err := k8s.Clientset.BatchV1beta1().CronJobs(namespace).List(metav1.ListOptions{})
	if err != nil {
		log.Panic(err.Error())
	}
	var cronList []types.CronJob
	for idx, crons := range cronJobs.Items {
		cronData := types.CronJob{
			SNo:        idx + 1,
			Name:       crons.Name,
			Schedule:   crons.Spec.Schedule,
			LinkFormat: getCronTabLinked(crons.Spec.Schedule),
			Namespace:  crons.Namespace,
		}
		cronList = append(cronList, cronData)
	}
	log.Infof("Total cronjob collected: %d ", len(cronList))
	return types.ServingData{
		CronJobLists: cronList,
		Namespace:    getScope(namespace),
	}
}

func getScope(namespace string) string {
	if scopeAvail {
		return "all"
	} else {
		return namespace
	}
}

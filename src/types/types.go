package types

// CronJobList holds the list of CronJob
type ServingData struct {
	CronJobLists []CronJob
	Namespace string
}

//CronJob holds the name, schedule and a formatted link for cronjob
type CronJob struct {
	SNo        int
	Name       string
	Schedule   string
	LinkFormat string
}

package google

// Regions is predefine google regions for deploy cloud run and database
//
var Regions = map[string]string{
	"us": "us-central1",
	"jp": "asia-northeast1",
	"be": "europe-west1",
}

const (
	MasterProject    = "piyuo-master"
	TestProject      = "piyuo-master-test"
	TestEmail        = "piyuo-master@gmail.com"
	StableProject    = "piyuo-stable"
	DefaultTaskQueue = "task"
	DefaultRegion    = "us-central1"
)

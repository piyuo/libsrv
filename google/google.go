package google

// Regions is predefine google regions for deploy cloud run and database
//
var Regions = map[string]string{
	"us": "us-central1",
	"jp": "asia-northeast1",
	"be": "europe-west1",
}

const (
	MyGithub         = "git@github.com:piyuo"
	MyDomain         = "piyuo.com"
	MasterBranch     = "master"
	MasterProject    = "piyuo-master"
	TestBranch       = "test"
	TestProject      = "piyuo-master-test"
	TestAccount      = "piyuo-master-test@piyuo-master-test.iam.gserviceaccount.com"
	TestEmail        = "piyuo-master@gmail.com"
	StableBranch     = "stable"
	StableProject    = "piyuo-stable"
	DefaultTaskQueue = "task"
	DefaultRegion    = "us-central1"
)

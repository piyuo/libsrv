package gcloud

import (
	"context"
	"fmt"
	"time"

	cloudtasks "cloud.google.com/go/cloudtasks/apiv2"
	"github.com/piyuo/libsrv/src/gaccount"
	"github.com/pkg/errors"
	"google.golang.org/api/option"
	tasks "google.golang.org/genproto/googleapis/cloud/tasks/v2"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const here = "gcloud"

const defaultQueueID = "tasks"

const defaultLocationID = "us-central1"

// Regions is predefine google regions for deploy cloud run and database
//
var Regions = map[string]string{
	"us": "us-central1",
	"jp": "asia-northeast1",
	"be": "europe-west1",
}

// testMode is true should return success, false return error, otherwise behave normal
//
var testMode *bool

// TestModeAlwaySuccess will let every function success
//
func TestModeAlwaySuccess() {
	t := true
	testMode = &t
}

// TestModeAlwayFail will let every function fail
//
func TestModeAlwayFail() {
	f := false
	testMode = &f
}

// TestModeBackNormal stop test mode and back to normal
//
func TestModeBackNormal() {
	testMode = nil
}

// CreateHTTPTask create google cloud task, if scheduleTime is nil mean now
//
//	err = gcloud.CreateHTTPTask(ctx, url,body,nil)
//
func CreateHTTPTask(ctx context.Context, url string, body []byte, scheduleTime *timestamppb.Timestamp) error {
	if testMode != nil {
		if *testMode {
			return nil
		}
		return errors.New("failed always")
	}

	//gcloud won't allow context deadline over 30 seconds
	ctx, cancel := context.WithTimeout(ctx, time.Second*20) // cloud flare dns call must completed in 15 seconds
	defer cancel()

	cred, err := gaccount.GlobalCredential(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to create global credential")
	}

	client, err := cloudtasks.NewClient(ctx, option.WithCredentials(cred))
	if err != nil {
		return errors.Wrap(err, "failed to create cloud tasks client")
	}

	// Build the Task queue path.
	queuePath := fmt.Sprintf("projects/%s/locations/%s/queues/%s", cred.ProjectID, defaultLocationID, defaultQueueID)

	// Build the Task payload.
	req := &tasks.CreateTaskRequest{
		Parent: queuePath,
		Task: &tasks.Task{
			ScheduleTime: scheduleTime,
			MessageType: &tasks.Task_HttpRequest{
				HttpRequest: &tasks.HttpRequest{
					HttpMethod: tasks.HttpMethod_POST,
					Url:        url,
				},
			},
		},
	}

	// Add a payload message if one is present.
	if body != nil {
		req.Task.GetHttpRequest().Body = body
	}

	_, err = client.CreateTask(ctx, req)
	if err != nil {
		return errors.Wrap(err, "failed to create task")
	}
	return nil
}

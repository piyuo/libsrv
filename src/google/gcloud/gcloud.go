package gcloud

import (
	"context"
	"fmt"
	"strings"
	"time"

	cloudtasks "cloud.google.com/go/cloudtasks/apiv2"
	"github.com/piyuo/libsrv/src/gaccount"
	"github.com/piyuo/libsrv/src/identifier"
	"github.com/pkg/errors"
	"google.golang.org/api/option"
	tasks "google.golang.org/genproto/googleapis/cloud/tasks/v2"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const here = "gcloud"

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

// CreateHTTPTask create task in us-central1, if scheduleTime is nil mean now, default deadline is 10 mins.
//
//	err = gcloud.CreateHTTPTask(ctx,"my-queue", url,body,nil)
//
func CreateHTTPTask(ctx context.Context, queueID, url string, body []byte, scheduleTime *timestamppb.Timestamp) error {
	if testMode != nil {
		if *testMode {
			return nil
		}
		return errors.New("failed always")
	}

	//gcloud won't allow context deadline over 30 seconds
	ctx, cancel := context.WithTimeout(ctx, time.Second*20)
	defer cancel()

	cred, err := gaccount.GlobalCredential(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to create global credential")
	}

	client, err := cloudtasks.NewClient(ctx, option.WithCredentials(cred))
	if err != nil {
		return errors.Wrap(err, "failed to create cloud tasks client")
	}
	taskID := identifier.UUID()
	if strings.Contains(url, "?") {
		url += "&"
	} else {
		url += "?"
	}
	url += "TaskID=" + taskID

	// Build the Task queue path.
	queuePath := fmt.Sprintf("projects/%s/locations/%s/queues/%s", cred.ProjectID, defaultLocationID, queueID)

	// Build the Task payload.
	req := &tasks.CreateTaskRequest{
		Parent: queuePath,
		Task: &tasks.Task{
			ScheduleTime: scheduleTime,
			//			DispatchDeadline: durationpb.New(deadline), // use default deadline 10 mins.
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

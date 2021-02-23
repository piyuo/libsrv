package gcloud

import (
	"context"
	"fmt"
	"time"

	cloudtasks "cloud.google.com/go/cloudtasks/apiv2"
	"github.com/pkg/errors"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	tasks "google.golang.org/genproto/googleapis/cloud/tasks/v2"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const here = "gcloud"

// CreateHTTPTask create http task on google cloud task with schedule time
//
//	err = gcloud.CreateHTTPTask(ctx, "my-project","us-central1","my-queue",3*time.Seconds,url,body)
//
func CreateHTTPTask(ctx context.Context, cred *google.Credentials, projectID, locationID, queueID, url string, schedule time.Time, body []byte) error {

	client, err := cloudtasks.NewClient(ctx, option.WithCredentials(cred))
	if err != nil {
		return errors.Wrap(err, "failed to create cloud tasks client")
	}

	// Build the Task queue path.
	queuePath := fmt.Sprintf("projects/%s/locations/%s/queues/%s", projectID, locationID, queueID)

	// Build the Task payload.
	req := &tasks.CreateTaskRequest{
		Parent: queuePath,
		Task: &tasks.Task{
			MessageType: &tasks.Task_HttpRequest{
				HttpRequest: &tasks.HttpRequest{
					HttpMethod: tasks.HttpMethod_POST,
					Url:        url,
				},
			},
			ScheduleTime: &timestamppb.Timestamp{Seconds: schedule.Unix()},
		},
	}

	// Add a payload message if one is present.
	req.Task.GetHttpRequest().Body = body

	_, err = client.CreateTask(ctx, req)
	if err != nil {
		return errors.Wrap(err, "failed to create task")
	}
	return nil
}

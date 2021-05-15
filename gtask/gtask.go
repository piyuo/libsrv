package gtask

import (
	"context"
	"fmt"
	"strings"
	"time"

	cloudtasks "cloud.google.com/go/cloudtasks/apiv2"
	"github.com/piyuo/libsrv/db"
	"github.com/piyuo/libsrv/gaccount"
	"github.com/piyuo/libsrv/gdb"
	"github.com/piyuo/libsrv/identifier"
	"github.com/pkg/errors"
	"google.golang.org/api/option"
	tasks "google.golang.org/genproto/googleapis/cloud/tasks/v2"
)

const defaultLocationID = "us-central1"

// Mock define key test flag
//
type Mock int8

const (
	// MockSuccess let function return nil
	//
	MockSuccess Mock = iota

	// MockError let function error
	//
	MockError
)

// New task in us-central1, if scheduleTime is nil mean now, default deadline is 10 mins. return task id if success
//
//	taskID, err = New(ctx,"my-queue", url,body,"my-task", 1800, 3)
//
func New(ctx context.Context, queueID, url string, body []byte, name string, duration, maxRetry int) (string, error) {
	if ctx.Value(MockSuccess) != nil {
		return "", nil
	}
	if ctx.Value(MockError) != nil {
		return "", errors.New("")
	}

	//gcloud won't allow context deadline over 30 seconds
	ctx, cancel := context.WithTimeout(ctx, time.Second*20)
	defer cancel()

	cred, err := gaccount.GlobalCredential(ctx)
	if err != nil {
		return "", errors.Wrap(err, "new global cred")
	}
	taskID := identifier.UUID()

	// create cloud task
	taskClient, err := cloudtasks.NewClient(ctx, option.WithCredentials(cred))
	if err != nil {
		return "", errors.Wrap(err, "new tasks client")
	}
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

	// not-exist use for test
	if !strings.Contains(url, "not-exists") {
		_, err = taskClient.CreateTask(ctx, req)
		if err != nil {
			return "", errors.Wrap(err, "create google task")
		}
	}

	// create task in database
	client, err := gdb.NewClient(ctx, cred)
	if err != nil {
		return "", errors.Wrap(err, "create task db client")
	}
	defer client.Close()

	task := newTask(taskID, name, duration, maxRetry)
	err = client.Set(ctx, task)
	if err != nil {
		return "", errors.Wrap(err, "set db task "+taskID)
	}

	return taskID, nil
}

// Lock task
//
//	ok, err := Lock(ctx, "task id")
//
func Lock(ctx context.Context, taskID string) error {
	if ctx.Value(MockSuccess) != nil {
		return nil
	}
	if ctx.Value(MockError) != nil {
		return errors.New("")
	}

	cred, err := gaccount.GlobalCredential(ctx)
	if err != nil {
		return errors.Wrap(err, "new global cred")
	}
	client, err := gdb.NewClient(ctx, cred)
	if err != nil {
		return errors.Wrap(err, "new db client")
	}
	defer client.Close()

	return client.Transaction(ctx, func(ctx context.Context, tx db.Transaction) error {
		obj, err := tx.Get(ctx, &Task{}, taskID)
		if err != nil {
			return errors.Wrap(err, "get task "+taskID)
		}
		if obj == nil {
			return errors.New("task not found " + taskID)
		}

		task := obj.(*Task)
		err = task.Lock()
		if err != nil {
			return err
		}

		err = tx.Set(ctx, task)
		if err != nil {
			return errors.Wrap(err, "set task "+taskID)
		}
		return nil
	})
}

// Delete task
//
//	err := Delete(ctx, "task id")
//
func Delete(ctx context.Context, taskID string) error {
	if ctx.Value(MockSuccess) != nil {
		return nil
	}
	if ctx.Value(MockError) != nil {
		return errors.New("")
	}

	cred, err := gaccount.GlobalCredential(ctx)
	if err != nil {
		return errors.Wrap(err, "new global cred")
	}
	client, err := gdb.NewClient(ctx, cred)
	if err != nil {
		return errors.Wrap(err, "new db client")
	}
	defer client.Close()

	task := &Task{}
	task.SetID(taskID)
	err = client.Delete(ctx, task)
	if err != nil {
		return errors.Wrap(err, "delete task "+taskID)
	}
	return nil
}

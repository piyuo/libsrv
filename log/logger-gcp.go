package log

import (
	"context"
	"fmt"

	"cloud.google.com/go/logging"
	gcp "github.com/piyuo/libsrv/gcp"
	"github.com/pkg/errors"
	"google.golang.org/api/option"
	mrpb "google.golang.org/genproto/googleapis/api/monitoredres"
)

type gcpLogger struct {
	Logger
	// client is google cloud platform logging client
	//
	client *logging.Client

	// logger is google cloud platform logging client
	//
	logger *logging.Logger
}

// NewGCPLogger return logger implement by google cloud platform
//
//	ctx := context.Background()
//	logger, err := NewGCPLogger(ctx)
//
func NewGCPLogger(ctx context.Context) (Logger, error) {
	cred, err := gcp.GlobalCredential(ctx)
	if err != nil {
		return nil, err
	}
	client, err := logging.NewClient(ctx, cred.ProjectID, option.WithCredentials(cred))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create stackdriver logging client")
	}
	logger := client.Logger(appName)

	return &gcpLogger{
		client: client,
		logger: logger,
	}, nil
}

// Close logger
//
//	defer logger.Close()
//
func (c *gcpLogger) Close() {
	if err := c.client.Close(); err != nil {
		fmt.Printf("!!! %v\n", err)
	}
}

// Write message and level to google cloud platform log
//
//	logger.write(ctx,"hi","app",DEBUG)
//
func (c *gcpLogger) Write(ctx context.Context, level Level, where, message string) {
	if message == "" {
		return
	}
	header, id := getHeader(ctx, where)
	fmt.Printf("%v%v (logged)\n", header, message)

	severity := logging.Info
	switch level {
	case WARNING:
		severity = logging.Warning
	case ALERT:
		severity = logging.Critical
	case DEBUG:
		severity = logging.Debug
	}

	entry := logging.Entry{
		Payload: header + message,
		Resource: &mrpb.MonitoredResource{
			Type: "project",
		},
		Severity: severity,
		Labels: map[string]string{
			"app": appName,
			"at":  where,
		},
	}
	if id != "" {
		entry.Labels["id"] = id
	}
	c.logger.Log(entry)
}

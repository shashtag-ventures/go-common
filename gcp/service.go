package gcp

import (
	"context"

	cloudbuild "cloud.google.com/go/cloudbuild/apiv1/v2"
	"cloud.google.com/go/logging/logadmin"
	run "cloud.google.com/go/run/apiv2"
	"cloud.google.com/go/storage"
)

type Service interface {
	// UploadSource uploads the source code zip to GCS and returns the object name/URI
	UploadSource(ctx context.Context, projectID string, zipContent []byte) (string, error)

	// TriggerBuild triggers a Cloud Build to build the container image
	// Returns the Build ID and the Image URI
	// Returns: buildID, logURL, imageURI, error
	TriggerBuild(ctx context.Context, karadaProjectID string, gcsSourceURI string, imageName string, onStart func(buildID, logURL string)) (string, string, string, error)

	// DeployService deploys the container image to Cloud Run
	// Returns the Service URL
	DeployService(ctx context.Context, projectID string, imageName string, serviceName string) (string, error)

	// GetBuildLogs fetches text logs for a build
	GetBuildLogs(ctx context.Context, buildID string) (string, error)

	// ListServices lists all Cloud Run services in the project and region
	ListServices(ctx context.Context) ([]string, error)

	// DeleteService deletes a Cloud Run service by name
	DeleteService(ctx context.Context, serviceName string) error

	// Close closes the clients
	Close() error
}

type gcpService struct {
	storageClient    *storage.Client
	buildClient      *cloudbuild.Client
	runClient        *run.ServicesClient
	projectID        string
	region           string
	bucketName       string
	artifactRegistry string
	logAdminClient   *logadmin.Client
}

func NewService(ctx context.Context, projectID, region, bucketName, artifactRegistry string) (Service, error) {
	sc, err := storage.NewClient(ctx)
	if err != nil {
		return nil, err
	}

	bc, err := cloudbuild.NewClient(ctx)
	if err != nil {
		return nil, err
	}

	rc, err := run.NewServicesClient(ctx)
	if err != nil {
		return nil, err
	}

	lac, err := logadmin.NewClient(ctx, projectID)
	if err != nil {
		return nil, err
	}

	return &gcpService{
		storageClient:    sc,
		buildClient:      bc,
		runClient:        rc,
		logAdminClient:   lac,
		projectID:        projectID,
		region:           region,
		bucketName:       bucketName,
		artifactRegistry: artifactRegistry,
	}, nil
}

func (s *gcpService) Close() error {
	// Best effort close
	_ = s.storageClient.Close()
	_ = s.buildClient.Close()
	_ = s.runClient.Close()
	_ = s.logAdminClient.Close()
	return nil
}

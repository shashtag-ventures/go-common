package gcp

import (
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	cloudbuildpb "cloud.google.com/go/cloudbuild/apiv1/v2/cloudbuildpb"
	"cloud.google.com/go/logging/logadmin"
	"google.golang.org/api/iterator"
)

func (s *gcpService) TriggerBuild(ctx context.Context, karadaProjectID string, gcsSourceURI string, imageName string, onStart func(buildID, logURL string)) (string, string, string, error) {
	// Apply a timeout to prevent the goroutine from hanging indefinitely
	ctx, cancel := context.WithTimeout(ctx, 5*time.Minute)
	defer cancel()

	// gcsSourceURI format: gs://bucket/path/to/object
	parts := strings.SplitN(strings.TrimPrefix(gcsSourceURI, "gs://"), "/", 2)
	if len(parts) != 2 {
		return "", "", "", fmt.Errorf("invalid gcs uri: %s", gcsSourceURI)
	}
	bucket, object := parts[0], parts[1]

	fullImageName := fmt.Sprintf("%s/%s", s.artifactRegistry, imageName)

	req := &cloudbuildpb.CreateBuildRequest{
		ProjectId: s.projectID,
		Build: &cloudbuildpb.Build{
			Source: &cloudbuildpb.Source{
				Source: &cloudbuildpb.Source_StorageSource{
					StorageSource: &cloudbuildpb.StorageSource{
						Bucket: bucket,
						Object: object,
					},
				},
			},
			Steps: []*cloudbuildpb.BuildStep{
				{
					Name: "gcr.io/cloud-builders/docker",
					Args: []string{"build", "-t", fullImageName, "."},
				},
				{
					Name: "gcr.io/cloud-builders/docker",
					Args: []string{"push", fullImageName},
				},
			},
			Images: []string{fullImageName},
		},
	}

	op, err := s.buildClient.CreateBuild(ctx, req)
	if err != nil {
		return "", "", "", fmt.Errorf("failed to create build: %w", err)
	}

	// Extract Build ID and Log URL from Metadata immediately
	var buildID, logURL string
	if meta, err := op.Metadata(); err == nil && meta.Build != nil {
		buildID = meta.Build.Id
		logURL = meta.Build.LogUrl
		if onStart != nil {
			onStart(buildID, logURL)
		}
	}

	// Wait for the build to complete
	resp, err := op.Wait(ctx)
	if err != nil {
		// Return available info even on failure
		return buildID, logURL, "", fmt.Errorf("build failed: %w", err)
	}

	return resp.Id, resp.LogUrl, fullImageName, nil
}

func (s *gcpService) GetBuildLogs(ctx context.Context, buildID string) (string, error) {
	// Apply a timeout to prevent hanging
	ctx, cancel := context.WithTimeout(ctx, 2*time.Minute)
	defer cancel()

	// 1. Get Build to check status and find logs bucket
	build, err := s.buildClient.GetBuild(ctx, &cloudbuildpb.GetBuildRequest{
		ProjectId: s.projectID,
		Id:        buildID,
	})
	if err != nil {
		return "", fmt.Errorf("failed to get build: %w", err)
	}

	// 2. If build is WORKING or QUEUED, fetch logs from Cloud Logging
	if build.Status == cloudbuildpb.Build_WORKING || build.Status == cloudbuildpb.Build_QUEUED {
		// Filter for logs related to this build
		filter := fmt.Sprintf(`resource.type="build" AND resource.labels.build_id="%s"`, buildID)

		// Use admin client to query logs
		// We want the most recent logs, but for simplicity let's just get the last 500 entries
		// In a real streaming scenario, we might use a cursor or timestamps.
		// For this MVP, fetching all/recent logs is okay as builds aren't huge.

		var logsBuilder strings.Builder
		it := s.logAdminClient.Entries(ctx, logadmin.Filter(filter), logadmin.NewestFirst())

		// Limit to prevent OOM on huge logs, say 1000 lines
		const maxEntries = 2000
		entries := make([]string, 0, maxEntries)

		for i := 0; i < maxEntries; i++ {
			entry, err := it.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				// If we fail to get live logs, we might just fall back or return what we have
				// But usually this means we can't read logs.
				// Let's just log error and continue to GCS fallback if possible, or return partial.
				// For now, return error to see what happens.
				return "", fmt.Errorf("failed to read live logs: %w", err)
			}

			// Payload can be string or struct/map.
			// Cloud Build logs are usually textPayload.
			if text, ok := entry.Payload.(string); ok {
				entries = append(entries, text)
			}
		}

		// Entries are NewestFirst, so reverse them to show in order
		for i := len(entries) - 1; i >= 0; i-- {
			logsBuilder.WriteString(entries[i])
			logsBuilder.WriteString("\n")
		}

		return logsBuilder.String(), nil
	}

	// 3. If build is finished (SUCCESS, FAILURE, CANCELLED), read from GCS
	// Parse Logs Bucket (gs://[bucket]/log-[id].txt)

	if build.LogsBucket == "" {
		return "", fmt.Errorf("no logs bucket found for build %s", buildID)
	}

	bucketName := strings.TrimPrefix(build.LogsBucket, "gs://")
	objectName := fmt.Sprintf("log-%s.txt", buildID)

	// 4. Read from GCS
	rc, err := s.storageClient.Bucket(bucketName).Object(objectName).NewReader(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to open log file: %w", err)
	}
	defer rc.Close()

	// Read content
	content, err := io.ReadAll(rc)
	if err != nil {
		return "", fmt.Errorf("failed to read logs: %w", err)
	}

	return string(content), nil
}

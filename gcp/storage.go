package gcp

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

func (s *gcpService) UploadSource(ctx context.Context, targetProjectID string, zipContent []byte) (string, error) {
	filename := fmt.Sprintf("source/%s/%s-%d.zip", targetProjectID, uuid.New().String(), time.Now().Unix())

	bucket := s.storageClient.Bucket(s.bucketName)
	obj := bucket.Object(filename)

	w := obj.NewWriter(ctx)
	if _, err := w.Write(zipContent); err != nil {
		return "", fmt.Errorf("failed to write zip to bucket: %w", err)
	}
	if err := w.Close(); err != nil {
		return "", fmt.Errorf("failed to close bucket writer: %w", err)
	}

	// Return format: gs://bucket/filename
	return fmt.Sprintf("gs://%s/%s", s.bucketName, filename), nil
}

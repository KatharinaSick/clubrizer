package storage

import (
	"context"
	"fmt"
	"io"

	gcs "cloud.google.com/go/storage"
	"github.com/google/uuid"
)

type Client struct {
	bucketName string
	gcs        *gcs.Client
}

// NewClient creates a GCS client using Application Default Credentials.
// Locally: run `gcloud auth application-default login`.
// On Cloud Run: the attached service account is used automatically.
func NewClient(ctx context.Context, bucketName string) (*Client, error) {
	gcsClient, err := gcs.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCS client: %w", err)
	}
	return &Client{bucketName: bucketName, gcs: gcsClient}, nil
}

// UploadProfilePicture uploads an image to GCS with a UUID-based object name and
// sets a public read ACL on the object. Returns the public URL.
// Note: per-object ACLs require uniform bucket-level access to be disabled on the bucket.
func (c *Client) UploadProfilePicture(ctx context.Context, contentType string, data io.Reader) (string, error) {
	objectName := uuid.New().String()

	obj := c.gcs.Bucket(c.bucketName).Object(objectName)
	wc := obj.NewWriter(ctx)
	wc.ContentType = contentType

	if _, err := io.Copy(wc, data); err != nil {
		return "", fmt.Errorf("failed to write to GCS: %w", err)
	}

	if err := wc.Close(); err != nil {
		return "", fmt.Errorf("failed to close GCS writer: %w", err)
	}

	if err := obj.ACL().Set(ctx, gcs.AllUsers, gcs.RoleReader); err != nil {
		return "", fmt.Errorf("failed to set public ACL on GCS object: %w", err)
	}

	return fmt.Sprintf("https://storage.googleapis.com/%s/%s", c.bucketName, objectName), nil
}

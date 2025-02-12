// client.go
package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
)

// Client wraps HTTP calls to an Archivista service.
type Client struct {
	baseURL    string
	httpClient *http.Client
}

// NewClient creates a new Archivista API client.
func NewClient(baseURL string) (*Client, error) {
	// Validate baseURL.
	_, err := url.ParseRequestURI(baseURL)
	if err != nil {
		return nil, err
	}
	return &Client{
		baseURL:    baseURL,
		httpClient: &http.Client{},
	}, nil
}

// UploadRequest describes a request to upload an artifact.
type UploadRequest struct {
	ArtifactType string `json:"artifactType"`
	Payload      []byte `json:"payload"`
	Signature    []byte `json:"signature"`
	KeyID        string `json:"keyID"`
}

// Upload sends an upload request to Archivista.
func (c *Client) Upload(ctx context.Context, req *UploadRequest) (*UploadResponse, error) {
	uploadURL, err := url.JoinPath(c.baseURL, "upload")
	if err != nil {
		return nil, err
	}
	bodyBytes, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, uploadURL, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(string(respBytes))
	}
	var uploadResp UploadResponse
	if err := json.Unmarshal(respBytes, &uploadResp); err != nil {
		return nil, err
	}
	return &uploadResp, nil
}

// Artifact represents a retrieved artifact from Archivista.
type Artifact struct {
	Payload   []byte `json:"payload"`
	Signature []byte `json:"signature"`
}

// GetArtifact retrieves an artifact by key.
func (c *Client) GetArtifact(ctx context.Context, key string) (*Artifact, error) {
	downloadURL, err := url.JoinPath(c.baseURL, "download", key)
	if err != nil {
		return nil, err
	}
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, downloadURL, nil)
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		respBytes, _ := io.ReadAll(resp.Body)
		return nil, errors.New(string(respBytes))
	}
	var artifact Artifact
	if err := json.NewDecoder(resp.Body).Decode(&artifact); err != nil {
		return nil, err
	}
	return &artifact, nil
}

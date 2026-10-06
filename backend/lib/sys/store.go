// Package sys holds the validator's system boundaries: where input files are
// read from and where the service writes its own output.
//
// S3 is disabled by default, and every input comes from the local
// filesystem and is synthetic. Real Medicare 835s are protected health
// information and live in
// S3 inside the CMS environment; the S3 reader exists for that deployment and
// is used only when the service or vdb is started with
// APP_SERVER_S3_ENABLED=true.
//
// S3 is read-only by design. The S3 type lists and reads objects and has no
// write, copy, or delete method, so nothing in this program can change or
// remove an input. Output (difference files, the database) never goes to S3.
package sys

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/credentials/ec2rolecreds"
	"github.com/aws/aws-sdk-go-v2/credentials/endpointcreds"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

// Reader reads input files. Names are relative to the filesystem root or the
// configured S3 prefix.
type Reader interface {
	Read(ctx context.Context, name string) ([]byte, error)
}

// FileSystem reads inputs and writes the service's own output under one root.
type FileSystem struct {
	Root *os.Root
}

func (s FileSystem) Read(ctx context.Context, name string) ([]byte, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	data, err := s.Root.ReadFile(name)
	if err != nil {
		return nil, fmt.Errorf("sys: reading file: %w", err)
	}
	return data, ctx.Err()
}

func (s FileSystem) Write(ctx context.Context, name string, data []byte) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if err := s.Root.MkdirAll(filepath.Dir(name), 0700); err != nil {
		return fmt.Errorf("sys: creating output directory: %w", err)
	}
	if err := s.Root.WriteFile(name, data, 0600); err != nil {
		return fmt.Errorf("sys: writing file: %w", err)
	}
	return nil
}

// S3Config is every S3 setting. The web service and vdb fill it from their
// APP_SERVER_S3_* configuration and, for task-role, from the four
// AWS_CONTAINER_* variables ECS sets in a task's environment. Each program
// reads those in the first lines of its run function and hands them over
// here, so the SDK never reads the environment itself.
type S3Config struct {
	Bucket          string
	Prefix          string
	Region          string
	Credentials     string // "static", "instance-role", or "task-role"
	AccessKeyID     string // static only
	SecretAccessKey string // static only
	SessionToken    string // static only; optional

	// task-role only: the container credentials endpoint and its
	// authorization token, as ECS sets them in a task's environment.
	ContainerCredentialsRelativeURI string // AWS_CONTAINER_CREDENTIALS_RELATIVE_URI: a path on the ECS endpoint
	ContainerCredentialsFullURI     string // AWS_CONTAINER_CREDENTIALS_FULL_URI: a whole URL; used when no relative URI is set
	ContainerAuthorizationToken     string // AWS_CONTAINER_AUTHORIZATION_TOKEN: optional
	ContainerAuthorizationTokenFile string // AWS_CONTAINER_AUTHORIZATION_TOKEN_FILE: optional; read once at startup, and used over the token value
}

// OpenS3 checks the settings, builds the client, and proves the credentials
// with one list request for at most one key under the prefix. The read-only
// policy (docs/delivery/iam-s3-read-only.json) allows listing only under the
// prefix, so HeadBucket, which it authorizes as a list without a prefix,
// would be denied. The client is built only from cfg: no AWS environment
// variables, shared config files, or default credential chain.
// instance-role reads credentials from the EC2 instance metadata service.
// task-role reads them from the container credentials endpoint ECS gives a
// task with a task role (the Container* fields); a wrong role fails the
// list request here, at startup.
func OpenS3(ctx context.Context, cfg S3Config) (S3, error) {
	if cfg.Bucket == "" {
		return S3{}, fmt.Errorf("sys: configuring S3: bucket is required when S3 is enabled")
	}
	if cfg.Region == "" {
		return S3{}, fmt.Errorf("sys: configuring S3: region is required")
	}
	var provider aws.CredentialsProvider
	switch cfg.Credentials {
	case "static":
		if cfg.AccessKeyID == "" || cfg.SecretAccessKey == "" {
			return S3{}, fmt.Errorf("sys: configuring S3: static credentials need an access key ID and secret")
		}
		provider = credentials.NewStaticCredentialsProvider(cfg.AccessKeyID, cfg.SecretAccessKey, cfg.SessionToken)
	case "instance-role":
		provider = aws.NewCredentialsCache(ec2rolecreds.New())
	case "task-role":
		var err error
		provider, err = taskRoleCredentials(cfg)
		if err != nil {
			return S3{}, err
		}
	default:
		return S3{}, fmt.Errorf("sys: configuring S3: credentials must be static, instance-role, or task-role")
	}
	client := s3.New(s3.Options{Region: cfg.Region, Credentials: provider})
	checkPrefix := keyPrefix(cfg.Prefix, ".")
	maxKeys := int32(1)
	_, err := client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{Bucket: &cfg.Bucket, Prefix: &checkPrefix, MaxKeys: &maxKeys})
	if err != nil {
		return S3{}, fmt.Errorf("sys: checking S3 access: %w", err)
	}
	return S3{Client: client, Bucket: cfg.Bucket, Prefix: cfg.Prefix}, nil
}

// ecsCredentialsEndpoint is where ECS answers a task's credential requests;
// AWS_CONTAINER_CREDENTIALS_RELATIVE_URI is a path on it.
const ecsCredentialsEndpoint = "http://169.254.170.2"

// taskRoleCredentials builds the task-role provider from the Container*
// fields: the SDK's endpoint credentials provider, fetching from
// taskRoleEndpoint and sending the authorization token, if there is one, as
// the Authorization header. When a token file is set, its whole contents
// are the token, read once here, and the token value is not used. The
// provider is wrapped in the SDK's cache, as instance-role is, so a request
// fetches credentials only when the cached ones have expired.
func taskRoleCredentials(cfg S3Config) (aws.CredentialsProvider, error) {
	endpoint, err := taskRoleEndpoint(cfg)
	if err != nil {
		return nil, err
	}
	token := cfg.ContainerAuthorizationToken
	if cfg.ContainerAuthorizationTokenFile != "" {
		data, err := os.ReadFile(cfg.ContainerAuthorizationTokenFile)
		if err != nil {
			return nil, fmt.Errorf("sys: configuring S3: reading AWS_CONTAINER_AUTHORIZATION_TOKEN_FILE: %w", err)
		}
		token = string(data)
	}
	return aws.NewCredentialsCache(endpointcreds.New(endpoint, func(o *endpointcreds.Options) {
		o.AuthorizationToken = token
	})), nil
}

// taskRoleEndpoint is the URL task-role fetches credentials from: the ECS
// endpoint plus the relative URI or, when no relative URI is set, the full
// URI as given. Neither set is a configuration error naming both variables.
func taskRoleEndpoint(cfg S3Config) (string, error) {
	if cfg.ContainerCredentialsRelativeURI != "" {
		return ecsCredentialsEndpoint + cfg.ContainerCredentialsRelativeURI, nil
	}
	if cfg.ContainerCredentialsFullURI != "" {
		return cfg.ContainerCredentialsFullURI, nil
	}
	return "", fmt.Errorf("sys: configuring S3: task-role needs AWS_CONTAINER_CREDENTIALS_RELATIVE_URI or AWS_CONTAINER_CREDENTIALS_FULL_URI in the environment (ECS sets one for a task that has a task role)")
}

// keyPrefix is the S3 key prefix for dir under the configured prefix,
// ending in "/", or empty for the whole bucket.
func keyPrefix(prefix, dir string) string {
	p := path.Join(prefix, dir) + "/"
	if p == "./" {
		return ""
	}
	return p
}

// S3 lists and reads input objects. It is disabled unless
// APP_SERVER_S3_ENABLED=true, and it deliberately has no method that writes,
// copies, or deletes.
type S3 struct {
	Client *s3.Client
	Bucket string
	Prefix string
}

func (s S3) Read(ctx context.Context, name string) ([]byte, error) {
	if !filepath.IsLocal(name) {
		return nil, fmt.Errorf("sys: reading object: name must be relative")
	}
	key := path.Join(s.Prefix, name)
	object, err := s.Client.GetObject(ctx, &s3.GetObjectInput{Bucket: &s.Bucket, Key: &key})
	if err != nil {
		var missing *types.NoSuchKey
		if errors.As(err, &missing) {
			err = errors.Join(fs.ErrNotExist, err)
		}
		return nil, fmt.Errorf("sys: reading object: %w", err)
	}
	data, readErr := io.ReadAll(object.Body)
	closeErr := object.Body.Close()
	if readErr != nil {
		return nil, fmt.Errorf("sys: reading object body: %w", readErr)
	}
	if closeErr != nil {
		return nil, fmt.Errorf("sys: closing object body: %w", closeErr)
	}
	return data, nil
}

// List returns the names of the objects directly under dir, relative to dir
// and sorted, as os.ReadDir does for a directory: objects in deeper
// "subdirectories", directory markers, and names starting with a dot are
// left out. dir is relative to the configured prefix; "." lists the prefix
// itself. Each name can be read with Read(path.Join(dir, name)).
func (s S3) List(ctx context.Context, dir string) ([]string, error) {
	if !filepath.IsLocal(dir) {
		return nil, fmt.Errorf("sys: listing objects: directory must be relative")
	}
	listPrefix := keyPrefix(s.Prefix, dir)
	delimiter := "/"
	input := s3.ListObjectsV2Input{Bucket: &s.Bucket, Prefix: &listPrefix, Delimiter: &delimiter}
	var names []string
	for {
		page, err := s.Client.ListObjectsV2(ctx, &input)
		if err != nil {
			return nil, fmt.Errorf("sys: listing objects: %w", err)
		}
		for _, object := range page.Contents {
			name := strings.TrimPrefix(aws.ToString(object.Key), listPrefix)
			if name == "" || strings.HasPrefix(name, ".") {
				continue
			}
			names = append(names, name)
		}
		if !aws.ToBool(page.IsTruncated) {
			break
		}
		input.ContinuationToken = page.NextContinuationToken
	}
	sort.Strings(names)
	return names, nil
}

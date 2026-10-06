package sys

import (
	"context"
	"encoding/xml"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"sync"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// fakeS3 answers ListObjectsV2 and GetObject for one bucket from a map of
// keys, two keys per list page so paging is exercised, and records the
// method of every request.
type fakeS3 struct {
	objects map[string]string
	mu      sync.Mutex
	methods []string
}

type listResult struct {
	XMLName               xml.Name `xml:"ListBucketResult"`
	IsTruncated           bool
	NextContinuationToken string   `xml:",omitempty"`
	Contents              []object `xml:"Contents"`
	CommonPrefixes        []prefix `xml:"CommonPrefixes"`
}

type object struct {
	Key  string
	Size int
}

type prefix struct {
	Prefix string
}

func (f *fakeS3) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	f.mu.Lock()
	f.methods = append(f.methods, r.Method)
	f.mu.Unlock()

	key := strings.TrimPrefix(r.URL.Path, "/bucket/")
	if r.URL.Query().Get("list-type") != "2" {
		body, ok := f.objects[key]
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`<Error><Code>NoSuchKey</Code></Error>`))
			return
		}
		w.Write([]byte(body))
		return
	}

	listPrefix := r.URL.Query().Get("prefix")
	delimiter := r.URL.Query().Get("delimiter")
	var keys []string
	common := map[string]bool{}
	for k := range f.objects {
		if !strings.HasPrefix(k, listPrefix) {
			continue
		}
		rest := strings.TrimPrefix(k, listPrefix)
		if i := strings.Index(rest, delimiter); delimiter != "" && i >= 0 && i < len(rest)-1 {
			common[listPrefix+rest[:i+1]] = true
			continue
		}
		keys = append(keys, k)
	}
	sort.Strings(keys)
	start, _ := strconv.Atoi(r.URL.Query().Get("continuation-token"))
	end := min(start+2, len(keys))
	result := listResult{}
	for _, k := range keys[start:end] {
		result.Contents = append(result.Contents, object{Key: k, Size: len(f.objects[k])})
	}
	if end < len(keys) {
		result.IsTruncated = true
		result.NextContinuationToken = strconv.Itoa(end)
	}
	for p := range common {
		result.CommonPrefixes = append(result.CommonPrefixes, prefix{Prefix: p})
	}
	w.Header().Set("Content-Type", "application/xml")
	xml.NewEncoder(w).Encode(result)
}

// TestS3ListAndRead lists a submission's objects under the configured prefix
// across list pages, reads one back, and checks that only reads were sent.
func TestS3ListAndRead(t *testing.T) {
	fake := &fakeS3{objects: map[string]string{
		"inputs/vendor-1/2026-10-02/":                  "",
		"inputs/vendor-1/2026-10-02/b.835":             "B",
		"inputs/vendor-1/2026-10-02/a.835":             "A",
		"inputs/vendor-1/2026-10-02/c.835":             "C",
		"inputs/vendor-1/2026-10-02/.hidden":           "H",
		"inputs/vendor-1/2026-10-02/deeper/d.835":      "D",
		"inputs/vendor-1/2026-10-02-late/e.835":        "E",
		"other/vendor-1/2026-10-02/f.835":              "F",
		"inputs/vendor-1/2026-10-02/deeper/nested.835": "N",
	}}
	server := httptest.NewServer(fake)
	defer server.Close()
	client := s3.New(s3.Options{
		Region:       "us-east-1",
		Credentials:  credentials.NewStaticCredentialsProvider("id", "secret", ""),
		BaseEndpoint: aws.String(server.URL),
		UsePathStyle: true,
	})
	store := S3{Client: client, Bucket: "bucket", Prefix: "inputs"}
	ctx := context.Background()

	names, err := store.List(ctx, "vendor-1/2026-10-02")
	if err != nil {
		t.Fatal(err)
	}
	if want := []string{"a.835", "b.835", "c.835"}; !reflect.DeepEqual(names, want) {
		t.Errorf("List(vendor-1/2026-10-02) = %q, want %q", names, want)
	}

	data, err := store.Read(ctx, path.Join("vendor-1/2026-10-02", names[1]))
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "B" {
		t.Errorf("Read(b.835) = %q, want %q", data, "B")
	}

	empty, err := store.List(ctx, "vendor-2")
	if err != nil || len(empty) != 0 {
		t.Errorf("List(vendor-2) = %q, %v, want no names and no error", empty, err)
	}

	for _, dir := range []string{"../other", "/inputs"} {
		if _, err := store.List(ctx, dir); err == nil {
			t.Errorf("List(%q) error = nil, want an error for a non-relative directory", dir)
		}
	}

	for _, method := range fake.methods {
		if method != http.MethodGet && method != http.MethodHead {
			t.Errorf("request method %s, want only GET or HEAD", method)
		}
	}
}

// TestTaskRoleEndpoint checks the URL task-role fetches credentials from:
// the relative URI on the ECS endpoint, the full URI as given, the relative
// URI first when both are set, and an error naming the variables when
// neither is set.
func TestTaskRoleEndpoint(t *testing.T) {
	tests := []struct {
		name    string
		cfg     S3Config
		want    string
		wantErr bool
	}{
		{"relative URI", S3Config{ContainerCredentialsRelativeURI: "/v2/credentials/0123"}, "http://169.254.170.2/v2/credentials/0123", false},
		{"full URI", S3Config{ContainerCredentialsFullURI: "http://localhost:8080/credentials"}, "http://localhost:8080/credentials", false},
		{"both", S3Config{ContainerCredentialsRelativeURI: "/v2/credentials/0123", ContainerCredentialsFullURI: "http://localhost:8080/credentials"}, "http://169.254.170.2/v2/credentials/0123", false},
		{"neither", S3Config{}, "", true},
	}
	for _, tt := range tests {
		got, err := taskRoleEndpoint(tt.cfg)
		if got != tt.want || (err != nil) != tt.wantErr {
			t.Errorf("%s: taskRoleEndpoint() = %q, %v, want %q, error %v", tt.name, got, err, tt.want, tt.wantErr)
		}
		if err != nil && !strings.Contains(err.Error(), "AWS_CONTAINER_CREDENTIALS_RELATIVE_URI") {
			t.Errorf("%s: taskRoleEndpoint() error = %v, want one naming AWS_CONTAINER_CREDENTIALS_RELATIVE_URI", tt.name, err)
		}
	}
}

// TestOpenS3TaskRoleNeedsEndpoint checks that task-role with none of the
// container endpoint variables set fails before any request is made.
func TestOpenS3TaskRoleNeedsEndpoint(t *testing.T) {
	_, err := OpenS3(context.Background(), S3Config{Bucket: "bucket", Region: "us-east-1", Credentials: "task-role"})
	if err == nil || !strings.Contains(err.Error(), "AWS_CONTAINER_CREDENTIALS_FULL_URI") {
		t.Errorf("OpenS3(task-role, no endpoint) error = %v, want one naming the endpoint variables", err)
	}
}

// TestTaskRoleCredentials fetches credentials through the task-role
// provider from a local stand-in for the container credentials endpoint,
// given as the full URI, and checks that the token file's contents, not
// the token value, went out as the Authorization header.
func TestTaskRoleCredentials(t *testing.T) {
	var mu sync.Mutex
	var authorization []string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		authorization = append(authorization, r.Header.Get("Authorization"))
		mu.Unlock()
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"AccessKeyId":"id","SecretAccessKey":"secret","Token":"session","Expiration":"2099-01-01T00:00:00Z"}`))
	}))
	defer server.Close()
	file := filepath.Join(t.TempDir(), "token")
	if err := os.WriteFile(file, []byte("from-file"), 0600); err != nil {
		t.Fatal(err)
	}

	provider, err := taskRoleCredentials(S3Config{
		ContainerCredentialsFullURI:     server.URL + "/credentials",
		ContainerAuthorizationToken:     "from-value",
		ContainerAuthorizationTokenFile: file,
	})
	if err != nil {
		t.Fatal(err)
	}
	creds, err := provider.Retrieve(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if got := creds.AccessKeyID + "/" + creds.SecretAccessKey + "/" + creds.SessionToken; got != "id/secret/session" {
		t.Errorf("Retrieve() = %q, want %q", got, "id/secret/session")
	}
	mu.Lock()
	defer mu.Unlock()
	if want := []string{"from-file"}; !reflect.DeepEqual(authorization, want) {
		t.Errorf("Authorization headers = %q, want %q", authorization, want)
	}

	_, err = taskRoleCredentials(S3Config{
		ContainerCredentialsFullURI:     server.URL + "/credentials",
		ContainerAuthorizationTokenFile: filepath.Join(t.TempDir(), "missing"),
	})
	if err == nil {
		t.Error("taskRoleCredentials(missing token file) error = nil, want an error")
	}
}

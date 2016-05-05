package web

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/hashicorp/go-version"
	"github.com/lgtmco/lgtm/model"
	"github.com/lgtmco/lgtm/remote"
	"math"
	"net/http"
	"strconv"
	"testing"
	"time"
	"github.com/lgtmco/lgtm/approval"
)

func TestHandleTimestampMillis(t *testing.T) {
	c := &model.Config{
		VersionFormat: "millis",
	}
	stamp, err := handleTimestamp(c)
	if err != nil {
		t.Errorf("didn't expect error: %v", err)
	}
	m := time.Now().UTC().Unix()
	m1, err := strconv.ParseInt(*stamp, 10, 64)
	if err != nil {
		t.Error(err)
	}
	if math.Abs(float64(m - m1)) > 100 {
		t.Errorf("shouldn't be that different: %d, %d", m, m1)
	}
}

func TestHandleTimestampBlank(t *testing.T) {
	c := &model.Config{
		VersionFormat: "",
	}
	stamp, err := handleTimestamp(c)
	if err != nil {
		t.Errorf("didn't expect error: %v", err)
	}
	//should be able to parse with rfc3339
	t2, err := time.Parse(time.RFC3339, *stamp)
	if err != nil {
		t.Error(err)
	}
	round := t2.Format(time.RFC3339)
	if round != *stamp {
		t.Errorf("Expected to be same, but wasn't %s, %s", *stamp, round)
	}
	fmt.Println(*stamp)
}

func TestHandleTimestamp3339(t *testing.T) {
	c := &model.Config{
		VersionFormat: "rfc3339",
	}
	stamp, err := handleTimestamp(c)
	if err != nil {
		t.Errorf("didn't expect error: %v", err)
	}
	//should be able to parse with rfc3339
	t2, err := time.Parse(time.RFC3339, *stamp)
	if err != nil {
		t.Error(err)
	}
	round := t2.Format(time.RFC3339)
	if round != *stamp {
		t.Errorf("Expected to be same, but wasn't %s, %s", *stamp, round)
	}
	fmt.Println(*stamp)
}

func TestHandleTimestampCustom(t *testing.T) {
	//01/02 03:04:05PM '06 -0700
	c := &model.Config{
		VersionFormat: "Jan 2 2006, 3:04:05 PM",
	}
	stamp, err := handleTimestamp(c)
	if err != nil {
		t.Errorf("didn't expect error: %v", err)
	}
	//should be able to parse with custom
	t2, err := time.Parse(c.VersionFormat, *stamp)
	if err != nil {
		t.Error(err)
	}
	round := t2.Format(c.VersionFormat)
	if round != *stamp {
		t.Errorf("Expected to be same, but wasn't %s, %s", *stamp, round)
	}
	fmt.Println(*stamp)
}

type myR struct {
	remote.Remote
}

func (m *myR) ListTags(u *model.User, r *model.Repo) ([]model.Tag, error) {
	return []model.Tag{
		"a",
		"0.1.0",
		"0.0.1",
	}, nil
}

func (m *myR) GetComments(u *model.User, r *model.Repo, num int) ([]*model.Comment, error) {
	return []*model.Comment{
		{
			Author: "test_guy",
			Body:   "LGTM 0.1.0",
		},
		{
			Author: "not_test_guy",
			Body:   "this is not an LGTM comment",
		},
		{
			Author: "not_test_guy",
			Body:   "LGTM",
		},
		{
			Author: "test_guy",
			Body:   "LGTM 0.1.0",
		},
		{
			Author: "test_guy2",
			Body:   "LGTM",
		},
		{
			Author: "test_guy3",
			Body:   "LGTM 0.0.1",
		},
	}, nil
}

func TestGetMaxVersionComment(t *testing.T) {
	c := &gin.Context{}

	remote.ToContext(c, &myR{})
	config := &model.Config{
		DoVersion: true,
		Pattern:   `(?i)LGTM\s*(\S*)`,
	}
	m := &model.Maintainer{
		People: map[string]*model.Person{
			"test_guy": &model.Person{
				Name: "test_guy",
			},
			"test_guy2": &model.Person{
				Name: "test_guy2",
			},
			"test_guy3": &model.Person{
				Name: "test_guy3",
			},
		},
	}
	i := model.Issue{
		Author: "test_guy",
	}
	comments := []*model.Comment{
		{
			Author: "test_guy",
			Body:   "LGTM 0.1.0",
		},
		{
			Author: "not_test_guy",
			Body:   "this is not an LGTM comment",
		},
		{
			Author: "not_test_guy",
			Body:   "LGTM",
		},
		{
			Author: "test_guy",
			Body:   "LGTM 0.1.0",
		},
		{
			Author: "test_guy2",
			Body:   "LGTM",
		},
		{
			Author: "test_guy3",
			Body:   "LGTM 0.0.1",
		},
	}
	alg, _ := approval.Lookup("simple")
	ver := getMaxVersionComment(config, m, i, comments, alg)
	if ver == nil {
		t.Fatalf("Got nil for version")
	}
	expected, _ := version.NewVersion("0.1.0")
	if !expected.Equal(ver) {
		t.Errorf("Expected %s, got %s", expected.String(), ver.String())
	}
}

func TestGetMaxVersionCommentBadPattern(t *testing.T) {
	c := &gin.Context{}

	remote.ToContext(c, &myR{})
	config := &model.Config{
		DoVersion: true,
		Pattern:   `?i)LGTM\s*(\S*)`,
	}
	m := &model.Maintainer{
		People: map[string]*model.Person{
			"test_guy": &model.Person{
				Name: "test_guy",
			},
			"test_guy2": &model.Person{
				Name: "test_guy2",
			},
			"test_guy3": &model.Person{
				Name: "test_guy3",
			},
		},
	}
	i := model.Issue{
		Author: "test_guy",
	}
	comments := []*model.Comment{
		{
			Author: "test_guy",
			Body:   "LGTM 0.1.0",
		},
		{
			Author: "not_test_guy",
			Body:   "this is not an LGTM comment",
		},
		{
			Author: "not_test_guy",
			Body:   "LGTM",
		},
		{
			Author: "test_guy",
			Body:   "not an approval comment",
		},
		{
			Author: "test_guy",
			Body:   "LGTM 0.1.0",
		},
		{
			Author: "test_guy2",
			Body:   "LGTM",
		},
		{
			Author: "test_guy3",
			Body:   "LGTM 0.0.1",
		},
	}
	alg, _ := approval.Lookup("simple")
	ver := getMaxVersionComment(config, m, i, comments, alg)
	if ver != nil {
		t.Fatalf("Should get nil for version")
	}
}

func TestGetMaxVersionCommentNoSelfApproval(t *testing.T) {
	c := &gin.Context{}

	remote.ToContext(c, &myR{})
	config := &model.Config{
		DoVersion:       true,
		Pattern:         `(?i)LGTM\s*(\S*)`,
		SelfApprovalOff: true,
	}
	m := &model.Maintainer{
		People: map[string]*model.Person{
			"test_guy": &model.Person{
				Name: "test_guy",
			},
			"test_guy2": &model.Person{
				Name: "test_guy2",
			},
			"test_guy3": &model.Person{
				Name: "test_guy3",
			},
		},
	}
	i := model.Issue{
		Author: "test_guy",
	}
	comments := []*model.Comment{
		{
			Author: "test_guy",
			Body:   "LGTM 0.1.0",
		},
		{
			Author: "not_test_guy",
			Body:   "this is not an LGTM comment",
		},
		{
			Author: "not_test_guy",
			Body:   "LGTM",
		},
		{
			Author: "test_guy",
			Body:   "LGTM 0.1.0",
		},
		{
			Author: "test_guy2",
			Body:   "LGTM",
		},
		{
			Author: "test_guy3",
			Body:   "LGTM 0.0.1",
		},
	}
	alg, _ := approval.Lookup("simple")
	ver := getMaxVersionComment(config, m, i, comments, alg)
	if ver == nil {
		t.Fatalf("Got nil for version")
	}
	expected, _ := version.NewVersion("0.0.1")
	if !expected.Equal(ver) {
		t.Errorf("Expected %s, got %s", expected.String(), ver.String())
	}
}

func TestGetMaxExistingTagFound(t *testing.T) {
	ver := getMaxExistingTag([]model.Tag{
		"a",
		"0.1.0",
		"0.0.1",
	})

	expected, _ := version.NewVersion("0.1.0")
	if !expected.Equal(ver) {
		t.Errorf("Expected %s, got %s", expected.String(), ver.String())
	}
}

func TestGetMaxExistingTagNotFound(t *testing.T) {
	ver := getMaxExistingTag([]model.Tag{
		"a",
		"b",
		"c",
	})

	expected, _ := version.NewVersion("0.0.0")
	if !expected.Equal(ver) {
		t.Errorf("Expected %s, got %s", expected.String(), ver.String())
	}
}

func TestHandleSemver(t *testing.T) {
	c := &gin.Context{}

	remote.ToContext(c, &myR{})
	config := &model.Config{
		DoVersion: true,
		Pattern:   `(?i)LGTM\s*(\S*)`,
		ApprovalAlg: "simple",
	}
	m := &model.Maintainer{
		People: map[string]*model.Person{
			"test_guy": &model.Person{
				Name: "test_guy",
			},
			"test_guy2": &model.Person{
				Name: "test_guy2",
			},
			"test_guy3": &model.Person{
				Name: "test_guy3",
			},
		},
	}
	user := &model.User{}
	repo := &model.Repo{}
	hook := &model.StatusHook{
		Repo: &model.Repo{
			Owner: "test_guy",
			Name:  "test_repo",
		},
	}
	pr := model.PullRequest{
		Issue: model.Issue{
			Author: "test_guy",
		},
	}
	ver, err := handleSemver(c, user, hook, pr, config, m, repo)
	if err != nil {
		t.Fatalf("unexpected error %s", err)
	}
	expected, _ := version.NewVersion("0.1.1")
	if expected.String() != *ver {
		t.Errorf("Expected %s, got %s", expected.String(), *ver)
	}
}

type myR2 struct {
	remote.Remote
}

func (m *myR2) ListTags(u *model.User, r *model.Repo) ([]model.Tag, error) {
	return []model.Tag{
		"a",
		"0.0.1",
		"0.0.2",
	}, nil
}

func (m *myR2) GetComments(u *model.User, r *model.Repo, num int) ([]*model.Comment, error) {
	return []*model.Comment{
		{
			Author: "test_guy",
			Body:   "LGTM 0.1.0",
		},
		{
			Author: "not_test_guy",
			Body:   "this is not an LGTM comment",
		},
		{
			Author: "not_test_guy",
			Body:   "LGTM",
		},
		{
			Author: "test_guy",
			Body:   "LGTM 0.1.0",
		},
		{
			Author: "test_guy2",
			Body:   "LGTM",
		},
		{
			Author: "test_guy3",
			Body:   "LGTM 0.0.1",
		},
	}, nil
}

func TestHandleSemver2(t *testing.T) {
	c := &gin.Context{}

	remote.ToContext(c, &myR2{})
	config := &model.Config{
		DoVersion: true,
		Pattern:   `(?i)LGTM\s*(\S*)`,
		ApprovalAlg: "simple",
	}
	m := &model.Maintainer{
		People: map[string]*model.Person{
			"test_guy": &model.Person{
				Name: "test_guy",
			},
			"test_guy2": &model.Person{
				Name: "test_guy2",
			},
			"test_guy3": &model.Person{
				Name: "test_guy3",
			},
		},
	}
	user := &model.User{}
	repo := &model.Repo{}
	hook := &model.StatusHook{
		Repo: &model.Repo{
			Owner: "test_guy",
			Name:  "test_repo",
		},
	}
	pr := model.PullRequest{
		Issue: model.Issue{
			Author: "test_guy",
		},
	}
	ver, err := handleSemver(c, user, hook, pr, config, m, repo)
	if err != nil {
		t.Fatalf("unexpected error %s", err)
	}
	expected, _ := version.NewVersion("0.1.0")
	if expected.String() != *ver {
		t.Errorf("Expected %s, got %s", expected.String(), *ver)
	}
}

type myR3 struct {
	remote.Remote
}

func (m *myR3) ListTags(u *model.User, r *model.Repo) ([]model.Tag, error) {
	return nil, errors.New("This is an error")
}

func (m *myR3) GetComments(u *model.User, r *model.Repo, num int) ([]*model.Comment, error) {
	return []*model.Comment{
		{
			Author: "test_guy",
			Body:   "LGTM 0.1.0",
		},
		{
			Author: "not_test_guy",
			Body:   "this is not an LGTM comment",
		},
		{
			Author: "not_test_guy",
			Body:   "LGTM",
		},
		{
			Author: "test_guy",
			Body:   "LGTM 0.1.0",
		},
		{
			Author: "test_guy2",
			Body:   "LGTM",
		},
		{
			Author: "test_guy3",
			Body:   "LGTM 0.0.1",
		},
	}, nil
}

func TestHandleSemver3(t *testing.T) {
	c := &gin.Context{}

	remote.ToContext(c, &myR3{})
	config := &model.Config{
		DoVersion: true,
		Pattern:   `(?i)LGTM\s*(\S*)`,
		ApprovalAlg: "simple",
	}
	m := &model.Maintainer{
		People: map[string]*model.Person{
			"test_guy": &model.Person{
				Name: "test_guy",
			},
			"test_guy2": &model.Person{
				Name: "test_guy2",
			},
			"test_guy3": &model.Person{
				Name: "test_guy3",
			},
		},
	}
	user := &model.User{}
	repo := &model.Repo{}
	hook := &model.StatusHook{
		Repo: &model.Repo{
			Owner: "test_guy",
			Name:  "test_repo",
		},
	}
	pr := model.PullRequest{
		Issue: model.Issue{
			Author: "test_guy",
		},
	}
	ver, err := handleSemver(c, user, hook, pr, config, m, repo)
	if err != nil {
		t.Fatalf("unexpected error %s", err)
	}
	expected, _ := version.NewVersion("0.1.0")
	if expected.String() != *ver {
		t.Errorf("Expected %s, got %s", expected.String(), *ver)
	}
}

type myR4 struct {
	remote.Remote
}

func (m *myR4) ListTags(u *model.User, r *model.Repo) ([]model.Tag, error) {
	return []model.Tag{
		"a",
		"0.0.1",
		"0.0.2",
	}, nil
}

func (m *myR4) GetComments(u *model.User, r *model.Repo, num int) ([]*model.Comment, error) {
	return nil, errors.New("This is an error")
}

type myRR struct {
	gin.ResponseWriter
}

func (mr *myRR) Header() http.Header {
	return http.Header{}
}

func (my *myRR) Write(p []byte) (n int, err error) {
	return len(p), nil
}

func TestHandleSemver4(t *testing.T) {
	c := &gin.Context{
		Writer: &myRR{},
	}

	remote.ToContext(c, &myR4{})
	config := &model.Config{
		DoVersion: true,
		Pattern:   `(?i)LGTM\s*(\S*)`,
		ApprovalAlg: "simple",
	}
	m := &model.Maintainer{
		People: map[string]*model.Person{
			"test_guy": &model.Person{
				Name: "test_guy",
			},
			"test_guy2": &model.Person{
				Name: "test_guy2",
			},
			"test_guy3": &model.Person{
				Name: "test_guy3",
			},
		},
	}
	user := &model.User{}
	repo := &model.Repo{}
	hook := &model.StatusHook{
		Repo: &model.Repo{
			Owner: "test_guy",
			Name:  "test_repo",
		},
	}
	pr := model.PullRequest{
		Issue: model.Issue{
			Author: "test_guy",
		},
	}
	_, err := handleSemver(c, user, hook, pr, config, m, repo)
	if err == nil {
		t.Fatalf("expected error, got none")
	}
}

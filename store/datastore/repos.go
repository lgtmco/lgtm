package datastore

import (
	"fmt"
	"strings"

	"github.com/lgtmco/lgtm/model"

	"github.com/russross/meddler"
)

func (db *datastore) GetRepo(id int64) (*model.Repo, error) {
	var repo = new(model.Repo)
	var err = meddler.Load(db, repoTable, repo, id)
	return repo, err
}

func (db *datastore) GetRepoSlug(slug string) (*model.Repo, error) {
	var repo = new(model.Repo)
	var err = meddler.QueryRow(db, repo, repoSlugQuery, slug)
	return repo, err
}

func (db *datastore) GetRepoMulti(slug ...string) ([]*model.Repo, error) {
	var repos = []*model.Repo{}
	var instr, params = toList(slug)
	var stmt = fmt.Sprintf(repoListQuery, instr)
	var err = meddler.QueryAll(db, &repos, stmt, params...)
	return repos, err
}

func (db *datastore) GetRepoOwner(owner string) ([]*model.Repo, error) {
	var repos = []*model.Repo{}
	var err = meddler.QueryAll(db, &repos, repoOwnerQuery, owner)
	return repos, err
}

func (db *datastore) CreateRepo(repo *model.Repo) error {
	return meddler.Insert(db, repoTable, repo)
}

func (db *datastore) UpdateRepo(repo *model.Repo) error {
	return meddler.Update(db, repoTable, repo)
}

func (db *datastore) DeleteRepo(repo *model.Repo) error {
	var _, err = db.Exec(repoDeleteStmt, repo.ID)
	return err
}

func toList(items []string) (string, []interface{}) {
	var size = len(items)
	if size > 990 {
		size = 990
		items = items[:990]
	}
	var qs = make([]string, size, size)
	var in = make([]interface{}, size, size)
	for i, item := range items {
		qs[i] = "?"
		in[i] = item
	}
	return strings.Join(qs, ","), in
}

const repoTable = "repos"

const repoSlugQuery = `
SELECT *
FROM repos
WHERE repo_slug = ?
LIMIT 1;
`

const repoOwnerQuery = `
SELECT *
FROM repos
WHERE repo_owner = ?
`

const repoListQuery = `
SELECT *
FROM repos
WHERE repo_slug IN (%s)
ORDER BY repo_slug
`

const repoDeleteStmt = `
DELETE FROM repos
WHERE repo_id = ?
`

// Copyright 2023 The Witness Contributors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package git

import (
	"crypto"
	"fmt"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/in-toto/go-witness/attestation"
	"github.com/in-toto/go-witness/cryptoutil"
)

const (
	Name    = "git"
	Type    = "https://witness.dev/attestations/git/v0.1"
	RunType = attestation.PreMaterialRunType
)

// This is a hacky way to create a compile time error in case the attestor
// doesn't implement the expected interfaces.
var (
	_ attestation.Attestor   = &Attestor{}
	_ attestation.Subjecter  = &Attestor{}
	_ attestation.BackReffer = &Attestor{}
)

func init() {
	attestation.RegisterAttestation(Name, Type, RunType, func() attestation.Attestor {
		return New()
	})
}

type Status struct {
	Staging  string `json:"staging,omitempty"`
	Worktree string `json:"worktree,omitempty"`
}

type Tag struct {
	Name         string `json:"name"`
	TaggerName   string `json:"taggername"`
	TaggerEmail  string `json:"taggeremail"`
	When         string `json:"when"`
	PGPSignature string `json:"pgpsignature"`
	Message      string `json:"message"`
}

type Attestor struct {
	CommitHash     string               `json:"commithash"`
	Author         string               `json:"author"`
	AuthorEmail    string               `json:"authoremail"`
	CommitterName  string               `json:"committername"`
	CommitterEmail string               `json:"committeremail"`
	CommitDate     string               `json:"commitdate"`
	CommitMessage  string               `json:"commitmessage"`
	Status         map[string]Status    `json:"status,omitempty"`
	CommitDigest   cryptoutil.DigestSet `json:"commitdigest,omitempty"`
	Signature      string               `json:"signature,omitempty"`
	ParentHashes   []string             `json:"parenthashes,omitempty"`
	TreeHash       string               `json:"treehash,omitempty"`
	Refs           []string             `json:"refs,omitempty"`
	Tags           []Tag                `json:"tags,omitempty"`
}

func New() *Attestor {
	return &Attestor{
		Status: make(map[string]Status),
	}
}

func (a *Attestor) Name() string {
	return Name
}

func (a *Attestor) Type() string {
	return Type
}

func (a *Attestor) RunType() attestation.RunType {
	return RunType
}

func (a *Attestor) Attest(ctx *attestation.AttestationContext) error {
	repo, err := git.PlainOpenWithOptions(ctx.WorkingDir(), &git.PlainOpenOptions{
		DetectDotGit: true,
	})

	if err != nil {
		return err
	}

	head, err := repo.Head()
	if err != nil {
		if strings.Contains(err.Error(), "reference not found") {
			return nil
		}
		return err
	}

	commit, err := repo.CommitObject(head.Hash())
	if err != nil {
		return err
	}

	a.CommitDigest = cryptoutil.DigestSet{
		{
			Hash:   crypto.SHA1,
			GitOID: false,
		}: commit.Hash.String(),
	}

	//get all the refs for the repo
	refs, err := repo.References()
	if err != nil {
		return err
	}

	//iterate over the refs and add them to the attestor
	err = refs.ForEach(func(ref *plumbing.Reference) error {
		//only add the ref if it points to the head
		if ref.Hash() != head.Hash() {
			return nil
		}

		//add the ref name to the attestor
		a.Refs = append(a.Refs, ref.Name().String())

		return nil
	})
	if err != nil {
		return err
	}

	a.CommitHash = head.Hash().String()
	a.Author = commit.Author.Name
	a.AuthorEmail = commit.Author.Email
	a.CommitterName = commit.Committer.Name
	a.CommitterEmail = commit.Committer.Email
	a.CommitDate = commit.Author.When.String()
	a.CommitMessage = commit.Message
	a.Signature = commit.PGPSignature

	for _, parent := range commit.ParentHashes {
		a.ParentHashes = append(a.ParentHashes, parent.String())
	}

	tags, err := repo.TagObjects()
	if err != nil {
		return fmt.Errorf("get tags error: %s", err)
	}

	var tagList []Tag

	err = tags.ForEach(func(t *object.Tag) error {

		//check if the tag points to the head
		if t.Target.String() != head.Hash().String() {
			return nil
		}

		tagList = append(tagList, Tag{
			Name:         t.Name,
			TaggerName:   t.Tagger.Name,
			TaggerEmail:  t.Tagger.Email,
			When:         t.Tagger.When.Format(time.RFC3339),
			PGPSignature: t.PGPSignature,
			Message:      t.Message,
		})
		return nil
	})

	if err != nil {
		return fmt.Errorf("iterate tags error: %s", err)
	}
	a.Tags = tagList

	a.TreeHash = commit.TreeHash.String()

	worktree, err := repo.Worktree()
	if err != nil {
		return err
	}

	status, err := worktree.Status()
	if err != nil {
		return err
	}

	for file, status := range status {
		if status.Worktree == git.Unmodified && status.Staging == git.Unmodified {
			continue
		}

		attestStatus := Status{
			Worktree: statusCodeString(status.Worktree),
			Staging:  statusCodeString(status.Staging),
		}

		a.Status[file] = attestStatus
	}

	return nil
}

func (a *Attestor) Subjects() map[string]cryptoutil.DigestSet {
	subjects := make(map[string]cryptoutil.DigestSet)
	hashes := []cryptoutil.DigestValue{{Hash: crypto.SHA256}}

	subjectName := fmt.Sprintf("commithash:%v", a.CommitHash)
	subjects[subjectName] = cryptoutil.DigestSet{
		{
			Hash:   crypto.SHA1,
			GitOID: false,
		}: a.CommitHash,
	}

	//add author email
	subjectName = fmt.Sprintf("authoremail:%v", a.AuthorEmail)
	ds, err := cryptoutil.CalculateDigestSetFromBytes([]byte(a.AuthorEmail), hashes)
	if err != nil {
		return nil
	}

	subjects[subjectName] = ds

	//add committer email
	subjectName = fmt.Sprintf("committeremail:%v", a.CommitterEmail)
	ds, err = cryptoutil.CalculateDigestSetFromBytes([]byte(a.CommitterEmail), hashes)
	if err != nil {
		return nil
	}

	subjects[subjectName] = ds

	//add parent hashes
	for _, parentHash := range a.ParentHashes {
		subjectName = fmt.Sprintf("parenthash:%v", parentHash)
		ds, err = cryptoutil.CalculateDigestSetFromBytes([]byte(parentHash), hashes)
		if err != nil {
			return nil
		}
		subjects[subjectName] = ds
	}

	return subjects
}

func (a *Attestor) BackRefs() map[string]cryptoutil.DigestSet {
	backrefs := make(map[string]cryptoutil.DigestSet)
	subjectName := fmt.Sprintf("commithash:%v", a.CommitHash)
	backrefs[subjectName] = cryptoutil.DigestSet{
		{
			Hash:   crypto.SHA1,
			GitOID: false,
		}: a.CommitHash,
	}
	return backrefs
}

func statusCodeString(statusCode git.StatusCode) string {
	switch statusCode {
	case git.Unmodified:
		return "unmodified"
	case git.Untracked:
		return "untracked"
	case git.Modified:
		return "modified"
	case git.Added:
		return "added"
	case git.Deleted:
		return "deleted"
	case git.Renamed:
		return "renamed"
	case git.Copied:
		return "copied"
	case git.UpdatedButUnmerged:
		return "updated"
	default:
		return string(statusCode)
	}
}

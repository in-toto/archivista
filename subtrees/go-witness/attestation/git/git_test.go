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
	"math/rand"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/stretchr/testify/require"
	"github.com/testifysec/go-witness/attestation"
)

func TestNew(t *testing.T) {
	attestor := New()
	require.NotNil(t, attestor, "Expected a new attestor")
	require.NotNil(t, attestor.Status, "Expected a map for Status")
}

func TestNameTypeRunType(t *testing.T) {
	attestor := New()
	require.Equal(t, Name, attestor.Name(), "Expected the attestor's name")
	require.Equal(t, Type, attestor.Type(), "Expected the attestor's type")
	require.Equal(t, RunType, attestor.RunType(), "Expected the attestor's run type")
}

func TestRun(t *testing.T) {
	attestor := New()

	_, dir, cleanup := createTestRepo(t)
	defer cleanup()

	ctx, err := attestation.NewContext([]attestation.Attestor{attestor}, attestation.WithWorkingDir(dir))
	require.NoError(t, err, "Expected no error from NewContext")

	err = ctx.RunAttestors()
	require.NoError(t, err, "Expected no error from RunAttestors")

	require.Empty(t, attestor.ParentHashes, "Expected the parent hashes to be set")

	createTestCommit(t, dir, "Test commit")
	createTestRefs(t, dir)
	createAnnotatedTagOnHead(t, dir)
	err = ctx.RunAttestors()

	// Check that the attestor has the expected values

	require.NoError(t, err, "Expected no error from attestation")
	require.NotEmpty(t, attestor.CommitHash, "Expected the commit hash to be set")
	require.NotEmpty(t, attestor.Author, "Expected the author to be set")
	require.NotEmpty(t, attestor.AuthorEmail, "Expected the author's email to be set")
	require.NotEmpty(t, attestor.CommitterName, "Expected the committer to be set")
	require.NotEmpty(t, attestor.CommitterEmail, "Expected the committer's email to be set")
	require.NotEmpty(t, attestor.CommitDate, "Expected the commit date to be set")
	require.NotEmpty(t, attestor.CommitMessage, "Expected the commit message to be set")
	require.NotEmpty(t, attestor.CommitDigest, "Expected the commit digest to be set")
	require.NotEmpty(t, attestor.TreeHash, "Expected the tree hash to be set")
	require.NotEmpty(t, attestor.ParentHashes, "Expected the parent hashes to be set")

	subjects := attestor.Subjects()
	require.NotNil(t, subjects, "Expected subjects to be non-nil")

	// Test for the existence of subjects
	require.Contains(t, subjects, fmt.Sprintf("commithash:%v", attestor.CommitHash), "Expected commithash subject to exist")
	require.Contains(t, subjects, fmt.Sprintf("authoremail:%v", attestor.AuthorEmail), "Expected authoremail subject to exist")
	require.Contains(t, subjects, fmt.Sprintf("committeremail:%v", attestor.CommitterEmail), "Expected committeremail subject to exist")

	for _, parentHash := range attestor.ParentHashes {
		subjectName := fmt.Sprintf("parenthash:%v", parentHash)
		require.Contains(t, subjects, subjectName, "Expected parent hash subject to exist")
	}

	backrefs := attestor.BackRefs()
	require.NotNil(t, backrefs, "Expected backrefs to be non-nil")

	subjectName := fmt.Sprintf("commithash:%v", attestor.CommitHash)
	require.Contains(t, backrefs, subjectName, "Expected commithash backref to exist")

	ds := backrefs[subjectName]
	require.NotNil(t, ds, "Expected a digest set for the commithash backref")

	// Test for the existence of a SHA1 digest in the digest set
	var found bool
	for d, v := range ds {
		if d.Hash == crypto.SHA1 {
			found = true
			require.Equal(t, d.GitOID, false, "Expected GitOID to be false")
			require.Equal(t, v, attestor.CommitHash, "Expected the correct value for the SHA1 digest")
		}
	}

	require.True(t, found, "Expected a SHA1 digest in the digest set")

	subRefName := "refs/tags/v1.0.0-lightweight"
	require.NoError(t, err, "Expected lightweight tag to exist")
	require.Contains(t, attestor.Refs, subRefName, "Expected lightweight tag ref to be attested")

	subRefName = "refs/heads/my-feature-branch"
	require.NoError(t, err, "Expected branch to exist")
	require.Contains(t, attestor.Refs, subRefName, "Expected branch ref to be attested")

	subRefName = "refs/heads/my-feature-branch@123"
	require.NoError(t, err, "Expected ref with special characters to exist")
	require.Contains(t, attestor.Refs, subRefName, "Expected ref with special characters to be attested")

	// Test the annotated tag contents
	tags := attestor.Tags
	require.NotNil(t, tags, "Expected tags to be non-nil")

	//we should have 1 tag
	require.Len(t, tags, 1, "Expected 1 tag")

	//get the tag object
	tagObject := tags[0]

	require.NoError(t, err, "Expected no error from getTagObject")
	require.Equal(t, "v1.0.0-test", tagObject.Name)
	require.Equal(t, "example tag message\n", tagObject.Message)
	require.Equal(t, "John Doe", tagObject.TaggerName)
	require.Equal(t, "tagger@example.com", tagObject.TaggerEmail)

}

func createTestRepo(t *testing.T) (*git.Repository, string, func()) {
	// Create a temporary directory for the test repository
	tmpDir, err := os.MkdirTemp("", "test-repo")
	require.NoError(t, err)

	// Initialize a new Git repository in the temporary directory
	repo, err := git.PlainInit(tmpDir, false)
	require.NoError(t, err)

	// Create a new file in the repository
	filePath := filepath.Join(tmpDir, "test.txt")
	file, err := os.Create(filePath)
	require.NoError(t, err)
	_, err = file.WriteString("Test file")
	require.NoError(t, err)
	err = file.Close()
	require.NoError(t, err)

	// Add the new file to the repository
	worktree, err := repo.Worktree()
	require.NoError(t, err)
	_, err = worktree.Add("test.txt")
	require.NoError(t, err)

	// Commit the new file to the repository
	_, err = worktree.Commit("Initial commit", &git.CommitOptions{
		Author: &object.Signature{
			Name:  "Test User",
			Email: "test@example.com",
			When:  time.Now(),
		},
	})
	require.NoError(t, err)

	// Return the test repository, the path to the test repository, and a cleanup function
	return repo, tmpDir, func() {
		err := os.RemoveAll(tmpDir)
		require.NoError(t, err)
	}
}
func createTestCommit(t *testing.T, repoPath string, message string) {
	// Open the Git repository
	repo, err := git.PlainOpen(repoPath)
	require.NoError(t, err)

	// Get the HEAD reference
	headRef, err := repo.Head()
	require.NoError(t, err)

	// Get the commit that the HEAD reference points to
	commit, err := repo.CommitObject(headRef.Hash())
	require.NoError(t, err)

	// Create a new file in the repository with a random string
	randStr := fmt.Sprintf("%d", rand.Int())
	filePath := filepath.Join(repoPath, "test.txt")
	file, err := os.Create(filePath)
	require.NoError(t, err)
	_, err = file.WriteString(randStr)
	require.NoError(t, err)
	err = file.Close()
	require.NoError(t, err)

	// Add the new file to the repository
	worktree, err := repo.Worktree()
	require.NoError(t, err)
	_, err = worktree.Add("test.txt")
	require.NoError(t, err)

	// Commit the new file to the repository use current commit as parent
	_, err = worktree.Commit(message, &git.CommitOptions{
		All:       false,
		Author:    &object.Signature{Name: "Test User", Email: "test@example.com", When: time.Now()},
		Committer: &object.Signature{Name: "Test User", Email: "test@example.com", When: time.Now()},
		Parents:   []plumbing.Hash{commit.Hash},
	})
	require.NoError(t, err)
}

func createTestRefs(t *testing.T, dir string) {
	// Open the Git repository
	repo, err := git.PlainOpen(dir)
	require.NoError(t, err)

	// Get the HEAD reference
	headRef, err := repo.Head()
	require.NoError(t, err)

	// Get the commit that the HEAD reference points to
	hash := headRef.Hash()

	// Create a new branch ref pointing to the specified commit hash
	branchRef := plumbing.NewBranchReferenceName("my-feature-branch")
	err = repo.Storer.SetReference(plumbing.NewHashReference(branchRef, hash))
	require.NoError(t, err)

	// Create a new lightweight tag pointing to the specified commit hash
	lightweightTagName := "v1.0.0-lightweight"
	err = repo.Storer.SetReference(plumbing.NewHashReference(plumbing.ReferenceName("refs/tags/"+lightweightTagName), hash))
	require.NoError(t, err)

	// Create a new ref with a special character in the name
	specialCharRef := plumbing.NewHashReference(plumbing.ReferenceName("refs/heads/my-feature-branch@123"), hash)
	err = repo.Storer.SetReference(specialCharRef)
	require.NoError(t, err)
}

func createAnnotatedTagOnHead(t *testing.T, path string) {
	// Open the Git repository.
	repo, err := git.PlainOpen(path)
	require.NoError(t, err)

	// Get the HEAD reference.
	headRef, err := repo.Head()
	require.NoError(t, err)

	// Get the commit that the HEAD reference points to.
	commit, err := repo.CommitObject(headRef.Hash())
	require.NoError(t, err)

	_, err = repo.CreateTag("v1.0.0-test", commit.Hash, &git.CreateTagOptions{
		Tagger: &object.Signature{
			Name: "John Doe",

			Email: "tagger@example.com",
			When:  time.Now(),
		},
		Message: "example tag message",
	})

	require.NoError(t, err)
}

# Git Attestor

The Git Attestor records the current state of the objects in the git repository, including untracked objects.

Both staged and unstaged states are recorded.

The Git Attestor assumes you are working from a git repository that has been initialized and has commits.

## Subjects

The attestor returns the SHA1 ([Secure Hash Algorithm 1](https://en.wikipedia.org/wiki/SHA-1)) git commit hash as a subject.

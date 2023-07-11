# Subtrees

This monorepo contains a mix of internal private code and public open source code which we make available to the public.

The public open source code has been added to the project as [Git Subtrees](https://www.atlassian.com/git/tutorials/git-subtree).

## List of our subtrees

| path-to-subtree               | name-of-subtree      | name-of-remote       |
| ----------------------------- | -------------------- | -------------------- |
| `subtrees/archivista`         | `archivista`         | `archivista`         |
| `subtrees/go-witness`         | `go-witness`         | `go-witness`         |
| `subtrees/witness`            | `witness`            | `witness`            |
| `subtrees/witness-run-action` | `witness-run-action` | `witness-run-action` |

## Getting Started

Assuming you have already gotten started with the root readme.md, then you can do the following:

1. `npm run remotes:add:all` this is a one time script to add all of our remotes. You'll only need to run it one time initially, and anytime we add more remote subtrees.
1. `npm run remotes:fetch:all` will fetch all the remote subtrees.

From there you can add more subtrees, update subtree subfolders, and/or update subtree parent repositories.

## Helper scripts

Review the scripts in our root `package.json` for helper scripts related to subtrees and remotes.

## Adding a subtree

You should add subtrees as remotes for easier management.

1. `git remote add -f <name-of-your-subtree> https://github.com/testifysec/some-great-open-source-project.git` with the name of your subtree and the https or ssh link to the git project. You can add this as a script in our root `package.json` to help other devs out!
1. `git subtree add --prefix <path-to-your-subtree> <name-of-your-subtree> main --squash` will create an initial commit for your subtree into the monorepo project. This should be a one-time command. Be careful to keep any subtree work completely seperate from any other commits.

## Updating a subtree with upstream changes

So, some awesome soul decided to help us out and contributed open source changes to one of our subtrees!

Now what?

1. `git fetch <name-of-your-subtree> main`
1. `git subtree pull --prefix <path-to-your-subtree> <name-of-your-subtree> main --squash`

## Contributing internal changes back upstream to the subtree's true git repository

We can freely commit our fixes to the sub-project in our local working directory now.

So, you touched a subtree! Great work, we just made our open source project better.

Now what?

When it’s time to contribute back to the upstream project, we need to fork the project and add it as another remote:

1. `git remote add <name-of-your-fork> ssh://git@github.com/testifysec/some-great-open-source-project.git` to fork your changes into the subtree's true parent repository
1. `git subtree push --prefix <path-to-your-subtree> <name-of-your-fork> main`

## Review the Docs

Please [review the docs](https://gist.github.com/SKempin/b7857a6ff6bddb05717cc17a44091202) on [git subtree](https://www.atlassian.com/git/tutorials/git-subtree).

## Automate as much as we can, assist your fellow devs

As we learn and grow, we'll want to master managing these subtrees hygeniecally. Always offer tools and shortcuts for your fellow developers, and if we can, automate as much as we can.

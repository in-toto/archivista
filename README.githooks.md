# Our githooks

We use[Husky](https://www.npmjs.com/package/husky) to make githooks easy!

We use `lint-staged` [to only run our linters on staged files.](https://www.npmjs.com/package/lint-staged).

We use `prepush-if-changed` [to only run our prepush hooks on things we change when pushing](https://www.npmjs.com/package/prepush-if-changed).

To help string this all together, [we orchestrate it with npm workspaces](https://docs.npmjs.com/cli/v7/using-npm/workspaces).

## What happens when you make commits

1. Make changes to contribute to our projects, as normal.
1. Stage your changes, as normal. 
1. Make your commit, our `pre-commit` hook will fire from husky, which should run anything we have configured to automate pre-commit, including `lint-staged` to *lint* any code which you touched.

**_Note:_**: Our linter scripts are configured to auto-fix lint issues if they can, so if you missed lint issues when you try to commit, you may see the linters auto-fix issues and restage them for you. This is to help automate the process to help developers, so they do not need to think about it as much.

This works by configuring our `precommit` script in our package.json, which our husky pre-commit config points to.

Linting works by configuring `lint-staged` to look for specific file changes and then upon detection execute specific shell commands. You can find this configuration in `.lintstagedrc.yaml`.

The idea is that if you touch a file, whether it be `node` or `python` or whatever, the proper linters will run as needed. If you didn't touch something, you shouldn't have to wait.

## What happens when you make pushes

1. Make commits in your working branch, as normal.
1. Push your branch, our `pre-push` hook will fire from husky, which should run anything we have configured to automate pre-push, including `prepush-if-changed` to *test* any code you touched.

This works by configuring our `prepush` script in our package.json, which our husky pre-push config points to.

Testing works by configuring `prepush-if-changed` to look for specific file changes and then upon detection execute specific shell commands. You can find this configuration in `.prepushrc.yaml`.

The idea is that if you touch a file, whether it be `node` or `go` or whatever, the proper tests will run as needed. If you didn't touch something, you shouldn't have to wait.

## Githooks are just part of the picture

Remember that anything we automate with githooks needs proper enforcement in our CI!

Githooks can easily be bypassed, so they are only meant to help developers from breaking builds.

## How to bypass githooks

Let's say you're working on a branch and you want to punt on linting and testing for a while, you can commit using the `--no-verify` flag.

You can also push using `--no-verify` but if you break a build, that's bad.
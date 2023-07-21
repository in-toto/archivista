# Judge API

This is the judge-api, the graphql layer for the judge platform.

## Getting Started

You can follow the getting started from the root folder to spin this api up in the local kube.

## Code gen

This project uses gqlgen and ent to generate go code for us.

`go generate ./... -v` will generate all of the code from within this directory.

`npm run gen:judge-api` will generate all of the code from the root folder of the monorepo.

## Information Architecture

### Projects

A project is a git repository somewhere.

### Policies

Policies are custom requirement specs that allow us to assert that certain things must have happened. If they didn't happen, we can empower users by failing the shell command,
allowing users to stop a build or deployment.

### Policy Decisions

Policy Decisions are records of when a `witness verify` policy was executed and what the decision made on the result was.

Policy Decisions belong to a digest on a project.

You can query for policy decisions from a Project or from a DigestID or Subject_Name.

This api supports the ability to post-back policy decisions from `witness verify` at the `/policy_decsisions/` post endpoint.

It accepts policy decisions [cloudevents](https://github.com/cloudevents/spec) with an attached policy_decision object:

```go
type PolicyDecision struct {
	id          uuid.UUID // the id for this policy decision (auto created)
	SubjectName string // the subject name that the policy decision belonged to when `witness verify` was executed
	DigestID    string // the digest that the policy decision belonged to when `witness verify` was executed
	Timestamp   time.Time // the time the policy decision was created 
	Decision    DecisionEnum // the decision, either allowed, denied, or skipped.
}
```
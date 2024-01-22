package api

type GraphQLError struct {
	Message string `json:"message"`
}

type GraphQLResponseGeneric[T any] struct {
	Data   T              `json:"data,omitempty"`
	Errors []GraphQLError `json:"errors,omitempty"`
}

type GraphQLRequestBodyGeneric[TVars any] struct {
	Query     string `json:"query"`
	Variables TVars  `json:"variables,omitempty"`
}

type RetrieveSubjectVars struct {
	Gitoid string `json:"gitoid"`
}

type SearchVars struct {
	Algorithm string `json:"algo"`
	Digest    string `json:"digest"`
}

type RetrieveSubjectResults struct {
	Subjects Subjects `json:"subjects"`
}

type Subjects struct {
	Edges []SubjectEdge `json:"edges"`
}

type SubjectEdge struct {
	Node SubjectNode `json:"node"`
}

type SubjectNode struct {
	Name           string          `json:"name"`
	SubjectDigests []SubjectDigest `json:"subjectDigests"`
}

type SubjectDigest struct {
	Algorithm string `json:"algorithm"`
	Value     string `json:"value"`
}

type SearchResults struct {
	Dsses DSSES `json:"dsses"`
}

type DSSES struct {
	Edges []SearchEdge `json:"edges"`
}

type SearchEdge struct {
	Node SearchNode `json:"node"`
}

type SearchNode struct {
	GitoidSha256 string    `json:"gitoidSha256"`
	Statement    Statement `json:"statement"`
}

type Statement struct {
	AttestationCollection AttestationCollection `json:"attestationCollections"`
}

type AttestationCollection struct {
	Name         string        `json:"name"`
	Attestations []Attestation `json:"attestations"`
}

type Attestation struct {
	Type string `json:"type"`
}

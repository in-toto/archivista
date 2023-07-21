package policydecision_test

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/testifysec/judge/judge-api/ent/enttest"

	policy_decision "github.com/testifysec/judge/judge-api/policy/policy_decision"
)

func setupRequestAndResponse(body string) (*http.Request, *httptest.ResponseRecorder) {
	request := httptest.NewRequest(http.MethodPost, "/policy-decision", bytes.NewReader([]byte(body)))
	responseWriter := httptest.NewRecorder()
	return request, responseWriter
}

// This test makes sure that if a policy decision is available inside a cloud event (like we expect) we are able to parse it
func TestPostPolicy_WithValidCloudEvent(t *testing.T) {
	// Arrange
	body := `{
		"type": "policy decision",
		"source": "test module",
		"specversion": "1.0",
		"data": {
			"digest_id": "123456",
			"subject_name": "test",
			"decision": "allowed"
		}
	}`
	request, responseWriter := setupRequestAndResponse(body)
	request.Header.Set("Content-Type", "application/cloudevents+json")

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&_fk=1")
	defer client.Close()

	// Act
	policy_decision.PostPolicy(responseWriter, request, client)

	// Assert
	response := responseWriter.Result()
	bodyBytes, _ := ioutil.ReadAll(response.Body)
	assert.Equal(t, http.StatusOK, response.StatusCode, "Expected OK status code")
	assert.Equal(t, "", string(bodyBytes))
}

// This test makes sure that if something isn't wrapped in a cloud event properly, we won't accept it
func TestPostPolicy_WithoutCloudEvent(t *testing.T) {
	// Arrange
	body := `{ "subjectName": "", "digestId": "1234", "decision": "" }`
	request, responseWriter := setupRequestAndResponse(body)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&_fk=1")
	defer client.Close()

	// Act
	policy_decision.PostPolicy(responseWriter, request, client)

	// Assert
	assert.Equal(t, http.StatusBadRequest, responseWriter.Code, "Expected BadRequest status code")
}

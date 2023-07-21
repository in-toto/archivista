package policydecision

import (
	"net/http"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/testifysec/judge/judge-api/ent"
)

func PostPolicy(w http.ResponseWriter, r *http.Request, database *ent.Client) {
	ctx := r.Context()

	event, err := cloudevents.NewEventFromHTTPRequest(r)
	if err != nil {
		http.Error(w, "Failed to parse CloudEvent", http.StatusBadRequest)
		return
	}

	// Extract the policy decision from the event data
	var pd ent.PolicyDecision
	if err := event.DataAs(&pd); err != nil {
		http.Error(w, "Failed to extract policy decision", http.StatusInternalServerError)
		return
	}

	// Store the policy decision in the database
	if _, err := database.PolicyDecision.
		Create().
		SetSubjectName(pd.SubjectName).
		SetDigestID(pd.DigestID).
		SetDecision(pd.Decision).
		Save(ctx); err != nil {
		http.Error(w, "Failed to store policy decision: "+err.Error()+pd.SubjectName, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

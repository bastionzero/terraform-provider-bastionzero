package bzpolicygen

import (
	"github.com/bastionzero/bastionzero-sdk-go/bastionzero/service/policies"
	"pgregory.net/rapid"
)

func SessionRecordingPolicyGen() *rapid.Generator[policies.SessionRecordingPolicy] {
	return rapid.Custom(func(t *rapid.T) policies.SessionRecordingPolicy {
		return policies.SessionRecordingPolicy{
			Policy:      PolicyGen().Draw(t, "BasePolicy"),
			RecordInput: rapid.Ptr(rapid.Bool(), true).Draw(t, "RecordInput"),
		}
	})
}

package slackevents

type URLVerification struct {
	Challenge string `json:"challenge"`
}

func DefaultURLVerificationHandler(verification *URLVerification) (string, error) {
	return verification.Challenge, nil
}

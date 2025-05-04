package security

type Payload struct {
	Tier string `json:"tier"`
}

type Token interface {
	Generate(sub string, payload Payload) (token string, err error)
	Validate(token string) (sub string, tokenPayload Payload, err error)
}

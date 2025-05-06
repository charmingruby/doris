package privilege

import "errors"

var ErrInvalidTier = errors.New("invalid tier")

const (
	API_KEY_TIER_ROOKIE  = "ROOKIE"
	API_KEY_TIER_PRO     = "PRO"
	API_KEY_TIER_MANAGER = "MANAGER"
	API_KEY_TIER_ADMIN   = "ADMIN"
)

func ValidateAPIKeyTier(tier string) error {
	tiers := map[string]string{
		API_KEY_TIER_ROOKIE:  API_KEY_TIER_ROOKIE,
		API_KEY_TIER_PRO:     API_KEY_TIER_PRO,
		API_KEY_TIER_MANAGER: API_KEY_TIER_MANAGER,
		API_KEY_TIER_ADMIN:   API_KEY_TIER_ADMIN,
	}

	if _, ok := tiers[tier]; !ok {
		return ErrInvalidTier
	}

	return nil
}

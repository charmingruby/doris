package privilege

import (
	"errors"

	"github.com/charmingruby/doris/lib/delivery/proto/gen/account"
)

var ErrInvalidTier = errors.New("invalid tier")

const (
	TIER_ROOKIE  = "ROOKIE"
	TIER_PRO     = "PRO"
	TIER_MANAGER = "MANAGER"
	TIER_ADMIN   = "ADMIN"
)

var validTiers = map[string]struct{}{
	TIER_ROOKIE:  {},
	TIER_PRO:     {},
	TIER_MANAGER: {},
	TIER_ADMIN:   {},
}

var stringToProtoTier = map[string]account.Tier{
	TIER_ROOKIE:  account.Tier_ROOKIE,
	TIER_PRO:     account.Tier_PRO,
	TIER_MANAGER: account.Tier_MANAGER,
	TIER_ADMIN:   account.Tier_ADMIN,
}

var protoTierToString = map[account.Tier]string{
	account.Tier_ROOKIE:  TIER_ROOKIE,
	account.Tier_PRO:     TIER_PRO,
	account.Tier_MANAGER: TIER_MANAGER,
	account.Tier_ADMIN:   TIER_ADMIN,
}

func IsTierValid(tier string) error {
	if _, ok := validTiers[tier]; !ok {
		return ErrInvalidTier
	}

	return nil
}

func MapTierFromProto(tier account.Tier) (string, error) {
	str, ok := protoTierToString[tier]
	if !ok {
		return "", ErrInvalidTier
	}

	return str, nil
}

func MapTierToProto(tier string) (account.Tier, error) {
	proto, ok := stringToProtoTier[tier]
	if !ok {
		return account.Tier_UNSPECIFIED, ErrInvalidTier
	}

	return proto, nil
}

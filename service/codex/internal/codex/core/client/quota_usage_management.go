package client

import "context"

type QuotaUsageManagement interface {
	CheckQuotaAvailability(ctx context.Context, correlationID, kind string, usage int) (bool, error)
	ConsumeQuota(ctx context.Context, correlationID, kind string, usage int) error
}

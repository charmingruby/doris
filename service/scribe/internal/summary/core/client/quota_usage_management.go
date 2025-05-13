package client

import "context"

type QuotaUsageManagement interface {
	HasQuotaRemaining(ctx context.Context, correlationID, kind string, usage int) error
	ConsumeQuota(ctx context.Context, correlationID, kind string, usage int) error
}

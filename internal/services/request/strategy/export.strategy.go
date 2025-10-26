package strategy

import "context"

type ExportStrategy interface {
	Export(ctx context.Context, key string) error
}

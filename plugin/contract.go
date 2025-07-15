package plugin

import "context"

type Plugin interface {
	Run(ctx context.Context, input any) (output any, err error)
}

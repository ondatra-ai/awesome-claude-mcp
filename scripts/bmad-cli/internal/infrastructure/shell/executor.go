package shell

import "context"

type Executor interface {
	Run(ctx context.Context, name string, args ...string) (string, error)
	RunWithStdin(ctx context.Context, name, stdin string, args ...string) (string, error)
}

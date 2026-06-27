package workload

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

// WithCustom adds an arbitrary workload to the lifecycle group.  execute runs
// the workload; interrupt is called when the group starts shutting down (either
// because another workload returned or because WithGracefulShutdown fired).
// Panics inside execute are caught and converted to errors.
func (s *Server) WithCustom(execute func() error, interrupt func(interruptReason error)) *Server {
	return s.add(execute, interrupt)
}

func (s *Server) add(execute func() error, interrupt func(interruptReason error)) *Server {
	s.group.Add(func() (err error) {
		defer func() {
			if r := recover(); r != nil {
				err = fmt.Errorf("panic: %v", r)
			}
		}()

		return execute()
	},
		interrupt,
	)

	return s
}

// WithGracefulShutdown adds a workload that listens for SIGINT/SIGTERM and
// returns nil when either signal is received, triggering an orderly shutdown of
// the whole group.  ctx cancellation also triggers shutdown with ctx.Err().
func (s *Server) WithGracefulShutdown(ctx context.Context) *Server {
	cancelCh := make(chan os.Signal, 1)

	return s.add(func() error {
		signal.Notify(cancelCh, os.Interrupt, syscall.SIGTERM)
		defer signal.Stop(cancelCh)
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-cancelCh:
			return nil
		}
	}, func(error) {
		close(cancelCh)
	})
}

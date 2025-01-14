package ctrl

import (
	"context"
	"fmt"
	"sync"

	"sigs.k8s.io/controller-runtime/pkg/log"

	"github.com/pomerium/pomerium/config"
	pomerium_cmd "github.com/pomerium/pomerium/pkg/cmd/pomerium"

	"github.com/pomerium/ingress-controller/model"
	"github.com/pomerium/ingress-controller/pomerium"
)

var (
	_ = pomerium.ConfigReconciler(new(Runner))
)

// Runner implements pomerium control loop
type Runner struct {
	src  *InMemoryConfigSource
	base config.Config
	sync.Once
	ready chan struct{}
}

// waitForConfig waits until initial configuration is available
func (r *Runner) waitForConfig(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-r.ready:
	}
	return nil
}

func (r *Runner) readyToRun() {
	close(r.ready)
}

// GetConfig returns current configuration snapshot
func (r *Runner) GetConfig() *config.Config {
	return r.src.GetConfig()
}

// SetConfig updates just the shared config settings
func (r *Runner) SetConfig(ctx context.Context, src *model.Config) (changes bool, err error) {
	dst := r.base.Clone()

	if err := Apply(ctx, dst.Options, src); err != nil {
		return false, fmt.Errorf("transform config: %w", err)
	}

	changed := r.src.SetConfig(ctx, dst)
	r.Once.Do(r.readyToRun)

	return changed, nil
}

// NewPomeriumRunner creates new pomerium command and control
func NewPomeriumRunner(base config.Config, listener config.ChangeListener) (*Runner, error) {
	return &Runner{
		base: base,
		src: &InMemoryConfigSource{
			listeners: []config.ChangeListener{listener},
		},
		ready: make(chan struct{}),
	}, nil
}

// Run starts pomerium once config is available
func (r *Runner) Run(ctx context.Context) error {
	if err := r.waitForConfig(ctx); err != nil {
		return fmt.Errorf("waiting for pomerium bootstrap config: %w", err)
	}

	log.FromContext(ctx).V(1).Info("got bootstrap config, starting pomerium...", "cfg", r.src.GetConfig())

	return pomerium_cmd.Run(ctx, r.src)
}

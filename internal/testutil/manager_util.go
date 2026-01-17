package testutil

import (
	"context"
	"sync"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
	"k8s.io/utils/ptr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/config"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

// ManagerUtil is a utility interface for managing a controller-runtime manager
// You can use this to start and stop the manager in a goroutine.
type ManagerUtil interface {
	Start()
	Stop()
	GetManager() manager.Manager
}

type managerUtilImpl struct {
	mgr      manager.Manager
	errCh    chan error
	stopCh   chan struct{}
	stopOnce sync.Once
}

func NewManagerUtil(ctxRoot context.Context, restConfig *rest.Config, schema *runtime.Scheme) ManagerUtil {
	mgr, err := ctrl.NewManager(restConfig, ctrl.Options{
		Scheme: schema,
		Controller: config.Controller{
			SkipNameValidation: ptr.To(true),
		},
	})
	if err != nil {
		panic(err)
	}

	return &managerUtilImpl{
		mgr:    mgr,
		errCh:  make(chan error),
		stopCh: make(chan struct{}),
	}
}

// Start starts the manager in a goroutine.
func (m *managerUtilImpl) Start() {
	go func() {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		go func() {
			<-m.stopCh
			cancel()
		}()

		m.errCh <- m.mgr.Start(ctx)
	}()
}

// Stop stops the manager and waits for it to stop.
func (m *managerUtilImpl) Stop() {
	m.stopOnce.Do(func() {
		close(m.stopCh)
	})
	<-m.errCh
}

func (m *managerUtilImpl) GetManager() manager.Manager {
	return m.mgr
}

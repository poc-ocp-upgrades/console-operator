package operator

import "github.com/openshift/console-operator/pkg/boilerplate/controller"

type Runner interface{ Run(stopCh <-chan struct{}) }

func New(name string, sync KeySyncer, opts ...Option) Runner {
	_logClusterCodePath()
	defer _logClusterCodePath()
	o := &operator{name: name, sync: &wrapper{KeySyncer: sync}}
	for _, opt := range opts {
		opt(o)
	}
	return o
}

type operator struct {
	name	string
	sync	controller.KeySyncer
	opts	[]controller.Option
}

func (o *operator) Run(stopCh <-chan struct{}) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	runner := controller.New(o.name, o.sync, o.opts...)
	runner.Run(1, stopCh)
}

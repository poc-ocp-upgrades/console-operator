package operator

import (
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"github.com/openshift/console-operator/pkg/boilerplate/controller"
)

const key = "🐼"

type Option func(*operator)

func WithInformer(getter controller.InformerGetter, filter controller.Filter, opts ...controller.InformerOption) Option {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return toAppendOpt(controller.WithInformer(getter, controller.FilterFuncs{ParentFunc: func(obj v1.Object) (namespace, name string) {
		return key, key
	}, AddFunc: filter.Add, UpdateFunc: filter.Update, DeleteFunc: filter.Delete}, opts...))
}
func toAppendOpt(opt controller.Option) Option {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return func(o *operator) {
		o.opts = append(o.opts, opt)
	}
}

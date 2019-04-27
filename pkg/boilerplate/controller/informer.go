package controller

type InformerOption func() informerOptionCase
type informerOptionCase int

const (
	syncDefault	informerOptionCase	= iota
	noSync
)

func WithNoSync() InformerOption {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return func() informerOptionCase {
		return noSync
	}
}
func withSync() InformerOption {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return func() informerOptionCase {
		return syncDefault
	}
}
func informerOptionToOption(opt InformerOption, getter InformerGetter) Option {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	switch opt() {
	case syncDefault:
		return WithInformerSynced(getter)
	case noSync:
		return func(*controller) {
		}
	default:
		panic(opt)
	}
}

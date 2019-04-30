package controller

type die string

func crash(i interface{}) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	mustDie, ok := i.(die)
	if ok {
		panic(string(mustDie))
	}
}

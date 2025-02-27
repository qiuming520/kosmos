package options

import (
	"time"

	"github.com/spf13/pflag"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/client-go/tools/leaderelection/resourcelock"
	componentbaseconfig "k8s.io/component-base/config"

	"github.com/kosmos.io/kosmos/pkg/utils"
)

var (
	defaultElectionLeaseDuration = metav1.Duration{Duration: 15 * time.Second}
	defaultElectionRenewDeadline = metav1.Duration{Duration: 10 * time.Second}
	defaultElectionRetryPeriod   = metav1.Duration{Duration: 2 * time.Second}
)

type Options struct {
	LeaderElection componentbaseconfig.LeaderElectionConfiguration
	utils.KubernetesOptions
}

// NewOptions builds a default elector options.
func NewOptions() *Options {
	return &Options{
		LeaderElection: componentbaseconfig.LeaderElectionConfiguration{
			LeaderElect:       true,
			ResourceLock:      resourcelock.LeasesResourceLock,
			ResourceNamespace: utils.DefaultNamespace,
			ResourceName:      "elector",
		},
	}
}

func (o *Options) Validate() field.ErrorList {
	errs := field.ErrorList{}
	newPath := field.NewPath("Options")
	if len(o.KubernetesOptions.ControlpanelKubeConfig) == 0 {
		errs = append(errs, field.Invalid(newPath.Child("controlpanelconfig"), o.ControlpanelKubeConfig, "controlpanelconfig path should not empty"))
	}
	return errs
}

func (o *Options) AddFlags(fs *pflag.FlagSet) {
	fs.Float32Var(&o.KubernetesOptions.QPS, "kube-qps", utils.DefaultKubeQPS, "QPS to use while talking with kube-apiserver.")
	fs.IntVar(&o.KubernetesOptions.Burst, "kube-burst", utils.DefaultKubeBurst, "Burst to use while talking with kube-apiserver.")
	fs.StringVar(&o.KubernetesOptions.KubeConfig, "kubeconfig", "", "Path to host kubeconfig file.")
	fs.StringVar(&o.KubernetesOptions.MasterURL, "master-url", "", "Used to generate kubeconfig for downloading, if not specified, will use host in kubeconfig.")
	fs.StringVar(&o.KubernetesOptions.ControlpanelKubeConfig, "controlpanel-kubeconfig", "", "Path to control plane kubeconfig file.")
	fs.StringVar(&o.KubernetesOptions.ControlpanelMasterURL, "controlpanel-master-url", "", "Used to generate host control plane kubeconfig for downloading, if not specified, will use host in controlpanel-kubeconfig.")
	fs.BoolVar(&o.LeaderElection.LeaderElect, "leader-elect", true, "Enable leader election, which must be true when running multi instances.")
	fs.StringVar(&o.LeaderElection.ResourceName, "leader-elect-resource-name", "elector", "The name of resource object that is used for locking during leader election.")
	fs.StringVar(&o.LeaderElection.ResourceNamespace, "leader-elect-resource-namespace", utils.DefaultNamespace, "The namespace of resource object that is used for locking during leader election.")
	fs.DurationVar(&o.LeaderElection.LeaseDuration.Duration, "leader-elect-lease-duration", defaultElectionLeaseDuration.Duration, ""+
		"The duration that non-leader candidates will wait after observing a leadership "+
		"renewal until attempting to acquire leadership of a led but unrenewed leader "+
		"slot. This is effectively the maximum duration that a leader can be stopped "+
		"before it is replaced by another candidate. This is only applicable if leader "+
		"election is enabled.")
	fs.DurationVar(&o.LeaderElection.RenewDeadline.Duration, "leader-elect-renew-deadline", defaultElectionRenewDeadline.Duration, ""+
		"The interval between attempts by the acting master to renew a leadership slot "+
		"before it stops leading. This must be less than or equal to the lease duration. "+
		"This is only applicable if leader election is enabled.")
	fs.DurationVar(&o.LeaderElection.RetryPeriod.Duration, "leader-elect-retry-period", defaultElectionRetryPeriod.Duration, ""+
		"The duration the clients should wait between attempting acquisition and renewal "+
		"of a leadership. This is only applicable if leader election is enabled.")
}

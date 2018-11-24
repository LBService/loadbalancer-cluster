package constants

import "time"

const (
	DefaultDialTimeout    = 5 * time.Second
	DefaultRequestTimeout = 5 * time.Second
	// DefaultBackupTimeout is the default maximal allowed time of the entire backup process.
	DefaultBackupTimeout    = 1 * time.Minute
	DefaultSnapshotInterval = 1800 * time.Second

	DefaultBackupPodHTTPPort = 19999

	OperatorRoot = "/var/tmp/lbaas-operator"

	EnvOperatorPodName      = "MY_POD_NAME"
	EnvOperatorPodNamespace = "MY_POD_NAMESPACE"
)

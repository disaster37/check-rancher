package rancherapi

type API interface {
	Cluster() ClusterAPI
	ETCDBackup() ETCDBackupAPI
}

type ClusterAPI interface {
	List() ([]*Cluster, error)
	GetByName(name string) (*Cluster, error)
}

type ETCDBackupAPI interface {
	List() ([]*ETCDBackup, error)
	ListByClusterID(clusterID string) ([]*ETCDBackup, error)
}

package ent

import (
	v1alpha "github.com/nokamoto/2pf23/pkg/api/ke/v1alpha"
)

func ClusterProto(x *Cluster) *v1alpha.Cluster {
	return &v1alpha.Cluster{
		Name:        x.Name,
		DisplayName: x.DisplayName,
		NumNodes:    x.NumNodes,
		MachineType: v1alpha.MachineType(x.MachineType),
	}
}

func ClusterCreateQuery(create *ClusterCreate, cluster *v1alpha.Cluster) *ClusterCreate {
	return create.SetName(cluster.GetName()).
		SetDisplayName(cluster.GetDisplayName()).
		SetNumNodes(cluster.GetNumNodes()).
		SetMachineType(int32(cluster.GetMachineType()))
}

func ClusterUpdateOneQuery(update *ClusterUpdateOne, cluster *v1alpha.Cluster) *ClusterUpdateOne {
	return update.SetDisplayName(cluster.GetDisplayName()).
		SetNumNodes(cluster.GetNumNodes()).
		SetMachineType(int32(cluster.GetMachineType()))
}

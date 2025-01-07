package controller

import (
	"github.com/tdevsin/keyforge/internal/cluster"
	"github.com/tdevsin/keyforge/internal/config"
	"github.com/tdevsin/keyforge/internal/proto"
	"go.uber.org/zap"
)

func GetClusterInfo(c *config.Config) (*proto.ClusterState, error) {
	var state proto.ClusterState
	c.ClusterInfo.GetClusterInfo().MapClusterStateToProto(&state)
	return &state, nil
}

func SetClusterInfo(c *config.Config, state *proto.ClusterState) error {
	c.Logger.Info("Setting cluster state", zap.Any("state", state))
	c.ClusterInfo.MergeClusterState(MapProtoToClusterInfo(state))
	return nil
}

func MapProtoToClusterInfo(state *proto.ClusterState) *cluster.ClusterInfo {
	ci := cluster.NewCluster("", 2)
	ci.Version = int(state.Version)
	ci.LastUpdated = state.LastUpdated.AsTime()
	ci.Nodes = make(map[string]cluster.Node)
	for _, node := range state.Nodes {
		ci.Nodes[node.Id] = cluster.Node{
			ID:      node.Id,
			Address: node.Address,
			Health: cluster.Health{
				LastChecked: node.Health.LastUpdated.AsTime(),
				Status:      cluster.Status(node.Health.Status),
			},
		}
	}
	return ci
}

package controller

import (
	"github.com/tdevsin/keyforge/internal/cluster"
	"github.com/tdevsin/keyforge/internal/config"
	"github.com/tdevsin/keyforge/internal/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func GetClusterInfo(c *config.Config) (*proto.ClusterState, error) {
	var state proto.ClusterState
	MapClusterStateToProto(&state, c)
	return &state, nil
}

func SetClusterInfo(c *config.Config, state *proto.ClusterState) error {
	c.ClusterInfo.MergeClusterState(MapProtoToClusterInfo(state))
	return nil
}

func MapClusterStateToProto(state *proto.ClusterState, c *config.Config) {
	localClusterInfo := c.ClusterInfo.GetClusterInfo()
	state.Version = int64(localClusterInfo.Version)
	state.LastUpdated = timestamppb.New(localClusterInfo.LastUpdated)
	state.Nodes = make([]*proto.Node, 0, len(localClusterInfo.Nodes))
	for _, node := range localClusterInfo.Nodes {
		state.Nodes = append(state.Nodes, &proto.Node{
			Id:      node.ID,
			Address: node.Address,
			Health: &proto.Health{
				LastUpdated: timestamppb.New(node.Health.LastChecked),
				Status:      proto.Status(node.Health.Status),
			},
		})
	}
}

func MapProtoToClusterInfo(state *proto.ClusterState) *cluster.ClusterInfo {
	ci := cluster.NewCluster()
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

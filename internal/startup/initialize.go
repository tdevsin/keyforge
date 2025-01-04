package startup

import (
	"context"

	"github.com/tdevsin/keyforge/internal/api/controller"
	"github.com/tdevsin/keyforge/internal/config"
	"github.com/tdevsin/keyforge/internal/proto"
	"github.com/tdevsin/keyforge/internal/utils"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
)

// StartNodeSetupInCluster initializes the node setup in the cluster and perform necessary operations
func StartNodeSetupInCluster(conf *config.Config, bootstrapNodeAddress string) error {
	// If no bootstrap node is passed, this is the first node in the cluster
	if utils.IsEmpty(bootstrapNodeAddress) {
		conf.Logger.Info("No bootstrap node provided. This is the first node in the cluster")

		// Update the versioning since this is first node in the cluster
		conf.ClusterInfo.IncrementVersion()
		return nil
	}

	// If bootstrap node is provided, join the cluster
	conf.Logger.Info("Joining existing cluster", zap.String("bootstrapNodeAddress", bootstrapNodeAddress))

	// Create client for calling bootstrap node
	conn, err := grpc.NewClient(bootstrapNodeAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	defer conn.Close()
	if err != nil {
		panic(err)
	}
	client := proto.NewClusterServiceClient(conn)

	// Get the cluster state from the bootstrap node
	clusterState, err := client.GetClusterState(context.TODO(), &emptypb.Empty{})
	if err != nil {
		panic(err)
	}

	// Merge the cluster state with the local cluster state
	conf.ClusterInfo.MergeClusterState(controller.MapProtoToClusterInfo(clusterState))

	// Update the versioning
	conf.ClusterInfo.IncrementVersion()

	// Send the updated cluster state to the bootstrap node
	var req proto.ClusterState
	conf.ClusterInfo.GetClusterInfo().MapClusterStateToProto(&req)
	_, err = client.SetClusterState(context.TODO(), &req)

	// Panic if state update fails
	if err != nil {
		panic(err)
	}

	conf.Logger.Info("Joined the cluster", zap.String("bootstrapNodeAddress", bootstrapNodeAddress))
	return nil
}

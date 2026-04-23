package sdk

import (
	"context"
	"crypto/tls"
	"fmt"
	pb "github.com/Royal17x/flagr/backend/gen/proto/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"time"
)

type grpcEvaluator struct {
	conn    *grpc.ClientConn
	client  pb.FlagServiceClient
	sdkKey  string
	timeout time.Duration
}

func newGRPCEvaluator(serverAddr, sdkKey string, timeout time.Duration, tlsEnabled bool) (*grpcEvaluator, error) {
	var creds grpc.DialOption
	if tlsEnabled {
		creds = grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{
			MinVersion: tls.VersionTLS12,
		}))
	} else {
		creds = grpc.WithTransportCredentials(insecure.NewCredentials())
	}

	conn, err := grpc.NewClient(serverAddr, creds)
	if err != nil {
		return nil, fmt.Errorf("grpc connect to %s: %w", serverAddr, err)
	}
	return &grpcEvaluator{
		conn:    conn,
		client:  pb.NewFlagServiceClient(conn),
		sdkKey:  sdkKey,
		timeout: timeout,
	}, nil
}
func (e *grpcEvaluator) Evaluate(ctx context.Context, req EvaluateRequest) (EvaluateResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, e.timeout)
	defer cancel()

	ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs(
		"X-SDK-Key", e.sdkKey,
	))

	resp, err := e.client.EvaluateFlag(ctx, &pb.EvaluateFlagRequest{
		FlagKey:       req.FlagKey,
		ProjectId:     req.ProjectID,
		EnvironmentId: req.EnvironmentID,
		Context:       req.Context,
	})
	if err != nil {
		return EvaluateResponse{}, fmt.Errorf("grpc.EvaluateFlag: %w", err)
	}
	return EvaluateResponse{
		Enabled:          resp.Enabled,
		Reason:           resp.Reason,
		EvaluationTimeMs: resp.EvaluationTimeMs,
	}, nil

}
func (e *grpcEvaluator) Close() error { return e.conn.Close() }

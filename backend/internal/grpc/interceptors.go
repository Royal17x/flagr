package grpc

import (
	"context"
	"log/slog"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/Royal17x/flagr/backend/internal/port"
)

type grpcContextKey string

const claimsKey grpcContextKey = "claims"

func UnaryLoggingInterceptor(
	ctx context.Context,
	req any,
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (any, error) {
	start := time.Now()
	resp, err := handler(ctx, req)

	slog.Info("grpc request",
		"method", info.FullMethod,
		"duration_ms", time.Since(start).Milliseconds(),
		"error", err,
	)
	return resp, err
}

func UnaryAuthInterceptor(authSvc port.AuthServiceInterface) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		if info.FullMethod == "/flagr.v1.FlagService/EvaluateFlag" ||
			info.FullMethod == "/flagr.v1.FlagService/EvaluateBatch" {
			return handler(ctx, req)
		}

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "missing metadata")
		}

		authHeader := md.Get("authorization")
		if len(authHeader) == 0 {
			return nil, status.Error(codes.Unauthenticated, "missing authorization")
		}

		parts := strings.SplitN(authHeader[0], " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			return nil, status.Error(codes.Unauthenticated, "invalid header")
		}
		claims, err := authSvc.ValidateAccessToken(parts[1])
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, "invalid credentials")
		}
		ctx = context.WithValue(ctx, claimsKey, claims)
		return handler(ctx, req)
	}
}

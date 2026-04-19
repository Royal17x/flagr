package grpc

import (
	"context"
	pb "github.com/Royal17x/flagr/backend/gen/proto/v1"
	"github.com/Royal17x/flagr/backend/internal/port"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log/slog"
	"sync"
	"sync/atomic"
	"time"
)

type FlagServer struct {
	pb.UnimplementedFlagServiceServer
	flagSvc port.FlagServiceInterface

	mu             sync.RWMutex
	subscribers    map[string][]chan *pb.FlagUpdate
	droppedUpdates atomic.Int64
}

func NewFlagServer(flagSvc port.FlagServiceInterface) *FlagServer {
	return &FlagServer{
		flagSvc:     flagSvc,
		subscribers: make(map[string][]chan *pb.FlagUpdate),
	}
}

func (s *FlagServer) EvaluateFlag(ctx context.Context, req *pb.EvaluateFlagRequest) (*pb.EvaluateFlagResponse, error) {
	start := time.Now()

	if req.FlagKey == "" {
		return nil, status.Error(codes.InvalidArgument, "flag_key is required")
	}

	projectID, err := uuid.Parse(req.ProjectId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid project_id")
	}

	environmentID, err := uuid.Parse(req.EnvironmentId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid environment_id")
	}

	enabled, err := s.flagSvc.EvaluateFlag(ctx, req.FlagKey, projectID, environmentID)
	if err != nil {
		return nil, domainErrToGRPC(err)
	}

	return &pb.EvaluateFlagResponse{
		Enabled:          enabled,
		FlagKey:          req.FlagKey,
		Reason:           "DB_LOOKUP",
		EvaluationTimeMs: time.Since(start).Milliseconds(),
	}, nil
}

func (s *FlagServer) EvaluateBatch(ctx context.Context, req *pb.EvaluateBatchRequest) (*pb.EvaluateBatchResponse, error) {
	responses := make([]*pb.EvaluateFlagResponse, 0, len(req.Requests))
	for _, request := range req.Requests {
		start := time.Now()

		errorResponse := func(reason string) *pb.EvaluateFlagResponse {
			return &pb.EvaluateFlagResponse{
				Enabled:          false,
				FlagKey:          request.FlagKey,
				Reason:           reason,
				EvaluationTimeMs: time.Since(start).Milliseconds(),
			}
		}

		if request.FlagKey == "" {
			responses = append(responses, errorResponse("ERROR: flag_key is required"))
			continue
		}

		projectID, err := uuid.Parse(request.ProjectId)
		if err != nil {
			responses = append(responses, errorResponse("ERROR: invalid project_id"))
			continue
		}

		environmentID, err := uuid.Parse(request.EnvironmentId)
		if err != nil {
			responses = append(responses, errorResponse("ERROR: invalid environment_id"))
			continue
		}

		enabled, err := s.flagSvc.EvaluateFlag(ctx, request.FlagKey, projectID, environmentID)
		if err != nil {
			responses = append(responses, errorResponse("ERROR"))
			continue
		}

		responses = append(responses, &pb.EvaluateFlagResponse{
			Enabled:          enabled,
			FlagKey:          request.FlagKey,
			Reason:           "DB_LOOKUP",
			EvaluationTimeMs: time.Since(start).Milliseconds(),
		})
	}
	return &pb.EvaluateBatchResponse{Responses: responses}, nil
}

func (s *FlagServer) WatchFlags(req *pb.WatchFlagsRequest, stream pb.FlagService_WatchFlagsServer) error {
	key := req.ProjectId + ":" + req.EnvironmentId
	subChan := make(chan *pb.FlagUpdate, 64)
	s.mu.Lock()
	s.subscribers[key] = append(s.subscribers[key], subChan)
	s.mu.Unlock()

	defer func() {
		s.mu.Lock()
		channels := s.subscribers[key]
		for i, ch := range channels {
			if ch == subChan {
				s.subscribers[key] = append(channels[:i], channels[i+1:]...)
				break
			}
		}
		s.mu.Unlock()
		close(subChan)
	}()

	for {
		select {
		case <-stream.Context().Done():
			return stream.Context().Err()
		case update, ok := <-subChan:
			if !ok {
				return nil
			}
			if err := stream.Send(update); err != nil {
				return err
			}
		}
	}
}

func (s *FlagServer) NotifySubscribers(projectID, envID string, update *pb.FlagUpdate) {
	key := projectID + ":" + envID
	s.mu.RLock()
	subs, ok := s.subscribers[key]
	if !ok || len(subs) == 0 {
		s.mu.RUnlock()
		return
	}

	channels := make([]chan *pb.FlagUpdate, len(subs))
	copy(channels, subs)
	s.mu.RUnlock()

	for _, ch := range channels {
		select {
		case ch <- update:
		default:
			// TODO : Prometheus counter
			s.droppedUpdates.Add(1)
			slog.Warn("kafka: dropped update for slow subscriber",
				"key", key,
				"total_dropped", s.droppedUpdates.Load(),
			)
		}
	}
}

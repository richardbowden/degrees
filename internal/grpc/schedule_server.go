package grpc

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/richardbowden/degrees/internal/dbpg"
	pb "github.com/richardbowden/degrees/internal/pb/degrees/v1"
	"github.com/richardbowden/degrees/internal/services"
)

type ScheduleServiceServer struct {
	pb.UnimplementedScheduleServiceServer
	scheduleSvc *services.ScheduleService
}

func NewScheduleServer(scheduleSvc *services.ScheduleService) *ScheduleServiceServer {
	return &ScheduleServiceServer{
		scheduleSvc: scheduleSvc,
	}
}

func (s *ScheduleServiceServer) GetScheduleConfig(ctx context.Context, req *pb.GetScheduleConfigRequest) (*pb.GetScheduleConfigResponse, error) {
	configs, err := s.scheduleSvc.GetScheduleConfig(ctx)
	if err != nil {
		return nil, ToGRPCError(err)
	}

	days := make([]*pb.ScheduleDay, len(configs))
	for i, c := range configs {
		days[i] = scheduleConfigToProto(&c)
	}

	return &pb.GetScheduleConfigResponse{
		Days: days,
	}, nil
}

func (s *ScheduleServiceServer) UpdateScheduleConfig(ctx context.Context, req *pb.UpdateScheduleConfigRequest) (*pb.UpdateScheduleConfigResponse, error) {
	openTime, err := parseTimeString(req.OpenTime)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid open_time format, expected HH:MM")
	}

	closeTime, err := parseTimeString(req.CloseTime)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid close_time format, expected HH:MM")
	}

	cfg, err := s.scheduleSvc.UpdateScheduleConfig(ctx, dbpg.UpdateScheduleConfigParams{
		DayOfWeek:     req.DayOfWeek,
		OpenTime:      openTime,
		CloseTime:     closeTime,
		IsOpen:        req.IsOpen,
		BufferMinutes: req.BufferMinutes,
	})
	if err != nil {
		return nil, ToGRPCError(err)
	}

	return &pb.UpdateScheduleConfigResponse{
		Day: scheduleConfigToProto(cfg),
	}, nil
}

func (s *ScheduleServiceServer) AddBlackout(ctx context.Context, req *pb.AddBlackoutRequest) (*pb.AddBlackoutResponse, error) {
	if req.Date == "" {
		return nil, status.Error(codes.InvalidArgument, "date is required")
	}

	blackout, err := s.scheduleSvc.AddBlackout(ctx, req.Date, req.Reason)
	if err != nil {
		return nil, ToGRPCError(err)
	}

	return &pb.AddBlackoutResponse{
		Blackout: &pb.Blackout{
			Id:     blackout.ID,
			Date:   formatPGDate(blackout.Date),
			Reason: blackout.Reason.String,
		},
	}, nil
}

func (s *ScheduleServiceServer) RemoveBlackout(ctx context.Context, req *pb.RemoveBlackoutRequest) (*pb.RemoveBlackoutResponse, error) {
	if req.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	err := s.scheduleSvc.RemoveBlackout(ctx, req.Id)
	if err != nil {
		return nil, ToGRPCError(err)
	}

	return &pb.RemoveBlackoutResponse{
		Success: true,
	}, nil
}

func scheduleConfigToProto(c *dbpg.ScheduleConfig) *pb.ScheduleDay {
	return &pb.ScheduleDay{
		Id:            c.ID,
		DayOfWeek:     c.DayOfWeek,
		OpenTime:      formatPGTime(c.OpenTime),
		CloseTime:     formatPGTime(c.CloseTime),
		IsOpen:        c.IsOpen,
		BufferMinutes: c.BufferMinutes,
	}
}

func parseTimeString(timeStr string) (pgtype.Time, error) {
	t, err := time.Parse("15:04", timeStr)
	if err != nil {
		return pgtype.Time{}, fmt.Errorf("invalid time: %s", timeStr)
	}
	microseconds := int64(t.Hour())*3600000000 + int64(t.Minute())*60000000
	return pgtype.Time{Microseconds: microseconds, Valid: true}, nil
}

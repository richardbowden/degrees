package grpc

import (
	"context"
	"encoding/json"
	"fmt"

	pb "github.com/richardbowden/degrees/internal/pb/degrees/v1"
	"github.com/richardbowden/degrees/internal/repos"
	"github.com/richardbowden/degrees/internal/services"
	"github.com/richardbowden/degrees/internal/settings"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/structpb"
)

// SettingsServiceServer implements the gRPC SettingsService
type SettingsServiceServer struct {
	pb.UnimplementedSettingsServiceServer
	settingsService *settings.Service
	settingsRepo    *repos.Settings
	authSvc         *services.AuthN
}

// NewSettingsServiceServer creates a new settings service gRPC server
func NewSettingsServiceServer(settingsService *settings.Service, settingsRepo *repos.Settings, authSvc *services.AuthN) *SettingsServiceServer {
	return &SettingsServiceServer{
		settingsService: settingsService,
		settingsRepo:    settingsRepo,
		authSvc:         authSvc,
	}
}

// GetSystemSettings retrieves all settings for a subsystem at system scope
func (s *SettingsServiceServer) GetSystemSettings(ctx context.Context, req *pb.GetSystemSettingsRequest) (*pb.GetSettingsResponse, error) {
	if req.Subsystem == "" {
		return nil, status.Error(codes.InvalidArgument, "subsystem is required")
	}

	settingsMap, err := s.settingsService.GetBySubsystem(ctx, req.Subsystem, settings.SystemScope())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get settings: %v", err)
	}

	// Convert map[string]any to map[string]*structpb.Value
	protoSettings, err := convertMapToProtoStruct(settingsMap)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to convert settings: %v", err)
	}

	return &pb.GetSettingsResponse{
		Subsystem: req.Subsystem,
		Settings:  protoSettings,
	}, nil
}

// SetSystemSetting sets a system-level setting
func (s *SettingsServiceServer) SetSystemSetting(ctx context.Context, req *pb.SetSystemSettingRequest) (*pb.SetSettingResponse, error) {
	if err := RequireSysop(ctx, s.authSvc); err != nil {
		return nil, err
	}

	if req.Subsystem == "" || req.Key == "" {
		return nil, status.Error(codes.InvalidArgument, "subsystem and key are required")
	}
	if req.Value == nil {
		return nil, status.Error(codes.InvalidArgument, "value is required")
	}

	// Get user ID from context (set by auth interceptor)
	userID := getUserIDFromGRPCContext(ctx)

	// Convert protobuf Value to Go value
	value, err := structpb.NewValue(req.Value.AsInterface())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid value: %v", err)
	}

	var description *string
	if req.Description != nil {
		description = req.Description
	}

	err = s.settingsService.SetSystem(ctx, req.Subsystem, req.Key, value.AsInterface(), description, userID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to set setting: %v", err)
	}

	// Build response setting
	setting := &pb.Setting{
		Scope:       pb.SettingScope_SETTING_SCOPE_SYSTEM,
		Subsystem:   req.Subsystem,
		Key:         req.Key,
		Value:       req.Value,
		Description: ptrToString(description),
	}

	return &pb.SetSettingResponse{
		Success: true,
		Setting: setting,
	}, nil
}

// GetOrganizationSettings retrieves all settings for a subsystem at organization scope
func (s *SettingsServiceServer) GetOrganizationSettings(ctx context.Context, req *pb.GetOrganizationSettingsRequest) (*pb.GetSettingsResponse, error) {
	if req.OrganizationId == 0 {
		return nil, status.Error(codes.InvalidArgument, "organization_id is required")
	}
	if req.Subsystem == "" {
		return nil, status.Error(codes.InvalidArgument, "subsystem is required")
	}

	settingsMap, err := s.settingsService.GetBySubsystem(ctx, req.Subsystem, settings.OrganizationScope(req.OrganizationId))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get settings: %v", err)
	}

	protoSettings, err := convertMapToProtoStruct(settingsMap)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to convert settings: %v", err)
	}

	return &pb.GetSettingsResponse{
		Subsystem: req.Subsystem,
		Settings:  protoSettings,
	}, nil
}

// SetOrganizationSetting sets an organization-level setting
func (s *SettingsServiceServer) SetOrganizationSetting(ctx context.Context, req *pb.SetOrganizationSettingRequest) (*pb.SetSettingResponse, error) {
	if err := RequireSysop(ctx, s.authSvc); err != nil {
		return nil, err
	}

	if req.OrganizationId == 0 {
		return nil, status.Error(codes.InvalidArgument, "organization_id is required")
	}
	if req.Subsystem == "" || req.Key == "" {
		return nil, status.Error(codes.InvalidArgument, "subsystem and key are required")
	}
	if req.Value == nil {
		return nil, status.Error(codes.InvalidArgument, "value is required")
	}

	userID := getUserIDFromGRPCContext(ctx)

	value, err := structpb.NewValue(req.Value.AsInterface())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid value: %v", err)
	}

	var description *string
	if req.Description != nil {
		description = req.Description
	}

	err = s.settingsService.SetOrganization(ctx, req.OrganizationId, req.Subsystem, req.Key, value.AsInterface(), description, userID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to set setting: %v", err)
	}

	setting := &pb.Setting{
		Scope:          pb.SettingScope_SETTING_SCOPE_ORGANIZATION,
		OrganizationId: req.OrganizationId,
		Subsystem:      req.Subsystem,
		Key:            req.Key,
		Value:          req.Value,
		Description:    ptrToString(description),
	}

	return &pb.SetSettingResponse{
		Success: true,
		Setting: setting,
	}, nil
}

// GetProjectSettings retrieves all settings for a subsystem at project scope
func (s *SettingsServiceServer) GetProjectSettings(ctx context.Context, req *pb.GetProjectSettingsRequest) (*pb.GetSettingsResponse, error) {
	if req.ProjectId == 0 {
		return nil, status.Error(codes.InvalidArgument, "project_id is required")
	}
	if req.Subsystem == "" {
		return nil, status.Error(codes.InvalidArgument, "subsystem is required")
	}

	settingsMap, err := s.settingsService.GetBySubsystem(ctx, req.Subsystem, settings.ProjectScope(req.ProjectId))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get settings: %v", err)
	}

	protoSettings, err := convertMapToProtoStruct(settingsMap)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to convert settings: %v", err)
	}

	return &pb.GetSettingsResponse{
		Subsystem: req.Subsystem,
		Settings:  protoSettings,
	}, nil
}

// SetProjectSetting sets a project-level setting
func (s *SettingsServiceServer) SetProjectSetting(ctx context.Context, req *pb.SetProjectSettingRequest) (*pb.SetSettingResponse, error) {
	if err := RequireSysop(ctx, s.authSvc); err != nil {
		return nil, err
	}

	if req.ProjectId == 0 {
		return nil, status.Error(codes.InvalidArgument, "project_id is required")
	}
	if req.Subsystem == "" || req.Key == "" {
		return nil, status.Error(codes.InvalidArgument, "subsystem and key are required")
	}
	if req.Value == nil {
		return nil, status.Error(codes.InvalidArgument, "value is required")
	}

	userID := getUserIDFromGRPCContext(ctx)

	value, err := structpb.NewValue(req.Value.AsInterface())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid value: %v", err)
	}

	var description *string
	if req.Description != nil {
		description = req.Description
	}

	err = s.settingsService.SetProject(ctx, req.ProjectId, req.Subsystem, req.Key, value.AsInterface(), description, userID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to set setting: %v", err)
	}

	setting := &pb.Setting{
		Scope:       pb.SettingScope_SETTING_SCOPE_PROJECT,
		ProjectId:   req.ProjectId,
		Subsystem:   req.Subsystem,
		Key:         req.Key,
		Value:       req.Value,
		Description: ptrToString(description),
	}

	return &pb.SetSettingResponse{
		Success: true,
		Setting: setting,
	}, nil
}

// GetUserSettings retrieves all settings for a subsystem at user scope
func (s *SettingsServiceServer) GetUserSettings(ctx context.Context, req *pb.GetUserSettingsRequest) (*pb.GetSettingsResponse, error) {
	if req.UserId == 0 {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}
	if req.Subsystem == "" {
		return nil, status.Error(codes.InvalidArgument, "subsystem is required")
	}

	settingsMap, err := s.settingsService.GetBySubsystem(ctx, req.Subsystem, settings.UserScope(req.UserId))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get settings: %v", err)
	}

	protoSettings, err := convertMapToProtoStruct(settingsMap)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to convert settings: %v", err)
	}

	return &pb.GetSettingsResponse{
		Subsystem: req.Subsystem,
		Settings:  protoSettings,
	}, nil
}

// SetUserSetting sets a user-level setting
func (s *SettingsServiceServer) SetUserSetting(ctx context.Context, req *pb.SetUserSettingRequest) (*pb.SetSettingResponse, error) {
	if err := RequireSysop(ctx, s.authSvc); err != nil {
		return nil, err
	}

	if req.UserId == 0 {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}
	if req.Subsystem == "" || req.Key == "" {
		return nil, status.Error(codes.InvalidArgument, "subsystem and key are required")
	}
	if req.Value == nil {
		return nil, status.Error(codes.InvalidArgument, "value is required")
	}

	userID := getUserIDFromGRPCContext(ctx)

	value, err := structpb.NewValue(req.Value.AsInterface())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid value: %v", err)
	}

	var description *string
	if req.Description != nil {
		description = req.Description
	}

	err = s.settingsService.SetUser(ctx, req.UserId, req.Subsystem, req.Key, value.AsInterface(), description, userID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to set setting: %v", err)
	}

	setting := &pb.Setting{
		Scope:       pb.SettingScope_SETTING_SCOPE_USER,
		UserId:      req.UserId,
		Subsystem:   req.Subsystem,
		Key:         req.Key,
		Value:       req.Value,
		Description: ptrToString(description),
	}

	return &pb.SetSettingResponse{
		Success: true,
		Setting: setting,
	}, nil
}

// DeleteSetting deletes a setting by ID
func (s *SettingsServiceServer) DeleteSetting(ctx context.Context, req *pb.DeleteSettingRequest) (*pb.DeleteSettingResponse, error) {
	if err := RequireSysop(ctx, s.authSvc); err != nil {
		return nil, err
	}

	if req.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	err := s.settingsService.Delete(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete setting: %v", err)
	}

	return &pb.DeleteSettingResponse{
		Success: true,
	}, nil
}

// ListAllSettings lists all settings (for admin interface)
func (s *SettingsServiceServer) ListAllSettings(ctx context.Context, req *pb.ListAllSettingsRequest) (*pb.ListAllSettingsResponse, error) {
	if err := RequireSysop(ctx, s.authSvc); err != nil {
		return nil, err
	}

	// Get all settings from repository via service
	domainSettings, err := s.settingsService.ListAll(ctx, s.settingsRepo)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list settings: %v", err)
	}

	// Apply filters and convert to proto
	var filteredSettings []*pb.Setting
	for _, setting := range domainSettings {
		// Apply subsystem filter
		if req.Subsystem != nil && setting.Subsystem != *req.Subsystem {
			continue
		}

		// Apply scope filter
		pbScope := convertScopeToPB(setting.Scope)
		if req.Scope != nil && pbScope != *req.Scope {
			continue
		}

		// Apply organization ID filter
		if req.OrganizationId != nil && (setting.OrganizationID == nil || *setting.OrganizationID != *req.OrganizationId) {
			continue
		}

		// Apply project ID filter
		if req.ProjectId != nil && (setting.ProjectID == nil || *setting.ProjectID != *req.ProjectId) {
			continue
		}

		// Apply user ID filter
		if req.UserId != nil && (setting.UserID == nil || *setting.UserID != *req.UserId) {
			continue
		}

		// Unmarshal JSON value
		var value any
		if err := json.Unmarshal(setting.Value, &value); err != nil {
			return nil, status.Errorf(codes.Internal, "failed to unmarshal value: %v", err)
		}

		protoValue, err := structpb.NewValue(value)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to convert value: %v", err)
		}

		// Convert to proto Setting
		pbSetting := &pb.Setting{
			Id:          setting.ID,
			Scope:       pbScope,
			Subsystem:   setting.Subsystem,
			Key:         setting.Key,
			Value:       protoValue,
			Description: setting.Description,
			CreatedAt:   setting.CreatedAt.Unix(),
			UpdatedAt:   setting.UpdatedAt.Unix(),
		}

		// Set optional fields
		if setting.OrganizationID != nil {
			pbSetting.OrganizationId = *setting.OrganizationID
		}
		if setting.ProjectID != nil {
			pbSetting.ProjectId = *setting.ProjectID
		}
		if setting.UserID != nil {
			pbSetting.UserId = *setting.UserID
		}
		if setting.UpdatedBy != nil {
			pbSetting.UpdatedBy = *setting.UpdatedBy
		}

		filteredSettings = append(filteredSettings, pbSetting)
	}

	return &pb.ListAllSettingsResponse{
		Settings: filteredSettings,
	}, nil
}

// GetSetting retrieves a single setting with hierarchical resolution
func (s *SettingsServiceServer) GetSetting(ctx context.Context, req *pb.GetSettingRequest) (*pb.GetSettingResponse, error) {
	if req.Subsystem == "" || req.Key == "" {
		return nil, status.Error(codes.InvalidArgument, "subsystem and key are required")
	}

	// Build scope from request
	scope := settings.ScopeContext{}
	if req.OrganizationId != nil {
		scope.OrganizationID = req.OrganizationId
	}
	if req.ProjectId != nil {
		scope.ProjectID = req.ProjectId
	}
	if req.UserId != nil {
		scope.UserID = req.UserId
	}

	// Get the setting value with metadata
	result, err := s.settingsService.GetWithMetadata(ctx, req.Subsystem, req.Key, scope)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "setting not found: %v", err)
	}

	// Convert JSON bytes to protobuf Value
	var value any
	if err := json.Unmarshal(result.Value, &value); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to unmarshal value: %v", err)
	}

	protoValue, err := structpb.NewValue(value)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to convert value: %v", err)
	}

	return &pb.GetSettingResponse{
		Value:         protoValue,
		ResolvedScope: convertScopeToPB(result.ResolvedScope),
		Description:   result.Description,
	}, nil
}

// Helper functions

func convertMapToProtoStruct(m map[string]any) (map[string]*structpb.Value, error) {
	result := make(map[string]*structpb.Value)
	for k, v := range m {
		protoValue, err := structpb.NewValue(v)
		if err != nil {
			return nil, fmt.Errorf("failed to convert value for key %s: %w", k, err)
		}
		result[k] = protoValue
	}
	return result, nil
}

func getUserIDFromGRPCContext(ctx context.Context) *int64 {
	// Extract user ID from context (set by auth interceptor)
	userID, ok := GetUserIDFromContext(ctx)
	if !ok {
		return nil
	}
	return &userID
}

func ptrToString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func convertScopeToPB(scope string) pb.SettingScope {
	switch scope {
	case "system":
		return pb.SettingScope_SETTING_SCOPE_SYSTEM
	case "organization":
		return pb.SettingScope_SETTING_SCOPE_ORGANIZATION
	case "project":
		return pb.SettingScope_SETTING_SCOPE_PROJECT
	case "user":
		return pb.SettingScope_SETTING_SCOPE_USER
	default:
		return pb.SettingScope_SETTING_SCOPE_UNSPECIFIED
	}
}


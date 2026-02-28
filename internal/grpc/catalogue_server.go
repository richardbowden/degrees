package grpc

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/richardbowden/degrees/internal/dbpg"
	pb "github.com/richardbowden/degrees/internal/pb/degrees/v1"
	"github.com/richardbowden/degrees/internal/services"
)

type CatalogueServiceServer struct {
	pb.UnimplementedCatalogueServiceServer
	catalogueSvc *services.CatalogueService
}

func NewCatalogueServiceServer(catalogueSvc *services.CatalogueService) *CatalogueServiceServer {
	return &CatalogueServiceServer{
		catalogueSvc: catalogueSvc,
	}
}

func (s *CatalogueServiceServer) ListCategories(ctx context.Context, req *pb.ListCategoriesRequest) (*pb.ListCategoriesResponse, error) {
	cats, err := s.catalogueSvc.ListCategories(ctx)
	if err != nil {
		return nil, ToGRPCError(err)
	}

	pbCats := make([]*pb.ServiceCategory, len(cats))
	for i, c := range cats {
		pbCats[i] = dbCategoryToPB(c)
	}

	return &pb.ListCategoriesResponse{Categories: pbCats}, nil
}

func (s *CatalogueServiceServer) AdminListServices(ctx context.Context, req *pb.ListCatalogueServicesRequest) (*pb.ListCatalogueServicesResponse, error) {
	userID, ok := GetUserIDFromContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "user not authenticated")
	}

	svcs, err := s.catalogueSvc.ListAllServices(ctx, userID)
	if err != nil {
		return nil, ToGRPCError(err)
	}

	pbSvcs := make([]*pb.DetailingService, len(svcs))
	for i, swt := range svcs {
		pbSvc := dbServiceToPB(swt.Service)
		pbSvc.PriceTiers = dbPriceTiersToPB(swt.Tiers)
		pbSvcs[i] = pbSvc
	}

	return &pb.ListCatalogueServicesResponse{Services: pbSvcs}, nil
}

func (s *CatalogueServiceServer) ListServices(ctx context.Context, req *pb.ListCatalogueServicesRequest) (*pb.ListCatalogueServicesResponse, error) {
	svcs, err := s.catalogueSvc.ListServices(ctx)
	if err != nil {
		return nil, ToGRPCError(err)
	}

	pbSvcs := make([]*pb.DetailingService, len(svcs))
	for i, swt := range svcs {
		pbSvc := dbServiceToPB(swt.Service)
		pbSvc.PriceTiers = dbPriceTiersToPB(swt.Tiers)
		pbSvcs[i] = pbSvc
	}

	return &pb.ListCatalogueServicesResponse{Services: pbSvcs}, nil
}

func (s *CatalogueServiceServer) GetService(ctx context.Context, req *pb.GetCatalogueServiceRequest) (*pb.GetCatalogueServiceResponse, error) {
	if req.Slug == "" {
		return nil, status.Error(codes.InvalidArgument, "slug is required")
	}

	svc, opts, tiers, err := s.catalogueSvc.GetServiceBySlug(ctx, req.Slug)
	if err != nil {
		return nil, ToGRPCError(err)
	}

	pbOpts := make([]*pb.DetailingServiceOption, len(opts))
	for i, opt := range opts {
		pbOpts[i] = dbServiceOptionToPB(opt)
	}

	pbSvc := &pb.DetailingService{
		Id:              svc.ID,
		CategoryId:      svc.CategoryID,
		Name:            svc.Name,
		Slug:            svc.Slug,
		Description:     svc.Description.String,
		ShortDesc:       svc.ShortDesc.String,
		BasePrice:       svc.BasePrice,
		DurationMinutes: svc.DurationMinutes,
		IsActive:        svc.IsActive,
		SortOrder:       svc.SortOrder,
		CategoryName:    svc.CategoryName,
		Options:         pbOpts,
		PriceTiers:      dbPriceTiersToPB(tiers),
	}
	if svc.CreatedAt.Valid {
		pbSvc.CreatedAt = timestamppb.New(svc.CreatedAt.Time)
	}
	if svc.UpdatedAt.Valid {
		pbSvc.UpdatedAt = timestamppb.New(svc.UpdatedAt.Time)
	}

	return &pb.GetCatalogueServiceResponse{Service: pbSvc}, nil
}

func (s *CatalogueServiceServer) CreateService(ctx context.Context, req *pb.CreateServiceRequest) (*pb.CreateServiceResponse, error) {
	if req.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "name is required")
	}
	if req.Slug == "" {
		return nil, status.Error(codes.InvalidArgument, "slug is required")
	}

	userID, ok := GetUserIDFromContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "user not authenticated")
	}

	params := dbpg.CreateServiceParams{
		CategoryID:      req.CategoryId,
		Name:            req.Name,
		Slug:            req.Slug,
		Description:     dbpg.StringToPGString(req.Description),
		ShortDesc:       dbpg.StringToPGString(req.ShortDesc),
		BasePrice:       req.BasePrice,
		DurationMinutes: req.DurationMinutes,
		IsActive:        req.IsActive,
		SortOrder:       req.SortOrder,
	}

	svc, err := s.catalogueSvc.CreateService(ctx, userID, params)
	if err != nil {
		return nil, ToGRPCError(err)
	}

	return &pb.CreateServiceResponse{Service: dbServiceToPB(svc)}, nil
}

func (s *CatalogueServiceServer) UpdateService(ctx context.Context, req *pb.UpdateServiceRequest) (*pb.UpdateServiceResponse, error) {
	if req.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	userID, ok := GetUserIDFromContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "user not authenticated")
	}

	params := dbpg.UpdateServiceParams{
		ID:              req.Id,
		CategoryID:      req.CategoryId,
		Name:            req.Name,
		Slug:            req.Slug,
		Description:     dbpg.StringToPGString(req.Description),
		ShortDesc:       dbpg.StringToPGString(req.ShortDesc),
		BasePrice:       req.BasePrice,
		DurationMinutes: req.DurationMinutes,
		IsActive:        req.IsActive,
		SortOrder:       req.SortOrder,
	}

	svc, err := s.catalogueSvc.UpdateService(ctx, userID, params)
	if err != nil {
		return nil, ToGRPCError(err)
	}

	return &pb.UpdateServiceResponse{Service: dbServiceToPB(svc)}, nil
}

func (s *CatalogueServiceServer) DeleteService(ctx context.Context, req *pb.DeleteServiceRequest) (*pb.DeleteServiceResponse, error) {
	if req.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	userID, ok := GetUserIDFromContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "user not authenticated")
	}

	err := s.catalogueSvc.DeleteService(ctx, userID, req.Id)
	if err != nil {
		return nil, ToGRPCError(err)
	}

	return &pb.DeleteServiceResponse{Success: true}, nil
}

func (s *CatalogueServiceServer) AddServiceOption(ctx context.Context, req *pb.AddServiceOptionRequest) (*pb.AddServiceOptionResponse, error) {
	if req.ServiceId == 0 {
		return nil, status.Error(codes.InvalidArgument, "service_id is required")
	}
	if req.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "name is required")
	}

	userID, ok := GetUserIDFromContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "user not authenticated")
	}

	params := dbpg.CreateServiceOptionParams{
		ServiceID:   req.ServiceId,
		Name:        req.Name,
		Description: dbpg.StringToPGString(req.Description),
		Price:       req.Price,
		IsActive:    req.IsActive,
		SortOrder:   req.SortOrder,
	}

	opt, err := s.catalogueSvc.AddServiceOption(ctx, userID, params)
	if err != nil {
		return nil, ToGRPCError(err)
	}

	return &pb.AddServiceOptionResponse{Option: dbServiceOptionToPB(opt)}, nil
}

func (s *CatalogueServiceServer) ListVehicleCategories(ctx context.Context, req *pb.ListVehicleCategoriesRequest) (*pb.ListVehicleCategoriesResponse, error) {
	cats, err := s.catalogueSvc.ListVehicleCategories(ctx)
	if err != nil {
		return nil, ToGRPCError(err)
	}

	pbCats := make([]*pb.VehicleCategory, len(cats))
	for i, c := range cats {
		pbCats[i] = dbVehicleCategoryToPB(c)
	}

	return &pb.ListVehicleCategoriesResponse{VehicleCategories: pbCats}, nil
}

func (s *CatalogueServiceServer) CreateVehicleCategory(ctx context.Context, req *pb.CreateVehicleCategoryRequest) (*pb.CreateVehicleCategoryResponse, error) {
	if req.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "name is required")
	}
	if req.Slug == "" {
		return nil, status.Error(codes.InvalidArgument, "slug is required")
	}

	userID, ok := GetUserIDFromContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "user not authenticated")
	}

	params := dbpg.CreateVehicleCategoryParams{
		Name:        req.Name,
		Slug:        req.Slug,
		Description: dbpg.StringToPGString(req.Description),
		SortOrder:   req.SortOrder,
	}

	vc, err := s.catalogueSvc.CreateVehicleCategory(ctx, userID, params)
	if err != nil {
		return nil, ToGRPCError(err)
	}

	return &pb.CreateVehicleCategoryResponse{VehicleCategory: dbVehicleCategoryToPB(vc)}, nil
}

func (s *CatalogueServiceServer) UpdateVehicleCategory(ctx context.Context, req *pb.UpdateVehicleCategoryRequest) (*pb.UpdateVehicleCategoryResponse, error) {
	if req.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	userID, ok := GetUserIDFromContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "user not authenticated")
	}

	params := dbpg.UpdateVehicleCategoryParams{
		ID:          req.Id,
		Name:        req.Name,
		Slug:        req.Slug,
		Description: dbpg.StringToPGString(req.Description),
		SortOrder:   req.SortOrder,
	}

	vc, err := s.catalogueSvc.UpdateVehicleCategory(ctx, userID, params)
	if err != nil {
		return nil, ToGRPCError(err)
	}

	return &pb.UpdateVehicleCategoryResponse{VehicleCategory: dbVehicleCategoryToPB(vc)}, nil
}

func (s *CatalogueServiceServer) DeleteVehicleCategory(ctx context.Context, req *pb.DeleteVehicleCategoryRequest) (*pb.DeleteVehicleCategoryResponse, error) {
	if req.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	userID, ok := GetUserIDFromContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "user not authenticated")
	}

	err := s.catalogueSvc.DeleteVehicleCategory(ctx, userID, req.Id)
	if err != nil {
		return nil, ToGRPCError(err)
	}

	return &pb.DeleteVehicleCategoryResponse{Success: true}, nil
}

func (s *CatalogueServiceServer) SetServicePriceTiers(ctx context.Context, req *pb.SetServicePriceTiersRequest) (*pb.SetServicePriceTiersResponse, error) {
	if req.ServiceId == 0 {
		return nil, status.Error(codes.InvalidArgument, "service_id is required")
	}

	userID, ok := GetUserIDFromContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "user not authenticated")
	}

	tiers := make([]dbpg.UpsertPriceTierParams, len(req.Tiers))
	for i, t := range req.Tiers {
		tiers[i] = dbpg.UpsertPriceTierParams{
			ServiceID:         req.ServiceId,
			VehicleCategoryID: t.VehicleCategoryId,
			Price:             t.Price,
		}
	}

	result, err := s.catalogueSvc.SetServicePriceTiers(ctx, userID, req.ServiceId, tiers)
	if err != nil {
		return nil, ToGRPCError(err)
	}

	return &pb.SetServicePriceTiersResponse{PriceTiers: dbPriceTiersToPB(result)}, nil
}

// Conversion helpers

func dbCategoryToPB(c dbpg.ServiceCategory) *pb.ServiceCategory {
	cat := &pb.ServiceCategory{
		Id:          c.ID,
		Name:        c.Name,
		Slug:        c.Slug,
		Description: c.Description.String,
		SortOrder:   c.SortOrder,
	}
	if c.CreatedAt.Valid {
		cat.CreatedAt = timestamppb.New(c.CreatedAt.Time)
	}
	if c.UpdatedAt.Valid {
		cat.UpdatedAt = timestamppb.New(c.UpdatedAt.Time)
	}
	return cat
}

func dbServiceToPB(s dbpg.Service) *pb.DetailingService {
	svc := &pb.DetailingService{
		Id:              s.ID,
		CategoryId:      s.CategoryID,
		Name:            s.Name,
		Slug:            s.Slug,
		Description:     s.Description.String,
		ShortDesc:       s.ShortDesc.String,
		BasePrice:       s.BasePrice,
		DurationMinutes: s.DurationMinutes,
		IsActive:        s.IsActive,
		SortOrder:       s.SortOrder,
	}
	if s.CreatedAt.Valid {
		svc.CreatedAt = timestamppb.New(s.CreatedAt.Time)
	}
	if s.UpdatedAt.Valid {
		svc.UpdatedAt = timestamppb.New(s.UpdatedAt.Time)
	}
	return svc
}

func dbServiceOptionToPB(o dbpg.ServiceOption) *pb.DetailingServiceOption {
	opt := &pb.DetailingServiceOption{
		Id:          o.ID,
		ServiceId:   o.ServiceID,
		Name:        o.Name,
		Description: o.Description.String,
		Price:       o.Price,
		IsActive:    o.IsActive,
		SortOrder:   o.SortOrder,
	}
	if o.CreatedAt.Valid {
		opt.CreatedAt = timestamppb.New(o.CreatedAt.Time)
	}
	return opt
}

func dbVehicleCategoryToPB(c dbpg.VehicleCategory) *pb.VehicleCategory {
	vc := &pb.VehicleCategory{
		Id:          c.ID,
		Name:        c.Name,
		Slug:        c.Slug,
		Description: c.Description.String,
		SortOrder:   c.SortOrder,
	}
	if c.CreatedAt.Valid {
		vc.CreatedAt = timestamppb.New(c.CreatedAt.Time)
	}
	if c.UpdatedAt.Valid {
		vc.UpdatedAt = timestamppb.New(c.UpdatedAt.Time)
	}
	return vc
}

func dbPriceTiersToPB(tiers []dbpg.ListPriceTiersByServiceRow) []*pb.ServicePriceTier {
	result := make([]*pb.ServicePriceTier, len(tiers))
	for i, t := range tiers {
		result[i] = &pb.ServicePriceTier{
			ServiceId:         t.ServiceID,
			VehicleCategoryId: t.VehicleCategoryID,
			CategoryName:      t.CategoryName,
			CategorySlug:      t.CategorySlug,
			Price:             t.Price,
		}
	}
	return result
}

package resourcelocalization

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	resourcesv1 "godev.bluffpointtech.com/i18nextservice/proto/resources"
	httpbody "google.golang.org/genproto/googleapis/api/httpbody"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type ResourceLocalizationService struct {
	resourcesv1.UnimplementedResourcesAPIServer
}

func NewResourceLocalizationService() *ResourceLocalizationService {
	return &ResourceLocalizationService{}
}

func (s *ResourceLocalizationService) GetResource(ctx context.Context, req *resourcesv1.GetResourceRequest) (*httpbody.HttpBody, error) {
	payload := map[string]any{
		"application": req.GetApplication(),
		"component":   req.GetComponent(),
		"page":        req.GetPage(),
		"resources":   map[string]any{},
	}
	b, _ := json.Marshal(payload)
	return &httpbody.HttpBody{
		ContentType: "application/json",
		Data:        b,
	}, nil
}

func (s *ResourceLocalizationService) GetResourceByUUID(ctx context.Context, req *resourcesv1.GetResourceByUUIDRequest) (*resourcesv1.PageResourceDetails, error) {
	now := timestamppb.New(time.Now())
	st := &structpb.Struct{Fields: map[string]*structpb.Value{}}
	return &resourcesv1.PageResourceDetails{
		Uuid:              req.GetUuid(),
		Resource:          st,
		CreationUID:       req.GetUuid(),
		CreationTimeStamp: now,
		RevisionUID:       req.GetUuid(),
		RevisionTimeStamp: now,
		Tags:              nil,
	}, nil
}

func (s *ResourceLocalizationService) CreatePageResource(ctx context.Context, req *resourcesv1.CreateResourceRequest) (*resourcesv1.PageResourceDetails, error) {
	id := uuid.NewString()
	now := timestamppb.New(time.Now())
	return &resourcesv1.PageResourceDetails{
		Uuid:              id,
		Resource:          req.GetResource(),
		CreationUID:       req.GetCreationUID(),
		CreationTimeStamp: now,
		RevisionUID:       id,
		RevisionTimeStamp: now,
		Tags:              req.GetTags(),
	}, nil
}

func (s *ResourceLocalizationService) UpdatePageResource(ctx context.Context, req *resourcesv1.UpdateResourceRequest) (*resourcesv1.PageResourceDetails, error) {
	now := timestamppb.New(time.Now())
	return &resourcesv1.PageResourceDetails{
		Uuid:              req.GetUuid(),
		Resource:          req.GetResource(),
		RevisionUID:       req.GetRevisionUID(),
		RevisionTimeStamp: now,
		Tags:              req.GetTags(),
	}, nil
}

func (s *ResourceLocalizationService) DeletePageResource(ctx context.Context, req *resourcesv1.DeleteResourceRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

var _ resourcesv1.ResourcesAPIServer = (*ResourceLocalizationService)(nil)

func (s *ResourceLocalizationService) SearchByTag(ctx context.Context, req *resourcesv1.SearchByTagRequest) (*resourcesv1.SearchResults, error) {
	return &resourcesv1.SearchResults{Items: nil}, nil
}

func (s *ResourceLocalizationService) SearchByText(ctx context.Context, req *resourcesv1.SearchByTextRequest) (*resourcesv1.SearchResults, error) {
	return &resourcesv1.SearchResults{Items: nil}, nil
}

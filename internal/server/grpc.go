package server

import (
	pb "blob-service/gen/blob/v1"
	"blob-service/internal/storage"
	"context"
	"database/sql"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type BlobServer struct {
	pb.UnimplementedBlobServiceServer
	repo *storage.BlobRepository
}

func NewBlobServer(repo *storage.BlobRepository) *BlobServer {
	return &BlobServer{repo: repo}
}

func (s *BlobServer) CreateBlob(ctx context.Context, req *pb.CreateBlobRequest) (*pb.CreateBlobResponse, error) {
	if req.Name == "" || req.Data == nil || req.ContentType == "" {
		return nil, status.Error(codes.InvalidArgument, "name, data, and content_type are required")
	}

	blob, err := s.repo.Create(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create blob: %v", err)
	}

	return &pb.CreateBlobResponse{Blob: blob}, nil
}

func (s *BlobServer) GetBlob(ctx context.Context, req *pb.GetBlobRequest) (*pb.GetBlobResponse, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	blob, err := s.repo.Get(ctx, req)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, status.Errorf(codes.NotFound, "blob not found")
	} else if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get blob: %v", err)
	}

	return &pb.GetBlobResponse{Blob: blob}, nil
}

func (s *BlobServer) UpdateBlob(ctx context.Context, req *pb.UpdateBlobRequest) (*pb.UpdateBlobResponse, error) {
	if req.Id == "" || req.Name == "" || req.Data == nil || req.ContentType == "" {
		return nil, status.Error(codes.InvalidArgument, "id, name, data, content_type are required")
	}

	blob, err := s.repo.Update(ctx, req)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, status.Errorf(codes.NotFound, "blob not found")
	} else if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update blob: %v", err)
	}

	return &pb.UpdateBlobResponse{Blob: blob}, nil
}

func (s *BlobServer) DeleteBlob(ctx context.Context, req *pb.DeleteBlobRequest) (*pb.DeleteBlobResponse, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	blob, err := s.repo.Delete(ctx, req)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, status.Errorf(codes.NotFound, "blob not found")
	} else if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete blob: %v", err)
	}

	return &pb.DeleteBlobResponse{Blob: blob}, nil
}

func (s *BlobServer) ListBlobs(ctx context.Context, req *pb.ListBlobsRequest) (*pb.ListBlobsResponse, error) {
	if req.Limit <= 0 {
		return nil, status.Error(codes.InvalidArgument, "limit must be greater than 0")
	}

	blobs, err := s.repo.List(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get blobs: %v", err)
	}

	return &pb.ListBlobsResponse{Blobs: blobs}, nil
}

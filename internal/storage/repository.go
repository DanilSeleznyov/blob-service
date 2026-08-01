package storage

import (
	pb "blob-service/gen/blob/v1"
	"context"
	"crypto/md5"
	"crypto/sha256"
	"database/sql"
	"encoding/json"
	"fmt"
	"path/filepath"
	"time"

	"github.com/lib/pq"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type BlobRepository struct {
	db *sql.DB
}

func NewBlobRepository(db *sql.DB) *BlobRepository {
	return &BlobRepository{db: db}
}

func (r *BlobRepository) Create(ctx context.Context, req *pb.CreateBlobRequest) (*pb.Blob, error) {
	query := `INSERT INTO blobs (name, content_type, extension, data, size, md5_hash, sha256_hash, tags)
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING *`

	var blob pb.Blob
	var createdAt, updatedAt time.Time
	var extension, md5Val string
	var sha256Null, categoryNull, ownerIdNull, accessLevelNull sql.NullString
	var metadata []byte

	ext := filepath.Ext(req.Name)
	md5Hash := fmt.Sprintf("%x", md5.Sum(req.Data))
	sha256Hash := fmt.Sprintf("%x", sha256.Sum256(req.Data))

	row := r.db.QueryRowContext(ctx, query,
		req.Name, req.ContentType,
		ext, req.Data, len(req.Data),
		md5Hash, sha256Hash,
		pq.Array(req.Tags),
	)
	err := row.Scan(
		&blob.Id, &blob.Name, &blob.ContentType,
		&extension, &blob.Data, &blob.Size,
		&md5Val, &sha256Null,
		&createdAt, &updatedAt,
		&categoryNull, pq.Array(&blob.Tags), &metadata,
		&ownerIdNull, &accessLevelNull,
	)

	if err != nil {
		return nil, err
	}

	if len(metadata) > 0 {
		_ = json.Unmarshal(metadata, &blob.Metadata)
	}

	blob.Extension = extension
	blob.Md5Hash = md5Val
	blob.Sha256Hash = sha256Null.String
	blob.Category = categoryNull.String
	blob.OwnerId = ownerIdNull.String
	blob.AccessLevel = accessLevelNull.String
	blob.CreatedAt = timestamppb.New(createdAt)
	blob.UpdatedAt = timestamppb.New(updatedAt)

	return &blob, nil
}

func (r *BlobRepository) Get(ctx context.Context, req *pb.GetBlobRequest) (*pb.Blob, error) {
	query := `SELECT * FROM blobs WHERE id = $1;`

	var blob pb.Blob
	var createdAt, updatedAt time.Time
	var extension, md5Val string
	var sha256Null, categoryNull, ownerIdNull, accessLevelNull sql.NullString
	var metadata []byte

	row := r.db.QueryRowContext(ctx, query, req.Id)
	err := row.Scan(
		&blob.Id, &blob.Name, &blob.ContentType,
		&extension, &blob.Data, &blob.Size,
		&md5Val, &sha256Null,
		&createdAt, &updatedAt,
		&categoryNull, pq.Array(&blob.Tags), &metadata,
		&ownerIdNull, &accessLevelNull,
	)

	if err != nil {
		return nil, err
	}

	if len(metadata) > 0 {
		_ = json.Unmarshal(metadata, &blob.Metadata)
	}

	blob.Extension = extension
	blob.Md5Hash = md5Val
	blob.Sha256Hash = sha256Null.String
	blob.Category = categoryNull.String
	blob.OwnerId = ownerIdNull.String
	blob.AccessLevel = accessLevelNull.String
	blob.CreatedAt = timestamppb.New(createdAt)
	blob.UpdatedAt = timestamppb.New(updatedAt)

	return &blob, nil
}

func (r *BlobRepository) Update(ctx context.Context, req *pb.UpdateBlobRequest) (*pb.Blob, error) {
	query := `UPDATE blobs SET name = $1, content_type = $2, data = $3, category = $4, tags = $5, metadata = $6, access_level = $7 WHERE id = $8 RETURNING *`

	var blob pb.Blob
	var createdAt, updatedAt time.Time
	var extension, md5Val string
	var sha256Null, categoryNull, ownerIdNull, accessLevelNull sql.NullString
	var metadata []byte

	var metadataJSON []byte
	if len(req.Metadata) > 0 {
		var err error
		metadataJSON, err = json.Marshal(req.Metadata)
		if err != nil {
			return nil, err
		}
	}

	row := r.db.QueryRowContext(ctx, query,
		req.Name, req.ContentType, req.Data, req.Category,
		pq.Array(req.Tags), metadataJSON, req.AccessLevel,
		req.Id,
	)
	err := row.Scan(
		&blob.Id, &blob.Name, &blob.ContentType,
		&extension, &blob.Data, &blob.Size,
		&md5Val, &sha256Null,
		&createdAt, &updatedAt,
		&categoryNull, pq.Array(&blob.Tags), &metadata,
		&ownerIdNull, &accessLevelNull,
	)

	if err != nil {
		return nil, err
	}

	if len(metadata) > 0 {
		_ = json.Unmarshal(metadata, &blob.Metadata)
	}

	blob.Extension = extension
	blob.Md5Hash = md5Val
	blob.Sha256Hash = sha256Null.String
	blob.Category = categoryNull.String
	blob.OwnerId = ownerIdNull.String
	blob.AccessLevel = accessLevelNull.String
	blob.CreatedAt = timestamppb.New(createdAt)
	blob.UpdatedAt = timestamppb.New(updatedAt)

	return &blob, nil
}

func (r *BlobRepository) Delete(ctx context.Context, req *pb.DeleteBlobRequest) (*pb.Blob, error) {
	query := `DELETE FROM blobs WHERE id = $1 RETURNING *`

	var blob pb.Blob
	var createdAt, updatedAt time.Time
	var extension, md5Val string
	var sha256Null, categoryNull, ownerIdNull, accessLevelNull sql.NullString
	var metadata []byte

	row := r.db.QueryRowContext(ctx, query, req.Id)
	err := row.Scan(
		&blob.Id, &blob.Name, &blob.ContentType,
		&extension, &blob.Data, &blob.Size,
		&md5Val, &sha256Null,
		&createdAt, &updatedAt,
		&categoryNull, pq.Array(&blob.Tags), &metadata,
		&ownerIdNull, &accessLevelNull,
	)

	if err != nil {
		return nil, err
	}

	if len(metadata) > 0 {
		_ = json.Unmarshal(metadata, &blob.Metadata)
	}

	blob.Extension = extension
	blob.Md5Hash = md5Val
	blob.Sha256Hash = sha256Null.String
	blob.Category = categoryNull.String
	blob.OwnerId = ownerIdNull.String
	blob.AccessLevel = accessLevelNull.String
	blob.CreatedAt = timestamppb.New(createdAt)
	blob.UpdatedAt = timestamppb.New(updatedAt)

	return &blob, nil
}

func (r *BlobRepository) List(ctx context.Context, req *pb.ListBlobsRequest) ([]*pb.Blob, error) {
	query := `SELECT * FROM blobs
			  WHERE ($1 = '' OR name ILIKE '%' || $1 || '%')
  			    AND ($2 = '' OR content_type = $2)
  			    AND (COALESCE(cardinality($3::varchar[]), 0) = 0 OR tags && $3::varchar[])
			  LIMIT NULLIF($4, 0)`

	rows, err := r.db.QueryContext(ctx, query, req.NameFilter, req.ContentTypeFilter, pq.Array(req.TagFilters), req.Limit)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var blobs []*pb.Blob

	for rows.Next() {
		var blob pb.Blob
		var createdAt, updatedAt time.Time
		var extension, md5Val string
		var sha256Null, categoryNull, ownerIdNull, accessLevelNull sql.NullString
		var metadata []byte

		err := rows.Scan(
			&blob.Id, &blob.Name, &blob.ContentType,
			&extension, &blob.Data, &blob.Size,
			&md5Val, &sha256Null,
			&createdAt, &updatedAt,
			&categoryNull, pq.Array(&blob.Tags), &metadata,
			&ownerIdNull, &accessLevelNull,
		)
		if err != nil {
			return nil, err
		}

		if len(metadata) > 0 {
			_ = json.Unmarshal(metadata, &blob.Metadata)
		}

		blob.Extension = extension
		blob.Md5Hash = md5Val
		blob.Sha256Hash = sha256Null.String
		blob.Category = categoryNull.String
		blob.OwnerId = ownerIdNull.String
		blob.AccessLevel = accessLevelNull.String
		blob.CreatedAt = timestamppb.New(createdAt)
		blob.UpdatedAt = timestamppb.New(updatedAt)

		blobs = append(blobs, &blob)
	}

	return blobs, nil
}

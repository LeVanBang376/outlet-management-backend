package service

import (
	"context"
	"database/sql"
	customerrors "magnolia-test-backend/internal/custom-errors"
	"magnolia-test-backend/internal/dto"
	"magnolia-test-backend/internal/model"
	"magnolia-test-backend/internal/repository"
	"time"
)

type FileService struct {
	db   *sql.DB
	repo *repository.FileRepository
}

func NewFileService(
	db *sql.DB,
	repo *repository.FileRepository,
) *FileService {
	return &FileService{
		db:   db,
		repo: repo,
	}
}

func (s *FileService) Create(
	ctx context.Context,
	req dto.CreateFileRequest,
) (*dto.FileResponse, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	file := model.File{
		ObjectKey:   req.ObjectKey,
		FileName:    req.FileName,
		ContentType: req.ContentType,
		Size:        req.Size,
		CreatedAt:   time.Now(),
	}

	res, err := s.repo.Create(ctx, tx, &file)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return res, nil
}

func (s *FileService) GetByID(
	ctx context.Context,
	id uint,
) (*dto.FileResponse, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	file, err := s.repo.FindByID(ctx, tx, id)
	if err != nil {
		return nil, err
	}

	if file == nil {
		return nil, customerrors.FileErrNotFound
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return dto.ToFileResponse(file), nil
}

func (s *FileService) GetAll(
	ctx context.Context,
) ([]*dto.FileResponse, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	files, err := s.repo.FindAll(ctx, tx)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return files, nil
}

func (s *FileService) Delete(
	ctx context.Context,
	id uint,
) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	file, err := s.repo.FindByID(ctx, tx, id)
	if err != nil {
		return err
	}

	if file == nil {
		return customerrors.FileErrNotFound
	}

	if err := s.repo.Delete(ctx, tx, id); err != nil {
		return err
	}

	return tx.Commit()
}

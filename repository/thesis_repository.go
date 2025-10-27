package repository

import (
	"context"
	"errors"
	"math"

	"github.com/Amierza/chat-service/dto"
	"github.com/Amierza/chat-service/entity"
	"github.com/Amierza/chat-service/response"
	"gorm.io/gorm"
)

type (
	IThesisRepository interface {
		// CREATE / POST

		// READ / GET
		GetThesisByID(ctx context.Context, tx *gorm.DB, id string) (*entity.Thesis, bool, error)
		GetAllThesesByLecturerIDWithPagination(ctx context.Context, tx *gorm.DB, pagination response.PaginationRequest, lecturerID string) (dto.ThesisPaginationRepositoryResponse, error)

		// UPDATE / PATCH
		UpdateThesis(ctx context.Context, tx *gorm.DB, thesis *entity.Thesis) error

		// DELETE / DELETE
	}

	thesisRepository struct {
		db *gorm.DB
	}
)

func NewThesisRepository(db *gorm.DB) *thesisRepository {
	return &thesisRepository{
		db: db,
	}
}

// CREATE / POST

// READ / GET
func (nr *thesisRepository) GetThesisByID(ctx context.Context, tx *gorm.DB, id string) (*entity.Thesis, bool, error) {
	if tx == nil {
		tx = nr.db
	}

	var thesis *entity.Thesis
	err := tx.WithContext(ctx).
		Preload("ThesisLogs").
		Preload("Sessions").
		Preload("Supervisors.Lecturer.StudyProgram.Faculty").
		Preload("Student.StudyProgram.Faculty").
		Where("id = ?", id).
		Take(&thesis).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return &entity.Thesis{}, false, nil
	}
	if err != nil {
		return &entity.Thesis{}, false, err
	}

	return thesis, true, nil
}
func (tr *thesisRepository) GetAllThesesByLecturerIDWithPagination(ctx context.Context, tx *gorm.DB, pagination response.PaginationRequest, lecturerID string) (dto.ThesisPaginationRepositoryResponse, error) {
	if tx == nil {
		tx = tr.db
	}

	var (
		theses []*entity.Thesis
		err    error
		count  int64
	)

	if pagination.PerPage == 0 {
		pagination.PerPage = 10
	}

	if pagination.Page == 0 {
		pagination.Page = 1
	}

	query := tx.WithContext(ctx).
		Preload("ThesisLogs").
		Preload("Sessions").
		Preload("Supervisors.Lecturer.StudyProgram.Faculty").
		Preload("Student.StudyProgram.Faculty").
		Joins("JOIN thesis_supervisors ts ON ts.thesis_id = theses.id").
		Where("ts.lecturer_id = ?", lecturerID).
		Model(&entity.Thesis{})

	if err := query.Order(`"created_at" DESC`).Find(&theses).Error; err != nil {
		return dto.ThesisPaginationRepositoryResponse{}, err
	}

	if err := query.Count(&count).Error; err != nil {
		return dto.ThesisPaginationRepositoryResponse{}, err
	}

	if err := query.Scopes(Paginate(pagination.Page, pagination.PerPage)).Find(&theses).Error; err != nil {
		return dto.ThesisPaginationRepositoryResponse{}, err
	}

	totalPage := int64(math.Ceil(float64(count) / float64(pagination.PerPage)))

	return dto.ThesisPaginationRepositoryResponse{
		Theses: theses,
		PaginationResponse: response.PaginationResponse{
			Page:    pagination.Page,
			PerPage: pagination.PerPage,
			MaxPage: totalPage,
			Count:   count,
		},
	}, err
}

// UPDATE / PATCH
func (sr *thesisRepository) UpdateThesis(ctx context.Context, tx *gorm.DB, thesis *entity.Thesis) error {
	if tx == nil {
		tx = sr.db
	}

	return tx.WithContext(ctx).Model(&entity.Thesis{}).Where("id = ?", thesis.ID).Updates(&thesis).Error
}

// DELETE / DELETE

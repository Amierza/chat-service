package repository

import (
	"context"
	"errors"
	"math"

	"github.com/Amierza/chat-service/constants"
	"github.com/Amierza/chat-service/dto"
	"github.com/Amierza/chat-service/entity"
	"github.com/Amierza/chat-service/response"
	"gorm.io/gorm"
)

type (
	IScheduleRepository interface {
		// CREATE / POST
		CreateSchedule(ctx context.Context, tx *gorm.DB, schedule *entity.Schedule) error

		// READ / GET
		GetThesisByID(ctx context.Context, tx *gorm.DB, thesisID string) (*entity.Thesis, bool, error)
		GetAllSchedulesByUserIDWithPagination(ctx context.Context, tx *gorm.DB, pagination response.PaginationRequest, role, userID string) (dto.SchedulePaginationRepositoryResponse, error)
		GetScheduleByID(ctx context.Context, tx *gorm.DB, id *string) (*entity.Schedule, bool, error)

		// UPDATE / PATCH
		UpdateSchedule(ctx context.Context, tx *gorm.DB, schedule *entity.Schedule) error
		UpdateScheduleStatus(ctx context.Context, tx *gorm.DB, id, status, approvedID string) error

		// DELETE / DELETE
		DeleteScheduleByID(ctx context.Context, tx *gorm.DB, id *string) error
	}

	scheduleRepository struct {
		db *gorm.DB
	}
)

func NewScheduleRepository(db *gorm.DB) *scheduleRepository {
	return &scheduleRepository{
		db: db,
	}
}

// CREATE / POST
func (sr *scheduleRepository) CreateSchedule(ctx context.Context, tx *gorm.DB, schedule *entity.Schedule) error {
	if tx == nil {
		tx = sr.db
	}

	return tx.WithContext(ctx).Create(&schedule).Error
}

// READ / GET
func (sr *scheduleRepository) GetThesisByID(ctx context.Context, tx *gorm.DB, thesisID string) (*entity.Thesis, bool, error) {
	if tx == nil {
		tx = sr.db
	}

	thesis := &entity.Thesis{}
	err := tx.WithContext(ctx).
		Preload("ThesisLogs").
		Preload("Sessions").
		Preload("Supervisors.Lecturer.StudyProgram.Faculty").
		Preload("Student.StudyProgram.Faculty").
		Where("id = ?", thesisID).
		Take(&thesis).Error
	if err != nil {
		return &entity.Thesis{}, false, err
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return &entity.Thesis{}, false, nil
	}

	return thesis, true, nil
}
func (sr *scheduleRepository) GetAllSchedulesByUserIDWithPagination(ctx context.Context, tx *gorm.DB, pagination response.PaginationRequest, role, userID string) (dto.SchedulePaginationRepositoryResponse, error) {
	if tx == nil {
		tx = sr.db
	}

	var (
		schedules []*entity.Schedule
		err       error
		count     int64
	)

	if pagination.PerPage == 0 {
		pagination.PerPage = 10
	}

	if pagination.Page == 0 {
		pagination.Page = 1
	}

	query := tx.WithContext(ctx).
		Preload("Thesis.Supervisors.Lecturer.StudyProgram.Faculty").
		Preload("Thesis.Student.StudyProgram.Faculty").
		Preload("CreatedBy.Student").
		Preload("ApprovedBy.Lecturer").
		Model(&entity.Schedule{})

	switch role {
	case constants.ENUM_ROLE_STUDENT:
		query = query.Where("schedules.created_by_id = ?", userID)

	case constants.ENUM_ROLE_LECTURER:
		query = query.
			Joins("JOIN theses t ON schedules.thesis_id = t.id").
			Joins("JOIN thesis_supervisors ts ON t.id = ts.thesis_id").
			Joins("JOIN lecturers l ON ts.lecturer_id = l.id").
			Joins("JOIN users u ON u.lecturer_id = l.id").
			Where("u.id = ?", userID)
	}

	if err := query.Order(`"created_at" DESC`).Find(&schedules).Error; err != nil {
		return dto.SchedulePaginationRepositoryResponse{}, err
	}

	if err := query.Count(&count).Error; err != nil {
		return dto.SchedulePaginationRepositoryResponse{}, err
	}

	if err := query.Scopes(Paginate(pagination.Page, pagination.PerPage)).Find(&schedules).Error; err != nil {
		return dto.SchedulePaginationRepositoryResponse{}, err
	}

	totalPage := int64(math.Ceil(float64(count) / float64(pagination.PerPage)))

	return dto.SchedulePaginationRepositoryResponse{
		Schedules: schedules,
		PaginationResponse: response.PaginationResponse{
			Page:    pagination.Page,
			PerPage: pagination.PerPage,
			MaxPage: totalPage,
			Count:   count,
		},
	}, err
}
func (sr *scheduleRepository) GetScheduleByID(ctx context.Context, tx *gorm.DB, id *string) (*entity.Schedule, bool, error) {
	if tx == nil {
		tx = sr.db
	}

	var schedule *entity.Schedule
	err := tx.WithContext(ctx).
		Preload("Thesis.Supervisors.Lecturer.StudyProgram.Faculty").
		Preload("Thesis.Student.StudyProgram.Faculty").
		Preload("CreatedBy.Student").
		Preload("ApprovedBy.Student").
		Preload("ApprovedBy.Lecturer").
		Where("id = ?", id).
		Take(&schedule).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return &entity.Schedule{}, false, nil
	}
	if err != nil {
		return &entity.Schedule{}, false, err
	}

	return schedule, true, nil
}

// UPDATE / PATCH
func (sr *scheduleRepository) UpdateSchedule(ctx context.Context, tx *gorm.DB, schedule *entity.Schedule) error {
	if tx == nil {
		tx = sr.db
	}

	return tx.WithContext(ctx).Model(&entity.Schedule{}).Where("id = ?", schedule.ID).Updates(&schedule).Error
}
func (sr *scheduleRepository) UpdateScheduleStatus(ctx context.Context, tx *gorm.DB, id, status, approvedID string) error {
	if tx == nil {
		tx = sr.db
	}

	err := tx.WithContext(ctx).
		Model(&entity.Schedule{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":         status,
			"approved_by_id": approvedID,
		}).Error

	if err != nil {
		return err
	}

	return nil
}

// DELETE / DELETE
func (sr *scheduleRepository) DeleteScheduleByID(ctx context.Context, tx *gorm.DB, id *string) error {
	if tx == nil {
		tx = sr.db
	}

	return tx.WithContext(ctx).Where("id = ?", id).Delete(&entity.Schedule{}).Error
}

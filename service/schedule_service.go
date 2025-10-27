package service

import (
	"context"
	"fmt"

	"github.com/Amierza/chat-service/constants"
	"github.com/Amierza/chat-service/dto"
	"github.com/Amierza/chat-service/entity"
	"github.com/Amierza/chat-service/jwt"
	"github.com/Amierza/chat-service/repository"
	"github.com/Amierza/chat-service/response"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type (
	IScheduleService interface {
		Create(ctx context.Context, req *dto.CreateScheduleRequest) (*dto.ScheduleResponse, error)
		GetAllWithPagination(ctx context.Context, req response.PaginationRequest) (dto.SchedulePaginationResponse, error)
		GetDetail(ctx context.Context, id *string) (*dto.ScheduleResponse, error)
		Update(ctx context.Context, req *dto.UpdateScheduleRequest) (*dto.ScheduleResponse, error)
		Approval(ctx context.Context, req *dto.ApprovalScheduleRequest) error
		Delete(ctx context.Context, id *string) error
	}

	scheduleService struct {
		scheduleRepo repository.IScheduleRepository
		userRepo     repository.IUserRepository
		logger       *zap.Logger
		jwt          jwt.IJWT
	}
)

func NewScheduleService(scheduleRepo repository.IScheduleRepository, userRepo repository.IUserRepository, logger *zap.Logger, jwt jwt.IJWT) *scheduleService {
	return &scheduleService{
		scheduleRepo: scheduleRepo,
		userRepo:     userRepo,
		logger:       logger,
		jwt:          jwt,
	}
}

func (ss *scheduleService) Create(ctx context.Context, req *dto.CreateScheduleRequest) (*dto.ScheduleResponse, error) {
	token := ctx.Value("Authorization").(string)
	userIDString, err := ss.jwt.GetUserIDByToken(token)
	if err != nil {
		ss.logger.Error("failed to get user id string from token",
			zap.Error(err),
		)
		return nil, err
	}

	user, found, err := ss.userRepo.GetUserByID(ctx, nil, userIDString)
	if err != nil {
		ss.logger.Error("failed to get user from user id string",
			zap.Error(err),
		)
		return nil, err
	}
	if !found {
		ss.logger.Warn("user not found",
			zap.String("user_id", userIDString),
		)
		return nil, dto.ErrNotFound
	}

	thesis, found, err := ss.scheduleRepo.GetThesisByID(ctx, nil, user.Student.Theses[0].ID.String())
	if err != nil {
		ss.logger.Error("failed to get thesis from thesis id",
			zap.Error(err),
		)
		return nil, err
	}
	if !found {
		ss.logger.Warn("thesis not found",
			zap.String("thesis_id", user.Student.Theses[0].ID.String()),
		)
		return nil, dto.ErrNotFound
	}

	newScheduleID := uuid.New()
	newSchedule := entity.Schedule{
		ID:          newScheduleID,
		ProposedAt:  req.ProposedAt,
		StartTime:   req.StartTime,
		EndTime:     req.EndTime,
		Status:      constants.ENUM_SCHEDULE_STATUS_PENDING,
		Description: req.Description,
		Location:    req.Location,
		ThesisID:    thesis.ID,
		CreatedByID: user.ID,
	}

	if err := ss.scheduleRepo.CreateSchedule(ctx, nil, &newSchedule); err != nil {
		ss.logger.Error("failed create schedule",
			zap.Error(err),
		)
		return nil, err
	}

	res := dto.ScheduleResponse{
		ID:          newSchedule.ID,
		ProposedAt:  newSchedule.ProposedAt,
		StartTime:   newSchedule.StartTime,
		EndTime:     newSchedule.EndTime,
		Status:      newSchedule.Status,
		Description: newSchedule.Description,
		Location:    newSchedule.Location,
		Thesis: dto.ThesisResponse{
			ID:          thesis.ID,
			Title:       thesis.Title,
			Description: thesis.Description,
			Progress:    thesis.Progress,
			Student: &dto.CustomUserResponse{
				ID:         thesis.Student.ID,
				Name:       thesis.Student.Name,
				Identifier: thesis.Student.Nim,
				Role:       constants.ENUM_ROLE_STUDENT,
			},
		},
		CreatedBy: dto.CustomUserResponse{
			ID:         user.ID,
			Name:       user.Student.Name,
			Identifier: user.Identifier,
			Role:       string(user.Role),
		},
	}

	for _, sup := range thesis.Supervisors {
		res.Thesis.Supervisors = append(res.Thesis.Supervisors, &dto.CustomUserResponse{
			ID:         sup.Lecturer.ID,
			Name:       sup.Lecturer.Name,
			Identifier: sup.Lecturer.Nip,
			Role:       constants.ENUM_ROLE_LECTURER,
		})
	}

	ss.logger.Info("success create schedule",
		zap.String("schedule_id", newSchedule.ID.String()),
	)

	return &res, nil
}

func (ss *scheduleService) GetAllWithPagination(ctx context.Context, req response.PaginationRequest) (dto.SchedulePaginationResponse, error) {
	token := ctx.Value("Authorization").(string)
	userIDString, err := ss.jwt.GetUserIDByToken(token)
	if err != nil {
		ss.logger.Error("failed to get user id string from token",
			zap.Error(err),
		)
		return dto.SchedulePaginationResponse{}, dto.ErrGetUserIDFromToken
	}
	user, found, err := ss.userRepo.GetUserByID(ctx, nil, userIDString)
	if err != nil {
		ss.logger.Error("failed to get user from user id string",
			zap.Error(err),
		)
		return dto.SchedulePaginationResponse{}, err
	}
	if !found {
		ss.logger.Warn("user not found",
			zap.String("user_id", userIDString),
		)
		return dto.SchedulePaginationResponse{}, dto.ErrNotFound
	}

	datas, err := ss.scheduleRepo.GetAllSchedulesByUserIDWithPagination(ctx, nil, req, string(user.Role), userIDString)
	if err != nil {
		ss.logger.Error("failed to get all schedules with pagination",
			zap.Error(err),
		)
		return dto.SchedulePaginationResponse{}, err
	}

	schedules := make([]*dto.ScheduleResponse, 0, len(datas.Schedules))
	for _, data := range datas.Schedules {
		schedule := dto.ScheduleResponse{
			ID:          data.ID,
			ProposedAt:  data.ProposedAt,
			StartTime:   data.StartTime,
			EndTime:     data.EndTime,
			Status:      data.Status,
			Description: data.Description,
			Location:    data.Location,
			Thesis: dto.ThesisResponse{
				ID:          data.Thesis.ID,
				Title:       data.Thesis.Title,
				Description: data.Thesis.Description,
				Progress:    data.Thesis.Progress,
				Student: &dto.CustomUserResponse{
					ID:         data.Thesis.Student.ID,
					Name:       data.Thesis.Student.Name,
					Identifier: data.Thesis.Student.Nim,
					Role:       constants.ENUM_ROLE_STUDENT,
				},
			},
			CreatedBy: dto.CustomUserResponse{
				ID:         data.CreatedBy.ID,
				Name:       data.CreatedBy.Student.Name,
				Identifier: data.CreatedBy.Identifier,
				Role:       string(data.CreatedBy.Role),
			},
		}

		for _, sup := range data.Thesis.Supervisors {
			schedule.Thesis.Supervisors = append(schedule.Thesis.Supervisors, &dto.CustomUserResponse{
				ID:         sup.Lecturer.ID,
				Name:       sup.Lecturer.Name,
				Identifier: sup.Lecturer.Nip,
				Role:       constants.ENUM_ROLE_LECTURER,
			})
		}

		if data.ApprovedBy != nil {
			schedule.ApprovedBy = &dto.CustomUserResponse{
				ID:         data.ApprovedBy.ID,
				Name:       data.ApprovedBy.Lecturer.Name,
				Identifier: data.ApprovedBy.Identifier,
				Role:       string(data.ApprovedBy.Role),
			}
		}

		schedules = append(schedules, &schedule)
	}
	ss.logger.Info("success get all schedules with pagination",
		zap.Int("count", len(datas.Schedules)),
	)

	return dto.SchedulePaginationResponse{
		Data: schedules,
		PaginationResponse: response.PaginationResponse{
			Page:    datas.Page,
			PerPage: datas.PerPage,
			MaxPage: datas.MaxPage,
			Count:   datas.Count,
		},
	}, nil
}

func (ss *scheduleService) GetDetail(ctx context.Context, id *string) (*dto.ScheduleResponse, error) {
	data, found, err := ss.scheduleRepo.GetScheduleByID(ctx, nil, id)
	if err != nil {
		ss.logger.Error("failed to get schedule by id",
			zap.String("id", *id),
			zap.Error(err),
		)
		return nil, err
	}
	if !found {
		ss.logger.Warn("schedule not found",
			zap.String("id", *id),
		)
		return nil, dto.ErrNotFound
	}

	schedule := dto.ScheduleResponse{
		ID:          data.ID,
		ProposedAt:  data.ProposedAt,
		StartTime:   data.StartTime,
		EndTime:     data.EndTime,
		Status:      data.Status,
		Description: data.Description,
		Location:    data.Location,
		Thesis: dto.ThesisResponse{
			ID:          data.Thesis.ID,
			Title:       data.Thesis.Title,
			Description: data.Thesis.Description,
			Progress:    data.Thesis.Progress,
			Student: &dto.CustomUserResponse{
				ID:         data.Thesis.Student.ID,
				Name:       data.Thesis.Student.Name,
				Identifier: data.Thesis.Student.Nim,
				Role:       constants.ENUM_ROLE_STUDENT,
			},
		},
		CreatedBy: dto.CustomUserResponse{
			ID:         data.CreatedBy.ID,
			Name:       data.CreatedBy.Student.Name,
			Identifier: data.CreatedBy.Identifier,
			Role:       string(data.CreatedBy.Role),
		},
	}

	for _, sup := range data.Thesis.Supervisors {
		schedule.Thesis.Supervisors = append(schedule.Thesis.Supervisors, &dto.CustomUserResponse{
			ID:         sup.Lecturer.ID,
			Name:       sup.Lecturer.Name,
			Identifier: sup.Lecturer.Nip,
			Role:       constants.ENUM_ROLE_LECTURER,
		})
	}

	if data.ApprovedBy != nil {
		schedule.ApprovedBy = &dto.CustomUserResponse{
			ID:         data.ApprovedBy.ID,
			Name:       data.ApprovedBy.Lecturer.Name,
			Identifier: data.ApprovedBy.Identifier,
			Role:       string(data.ApprovedBy.Role),
		}
	}

	ss.logger.Info("success get detail schedule",
		zap.String("id", *id),
	)

	return &schedule, nil
}

func (ss *scheduleService) Update(ctx context.Context, req *dto.UpdateScheduleRequest) (*dto.ScheduleResponse, error) {
	token := ctx.Value("Authorization").(string)
	userIDString, err := ss.jwt.GetUserIDByToken(token)
	if err != nil {
		ss.logger.Error("failed to get user id string from token",
			zap.Error(err),
		)
		return nil, err
	}

	user, found, err := ss.userRepo.GetUserByID(ctx, nil, userIDString)
	if err != nil {
		ss.logger.Error("failed to get user from user id string",
			zap.Error(err),
		)
		return nil, err
	}
	if !found {
		ss.logger.Warn("user not found",
			zap.String("user_id", userIDString),
		)
		return nil, dto.ErrNotFound
	}

	schedule, found, err := ss.scheduleRepo.GetScheduleByID(ctx, nil, &req.ID)
	if err != nil {
		ss.logger.Error("failed to get schedule by id",
			zap.String("id", req.ID),
			zap.Error(err),
		)
		return nil, err
	}
	if !found {
		ss.logger.Warn("schedule not found",
			zap.String("id", req.ID),
		)
		return nil, dto.ErrNotFound
	}

	_, found, err = ss.scheduleRepo.GetThesisByID(ctx, nil, schedule.Thesis.ID.String())
	if err != nil {
		ss.logger.Error("failed to get thesis from thesis id",
			zap.Error(err),
		)
		return nil, err
	}
	if !found {
		ss.logger.Warn("thesis not found",
			zap.String("thesis_id", user.Student.Theses[0].ID.String()),
		)
		return nil, dto.ErrNotFound
	}

	schedule.ProposedAt = req.ProposedAt
	schedule.StartTime = req.StartTime
	schedule.EndTime = req.EndTime
	schedule.Description = req.Description
	schedule.Location = req.Location

	if err := ss.scheduleRepo.UpdateSchedule(ctx, nil, schedule); err != nil {
		ss.logger.Error("failed update schedule",
			zap.Error(err),
		)
		return nil, err
	}

	res := dto.ScheduleResponse{
		ID:          schedule.ID,
		ProposedAt:  schedule.ProposedAt,
		StartTime:   schedule.StartTime,
		EndTime:     schedule.EndTime,
		Status:      schedule.Status,
		Description: schedule.Description,
		Location:    schedule.Location,
		Thesis: dto.ThesisResponse{
			ID:          schedule.Thesis.ID,
			Title:       schedule.Thesis.Title,
			Description: schedule.Thesis.Description,
			Progress:    schedule.Thesis.Progress,
			Student: &dto.CustomUserResponse{
				ID:         schedule.Thesis.Student.ID,
				Name:       schedule.Thesis.Student.Name,
				Identifier: schedule.Thesis.Student.Nim,
				Role:       constants.ENUM_ROLE_STUDENT,
			},
		},
		CreatedBy: dto.CustomUserResponse{
			ID:         user.ID,
			Name:       user.Student.Name,
			Identifier: user.Identifier,
			Role:       string(user.Role),
		},
	}

	for _, sup := range schedule.Thesis.Supervisors {
		res.Thesis.Supervisors = append(res.Thesis.Supervisors, &dto.CustomUserResponse{
			ID:         sup.Lecturer.ID,
			Name:       sup.Lecturer.Name,
			Identifier: sup.Lecturer.Nip,
			Role:       constants.ENUM_ROLE_LECTURER,
		})
	}

	ss.logger.Info("success update schedule",
		zap.String("schedule_id", schedule.ID.String()),
	)

	return &res, nil
}

func (ss *scheduleService) Approval(ctx context.Context, req *dto.ApprovalScheduleRequest) error {
	_, found, err := ss.scheduleRepo.GetScheduleByID(ctx, nil, &req.ID)
	if err != nil {
		ss.logger.Error("failed to get schedule by id",
			zap.String("id", req.ID),
			zap.Error(err),
		)
		return err
	}
	if !found {
		ss.logger.Warn("schedule not found",
			zap.String("id", req.ID),
		)
		return dto.ErrNotFound
	}

	token := ctx.Value("Authorization").(string)
	userIDString, err := ss.jwt.GetUserIDByToken(token)
	if err != nil {
		ss.logger.Error("failed to get user id string from token",
			zap.Error(err),
		)
		return err
	}
	user, found, err := ss.userRepo.GetUserByID(ctx, nil, userIDString)
	if err != nil {
		ss.logger.Error("failed to get user from user id string",
			zap.Error(err),
		)
		return err
	}
	if !found {
		ss.logger.Warn("user not found",
			zap.String("user_id", userIDString),
		)
		return dto.ErrNotFound
	}

	if user.Role == constants.ENUM_ROLE_STUDENT {
		ss.logger.Warn("student cannot approval",
			zap.String("user_id", userIDString),
		)
		return fmt.Errorf("access denied")
	}

	newStatus := req.Status
	if !entity.IsValidScheduleStatus(newStatus) {
		ss.logger.Warn("schedule status not found",
			zap.String("id", req.ID),
			zap.String("status", string(req.Status)),
		)
		return fmt.Errorf("invalid status schedule: %s", newStatus)
	}

	if err := ss.scheduleRepo.UpdateScheduleStatus(ctx, nil, req.ID, string(newStatus), userIDString); err != nil {
		ss.logger.Error("failed to update schedule status by id",
			zap.String("id", req.ID),
			zap.String("status", string(newStatus)),
			zap.Error(err),
		)
		return err
	}

	ss.logger.Info("success update schedule status by id",
		zap.String("id", req.ID),
		zap.String("status", string(newStatus)),
	)

	return nil
}

func (ss *scheduleService) Delete(ctx context.Context, id *string) error {
	_, found, err := ss.scheduleRepo.GetScheduleByID(ctx, nil, id)
	if err != nil {
		ss.logger.Error("failed to get schedule by id before delete",
			zap.String("id", *id),
			zap.Error(err),
		)
		return err
	}
	if !found {
		ss.logger.Warn("schedule not found for delete",
			zap.String("id", *id),
		)
		return dto.ErrNotFound
	}

	// delete schedule
	err = ss.scheduleRepo.DeleteScheduleByID(ctx, nil, id)
	if err != nil {
		ss.logger.Error("failed to delete schedule",
			zap.String("id", *id),
			zap.Error(err),
		)
		return err
	}

	ss.logger.Info("success delete schedule",
		zap.String("id", *id),
	)
	return nil
}

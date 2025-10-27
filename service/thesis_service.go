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
	"go.uber.org/zap"
)

type (
	IThesisService interface {
		GetDetail(ctx context.Context, id string) (*dto.ThesisResponse, error)
		Update(ctx context.Context, req *dto.UpdateThesisRequest) (*dto.ThesisResponse, error)
		GetAllByLecturer(ctx context.Context, req response.PaginationRequest, lecturerID string) (dto.ThesisPaginationResponse, error)
	}

	thesisService struct {
		thesisRepo repository.IThesisRepository
		userRepo   repository.IUserRepository
		logger     *zap.Logger
		jwt        jwt.IJWT
	}
)

func NewThesisService(thesisRepo repository.IThesisRepository, userRepo repository.IUserRepository, logger *zap.Logger, jwt jwt.IJWT) *thesisService {
	return &thesisService{
		thesisRepo: thesisRepo,
		userRepo:   userRepo,
		logger:     logger,
		jwt:        jwt,
	}
}

func (us *thesisService) GetDetail(ctx context.Context, id string) (*dto.ThesisResponse, error) {
	data, found, err := us.thesisRepo.GetThesisByID(ctx, nil, id)
	if err != nil {
		us.logger.Error("failed to get thesis by id",
			zap.String("id", id),
			zap.Error(err),
		)
		return nil, dto.ErrGetThesisByID
	}
	if !found {
		us.logger.Warn("thesis not found",
			zap.String("id", id),
		)
		return nil, dto.ErrNotFound
	}

	thesis := &dto.ThesisResponse{
		ID:          data.ID,
		Title:       data.Title,
		Description: data.Description,
		Progress:    data.Progress,
		Student: &dto.CustomUserResponse{
			ID:         data.Student.ID,
			Name:       data.Student.Name,
			Identifier: data.Student.Nim,
			Role:       constants.ENUM_ROLE_STUDENT,
		},
	}
	for _, sup := range data.Supervisors {
		thesis.Supervisors = append(thesis.Supervisors, &dto.CustomUserResponse{
			ID:         sup.Lecturer.ID,
			Name:       sup.Lecturer.Name,
			Identifier: sup.Lecturer.Nip,
			Role:       constants.ENUM_ROLE_LECTURER,
		})
	}
	us.logger.Info("success get detail thesis",
		zap.String("id", id),
	)

	return thesis, nil
}

func (ts *thesisService) Update(ctx context.Context, req *dto.UpdateThesisRequest) (*dto.ThesisResponse, error) {
	token := ctx.Value("Authorization").(string)
	userIDString, err := ts.jwt.GetUserIDByToken(token)
	if err != nil {
		ts.logger.Error("failed to get user id string from token",
			zap.Error(err),
		)
		return nil, err
	}

	user, found, err := ts.userRepo.GetUserByID(ctx, nil, userIDString)
	if err != nil {
		ts.logger.Error("failed to get user from user id string",
			zap.Error(err),
		)
		return nil, err
	}
	if !found {
		ts.logger.Warn("user not found",
			zap.String("user_id", userIDString),
		)
		return nil, dto.ErrNotFound
	}

	if user.Role == constants.ENUM_ROLE_LECTURER || user.Role == constants.ENUM_ROLE_PRIMARY_LECTURER || user.Role == constants.ENUM_ROLE_SECONDARY_LECTURER {
		ts.logger.Warn("lecturer cannot update the thesis",
			zap.String("user_id", userIDString),
			zap.String("role", string(user.Role)),
		)
		return nil, fmt.Errorf("lecturer cannot update thesis")
	}

	existingThesis, found, err := ts.thesisRepo.GetThesisByID(ctx, nil, req.ID)
	if err != nil {
		ts.logger.Error("failed to get thesis by id",
			zap.String("id", req.ID),
			zap.Error(err),
		)
		return nil, dto.ErrGetThesisByID
	}
	if !found {
		ts.logger.Warn("thesis not found",
			zap.String("id", req.ID),
		)
		return nil, dto.ErrNotFound
	}

	if !entity.IsValidProgress(req.Progress) {
		ts.logger.Warn("invalid thesis progress",
			zap.String("id", req.ID),
			zap.String("progress", string(req.Progress)),
		)
		return nil, fmt.Errorf("invalid thesis progress")
	}

	existingThesis.Title = req.Title
	existingThesis.Description = req.Description
	existingThesis.Progress = req.Progress
	existingThesis.StudentID = req.StudentID

	if err := ts.thesisRepo.UpdateThesis(ctx, nil, existingThesis); err != nil {
		ts.logger.Warn("failed update thesis",
			zap.String("id", req.ID),
		)
		return nil, fmt.Errorf("failed update thesis")
	}

	res := &dto.ThesisResponse{
		ID:          existingThesis.ID,
		Title:       existingThesis.Title,
		Description: existingThesis.Description,
		Progress:    existingThesis.Progress,
		Student: &dto.CustomUserResponse{
			ID:         existingThesis.Student.ID,
			Name:       existingThesis.Student.Name,
			Identifier: existingThesis.Student.Nim,
			Role:       constants.ENUM_ROLE_STUDENT,
		},
	}
	for _, sup := range existingThesis.Supervisors {
		res.Supervisors = append(res.Supervisors, &dto.CustomUserResponse{
			ID:         sup.Lecturer.ID,
			Name:       sup.Lecturer.Name,
			Identifier: sup.Lecturer.Nip,
			Role:       constants.ENUM_ROLE_LECTURER,
		})
	}
	ts.logger.Info("success update thesis",
		zap.String("id", req.ID),
	)

	return res, nil
}

func (ts *thesisService) GetAllByLecturer(ctx context.Context, req response.PaginationRequest, lecturerID string) (dto.ThesisPaginationResponse, error) {
	token := ctx.Value("Authorization").(string)
	userIDString, err := ts.jwt.GetUserIDByToken(token)
	if err != nil {
		ts.logger.Error("failed to get user id string from token",
			zap.Error(err),
		)
		return dto.ThesisPaginationResponse{}, err
	}

	user, found, err := ts.userRepo.GetUserByID(ctx, nil, userIDString)
	if err != nil {
		ts.logger.Error("failed to get user from user id string",
			zap.Error(err),
		)
		return dto.ThesisPaginationResponse{}, err
	}
	if !found {
		ts.logger.Warn("user not found",
			zap.String("user_id", userIDString),
		)
		return dto.ThesisPaginationResponse{}, dto.ErrNotFound
	}

	if user.Role == constants.ENUM_ROLE_STUDENT {
		ts.logger.Warn("student cannot read all lecturer thesis",
			zap.String("user_id", userIDString),
			zap.String("role", string(user.Role)),
		)
		return dto.ThesisPaginationResponse{}, fmt.Errorf("student cannot read all lecturer thesis")
	}

	datas, err := ts.thesisRepo.GetAllThesesByLecturerIDWithPagination(ctx, nil, req, lecturerID)
	if err != nil {
		ts.logger.Error("failed to get all theses with pagination",
			zap.Error(err),
		)
		return dto.ThesisPaginationResponse{}, err
	}

	theses := make([]*dto.ThesisResponse, 0, len(datas.Theses))
	for _, data := range datas.Theses {
		thesis := dto.ThesisResponse{
			ID:          data.ID,
			Title:       data.Title,
			Description: data.Description,
			Progress:    data.Progress,
			Student: &dto.CustomUserResponse{
				ID:         data.Student.ID,
				Name:       data.Student.Name,
				Identifier: data.Student.Nim,
				Role:       constants.ENUM_ROLE_STUDENT,
			},
		}
		for _, sup := range data.Supervisors {
			thesis.Supervisors = append(thesis.Supervisors, &dto.CustomUserResponse{
				ID:         sup.Lecturer.ID,
				Name:       sup.Lecturer.Name,
				Identifier: sup.Lecturer.Nip,
				Role:       constants.ENUM_ROLE_LECTURER,
			})
		}

		theses = append(theses, &thesis)
	}
	ts.logger.Info("success get all theses by lecturer id with pagination",
		zap.Int("count", len(datas.Theses)),
	)

	return dto.ThesisPaginationResponse{
		Data: theses,
		PaginationResponse: response.PaginationResponse{
			Page:    datas.Page,
			PerPage: datas.PerPage,
			MaxPage: datas.MaxPage,
			Count:   datas.Count,
		},
	}, nil
}

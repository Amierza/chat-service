package handler

import (
	"fmt"
	"net/http"

	"github.com/Amierza/chat-service/dto"
	"github.com/Amierza/chat-service/response"
	"github.com/Amierza/chat-service/service"
	"github.com/gin-gonic/gin"
)

type (
	IScheduleHandler interface {
		Create(ctx *gin.Context)
		GetAll(ctx *gin.Context)
		GetDetail(ctx *gin.Context)
		Update(ctx *gin.Context)
		Approval(ctx *gin.Context)
		Delete(ctx *gin.Context)
	}

	scheduleHandler struct {
		scheduleService service.IScheduleService
	}
)

func NewScheduleHandler(scheduleService service.IScheduleService) *scheduleHandler {
	return &scheduleHandler{
		scheduleService: scheduleService,
	}
}

func (sh *scheduleHandler) Create(ctx *gin.Context) {
	var payload *dto.CreateScheduleRequest
	if err := ctx.ShouldBind(&payload); err != nil {
		res := response.BuildResponseFailed(dto.MESSAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	result, err := sh.scheduleService.Create(ctx, payload)
	if err != nil {
		status := mapErrorToStatus(err)
		res := response.BuildResponseFailed(fmt.Sprintf("%s schedules", dto.FAILED_CREATE), err.Error(), nil)
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	res := response.BuildResponseSuccess(fmt.Sprintf("%s schedules", dto.SUCCESS_CREATE), result)
	ctx.JSON(http.StatusOK, res)
}

func (sh *scheduleHandler) GetAll(ctx *gin.Context) {
	var pagination response.PaginationRequest
	if err := ctx.ShouldBindQuery(&pagination); err != nil {
		res := response.BuildResponseFailed(dto.MESSAGE_FAILED_INVALID_QUERY_PARAMS, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	result, err := sh.scheduleService.GetAllWithPagination(ctx, pagination)
	if err != nil {
		status := mapErrorToStatus(err)
		res := response.BuildResponseFailed(fmt.Sprintf("%s schedules", dto.FAILED_GET_ALL), err.Error(), nil)
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	res := response.Response{
		Status:   true,
		Messsage: fmt.Sprintf("%s schedules", dto.SUCCESS_GET_ALL),
		Data:     result.Data,
		Meta:     result.PaginationResponse,
	}
	ctx.JSON(http.StatusOK, res)
}

func (sh *scheduleHandler) GetDetail(ctx *gin.Context) {
	idStr := ctx.Param("id")
	result, err := sh.scheduleService.GetDetail(ctx, &idStr)
	if err != nil {
		status := mapErrorToStatus(err)
		res := response.BuildResponseFailed(fmt.Sprintf("%s schedules", dto.FAILED_GET_DETAIL), err.Error(), nil)
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	res := response.BuildResponseSuccess(fmt.Sprintf("%s schedules", dto.SUCCESS_GET_DETAIL), result)
	ctx.JSON(http.StatusOK, res)
}

func (sh *scheduleHandler) Update(ctx *gin.Context) {
	payload := &dto.UpdateScheduleRequest{}
	idStr := ctx.Param("id")
	payload.ID = idStr
	if err := ctx.ShouldBind(&payload); err != nil {
		res := response.BuildResponseFailed(dto.MESSAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	result, err := sh.scheduleService.Update(ctx, payload)
	if err != nil {
		status := mapErrorToStatus(err)
		res := response.BuildResponseFailed(fmt.Sprintf("%s schedules", dto.FAILED_UPDATE), err.Error(), nil)
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	res := response.BuildResponseSuccess(fmt.Sprintf("%s schedules", dto.SUCCESS_UPDATE), result)
	ctx.JSON(http.StatusOK, res)
}

func (sh *scheduleHandler) Approval(ctx *gin.Context) {
	payload := &dto.ApprovalScheduleRequest{}
	idStr := ctx.Param("id")
	payload.ID = idStr
	if err := ctx.ShouldBind(&payload); err != nil {
		res := response.BuildResponseFailed(dto.MESSAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	err := sh.scheduleService.Approval(ctx, payload)
	if err != nil {
		status := mapErrorToStatus(err)
		res := response.BuildResponseFailed("failed to update approval schedules", err.Error(), nil)
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	res := response.BuildResponseSuccess("success to update approval schedules", nil)
	ctx.JSON(http.StatusOK, res)
}

func (sh *scheduleHandler) Delete(ctx *gin.Context) {
	idStr := ctx.Param("id")
	err := sh.scheduleService.Delete(ctx, &idStr)
	if err != nil {
		status := mapErrorToStatus(err)
		res := response.BuildResponseFailed(fmt.Sprintf("%s schedules", dto.FAILED_DELETE), err.Error(), nil)
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	res := response.BuildResponseSuccess(fmt.Sprintf("%s schedules", dto.SUCCESS_DELETE), nil)
	ctx.JSON(http.StatusOK, res)
}

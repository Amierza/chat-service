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
	IThesisHandler interface {
		GetDetail(ctx *gin.Context)
		Update(ctx *gin.Context)
		GetAllByLecturer(ctx *gin.Context)
	}

	thesisHandler struct {
		thesisService service.IThesisService
	}
)

func NewThesisHandler(thesisService service.IThesisService) *thesisHandler {
	return &thesisHandler{
		thesisService: thesisService,
	}
}

func (th *thesisHandler) GetDetail(ctx *gin.Context) {
	idStr := ctx.Param("id")
	result, err := th.thesisService.GetDetail(ctx, idStr)
	if err != nil {
		status := mapErrorToStatus(err)
		res := response.BuildResponseFailed(fmt.Sprintf("%s thesis", dto.FAILED_GET_DETAIL), err.Error(), nil)
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	res := response.BuildResponseSuccess(fmt.Sprintf("%s thesis", dto.SUCCESS_GET_DETAIL), result)
	ctx.JSON(http.StatusOK, res)
}

func (th *thesisHandler) Update(ctx *gin.Context) {
	payload := &dto.UpdateThesisRequest{}
	idStr := ctx.Param("id")
	payload.ID = idStr
	if err := ctx.ShouldBind(&payload); err != nil {
		res := response.BuildResponseFailed(dto.MESSAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	result, err := th.thesisService.Update(ctx, payload)
	if err != nil {
		status := mapErrorToStatus(err)
		res := response.BuildResponseFailed(fmt.Sprintf("%s thesis", dto.FAILED_UPDATE), err.Error(), nil)
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	res := response.BuildResponseSuccess(fmt.Sprintf("%s thesis", dto.SUCCESS_UPDATE), result)
	ctx.JSON(http.StatusOK, res)
}

func (th *thesisHandler) GetAllByLecturer(ctx *gin.Context) {
	var pagination response.PaginationRequest
	if err := ctx.ShouldBindQuery(&pagination); err != nil {
		res := response.BuildResponseFailed(dto.MESSAGE_FAILED_INVALID_QUERY_PARAMS, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	idStr := ctx.Param("lecturer_id")
	result, err := th.thesisService.GetAllByLecturer(ctx, pagination, idStr)
	if err != nil {
		status := mapErrorToStatus(err)
		res := response.BuildResponseFailed(fmt.Sprintf("%s lecturer thesis", dto.FAILED_GET_ALL), err.Error(), nil)
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	res := response.Response{
		Status:   true,
		Messsage: fmt.Sprintf("%s theses", dto.SUCCESS_GET_ALL),
		Data:     result.Data,
		Meta:     result.PaginationResponse,
	}
	ctx.JSON(http.StatusOK, res)
}

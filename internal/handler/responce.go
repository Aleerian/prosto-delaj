package handler

import (
	"github.com/sirupsen/logrus"
	"prosto-delaj-api/internal/exceptions"

	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) sendResponseSuccess(c *gin.Context, successResponse any, err exceptions.ErrorCode) {
	if successResponse == nil {
		if err == exceptions.NoContent {
			c.Status(http.StatusNoContent)
			return
		}
		code, response := getFailedResponse(err)
		c.AbortWithStatusJSON(code, struct {
			Code    string `json:"code"`
			Message any    `json:"message"`
		}{
			Code:    response.ErrorCode.String(),
			Message: response.Message,
		})
		return
	}
	c.JSON(http.StatusOK, successResponse)
}

func (h *Handler) sendResponseCreated(c *gin.Context, successResponse any, err exceptions.ErrorCode) {
	if successResponse == nil {
		if err == exceptions.NoContent {
			c.Status(http.StatusNoContent)
			return
		}
		code, response := getFailedResponse(err)
		c.AbortWithStatusJSON(code, struct {
			Code    string `json:"code"`
			Message any    `json:"message"`
		}{
			Code:    response.ErrorCode.String(),
			Message: response.Message,
		})
		return
	}
	c.JSON(http.StatusCreated, successResponse)
}

func (h *Handler) sendResponseError(c *gin.Context, err exceptions.ErrorCode) {
	code, response := getFailedResponse(err)
	c.AbortWithStatusJSON(code, struct {
		Code    string `json:"code"`
		Message any    `json:"message"`
	}{
		Code:    response.ErrorCode.String(),
		Message: response.Message,
	})
}

// getFailedResponse Возвращает http code status
func getFailedResponse(err exceptions.ErrorCode) (int, exceptions.FailedResponseBody) {
	failedResponse, isFound := exceptions.ErrorCodeToFailedResponse[err]
	if !isFound {
		logrus.Error("the specified error code not found")
		return http.StatusInternalServerError, exceptions.FailedResponseBody{
			ErrorCode: err,
			Message:   exceptions.ErrorServer.Error(),
		}
	}

	return int(failedResponse.HttpCode), exceptions.FailedResponseBody{
		ErrorCode: err,
		Message:   failedResponse.Message,
	}
}

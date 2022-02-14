package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hsjsjsj009/kubeEP/kubeEP-BE/internal/constant"
	"github.com/hsjsjsj009/kubeEP/kubeEP-BE/internal/entity/response"
	"net/http"
)

type baseHandler struct {
}

func (h baseHandler) errorResponse(c *fiber.Ctx, data interface{}) error {
	return c.Status(http.StatusBadRequest).JSON(&response.Base{
		Status: constant.Error,
		Data:   data,
	})
}

func (h baseHandler) errorResponseWithCode(c *fiber.Ctx, data interface{}, code int) error {
	return c.Status(http.StatusBadRequest).JSON(&response.Base{
		Code:   code,
		Status: constant.Error,
		Data:   data,
	})
}

func (h baseHandler) successResponse(c *fiber.Ctx, data interface{}) error {
	return c.JSON(&response.Base{
		Data:   data,
		Status: constant.Success,
	})
}

func (h baseHandler) failResponse(c *fiber.Ctx, data interface{}) error {
	return c.JSON(&response.Base{
		Data:   data,
		Status: constant.Fail,
	})
}

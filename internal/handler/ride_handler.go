package handler

import (
	"github.com/Furkanturan8/motorbike-rental-backend-v2/internal/dto"
	"github.com/Furkanturan8/motorbike-rental-backend-v2/internal/model"
	"github.com/Furkanturan8/motorbike-rental-backend-v2/internal/service"
	"github.com/Furkanturan8/motorbike-rental-backend-v2/pkg/errorx"
	"github.com/Furkanturan8/motorbike-rental-backend-v2/pkg/response"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

type RideHandler struct {
	service *service.RideService
}

func NewRideHandler(s *service.RideService) *RideHandler {
	return &RideHandler{service: s}
}

func (h *RideHandler) Create(c *fiber.Ctx) error {
	var req dto.CreateRideRequest
	if err := c.BodyParser(&req); err != nil {
		return errorx.Wrap(errorx.ErrInvalidRequest, err)
	}

	ride := req.ToDBModel(model.Ride{})

	if err := h.service.Create(c.Context(), &ride); err != nil {
		return errorx.WithDetails(errorx.ErrInternal, err.Error())
	}

	return response.Success(c, nil, "Ride başarıyla oluşturuldu")
}

func (h *RideHandler) GetByID(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return errorx.Wrap(errorx.ErrInvalidRequest, err)
	}

	resp, err := h.service.GetByID(c.Context(), int64(id))
	if err != nil {
		return errorx.WithDetails(errorx.ErrNotFound, "Ride bulunamadı")
	}

	ride := dto.RideResponse{}.ToResponseModel(*resp)

	return response.Success(c, ride)
}

func (h *RideHandler) Update(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return errorx.Wrap(errorx.ErrInvalidRequest, err)
	}

	var req dto.UpdateRideRequest
	if err = c.BodyParser(&req); err != nil {
		return errorx.Wrap(errorx.ErrInvalidRequest, err)
	}

	_, err = h.service.GetByID(c.Context(), int64(id))
	if err != nil {
		return err
	}

	ride := req.ToDBModel(model.Ride{})

	if err = h.service.Update(c.Context(), ride); err != nil {
		return errorx.WithDetails(errorx.ErrInternal, err.Error())
	}

	return response.Success(c, nil, "Ride başarıyla güncellendi")
}

func (h *RideHandler) Delete(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return errorx.Wrap(errorx.ErrInvalidRequest, err)
	}

	if err = h.service.Delete(c.Context(), int64(id)); err != nil {
		return errorx.WithDetails(errorx.ErrInternal, err.Error())
	}

	return response.Success(c, nil, "Ride başarıyla silindi")
}

func (h *RideHandler) List(c *fiber.Ctx) error {
	resp, err := h.service.List(c.Context())
	if err != nil {
		return errorx.WithDetails(errorx.ErrInternal, err.Error())
	}

	rides := make([]dto.RideResponse, len(resp))
	for i, item := range resp {
		rides[i] = dto.RideResponse{}.ToResponseModel(item)
	}
	return response.Success(c, rides)
}

func (h *RideHandler) ListRideByUserID(c *fiber.Ctx) error {
	param := c.Params("userID")
	userID, err := strconv.Atoi(param)
	if err != nil {
		return errorx.WithDetails(errorx.ErrInternal, "Geçersiz kullanıcı kimliği")
	}

	resp, err := h.service.GetByUserID(c.Context(), int64(userID))
	if err != nil {
		return errorx.WithDetails(errorx.ErrInternal, err.Error())
	}

	rides := make([]dto.RideResponse, len(resp))
	for i, item := range resp {
		rides[i] = dto.RideResponse{}.ToResponseModel(item)
	}

	return response.Success(c, rides)
}

package handler

import (
	"github.com/Furkanturan8/motorbike-rental-backend-v2/internal/dto"
	"github.com/Furkanturan8/motorbike-rental-backend-v2/internal/model"
	"github.com/Furkanturan8/motorbike-rental-backend-v2/internal/service"
	"github.com/Furkanturan8/motorbike-rental-backend-v2/pkg/errorx"
	"github.com/Furkanturan8/motorbike-rental-backend-v2/pkg/response"
	"github.com/gofiber/fiber/v2"
)

type MotorbikeHandler struct {
	service *service.MotorbikeService
}

func NewMotorbikeHandler(s *service.MotorbikeService) *MotorbikeHandler {
	return &MotorbikeHandler{service: s}
}

func (h *MotorbikeHandler) Create(c *fiber.Ctx) error {
	var req dto.CreateMotorbikeRequest
	if err := c.BodyParser(&req); err != nil {
		return errorx.WrapErr(errorx.ErrInvalidRequest, err)
	}

	motorbike := req.ToDBModel(model.Motorbike{})

	if err := h.service.Create(c.Context(), &motorbike); err != nil {
		return errorx.WrapErr(errorx.ErrInternal, err)
	}

	return response.Success(c, nil, "Motorbike başarıyla oluşturuldu")
}

func (h *MotorbikeHandler) GetByID(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return errorx.WrapErr(errorx.ErrInvalidRequest, err)
	}

	resp, err := h.service.GetByID(c.Context(), int64(id))
	if err != nil {
		return errorx.WrapMsg(errorx.ErrNotFound, "Motorbike bulunamadı")
	}

	motorbike := dto.MotorbikeResponse{}.ToResponseModel(*resp)

	return response.Success(c, motorbike)
}

func (h *MotorbikeHandler) Update(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return errorx.WrapErr(errorx.ErrInvalidRequest, err)
	}

	var req dto.UpdateMotorbikeRequest
	if err = c.BodyParser(&req); err != nil {
		return errorx.WrapErr(errorx.ErrInvalidRequest, err)
	}

	currentMotorbike, err := h.service.GetByID(c.Context(), int64(id))
	if err != nil {
		return err
	}

	updatedMotorbike := req.ToDBModel(*currentMotorbike)

	if err = h.service.Update(c.Context(), updatedMotorbike); err != nil {
		return errorx.WrapErr(errorx.ErrInternal, err)
	}

	return response.Success(c, nil, "Motorbike başarıyla güncellendi")
}

func (h *MotorbikeHandler) Delete(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return errorx.WrapErr(errorx.ErrInvalidRequest, err)
	}

	if err = h.service.Delete(c.Context(), int64(id)); err != nil {
		return errorx.WrapErr(errorx.ErrInternal, err)
	}

	return response.Success(c, nil, "Motorbike başarıyla silindi")
}

// todo: fotoları repoda preload/ilişkili şekilde getir.
func (h *MotorbikeHandler) List(c *fiber.Ctx) error {
	resp, err := h.service.List(c.Context())
	if err != nil {
		return errorx.WrapErr(errorx.ErrInternal, err)
	}

	motorbikes := make([]dto.MotorbikeResponse, len(resp))
	for i, item := range resp {
		motorbikes[i] = dto.MotorbikeResponse{}.ToResponseModel(item)
	}
	return response.Success(c, motorbikes)
}

func (h *MotorbikeHandler) GetAvailableMotors(c *fiber.Ctx) error {
	resp, err := h.service.GetMotorsForStatus(c.Context(), string(model.BikeAvailable))
	if err != nil {
		return errorx.WrapErr(errorx.ErrInternal, err)
	}

	if len(resp) == 0 {
		return response.Success(c, nil, "Müsait motor bulunamadı!")
	}

	motorbikes := make([]dto.MotorbikeResponse, len(resp))
	for i, item := range resp {
		motorbikes[i] = dto.MotorbikeResponse{}.ToResponseModel(item)
	}
	return response.Success(c, motorbikes)
}

func (h *MotorbikeHandler) GetMaintenanceMotors(c *fiber.Ctx) error {
	resp, err := h.service.GetMotorsForStatus(c.Context(), string(model.BikeInMaintenance))
	if err != nil {
		return errorx.WrapErr(errorx.ErrInternal, err)
	}

	if len(resp) == 0 {
		return response.Success(c, nil, "Bakımda motor yok!")
	}

	motorbikes := make([]dto.MotorbikeResponse, len(resp))
	for i, item := range resp {
		motorbikes[i] = dto.MotorbikeResponse{}.ToResponseModel(item)
	}

	return response.Success(c, motorbikes)
}

func (h *MotorbikeHandler) GetRentedMotors(c *fiber.Ctx) error {
	return nil
}

func (h *MotorbikeHandler) GetPhotosByID(c *fiber.Ctx) error {
	return nil
}

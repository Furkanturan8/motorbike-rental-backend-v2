package handler

import (
	"github.com/Furkanturan8/motorbike-rental-backend-v2/internal/dto"
	"github.com/Furkanturan8/motorbike-rental-backend-v2/internal/model"
	"github.com/Furkanturan8/motorbike-rental-backend-v2/internal/service"
	"github.com/Furkanturan8/motorbike-rental-backend-v2/pkg/errorx"
	"github.com/Furkanturan8/motorbike-rental-backend-v2/pkg/response"
	"github.com/gofiber/fiber/v2"
)

type BluetoothConnectionHandler struct {
	service *service.BluetoothConnectionService
}

func NewBluetoothConnectionHandler(s *service.BluetoothConnectionService) *BluetoothConnectionHandler {
	return &BluetoothConnectionHandler{service: s}
}

func (h *BluetoothConnectionHandler) Create(c *fiber.Ctx) error {
	var req dto.CreateBluetoothConnectionRequest
	if err := c.BodyParser(&req); err != nil {
		return errorx.WrapErr(errorx.ErrInvalidRequest, err)
	}

	bluetoothConnection := req.ToDBModel(model.BluetoothConnection{})

	if err := h.service.Create(c.Context(), &bluetoothConnection); err != nil {
		return errorx.WrapErr(errorx.ErrInternal, err)
	}

	return response.Success(c, nil, "Bluetooth_Connection başarıyla oluşturuldu")
}

func (h *BluetoothConnectionHandler) GetByID(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return errorx.WrapErr(errorx.ErrInvalidRequest, err)
	}

	resp, err := h.service.GetByID(c.Context(), int64(id))
	if err != nil {
		return errorx.WrapMsg(errorx.ErrNotFound, "Bluetooth_Connection bulunamadı")
	}

	bluetoothConnection := dto.BluetoothConnectionResponse{}.ToResponseModel(*resp)

	return response.Success(c, bluetoothConnection)
}

func (h *BluetoothConnectionHandler) Update(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return errorx.WrapErr(errorx.ErrInvalidRequest, err)
	}

	var req dto.UpdateBluetoothConnectionRequest
	if err = c.BodyParser(&req); err != nil {
		return errorx.WrapErr(errorx.ErrInvalidRequest, err)
	}

	_, err = h.service.GetByID(c.Context(), int64(id))
	if err != nil {
		return err
	}

	bluetoothConnection := req.ToDBModel(model.BluetoothConnection{})

	if err = h.service.Update(c.Context(), bluetoothConnection); err != nil {
		return errorx.WrapErr(errorx.ErrInternal, err)
	}

	return response.Success(c, nil, "Bluetooth_Connection başarıyla güncellendi")
}

func (h *BluetoothConnectionHandler) Delete(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return errorx.WrapErr(errorx.ErrInvalidRequest, err)
	}

	if err = h.service.Delete(c.Context(), int64(id)); err != nil {
		return errorx.WrapErr(errorx.ErrInternal, err)
	}

	return response.Success(c, nil, "Bluetooth_Connection başarıyla silindi")
}

func (h *BluetoothConnectionHandler) List(c *fiber.Ctx) error {
	resp, err := h.service.List(c.Context())
	if err != nil {
		return errorx.WrapErr(errorx.ErrInternal, err)
	}

	bluetoothConnection := make([]dto.BluetoothConnectionResponse, len(resp))
	for i, item := range resp {
		bluetoothConnection[i] = dto.BluetoothConnectionResponse{}.ToResponseModel(item)
	}
	return response.Success(c, bluetoothConnection)
}

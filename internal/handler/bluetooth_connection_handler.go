package handler

import (
	"github.com/Furkanturan8/motorbike-rental-backend-v2/internal/dto"
	"github.com/Furkanturan8/motorbike-rental-backend-v2/internal/model"
	"github.com/Furkanturan8/motorbike-rental-backend-v2/internal/service"
	"github.com/Furkanturan8/motorbike-rental-backend-v2/pkg/errorx"
	"github.com/Furkanturan8/motorbike-rental-backend-v2/pkg/response"
	"github.com/gofiber/fiber/v2"
	"time"
)

type BluetoothConnectionHandler struct {
	service          *service.BluetoothConnectionService
	motorbikeService *service.MotorbikeService
}

func NewBluetoothConnectionHandler(s *service.BluetoothConnectionService, m *service.MotorbikeService) *BluetoothConnectionHandler {
	return &BluetoothConnectionHandler{service: s, motorbikeService: m}
}

func (h *BluetoothConnectionHandler) GetMyConnections(ctx *fiber.Ctx) error {
	userID := ctx.Locals("userID").(int64)
	if userID == 0 {
		return errorx.WrapMsg(errorx.ErrInvalidRequest, " Kullanıcı bulunamadı")
	}

	resp, err := h.service.GetByUserID(ctx.Context(), userID)
	if err != nil {
		return errorx.WrapErr(errorx.ErrInternal, err)
	}

	bluetoothConnection := make([]dto.BluetoothConnectionResponse, len(resp))
	for i, item := range resp {
		bluetoothConnection[i] = dto.BluetoothConnectionResponse{}.ToResponseModel(item)
	}
	return response.Success(ctx, bluetoothConnection)
}

func (h *BluetoothConnectionHandler) Connect(ctx *fiber.Ctx) error {
	var req dto.ConnectRequest
	if err := ctx.BodyParser(&req); err != nil {
		return errorx.WrapErr(errorx.ErrInvalidRequest, err)
	}

	userID := ctx.Locals("userID").(int64)
	if userID == 0 {
		return errorx.WrapMsg(errorx.ErrInvalidRequest, " Kullanıcı bulunamadı")
	}

	motor, err := h.motorbikeService.GetByID(ctx.Context(), req.MotorbikeID)
	if err != nil {
		return errorx.WrapMsg(errorx.ErrInternal, "Motorbisiklet bulunamadı")
	}

	// Motorbike'ın durumu 'Available' mı kontrol et
	if motor.Status != model.BikeAvailable {
		return errorx.WrapMsg(errorx.ErrInvalidRequest, "Bu Motorbisiklet şu anda müsait değil!")
	}

	var bluetoothConnection model.BluetoothConnection
	bluetoothConnectionOK := true // todo: burada ardinuo bağlantısı kontrol edilecek. Ardinuo dan gelen true olsun şimdilik
	if bluetoothConnectionOK {
		bluetoothConnection = req.ToDBModel(model.BluetoothConnection{})
		// bağlandığı için lock statusu değiştir
		motor.LockStatus = model.Unlocked
	}

	bluetoothConnection.UserID = userID

	if err = h.service.Create(ctx.Context(), &bluetoothConnection); err != nil {
		return errorx.WrapErr(errorx.ErrInternal, err)
	}

	motor.Status = model.BikeRented

	err = h.motorbikeService.Update(ctx.Context(), *motor)
	if err != nil {
		return errorx.WrapMsg(errorx.ErrInternal, "Motor status güncellenirken hata oluştu!")
	}

	return response.Success(ctx, nil, "Bluetooth bağlantısı başarıyla kuruldu")
}

func (h *BluetoothConnectionHandler) Create(ctx *fiber.Ctx) error {
	var req dto.CreateBluetoothConnectionRequest
	if err := ctx.BodyParser(&req); err != nil {
		return errorx.WrapErr(errorx.ErrInvalidRequest, err)
	}

	bluetoothConnection := req.ToDBModel(model.BluetoothConnection{})

	if err := h.service.Create(ctx.Context(), &bluetoothConnection); err != nil {
		return errorx.WrapErr(errorx.ErrInternal, err)
	}

	return response.Success(ctx, nil, "Bluetooth_Connection başarıyla oluşturuldu")
}

func (h *BluetoothConnectionHandler) GetByID(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")
	if err != nil {
		return errorx.WrapErr(errorx.ErrInvalidRequest, err)
	}

	resp, err := h.service.GetByID(ctx.Context(), int64(id))
	if err != nil {
		return errorx.WrapMsg(errorx.ErrNotFound, "Bluetooth_Connection bulunamadı")
	}

	bluetoothConnection := dto.BluetoothConnectionResponse{}.ToResponseModel(*resp)

	return response.Success(ctx, bluetoothConnection)
}

func (h *BluetoothConnectionHandler) Update(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")
	if err != nil {
		return errorx.WrapErr(errorx.ErrInvalidRequest, err)
	}

	var req dto.UpdateBluetoothConnectionRequest
	if err = ctx.BodyParser(&req); err != nil {
		return errorx.WrapErr(errorx.ErrInvalidRequest, err)
	}

	_, err = h.service.GetByID(ctx.Context(), int64(id))
	if err != nil {
		return err
	}

	bluetoothConnection := req.ToDBModel(model.BluetoothConnection{})

	if err = h.service.Update(ctx.Context(), bluetoothConnection); err != nil {
		return errorx.WrapErr(errorx.ErrInternal, err)
	}

	return response.Success(ctx, nil, "Bluetooth_Connection başarıyla güncellendi")
}

func (h *BluetoothConnectionHandler) Delete(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")
	if err != nil {
		return errorx.WrapErr(errorx.ErrInvalidRequest, err)
	}

	if err = h.service.Delete(ctx.Context(), int64(id)); err != nil {
		return errorx.WrapErr(errorx.ErrInternal, err)
	}

	return response.Success(ctx, nil, "Bluetooth_Connection başarıyla silindi")
}

func (h *BluetoothConnectionHandler) List(ctx *fiber.Ctx) error {
	resp, err := h.service.List(ctx.Context())
	if err != nil {
		return errorx.WrapErr(errorx.ErrInternal, err)
	}

	bluetoothConnection := make([]dto.BluetoothConnectionResponse, len(resp))
	for i, item := range resp {
		bluetoothConnection[i] = dto.BluetoothConnectionResponse{}.ToResponseModel(item)
	}
	return response.Success(ctx, bluetoothConnection)
}

func (h *BluetoothConnectionHandler) Disconnect(ctx *fiber.Ctx) error {
	var req dto.ConnectRequest
	if err := ctx.BodyParser(&req); err != nil {
		return errorx.WrapErr(errorx.ErrInvalidRequest, err)
	}

	userID := ctx.Locals("userID").(int64)
	if userID == 0 {
		return errorx.WrapMsg(errorx.ErrInvalidRequest, " Kullanıcı bulunamadı")
	}

	connection, err := h.service.GetByMotorbikeID(ctx.Context(), req.MotorbikeID)
	if err != nil {
		return errorx.WrapMsg(errorx.ErrInternal, "Bluetooth bağlantısı bulunamadı")
	}

	// zaten bağlantı koptuysa..
	if connection.DisconnectedAt != nil {
		return errorx.WrapMsg(errorx.ErrInvalidRequest, "Zaten bağlantı kopmuş!")
	}

	bluetoothDisconnectOK := true // todo: bağlantı kesildi mi?? burada ardinuo bağlantısı kontrol edilecek. Ardinuo dan gelen true olsun şimdilik
	if bluetoothDisconnectOK {
		now := time.Now()
		connection.DisconnectedAt = &now
		if err = h.service.Update(ctx.Context(), *connection); err != nil {
			return errorx.WrapMsg(errorx.ErrInternal, "Bağlantı kesilemedi!")
		}
	}

	motor, err := h.motorbikeService.GetByID(ctx.Context(), connection.MotorbikeID)
	if err != nil {
		return errorx.WrapErr(errorx.ErrInternal, err)
	}

	motor.Status = model.BikeAvailable
	motor.LockStatus = model.Locked

	err = h.motorbikeService.Update(ctx.Context(), *motor)
	if err != nil {
		return errorx.WrapMsg(errorx.ErrInternal, "Motor status güncellenirken hata oluştu!")
	}

	return response.Success(ctx, nil, "Bluetooth bağlantısı başarıyla kesildi")
}

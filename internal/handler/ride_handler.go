package handler

import (
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/Furkanturan8/motorbike-rental-backend-v2/internal/dto"
	"github.com/Furkanturan8/motorbike-rental-backend-v2/internal/model"
	"github.com/Furkanturan8/motorbike-rental-backend-v2/internal/service"
	"github.com/Furkanturan8/motorbike-rental-backend-v2/pkg/errorx"
	"github.com/Furkanturan8/motorbike-rental-backend-v2/pkg/response"
	"github.com/gofiber/fiber/v2"
)

type RideHandler struct {
	rideService *service.RideService
}

func NewRideHandler(rideService *service.RideService) *RideHandler {
	return &RideHandler{
		rideService: rideService,
	}
}

func (h *RideHandler) Create(c *fiber.Ctx) error {
	var req dto.CreateRideRequest
	if err := c.BodyParser(&req); err != nil {
		return errorx.WrapErr(errorx.ErrInvalidRequest, err)
	}

	ride := req.ToDBModel(model.Ride{})

	if err := h.rideService.Create(c.Context(), &ride); err != nil {
		return errorx.WrapErr(errorx.ErrInternal, err)
	}

	return response.Success(c, nil, "Ride başarıyla oluşturuldu")
}

func (h *RideHandler) GetByID(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return errorx.WrapErr(errorx.ErrInvalidRequest, err)
	}

	resp, err := h.rideService.GetByID(c.Context(), int64(id))
	if err != nil {
		return errorx.WrapMsg(errorx.ErrNotFound, "Ride bulunamadı")
	}

	ride := dto.RideResponse{}.ToResponseModel(*resp)

	return response.Success(c, ride)
}

func (h *RideHandler) Update(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return errorx.WrapErr(errorx.ErrInvalidRequest, err)
	}

	var req dto.UpdateRideRequest
	if err = c.BodyParser(&req); err != nil {
		return errorx.WrapErr(errorx.ErrInvalidRequest, err)
	}

	currentRide, err := h.rideService.GetByID(c.Context(), int64(id))
	if err != nil {
		return err
	}

	ride := req.ToDBModel(*currentRide)

	if err = h.rideService.Update(c.Context(), ride); err != nil {
		return errorx.WrapErr(errorx.ErrInternal, err)
	}

	return response.Success(c, nil, "Ride başarıyla güncellendi")
}

func (h *RideHandler) Delete(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return errorx.WrapErr(errorx.ErrInvalidRequest, err)
	}

	if err = h.rideService.Delete(c.Context(), int64(id)); err != nil {
		return errorx.WrapErr(errorx.ErrInternal, err)
	}

	return response.Success(c, nil, "Ride başarıyla silindi")
}

func (h *RideHandler) List(c *fiber.Ctx) error {
	resp, err := h.rideService.List(c.Context())
	if err != nil {
		return errorx.WrapErr(errorx.ErrInternal, err)
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
		return errorx.WrapMsg(errorx.ErrInternal, "Geçersiz kullanıcı kimliği")
	}

	resp, err := h.rideService.GetByUserID(c.Context(), int64(userID))
	if err != nil {
		return errorx.WrapErr(errorx.ErrInternal, err)
	}

	rides := make([]dto.RideResponse, len(resp))
	for i, item := range resp {
		rides[i] = dto.RideResponse{}.ToResponseModel(item)
	}

	return response.Success(c, rides)
}

func (h *RideHandler) ListMyRides(c *fiber.Ctx) error {
	userID := c.Locals("userID").(int64)

	resp, err := h.rideService.GetByUserID(c.Context(), userID)
	if err != nil {
		return errorx.WrapErr(errorx.ErrInternal, err)
	}

	rides := make([]dto.RideResponse, len(resp))
	for i, item := range resp {
		rides[i] = dto.RideResponse{}.ToResponseModel(item)
	}

	return response.Success(c, rides)
}

func (h *RideHandler) ListRideByMotorbikeID(c *fiber.Ctx) error {
	param := c.Params("motorbikeID")
	motorbikeID, err := strconv.Atoi(param)
	if err != nil {
		return errorx.WrapMsg(errorx.ErrInternal, "Geçersiz motorbike kimliği")
	}

	resp, err := h.rideService.GetByMotorbikeID(c.Context(), int64(motorbikeID))
	if err != nil {
		return errorx.WrapErr(errorx.ErrInternal, err)
	}

	rides := make([]dto.RideResponse, len(resp))
	for i, item := range resp {
		rides[i] = dto.RideResponse{}.ToResponseModel(item)
	}

	return response.Success(c, rides)
}

func (h *RideHandler) FinishRide(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")
	if err != nil {
		return errorx.WrapErr(errorx.ErrInvalidRequest, err)
	}

	userID := ctx.Locals("userID").(int64)

	ride, err := h.rideService.FinishRide(ctx.Context(), int64(id), userID)
	if err != nil {
		return err // zaten wrap edilmiş şekilde dönüyor
	}

	return response.Success(ctx, dto.RideResponse{}.ToResponseModel(*ride), "Sürüş Bitirildi! Ücret: "+strconv.FormatFloat(ride.Cost, 'f', 2, 64)+" TL")
}

func (h *RideHandler) AddRidePhoto(ctx *fiber.Ctx) error {
	rideID, err := ctx.ParamsInt("id")
	if err != nil {
		return errorx.WrapErr(errorx.ErrInvalidRequest, err)
	}

	photo, err := ctx.FormFile("photo")
	if err != nil {
		return errorx.WrapMsg(errorx.ErrInvalidRequest, "Fotoğraf yüklenemedi")
	}

	// Dosya kaydedileceği yolu oluştur
	dir := "uploads/rides"
	if err = os.MkdirAll(dir, os.ModePerm); err != nil {
		return errorx.WrapMsg(errorx.ErrInternal, "Yükleme klasörü oluşturulamadı")
	}
	filePath := filepath.Join(dir, "ride_id_"+strconv.Itoa(rideID)+"_name_"+photo.Filename)

	if err = ctx.SaveFile(photo, filePath); err != nil {
		return errorx.WrapMsg(errorx.ErrInternal, "Fotoğraf kaydedilemedi")
	}

	// Motorun kilitlenip kilitlenmediğini kontrol et
	if err = h.rideService.HandleAfterPhotoUpload(ctx.Context(), rideID); err != nil {
		return err
	}

	return response.Success(ctx, nil, "Fotoğraf yüklendi ve motor bağlantısı kesildi.")
}

func (h *RideHandler) ListByDateRange(ctx *fiber.Ctx) error {
	startTimeStr := ctx.Query("start_time")
	endTimeStr := ctx.Query("end_time")

	if startTimeStr == "" || endTimeStr == "" {
		return errorx.WrapMsg(errorx.ErrInvalidRequest, "start_time ve end_time parametreleri zorunludur")
	}

	startTime, err := time.Parse("2006-01-02", startTimeStr)
	if err != nil {
		return errorx.WrapMsg(errorx.ErrInvalidRequest, "start_time formatı geçersiz. Beklenen format: YYYY-MM-DD")
	}

	endTime, err := time.Parse("2006-01-02", endTimeStr)
	if err != nil {
		return errorx.WrapMsg(errorx.ErrInvalidRequest, "end_time formatı geçersiz. Beklenen format: YYYY-MM-DD")
	}

	// Bitiş zamanını günün sonuna al (23:59:59) dahil etmek için
	endTime = endTime.Add(time.Hour*23 + time.Minute*59 + time.Second*59)

	rides, err := h.rideService.ListByDateRange(ctx.Context(), startTime, endTime)
	if err != nil {
		return errorx.WrapErr(errorx.ErrInternal, err)
	}

	resp := make([]dto.RideResponse, len(*rides))
	for i, r := range *rides {
		resp[i] = dto.RideResponse{}.ToResponseModel(r)
	}

	return response.Success(ctx, resp, "Başarıyla getirildi")
}

package response

import (
	"github.com/gofiber/fiber/v2"
	"reflect"
)

const (
	StatusOK = fiber.StatusOK
)

// Response yapısı
type Response struct {
	Success   bool        `json:"success"`
	Data      interface{} `json:"data,omitempty"`
	Message   interface{} `json:"message,omitempty"`
	DataCount int         `json:"data_count,omitempty"`
}

// Başarılı yanıt oluşturmak için yardımcı fonksiyonlar
func Success(c *fiber.Ctx, data interface{}, message ...string) error {
	var msg interface{}
	if len(message) > 0 {
		msg = message[0]
	}

	resp := Response{
		Success: true,
		Data:    data,
		Message: msg,
	}

	// Eğer data bir slice ise, count al
	if data != nil && reflect.TypeOf(data).Kind() == reflect.Slice {
		resp.DataCount = reflect.ValueOf(data).Len()
	}

	return c.Status(StatusOK).JSON(resp)
}

// Başarılı yanıt - veri olmadan
func SuccessNoData(c *fiber.Ctx) error {
	return c.Status(StatusOK).JSON(Response{
		Success: true,
	})
}

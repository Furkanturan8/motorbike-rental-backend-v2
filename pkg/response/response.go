package response

import "github.com/gofiber/fiber/v2"

const (
	StatusOK = fiber.StatusOK
)

// Response yapısı
type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Message interface{} `json:"message,omitempty"`
}

// Başarılı yanıt oluşturmak için yardımcı fonksiyonlar
func Success(c *fiber.Ctx, data interface{}, message ...string) error {
	var msg interface{}
	if len(message) > 0 {
		msg = message[0]
	}

	return c.Status(StatusOK).JSON(Response{
		Success: true,
		Data:    data,
		Message: msg,
	})
}

// Başarılı yanıt - veri olmadan
func SuccessNoData(c *fiber.Ctx) error {
	return c.Status(StatusOK).JSON(Response{
		Success: true,
	})
}

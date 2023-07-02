package validation

import (
	"fmt"

	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

var DecodedToken jwt.MapClaims

func Authenticator(req *fiber.Ctx) error {
	headers := req.GetReqHeaders()
	if headers["Authorization"] == "" {
		return req.Status(400).JSON(fiber.Map{
			"msg": "authorization header is required!",
		})
	} else {
		rawToken := headers["Authorization"]
		headerToken := strings.Split(rawToken, " ")
		token, err := verifyJwt(headerToken[1])
		if err == nil {
			DecodedToken = token
			return req.Next()
		}

		fmt.Println(err)
		return req.Status(401).JSON(fiber.Map{
			"msg": "user authentication failed!",
		})
	}
}

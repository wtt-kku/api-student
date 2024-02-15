package middleware

import (
	"errors"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"

	"student_score/utils"
)

func GenerateJWT(userId string) (string, error) {
	timeExp, err := strconv.Atoi(os.Getenv("ACCESS_TOKEN_EXP_TIME_MIN"))
	if err != nil {
		return "", err
	}

	currentTime := time.Now()
	expirationDuration := time.Duration(timeExp) * time.Minute

	// Create the claims for the access token
	accessTokenClaims := &jwt.RegisteredClaims{
		Subject:   userId,
		ExpiresAt: jwt.NewNumericDate(currentTime.Add(expirationDuration)),
		IssuedAt:  jwt.NewNumericDate(currentTime),
	}

	// Generate access token
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims)
	accessTokenString, err := accessToken.SignedString([]byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		return "", err
	}

	return accessTokenString, nil
}

func JWTMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return c.JSON(http.StatusUnauthorized, utils.Response{
				Result:  false,
				Code:    utils.CommonRespCode["UNAUTHORIZED"].Code,
				Message: utils.CommonRespCode["UNAUTHORIZED"].Message,
			})
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			return c.JSON(http.StatusUnauthorized, utils.Response{
				Result:  false,
				Code:    utils.CommonRespCode["INVALID_TOKEN"].Code,
				Message: utils.CommonRespCode["INVALID_TOKEN"].Message,
			})
		}

		keyFunc := func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("SECRET_KEY")), nil
		}

		claims := new(jwt.RegisteredClaims)

		token, err := jwt.ParseWithClaims(tokenString, claims, keyFunc)
		if err != nil {
			if errors.Is(err, jwt.ErrTokenExpired) {
				return c.JSON(http.StatusUnauthorized, utils.Response{
					Result:  false,
					Code:    utils.CommonRespCode["TOKEN_EXPIRED"].Code,
					Message: utils.CommonRespCode["TOKEN_EXPIRED"].Message,
				})
			} else {
				return c.JSON(http.StatusUnauthorized, utils.Response{
					Result:  false,
					Code:    utils.CommonRespCode["INVALID_TOKEN"].Code,
					Message: utils.CommonRespCode["INVALID_TOKEN"].Message,
				})
			}
		}

		if !token.Valid || len(claims.Subject) == 0 {
			return c.JSON(http.StatusUnauthorized, utils.Response{
				Result:  false,
				Code:    utils.CommonRespCode["INVALID_TOKEN"].Code,
				Message: utils.CommonRespCode["INVALID_TOKEN"].Message,
			})
		}

		c.Set("userId", claims.Subject)
		return next(c)
	}
}

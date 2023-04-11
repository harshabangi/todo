package service_echo

import (
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/harsha-aqfer/todo/pkg"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strings"
	"time"
)

const tokenExpirySec = 24 * 3600 // 2 hours

type Claims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

func signUp(c echo.Context) error {
	s := c.Get("service").(*Service)

	var req pkg.User
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := req.Validate(); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	hashedPassword, err := hashPassword(req.Password)
	if err != nil {
		return err
	}
	req.Password = hashedPassword

	if err = s.db.User.CreateUser(&req); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, pkg.NewMsgResp("Successfully signed up!"))
}

func signIn(c echo.Context) error {
	s := c.Get("service").(*Service)

	var req pkg.User
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	user, err := s.db.User.GetUser(req.Email)
	if err != nil {
		return err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, pkg.NewMsgResp("incorrect password"))
	}

	token, err := generateToken(user.Email, s.conf.SigningKey)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, &pkg.Token{ExpiresIn: tokenExpirySec, JWTToken: token})
}

func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func mkJwtToken(signKey []byte, claims jwt.Claims) (string, error) {
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(signKey)
}

func generateToken(email string, signingKey string) (string, error) {
	now := time.Now().Unix()

	claims := Claims{
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  now,
			ExpiresAt: now + tokenExpirySec,
		},
		Email: email,
	}

	return mkJwtToken([]byte(signingKey), claims)
}

type SecurityContext struct {
	Email  string
	UserID int64
}

func IsAuthorized(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		const authBearer = "Bearer"

		var (
			s = c.Get("service").(*Service)
			h = c.Request().Header.Get("Authorization")
		)

		if !strings.HasPrefix(h, authBearer) {
			return echo.NewHTTPError(http.StatusUnauthorized)
		}

		var (
			authToken = strings.Trim(h[len(authBearer):], " ")
			claims    = &Claims{}
		)

		tkn, err := jwt.ParseWithClaims(authToken, claims, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("error parsing the token")
			}
			return []byte(s.conf.SigningKey), nil
		})

		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				return echo.NewHTTPError(http.StatusUnauthorized)
			}
			return echo.NewHTTPError(http.StatusBadRequest)
		}

		if !tkn.Valid {
			return echo.NewHTTPError(http.StatusUnauthorized)
		}

		userID, err := s.db.User.GetUserID(claims.Email)
		if err != nil {
			return err
		}

		c.Set("security_context", &SecurityContext{Email: claims.Email, UserID: userID})
		return next(c)
	}
}

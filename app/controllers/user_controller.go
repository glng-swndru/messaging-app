package controllers

import (
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/kooroshh/fiber-boostrap/app/models"
	"github.com/kooroshh/fiber-boostrap/app/repository"
	jwttoken "github.com/kooroshh/fiber-boostrap/pkg/jwt_token"
	"github.com/kooroshh/fiber-boostrap/pkg/response"
	"go.elastic.co/apm"
	"golang.org/x/crypto/bcrypt"
)

func Register(ctx *fiber.Ctx) error {
	span, spanCtx := apm.StartSpan(ctx.Context(), "Register", "controller")
	defer span.End()

	user := new(models.User)
	err := ctx.BodyParser(user)
	if err != nil {
		errResponse := fmt.Errorf("failed to parse request: %v", err)
		log.Println(errResponse)
		return response.SendFailureResponse(ctx, fiber.StatusInternalServerError, errResponse.Error(), nil)
	}
	err = user.Validate()
	if err != nil {
		errResponse := fmt.Errorf("failed to validate request: %v", err)
		log.Println(errResponse)
		return response.SendFailureResponse(ctx, fiber.StatusInternalServerError, errResponse.Error(), nil)
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		errResponse := fmt.Errorf("failed to encrypt the password: %v", err)
		log.Println(errResponse)
		return response.SendFailureResponse(ctx, fiber.StatusInternalServerError, errResponse.Error(), nil)
	}
	user.Password = string(hashPassword)

	err = repository.InsertNewUser(spanCtx, user)
	if err != nil {
		errResponse := fmt.Errorf("failed to insert new user: %v", err)
		log.Println(errResponse)
		return response.SendFailureResponse(ctx, fiber.StatusInternalServerError, errResponse.Error(), nil)
	}

	resp := user
	resp.Password = ""

	return response.SendSuccessResponse(ctx, resp)
}

func Login(ctx *fiber.Ctx) error {
	span, spanCtx := apm.StartSpan(ctx.Context(), "Login", "controller")
	defer span.End()

	// parsing request and validation request
	var (
		loginReq = new(models.LoginRequest)
		resp     models.LoginResponse
		now      = time.Now()
	)

	err := ctx.BodyParser(loginReq)
	if err != nil {
		errResponse := fmt.Errorf("failed to parse request: %v", err)
		log.Println(errResponse)
		return response.SendFailureResponse(ctx, fiber.StatusInternalServerError, errResponse.Error(), nil)
	}
	err = loginReq.Validate()
	if err != nil {
		errResponse := fmt.Errorf("failed to validate request: %v", err)
		log.Println(errResponse)
		return response.SendFailureResponse(ctx, fiber.StatusInternalServerError, errResponse.Error(), nil)
	}

	user, err := repository.GetUserByUsername(spanCtx, loginReq.Username)
	if err != nil {
		errResponse := fmt.Errorf("failed to get username: %v", err)
		log.Println(errResponse)
		return response.SendFailureResponse(ctx, fiber.StatusNotFound, "username/password salah", nil)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginReq.Password))
	if err != nil {
		errResponse := fmt.Errorf("failed to check password: %v", err)
		log.Println(errResponse)
		return response.SendFailureResponse(ctx, fiber.StatusNotFound, "username/password salah", nil)
	}

	token, err := jwttoken.GenerateToken(spanCtx, user.Username, user.Fullname, "token", now)
	if err != nil {
		errResponse := fmt.Errorf("failed to generate token: %v", err)
		log.Println(errResponse)
		return response.SendFailureResponse(ctx, fiber.StatusNotFound, "an error occurred in the system", nil)
	}

	refreshToken, err := jwttoken.GenerateToken(spanCtx, user.Username, user.Fullname, "refresh token", now)
	if err != nil {
		errResponse := fmt.Errorf("failed to refresh token: %v", err)
		log.Println(errResponse)
		return response.SendFailureResponse(ctx, fiber.StatusNotFound, "an error occurred in the system", nil)
	}

	userSession := &models.UserSession{
		UserID:              user.ID,
		Token:               token,
		RefreshToken:        refreshToken,
		TokenExpired:        now.Add(jwttoken.MapTypeToken["token"]),
		RefreshTokenExpired: now.Add(jwttoken.MapTypeToken["refresh_token"]),
	}

	err = repository.InsertNewUserSession(spanCtx, userSession)
	if err != nil {
		errResponse := fmt.Errorf("failed insert user session: %v", err)
		log.Println(errResponse)
		return response.SendFailureResponse(ctx, fiber.StatusInternalServerError, "terjadi kesalahan pada sistem", nil)
	}

	resp.Username = user.Username
	resp.FullName = user.Fullname
	resp.Token = token
	resp.RefreshToken = refreshToken

	return response.SendSuccessResponse(ctx, resp)
}

func Logout(ctx *fiber.Ctx) error {
	span, spanCtx := apm.StartSpan(ctx.Context(), "Logout", "controller")
	defer span.End()

	token := ctx.Get("Authorization")
	err := repository.DeleteUserSessionByToken(spanCtx, token)
	if err != nil {
		errResponse := fmt.Errorf("failed delete user session: %v", err)
		log.Println(errResponse)
		return response.SendFailureResponse(ctx, fiber.StatusInternalServerError, "terjadi kesalahan pada sistem", nil)
	}

	return response.SendSuccessResponse(ctx, nil)
}
func RefreshToken(ctx *fiber.Ctx) error {
	span, spanCtx := apm.StartSpan(ctx.Context(), "RefreshToken", "controller")
	defer span.End()

	now := time.Now()
	refreshToken := ctx.Get("Authorization")
	userName := ctx.Locals("username").(string)
	fullName := ctx.Locals("full_name").(string)

	token, err := jwttoken.GenerateToken(spanCtx, userName, fullName, "token", now)
	if err != nil {
		errResponse := fmt.Errorf("failed to generate token: %v", err)
		log.Println(errResponse)
		return response.SendFailureResponse(ctx, fiber.StatusInternalServerError, "terjadi kesalahan pada sistem", nil)
	}

	err = repository.UpdateUserSessionToken(spanCtx, token, now.Add(jwttoken.MapTypeToken["token"]), refreshToken)
	if err != nil {
		errResponse := fmt.Errorf("failed to update token: %v", err)
		log.Println(errResponse)
		return response.SendFailureResponse(ctx, fiber.StatusInternalServerError, "terjadi kesalahan pada sistem", nil)
	}

	return response.SendSuccessResponse(ctx, fiber.Map{
		"token": token,
	})
}

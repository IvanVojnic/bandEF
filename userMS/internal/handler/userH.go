package handler

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"net/http"
	"userMS/models"
)

type RequestBody struct {
	SenderID     uuid.UUID     `json:"senderID"`
	ReceiverID   uuid.UUID     `json:"receiverID"`
	UserEmail    string        `json:"userEmail"`
	UsersInvited []models.User `json:"usersInvited"`
	statusID     int           `json:"statusID"`
}

func (h *Handler) AcceptFriendsRequest(ctx echo.Context) error {
	var reqBody RequestBody
	err := ctx.Bind(&reqBody)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"Error Bind json while creating book": err,
			"book":                                reqBody,
		}).Info("Bind json")
		return echo.NewHTTPError(http.StatusInternalServerError, "data not correct")
	}
	userID := ctx.Get("user_id").(uuid.UUID)
	err = h.userS.AcceptFriendsRequest(ctx.Request().Context(), reqBody.SenderID, userID)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"Error while accepting request": err,
		}).Info("Accept request")
		return echo.NewHTTPError(http.StatusBadRequest, "book creating failed")
	}
	return ctx.String(http.StatusOK, "request accepted")
}

func (h *Handler) GetFriends(ctx echo.Context) error {
	userID := ctx.Get("user_id").(uuid.UUID)
	rooms, err := h.userS.GetFriends(ctx.Request().Context(), userID)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"Error get rooms": err,
			"rooms":           rooms,
		}).Info("Get rooms request")
		return echo.NewHTTPError(http.StatusBadRequest, "book updating failed")
	}
	return ctx.JSON(http.StatusOK, rooms)
}

func (h *Handler) SendFriendsRequest(ctx echo.Context) error {
	var reqBody RequestBody
	errBind := ctx.Bind(&reqBody)
	if errBind != nil {
		logrus.WithFields(logrus.Fields{
			"Error Bind json while send request to be a friend to another user": errBind,
			"reqBody": reqBody,
		}).Info("Bind json")
		return echo.NewHTTPError(http.StatusInternalServerError, "data not correct")
	}
	userID := ctx.Get("user_id").(uuid.UUID)
	err := h.userS.SendFriendsRequest(ctx.Request().Context(), userID, reqBody.ReceiverID)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"Error while send friends request": err,
		}).Info("send friends request")
		return echo.NewHTTPError(http.StatusBadRequest, "smth went wrong")
	}
	return ctx.String(http.StatusOK, "request sent")
}

func (h *Handler) FindUser(ctx echo.Context) error {
	var reqBody RequestBody
	errBind := ctx.Bind(&reqBody)
	if errBind != nil {
		logrus.WithFields(logrus.Fields{
			"Error while send friends request": errBind,
		}).Info("send friends request")
		return echo.NewHTTPError(http.StatusBadRequest, "wrong data")
	}
	user, err := h.userS.FindUser(ctx.Request().Context(), reqBody.UserEmail)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"Error find user": err,
			"user":            user,
		}).Info("GET find user")
		return echo.NewHTTPError(http.StatusBadRequest)
	}
	return ctx.JSON(http.StatusOK, user)
}

func (h *Handler) GetRequest(ctx echo.Context) error {
	userID := ctx.Get("user_id").(uuid.UUID)
	users, err := h.userS.GetRequest(ctx.Request().Context(), userID)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"Error get users who sent request to be a friends": err,
			"users": users,
		}).Info("GET users request")
		return echo.NewHTTPError(http.StatusBadRequest, "cannot get users")
	}
	return ctx.JSON(http.StatusOK, users)
}

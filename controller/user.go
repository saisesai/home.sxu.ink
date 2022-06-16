package controller

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/saisesai/home.sxu.ink/log"
	"github.com/saisesai/home.sxu.ink/model"
	"gorm.io/gorm"
	"net/http"
	"time"
)

func HandleUserLogin(ctx *gin.Context) {
	var err error
	// get username and password
	username := ctx.PostForm("username")
	password := ctx.PostForm("password")
	// build setup logger
	L := log.WithField("ip", ctx.RemoteIP()).
		WithField("username", username).
		WithField("password", password)
	// check data format
	if username == "" || password == "" {
		L.Warningln(ErrInvalidDataFormat)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": ErrInvalidDataFormat.Error(),
		})
		return
	}
	// query
	user, err := model.FindUserByUsername(username)
	// internal error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		L.WithError(err).Errorln("failed to find user!")
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
		return
	}
	// user not found or wrong password
	if errors.Is(err, gorm.ErrRecordNotFound) || password != user.Password {
		L.Warningln("user not found or wrong password!")
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "user not found or wrong password",
		})
		return
	}
	// build jwt
	tokenString, err := BuildMapClaimsJwt(jwt.MapClaims{
		"username": username,
		"expire":   time.Now().Add(time.Hour * 24 * 7).Unix(), // expire in 7 days
	})
	if err != nil {
		L.WithError(err).Errorln("failed to build jwt!")
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
		return
	}
	// ok
	ctx.JSON(http.StatusOK, gin.H{
		"msg":   "ok",
		"token": tokenString,
	})
}

func HandleUserAddDevice(ctx *gin.Context) {
	var err error
	L := log.WithField("ip", ctx.RemoteIP())
	// parse jwt
	claims, err := ParseMapClaimsJwtHeader(ctx)
	if err != nil {
		L.Debugln(err)
		L.Warningln("invalid token!")
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"msg": "invalid token",
		})
		return
	}
	// parse body
	var data map[string]string
	err = ctx.BindJSON(&data)
	if err != nil {
		L.Debugln(err)
		L.Warningln(ErrInvalidDataFormat)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": ErrInvalidDataFormat.Error(),
		})
		return
	}
	// get data
	username := claims["username"].(string)
	L = L.WithField("username", username)
	deviceType, deviceTypeOk := data["type"]
	deviceClientId, deviceClientIdOk := data["client_id"]
	deviceRemark, deviceRemarkOk := data["remark"]
	// check data format
	if !deviceTypeOk || !deviceClientIdOk || !deviceRemarkOk ||
		(deviceType != "relay") || deviceClientId == "" || deviceRemark == "" {
		L.Warningln(ErrInvalidDataFormat)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": ErrInvalidDataFormat.Error(),
		})
		return
	}
	// find user
	user, err := model.FindUserByUsername(username)
	if err != nil {
		L.WithError(err).Errorln("failed to find user!")
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
		return
	}
	// add device client id to user
	err = user.AddDevice(deviceClientId)
	if err != nil {
		if errors.Is(err, model.ErrInvalidDataFormat) || errors.Is(err, model.ErrUserDeviceAlreadyExist) {
			L.Warningln(ErrInvalidDataFormat)
			ctx.JSON(http.StatusBadRequest, gin.H{
				"msg": err.Error(),
			})
			return
		}
		L.WithError(err).Errorln("failed to add device to user!")
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
		return
	}
	// add remark
	err = user.SetDeviceRemark(deviceClientId, deviceRemark)
	if err != nil {
		L.Warningln(err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": err.Error(),
		})
		return
	}
	// add device
	switch deviceType {
	case "relay":
		// add device type
		err = user.SetDeviceType(deviceClientId, "relay")
		if err != nil {
			L.Warningln(err)
			ctx.JSON(http.StatusBadRequest, gin.H{
				"msg": err.Error(),
			})
			return
		}
		HandleRelayAdd(L, ctx, deviceClientId)
	default:
		L.Warningln("unknown device type: " + deviceType)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "unknown device type: " + deviceType,
		})
	}
}

func HandleUserDeleteDevice(ctx *gin.Context) {
	var err error
	L := log.WithField("ip", ctx.RemoteIP())
	// check jwt
	claims, err := ParseMapClaimsJwtHeader(ctx)
	if err != nil {
		L.Debugln(err)
		L.Warningln("invalid token!")
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"msg": "invalid token",
		})
		return
	}
	// parse body
	var data map[string]string
	err = ctx.BindJSON(&data)
	if err != nil {
		L.Debugln(err)
		L.Warningln(ErrInvalidDataFormat)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": ErrInvalidDataFormat.Error(),
		})
		return
	}
	// get data
	username := claims["username"].(string)
	L = L.WithField("username", username)
	deviceType, deviceTypeOk := data["type"]
	deviceClientId, deviceClientIdOk := data["client_id"]
	// check data format
	if !deviceTypeOk || !deviceClientIdOk || (deviceType != "relay") || deviceClientId == "" {
		L.Warningln(ErrInvalidDataFormat)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": ErrInvalidDataFormat.Error(),
		})
		return
	}
	// find user
	user, err := model.FindUserByUsername(username)
	if err != nil {
		L.WithError(err).Errorln("failed to find user!")
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
		return
	}
	// remove device type
	err = user.DeleteDeviceType(deviceClientId)
	if err != nil {
		L.Warningln(err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": err.Error(),
		})
		return
	}
	// remove device remarks
	err = user.DeleteDeviceRemark(deviceClientId)
	if err != nil {
		L.Warningln(err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": err.Error(),
		})
		return
	}
	// remove device client id from user
	err = user.DeleteDevice(deviceClientId)
	if err != nil {
		// check if user owns device
		if errors.Is(err, model.ErrUserDeviceNotExist) {
			L.Warningln(model.ErrUserDeviceNotExist)
			ctx.JSON(http.StatusBadRequest, gin.H{
				"msg": model.ErrUserDeviceNotExist.Error(),
			})
			return
		}
		// internal error
		L.WithError(err).Errorln("failed to delete user!")
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
		return
	}
	// delete device
	switch deviceType {
	case "relay":
		HandleRelayDelete(L, ctx, deviceClientId)
	default:
		L.Warningln("unknown device type: " + deviceType)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "unknown device type: " + deviceType,
		})
	}
}

func HandleUserModifyDevice(ctx *gin.Context) {
	var err error
	L := log.WithField("ip", ctx.RemoteIP())
	// check jwt
	claims, err := ParseMapClaimsJwtHeader(ctx)
	if err != nil {
		L.Debugln(err)
		L.Warningln("invalid token!")
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"msg": "invalid token",
		})
		return
	}
	username := claims["username"].(string)
	L = L.WithField("username", username)
	// parse body
	var data map[string]string
	err = ctx.BindJSON(&data)
	if err != nil {
		L.Debugln(err)
		L.Warningln(ErrInvalidDataFormat)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": ErrInvalidDataFormat.Error(),
		})
		return
	}
	// get data
	remark, remarkOk := data["remark"]
	deviceClientId, deviceClientIdOk := data["client_id"]
	if remark == "" || !remarkOk ||
		deviceClientId == "" || !deviceClientIdOk {
		L.Warningln(ErrInvalidDataFormat)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": ErrInvalidDataFormat.Error(),
		})
		return
	}
	// find user
	user, err := model.FindUserByUsername(username)
	if err != nil {
		L.WithError(err).Errorln("failed to find user!")
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
		return
	}
	// modify
	err = user.SetDeviceRemark(deviceClientId, remark)
	if err != nil {
		L.Warningln(err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": err.Error(),
		})
		return
	}
	// ok
	ctx.JSON(http.StatusOK, gin.H{
		"msg": "ok",
	})
}

func HandleUserGetDevices(ctx *gin.Context) {
	var err error
	L := log.WithField("ip", ctx.RemoteIP())
	// check jwt
	claims, err := ParseMapClaimsJwtHeader(ctx)
	if err != nil {
		L.Debugln(err)
		L.Warningln("invalid token!")
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"msg": "invalid token",
		})
		return
	}
	username := claims["username"].(string)
	L = L.WithField("username", username)
	// return devices
	user, err := model.FindUserByUsername(username)
	if err != nil {
		L.WithError(err).Errorln("failed to find user!")
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
		return
	}
	// return client_id, device remarks and types
	var out []map[string]string
	devices := user.GetDevices()
	types := user.GetDeviceTypes()
	remarks := user.GetDeviceRemarks()
	for _, v := range devices {
		out = append(out, map[string]string{
			"type":      types[v],
			"client_id": v,
			"remark":    remarks[v],
		})
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg":     "ok",
		"devices": out,
	})
}

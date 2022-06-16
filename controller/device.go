package controller

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/saisesai/home.sxu.ink/log"
	"github.com/saisesai/home.sxu.ink/model"
	"github.com/saisesai/home.sxu.ink/mqtt"
	"gorm.io/gorm"
	"net/http"
)

func HandleDeviceInfo(ctx *gin.Context) {
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
	username := claims["username"].(string)
	L = L.WithField("username", username)
	// get data
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
	// check owner
	var f = false
	userDevices := user.GetDevices()
	for _, v := range userDevices {
		if v == deviceClientId {
			f = true
		}
	}
	if !f {
		L.Warningln("forbid to get device!")
		ctx.JSON(http.StatusForbidden, gin.H{
			"msg": "forbidden",
		})
		return
	}
	// find device
	switch deviceType {
	case "relay":
		relay, err := model.FindRelayByClientId(deviceClientId)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) { // not found
				L.WithField("client_id", deviceClientId).Warningln("rely not exist!")
				ctx.JSON(http.StatusBadRequest, gin.H{
					"msg": ErrRelayNotExist.Error(),
				})
				return
			}
			// internal error
			L.WithError(err).Errorln("failed to find relay:", relay.ClientId)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"msg": err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"msg":    "ok",
			"sv":     relay.SoftwareVersion,
			"hv":     relay.HardwareVersion,
			"online": relay.Online,
			"on":     relay.On,
		})
	}
}

func HandleDeviceCmd(ctx *gin.Context) {
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
	username := claims["username"].(string)
	L = L.WithField("username", username)
	deviceType, deviceTypeOk := data["type"]
	deviceClientId, deviceClientIdOk := data["client_id"]
	deviceCmd, deviceCmdOk := data["cmd"]
	if deviceType != "relay" || !deviceTypeOk ||
		deviceClientId == "" || !deviceClientIdOk ||
		deviceCmd == "" || !deviceCmdOk {
		L.Warningln(ErrInvalidDataFormat)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": ErrInvalidDataFormat.Error(),
		})
		return
	}
	switch deviceType {
	case "relay":
		err = mqtt.SendRelayCmd(deviceClientId, deviceCmd)
		if err != nil {
			if errors.Is(err, mqtt.ErrDeviceNotExist) || errors.Is(err, mqtt.ErrInvalidCmd) {
				L.WithField("client_id", deviceClientId).Warningln(err)
				ctx.JSON(http.StatusBadRequest, gin.H{
					"msg": err.Error(),
				})
				return
			}
			// internal error
			L.Errorln(err)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"msg": err.Error(),
			})
			return
		}
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg": "ok",
	})
}

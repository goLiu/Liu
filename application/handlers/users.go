package handlers

import (
	"community-governance/application/models"
	"community-governance/application/utils"
	dbMod "community-governance/db/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

const jwtSecretKey = "community-secret-key"

func Login(c *gin.Context) {
	var info models.LoginInfo
	err := c.ShouldBind(&info)
	if err != nil {
		c.JSON(http.StatusBadRequest, "传入的参数不合法:"+err.Error())
	}
	id, userType, err := dbMod.Login(info.IDNumber, info.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "登录失败:"+err.Error())
		return
	}
	if id == "" {
		c.JSON(http.StatusInternalServerError, "登录失败,用户不存在")
		return
	}
	token, err := utils.GenerateJWT(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "生成token失败:"+err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token, "type": userType, "data": "登录成功"})
}

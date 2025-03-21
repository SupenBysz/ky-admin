package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/SupenBysz/ky-admin/internal/model"
	"github.com/SupenBysz/ky-admin/internal/service"
)

// UserController 用户控制器
type UserController struct {
	userService service.UserService
}

// NewUserController 创建用户控制器
func NewUserController(userService service.UserService) *UserController {
	return &UserController{
		userService: userService,
	}
}

// GetUsers 获取用户列表
func (c *UserController) GetUsers(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "10"))

	users, total, err := c.userService.GetUsers(page, pageSize)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": 500,
			"msg":  "获取用户列表失败：" + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "success",
		"data": gin.H{
			"total": total,
			"list":  users,
		},
	})
}

// GetUserByID 根据ID获取用户
func (c *UserController) GetUserByID(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": 400,
			"msg":  "无效的用户ID",
		})
		return
	}

	user, err := c.userService.GetByID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": 500,
			"msg":  "获取用户失败：" + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "success",
		"data": user,
	})
}

// CreateUser 创建用户
func (c *UserController) CreateUser(ctx *gin.Context) {
	var user model.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": 400,
			"msg":  "无效的请求参数：" + err.Error(),
		})
		return
	}

	if err := c.userService.CreateUser(&user); err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": 500,
			"msg":  "创建用户失败：" + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "success",
		"data": user.ID,
	})
}

// UpdateUser 更新用户
func (c *UserController) UpdateUser(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": 400,
			"msg":  "无效的用户ID",
		})
		return
	}

	var user model.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": 400,
			"msg":  "无效的请求参数：" + err.Error(),
		})
		return
	}

	user.ID = uint(id)
	if err := c.userService.UpdateUser(&user); err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": 500,
			"msg":  "更新用户失败：" + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "success",
	})
}

// DeleteUser 删除用户
func (c *UserController) DeleteUser(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": 400,
			"msg":  "无效的用户ID",
		})
		return
	}

	if err := c.userService.DeleteUser(uint(id)); err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": 500,
			"msg":  "删除用户失败：" + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "success",
	})
}

package handlers

import (
	"log"
	"net/http"
	"regexp"
	"smartwatch-server/api/auth"
	"smartwatch-server/api/models"
	"smartwatch-server/api/repository"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// AuthHandler 认证处理器
type AuthHandler struct {
	repo *repository.AuthRepository
}

// NewAuthHandler 创建认证处理器
func NewAuthHandler(repo *repository.AuthRepository) *AuthHandler {
	return &AuthHandler{repo: repo}
}

// RegisterReq 注册请求
type RegisterReq struct {
	Username        string `json:"username" binding:"required"`
	Password        string `json:"password" binding:"required"`
	ConfirmPassword string `json:"confirm_password" binding:"required"`
	Phone           string `json:"phone"`
	Email           string `json:"email"`
}

// LoginReq 登录请求
type LoginReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// AuthResp 认证响应
type AuthResp struct {
	Token     string       `json:"token"`
	ExpiresAt int64        `json:"expires_at"`
	User      *UserInfoDTO `json:"user"`
}

// UserInfoDTO 用户信息（不含敏感字段）
type UserInfoDTO struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Phone    string `json:"phone,omitempty"`
	Email    string `json:"email,omitempty"`
	Role     string `json:"role"`
}

var (
	usernameRegex = regexp.MustCompile(`^[a-zA-Z0-9_]{3,32}$`)
	phoneRegex    = regexp.MustCompile(`^1[3-9]\d{9}$`)
	emailRegex    = regexp.MustCompile(`^[^\s@]+@[^\s@]+\.[^\s@]+$`)
)

// Register 注册
func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误: " + err.Error()})
		return
	}

	// 校验
	if !usernameRegex.MatchString(req.Username) {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "用户名需 3-32 位字母、数字或下划线"})
		return
	}
	if len(req.Password) < 6 || len(req.Password) > 32 {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "密码长度需 6-32 位"})
		return
	}
	if req.Password != req.ConfirmPassword {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "两次密码不一致"})
		return
	}
	if req.Phone != "" && !phoneRegex.MatchString(req.Phone) {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "手机号格式不正确"})
		return
	}
	if req.Email != "" && !emailRegex.MatchString(req.Email) {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "邮箱格式不正确"})
		return
	}

	// 查重
	exists, err := h.repo.ExistsUsername(req.Username)
	if err != nil {
		log.Printf("注册查重失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "服务异常"})
		return
	}
	if exists {
		c.JSON(http.StatusConflict, gin.H{"code": 409, "message": "用户名已存在"})
		return
	}
	if req.Phone != "" {
		exists, _ = h.repo.ExistsPhone(req.Phone)
		if exists {
			c.JSON(http.StatusConflict, gin.H{"code": 409, "message": "手机号已注册"})
			return
		}
	}
	if req.Email != "" {
		exists, _ = h.repo.ExistsEmail(req.Email)
		if exists {
			c.JSON(http.StatusConflict, gin.H{"code": 409, "message": "邮箱已注册"})
			return
		}
	}

	// 密码哈希
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("密码哈希失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "服务异常"})
		return
	}

	u := &models.AdminUser{
		Username:     req.Username,
		PasswordHash: string(hash),
		Role:         "admin",
		Status:       1,
	}
	if req.Phone != "" {
		u.Phone = &req.Phone
	}
	if req.Email != "" {
		u.Email = &req.Email
	}

	if err := h.repo.CreateUser(u); err != nil {
		log.Printf("创建用户失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "注册失败"})
		return
	}

	token, expiresAt, err := auth.IssueToken(u)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "生成令牌失败"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"code":    0,
		"message": "注册成功",
		"data": AuthResp{
			Token:     token,
			ExpiresAt: expiresAt,
			User:      toUserInfoDTO(u),
		},
	})
}

// Login 密码登录
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	var u *models.AdminUser
	var err error
	if emailRegex.MatchString(req.Username) {
		u, err = h.repo.GetByEmail(req.Username)
	} else if phoneRegex.MatchString(req.Username) {
		u, err = h.repo.GetByPhone(req.Username)
	} else {
		u, err = h.repo.GetByUsername(req.Username)
	}

	if err == gorm.ErrRecordNotFound || u == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "用户名或密码错误"})
		return
	}
	if err != nil {
		log.Printf("登录查询失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "服务异常"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "用户名或密码错误"})
		return
	}

	token, expiresAt, err := auth.IssueToken(u)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "生成令牌失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "登录成功",
		"data": AuthResp{
			Token:     token,
			ExpiresAt: expiresAt,
			User:      toUserInfoDTO(u),
		},
	})
}

// GetMe 获取当前用户信息（需鉴权）
func (h *AuthHandler) GetMe(c *gin.Context) {
	userID, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "未登录"})
		return
	}
	username, _ := c.Get("username")
	role, _ := c.Get("role")

	// 从 DB 查完整信息（含 phone/email）
	u, err := h.repo.GetByUsername(username.(string))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"data": UserInfoDTO{
				ID:       userID.(uint),
				Username: username.(string),
				Role:     safeStr(role),
			},
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": toUserInfoDTO(u),
	})
}

func toUserInfoDTO(u *models.AdminUser) *UserInfoDTO {
	dto := &UserInfoDTO{ID: u.ID, Username: u.Username, Role: u.Role}
	if u.Phone != nil {
		dto.Phone = *u.Phone
	}
	if u.Email != nil {
		dto.Email = *u.Email
	}
	return dto
}

func safeStr(v interface{}) string {
	if s, ok := v.(string); ok {
		return s
	}
	return ""
}

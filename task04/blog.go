package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// ======================== 模型定义 ========================

type User struct {
	gorm.Model
	Username string
	Password string
	Email    string
	Posts    []Post
}

type Post struct {
	gorm.Model
	UserID   uint
	Title    string
	Content  string
	Comments []Comment
}

type Comment struct {
	gorm.Model
	Content string
	PostID  uint
	UserID  uint
}

// ======================== 请求体结构 ========================

type UserRegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Email    string `json:"email" binding:"required"`
}

type UserLoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type CreatePostRequest struct {
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
}

type UpdatePostRequest struct {
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
}

type AddCommentRequest struct {
	Content string `json:"content" binding:"required"`
}

// ======================== 响应结构与工具 ========================

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    200,
		Message: "success",
		Data:    data,
	})
}

func Fail(c *gin.Context, code int, message string, err error) {
	msg := message
	if err != nil {
		msg = fmt.Sprintf("%s: %v", message, err.Error())
	}
	c.JSON(code, Response{
		Code:    code,
		Message: msg,
	})
	c.Abort()
}

// ======================== 初始化与日志 ========================

var (
	jwtSecret = []byte(os.Getenv("BLOG_SECRET"))
	logger    *zap.Logger
)

func init() {
	var err error
	logger, err = zap.NewProduction()
	if err != nil {
		panic(err)
	}
}

func Connect() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("blog.db"), &gorm.Config{})
	if err != nil {
		logger.Fatal("数据库连接失败", zap.Error(err))
	}
	return db
}

// ======================== 主程序入口 ========================

func main() {
	db := Connect()
	db.AutoMigrate(&User{}, &Post{}, &Comment{})

	r := gin.Default()

	// 全局请求日志中间件
	r.Use(func(c *gin.Context) {
		start := time.Now()
		c.Next()
		duration := time.Since(start)
		logger.Info("请求日志",
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.Int("status", c.Writer.Status()),
			zap.Duration("duration", duration),
			zap.String("client_ip", c.ClientIP()),
		)
	})

	// 注册与登录
	Register(r, db)
	Login(r, db)

	group := r.Group("/")
	group.Use(JwtAuthMiddleware)

	// ============= 文章管理 =============

	// 新增文章
	group.POST("/post", func(c *gin.Context) {
		var req CreatePostRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			Fail(c, http.StatusBadRequest, "参数绑定错误", err)
			return
		}
		userID, _ := c.Get("userID")
		post := Post{UserID: userID.(uint), Title: req.Title, Content: req.Content}
		if err := db.Create(&post).Error; err != nil {
			Fail(c, http.StatusInternalServerError, "创建文章失败", err)
			return
		}
		Success(c, gin.H{"id": post.ID})
	})

	// 删除文章
	group.DELETE("/post/:id", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			Fail(c, http.StatusBadRequest, "无效的 id 参数", err)
			return
		}

		userID, _ := c.Get("userID")
		var post Post
		if err := db.First(&post, id).Error; err != nil {
			Fail(c, http.StatusNotFound, "文章不存在", err)
			return
		}

		if post.UserID != userID {
			Fail(c, http.StatusForbidden, "只能删除自己的文章", nil)
			return
		}

		db.Delete(&post)
		Success(c, "删除成功")
	})

	// 更新文章
	group.PATCH("/post/:id", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			Fail(c, http.StatusBadRequest, "无效的 id 参数", err)
			return
		}

		var req UpdatePostRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			Fail(c, http.StatusBadRequest, "参数绑定错误", err)
			return
		}

		userID, _ := c.Get("userID")
		var post Post
		if err := db.First(&post, id).Error; err != nil {
			Fail(c, http.StatusNotFound, "文章不存在", err)
			return
		}

		if post.UserID != userID {
			Fail(c, http.StatusForbidden, "只能更新自己的文章", nil)
			return
		}

		db.Model(&post).Updates(Post{Title: req.Title, Content: req.Content})
		Success(c, "更新成功")
	})

	// 获取文章列表（支持 ?id= 查询单个）
	r.GET("/posts", func(c *gin.Context) {
		query := db.Preload("Comments")
		if id := c.Query("id"); id != "" {
			query = query.Where("id = ?", id)
		}
		var posts []Post
		if err := query.Find(&posts).Error; err != nil {
			Fail(c, http.StatusInternalServerError, "查询文章失败", err)
			return
		}
		Success(c, posts)
	})

	// ============= 评论管理 =============

	// 添加评论
	group.POST("/:id/comment", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			Fail(c, http.StatusBadRequest, "无效的 id 参数", err)
			return
		}

		var post Post
		if err := db.First(&post, id).Error; err != nil {
			Fail(c, http.StatusNotFound, "文章不存在", err)
			return
		}

		var req AddCommentRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			Fail(c, http.StatusBadRequest, "参数绑定错误", err)
			return
		}

		userID, _ := c.Get("userID")
		comment := Comment{UserID: userID.(uint), PostID: uint(id), Content: req.Content}
		if err := db.Create(&comment).Error; err != nil {
			Fail(c, http.StatusInternalServerError, "创建评论失败", err)
			return
		}
		Success(c, gin.H{"id": comment.ID})
	})

	// 获取某篇文章的评论
	r.GET("/:id/comments", func(c *gin.Context) {
		id := c.Param("id")
		var comments []Comment
		if err := db.Where("post_id = ?", id).Find(&comments).Error; err != nil {
			Fail(c, http.StatusInternalServerError, "查询评论失败", err)
			return
		}
		Success(c, comments)
	})

	r.Run(":8080")
}

// ======================== 用户注册与登录 ========================

func Register(r *gin.Engine, db *gorm.DB) {
	r.POST("/register", func(c *gin.Context) {
		var req UserRegisterRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			Fail(c, http.StatusBadRequest, "参数绑定错误", err)
			return
		}

		var exists User
		db.Where("username = ?", req.Username).First(&exists)
		if exists.ID != 0 {
			Fail(c, http.StatusConflict, "用户名已存在", nil)
			return
		}

		hash, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		user := User{Username: req.Username, Password: string(hash), Email: req.Email}
		if err := db.Create(&user).Error; err != nil {
			Fail(c, http.StatusInternalServerError, "注册失败", err)
			return
		}

		Success(c, "注册成功")
	})
}

func Login(r *gin.Engine, db *gorm.DB) {
	r.POST("/login", func(c *gin.Context) {
		var req UserLoginRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			Fail(c, http.StatusBadRequest, "参数绑定错误", err)
			return
		}

		var user User
		db.Where("username = ?", req.Username).First(&user)
		if user.ID == 0 {
			Fail(c, http.StatusNotFound, "用户不存在", nil)
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
			Fail(c, http.StatusUnauthorized, "用户名或密码错误", err)
			return
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"userID":   user.ID,
			"username": user.Username,
			"exp":      time.Now().Add(24 * time.Hour).Unix(),
		})

		tokenString, err := token.SignedString(jwtSecret)
		if err != nil {
			Fail(c, http.StatusInternalServerError, "生成 token 失败", err)
			return
		}

		Success(c, gin.H{"token": tokenString})
	})
}

// ======================== JWT 鉴权中间件 ========================

func JwtAuthMiddleware(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		Fail(c, http.StatusUnauthorized, "缺少 Authorization 请求头", nil)
		return
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		Fail(c, http.StatusUnauthorized, "请求头格式错误，应为 Bearer <token>", nil)
		return
	}

	token, err := jwt.Parse(parts[1], func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})
	if err != nil {
		Fail(c, http.StatusUnauthorized, "token 解析失败", err)
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if exp, ok := claims["exp"].(float64); ok && time.Now().Unix() > int64(exp) {
			Fail(c, http.StatusUnauthorized, "token 已过期", nil)
			return
		}
		c.Set("username", claims["username"].(string))
		c.Set("userID", uint(claims["userID"].(float64)))
		c.Next()
		return
	}

	Fail(c, http.StatusUnauthorized, "无效的 token", nil)
}

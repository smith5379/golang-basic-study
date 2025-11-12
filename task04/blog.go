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
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// 用户
type User struct {
	gorm.Model
	Username string
	Password string
	Email    string
	Posts    []Post
}

// 博客文章
type Post struct {
	gorm.Model
	UserID   uint
	Title    string
	Content  string
	Comments []Comment
}

// 博客评论
type Comment struct {
	gorm.Model
	Content string
	PostID  uint
	UserID  uint
}

type UserRegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Email    string `json:"email" binding:"required"`
}

type UserLoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// 新增文章
type CreatePostRequest struct {
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
}

// 更新文章
type UpdatePostRequest struct {
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
}

func Connect() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("blog.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	return db
}

var jwtSecret = []byte(os.Getenv("BLOG_SECRET"))

func main() {
	db := Connect()
	db.AutoMigrate(&User{}, &Post{}, &Comment{})
	r := gin.Default()

	//用户注册
	Register(r, db)

	//用户登录
	Login(r, db)

	group := r.Group("/")
	group.Use(JwtAuthMiddleware)

	//新增文章
	group.POST("/post", func(c *gin.Context) {
		var createPost CreatePostRequest
		if err := c.ShouldBindBodyWithJSON(&createPost); err != nil {
			c.JSON(400, gin.H{"message": "参数绑定错误或缺失必传参数", "error": err.Error()})
			return
		}
		userID, _ := c.Get("userID")
		db.Create(&Post{UserID: userID.(uint), Title: createPost.Title, Content: createPost.Content})
		c.JSON(200, gin.H{
			"message": "创建文章成功",
		})
	})

	//删除文章
	group.DELETE("/post/:id", func(c *gin.Context) {
		idStr := c.Param("id")
		if idStr == "" {
			c.JSON(400, gin.H{"message": "缺少必要的路径参数 id"})
			c.Abort()
			return
		}

		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(400, gin.H{"message": "无效的id参数"})
		}

		userID, _ := c.Get("userID")

		var existsPost Post
		db.First(&existsPost, id)
		if existsPost.ID == 0 {
			c.JSON(500, gin.H{
				"message": "文章不存在",
			})
			return
		}

		if existsPost.UserID != userID {
			c.JSON(500, gin.H{
				"message": "只能删除自己的文章",
			})
			return
		}

		db.Delete(&existsPost)
		c.JSON(200, gin.H{
			"message": "删除文章成功",
		})
	})

	//更新文章
	group.PATCH("/post/:id", func(c *gin.Context) {
		idStr := c.Param("id")
		if idStr == "" {
			c.JSON(400, gin.H{"message": "缺少必要的路径参数 id"})
			c.Abort()
			return
		}

		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(400, gin.H{"message": "无效的id参数"})
		}

		userID, _ := c.Get("userID")

		var existsPost Post
		db.First(&existsPost, id)
		if existsPost.ID == 0 {
			c.JSON(500, gin.H{
				"message": "文章不存在",
			})
			return
		}

		if existsPost.UserID != userID {
			c.JSON(500, gin.H{
				"message": "只能更新自己的文章",
			})
			return
		}
		var updatePost UpdatePostRequest
		if err := c.ShouldBindBodyWithJSON(&updatePost); err != nil {
			c.JSON(400, gin.H{"message": "参数绑定错误或缺失必传参数", "error": err.Error()})
			return
		}

		existsPost.Title = updatePost.Title
		existsPost.Content = updatePost.Content

		db.Model(&Post{}).Where("id = ?", id).Updates(&Post{Title: updatePost.Title, Content: updatePost.Content})
		c.JSON(200, gin.H{
			"message": "更新文章成功",
		})
	})

	//获取所有的文章列表
	r.GET("/posts", func(c *gin.Context) {

		var posts []Post
		db.Model(&Post{}).Preload("Comments").Find(&posts)
		for post := range posts {
			fmt.Println(post)
		}

		c.JSON(200, gin.H{
			"data": posts,
		})
	})

	r.Run(":8080")
}

// 注册
func Register(r *gin.Engine, db *gorm.DB) {
	r.POST("/register", func(c *gin.Context) {
		var userRegisterRequest UserRegisterRequest
		if err := c.ShouldBindBodyWithJSON(&userRegisterRequest); err != nil {
			c.JSON(400, gin.H{"message": "参数绑定错误", "error": err.Error()})
			return
		}
		//判断用户名是否存在
		var existsUser User
		db.Model(&User{}).Where("username = ?", userRegisterRequest.Username).First(&existsUser)
		if existsUser.ID != 0 {
			c.JSON(500, gin.H{
				"message": "用户名称已存在",
			})
			return
		}
		//密码加密
		hash, err := bcrypt.GenerateFromPassword([]byte(userRegisterRequest.Password), bcrypt.DefaultCost)
		if err != nil {
			panic(err)
		}

		db.Create(&User{Username: userRegisterRequest.Username, Password: string(hash), Email: userRegisterRequest.Email})
		c.JSON(200, gin.H{
			"message": "注册成功",
		})
	})
}

// 登录
func Login(r *gin.Engine, db *gorm.DB) {
	r.POST("/login", func(c *gin.Context) {
		var userLoginRequest UserLoginRequest
		if err := c.ShouldBindBodyWithJSON(&userLoginRequest); err != nil {
			c.JSON(400, gin.H{"message": "参数绑定错误", "error": err.Error()})
			return
		}
		//判断用户名是否存在
		var existsUser User
		db.Model(&User{}).Where("username = ?", userLoginRequest.Username).First(&existsUser)
		if existsUser.ID == 0 {
			c.JSON(500, gin.H{
				"message": "用户不存在",
			})
			return
		}
		//密码校验
		err := bcrypt.CompareHashAndPassword([]byte(existsUser.Password), []byte(userLoginRequest.Password))
		if err != nil {
			c.JSON(500, gin.H{
				"message": "用户与密码不匹配",
			})
			return
		}

		// 生成jwt token
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"userID":   existsUser.ID,
			"username": existsUser.Username,
			"exp":      time.Now().Add(time.Hour * 24).Unix(),
		})

		tokenString, err := token.SignedString(jwtSecret)
		if err != nil {
			c.JSON(500, gin.H{
				"message": "生成jwt token失败",
			})
			return
		}
		c.JSON(200, gin.H{
			"message": "登录成功",
			"token":   tokenString,
		})
	})
}

func JwtAuthMiddleware(c *gin.Context) {
	jwtToken := c.GetHeader("Authorization")
	if jwtToken == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "请求头缺失",
		})
		c.Abort()
		return
	}

	parts := strings.Split(jwtToken, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "请求头格式错误",
		})
		c.Abort()
		return
	}

	token, err := jwt.Parse(parts[1], func(token *jwt.Token) (interface{}, error) {
		// 校验算法是否符合预期
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "非法的token" + err.Error(),
		})
		c.Abort()
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// 校验是否过期
		if exp, ok := claims["exp"].(float64); ok {
			if time.Now().Unix() > int64(exp) {
				c.JSON(http.StatusUnauthorized, gin.H{"message": "token 过期"})
				c.Abort()
				return
			}
		}

		c.Set("username", claims["username"].(string))
		c.Set("userID", uint(claims["userID"].(float64)))
		c.Next()
		return
	}

}

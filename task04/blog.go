package main

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
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

type UserLoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

func Connect() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("blog.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	return db
}

func main() {
	db := Connect()
	db.AutoMigrate(&User{}, &Post{}, &Comment{})
	r := gin.Default()

	//用户注册
	register(r, db)

	//用户登录

	r.GET("/posts", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.Run(":8080")
}

// 注册
func register(r *gin.Engine, db *gorm.DB) {
	r.POST("/register", func(c *gin.Context) {
		var userLoginRequest UserLoginRequest
		err := c.MustBindWith(&userLoginRequest, binding.JSON)
		if err != nil {
			return
		}
		//判断用户名是否存在
		var existsUser User
		db.Model(&User{}).Where("username = ?", userLoginRequest.Username).First(&existsUser)
		if existsUser.ID != 0 {
			c.JSON(500, gin.H{
				"message": "username already exists",
			})
			return
		}
		//密码加密
		hash, err := bcrypt.GenerateFromPassword([]byte(userLoginRequest.Password), bcrypt.DefaultCost)
		if err != nil {
			panic(err)
		}

		db.Create(&User{Username: userLoginRequest.Username, Password: string(hash), Email: userLoginRequest.Email})
		c.JSON(200, gin.H{
			"message": "register success",
		})
	})
}

// 登录
func login(r *gin.Engine, db *gorm.DB) {
	r.POST("/register", func(c *gin.Context) {
		var userLoginRequest UserLoginRequest
		err := c.MustBindWith(&userLoginRequest, binding.JSON)
		if err != nil {
			return
		}
		//判断用户名是否存在
		var existsUser User
		db.Model(&User{}).Where("username = ?", userLoginRequest.Username).First(&existsUser)
		if existsUser.ID != 0 {
			c.JSON(500, gin.H{
				"message": "username already exists",
			})
			return
		}
		//密码加密
		hash, err := bcrypt.GenerateFromPassword([]byte(userLoginRequest.Password), bcrypt.DefaultCost)
		if err != nil {
			panic(err)
		}

		db.Create(&User{Username: userLoginRequest.Username, Password: string(hash), Email: userLoginRequest.Email})
		c.JSON(200, gin.H{
			"message": "register success",
		})
	})
}

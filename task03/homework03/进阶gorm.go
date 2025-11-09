package homework03

import (
	"fmt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name      string // 用户名称
	PostCount uint   // 文章数量
	Posts     []Post
}

type Post struct {
	gorm.Model
	UserID        uint
	Title         string // 文档标题
	Content       string // 文档标题
	Comments      []Comment
	CommentStatus string //文章评论状态
}

func (p *Post) AfterCreate(tx *gorm.DB) (err error) {
	err = tx.Model(&User{}).Where("id = ?", p.UserID).
		UpdateColumn("post_count", gorm.Expr("post_count + ?", 1)).Error
	return
}

type Comment struct {
	gorm.Model
	PostID  uint
	UserID  uint
	Content string // 评论内容
}

func (c *Comment) AfterDelete(tx *gorm.DB) (err error) {
	fmt.Println("c: ", c)
	var count int64
	tx.Model(&Comment{}).Where("post_id = ?", c.PostID).
		Count(&count)
	if count == 0 {
		tx.Model(&Post{}).Where("id = ?", c.PostID).UpdateColumn("comment_status", "无评论")
	}

	return err
}

func Run(db *gorm.DB) {

	db.AutoMigrate(&User{}, &Post{}, &Comment{})

	//db.Create(&User{
	//	Name: "zhangsan",
	//	Posts: []Post{
	//		{Title: "zhangsan的第一篇博客", Content: "zhangsan的第一篇博客的内容", Comments: []Comment{{UserID: 3, Content: "zhangsan的第一篇博客的评论1"}}},
	//		{Title: "zhangsan的第二篇博客", Content: "zhangsan的第二篇博客的内容", Comments: []Comment{{UserID: 4, Content: "zhangsan的第二篇博客的评论1"}, {UserID: 5, Content: "zhangsan的第二篇博客的评论2"}}},
	//	},
	//})

	// 使用Gorm查询某个用户发布的所有文章及其对应的评论信息。
	//var user User
	//db.Model(&User{}).Preload("Posts.Comments").First(&user, 1)
	//fmt.Printf("%+v\n", user)
	//for _, post := range user.Posts {
	//	fmt.Printf("文章: %s\n", post.Content)
	//	for _, comment := range post.Comments {
	//		fmt.Printf("  评论: %s\n", comment.Content)
	//	}
	//}

	// 使用Gorm查询评论数量最多的文章信息。
	//var post Post
	//db.Model(&Post{}).
	//	Joins("left join comments on posts.id = comments.post_id").
	//	Select("posts.id, count(comments.id) as comment_count").
	//	Group("posts.id").Order("comment_count desc").Limit(1).Scan(&post)
	//
	//db.Preload("Comments").First(&post, post.ID)
	//
	//fmt.Printf("评论最多的文章：%s \n", post.Title)
	//for _, com := range post.Comments {
	//	fmt.Printf("评论：%s \n", com.Content)
	//}

	// 删除第一篇文章的评论
	//var p Post
	//db.First(&p, 1)
	//fmt.Printf("before, title: %s ,id:%d, state:%s \n ", p.Title, p.ID, p.CommentStatus)
	//
	////删除评论, 测试钩子函数
	//var comments []Comment
	//db.Where("post_id = ?", 1).Find(&comments)
	//fmt.Println("del:", comments)
	//for _, comment := range comments {
	//	db.Delete(&comment)
	//}
	//
	//db.First(&p, 1)
	//fmt.Printf("after, title: %s ,id:%d, state:%s \n ", p.Title, p.ID, p.CommentStatus)

}

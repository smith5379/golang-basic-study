package main

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type Book struct {
	ID     uint64
	Title  string
	Author string
	Price  float64
}

func main() {

	db, err := sqlx.Connect("sqlite3", "test.db")
	if err != nil {
		fmt.Printf("连接数据库失败: %v", err)
	}
	defer db.Close()

	createSql := `
			create table if not exists books (
	        id          integer primary key autoincrement,
	        title       varchar(100),
	    	author      varchar(100),
	    	price       real )`

	_, err = db.Exec(createSql)
	if err != nil {
		fmt.Printf("无法创建表: %v", err)
	}

	books := []Book{
		{Title: "Alice", Author: "张三", Price: 20},
		{Title: "Bob", Author: "李四", Price: 30},
		{Title: "Charlie", Author: "王五", Price: 55},
		{Title: "David", Author: "钱六", Price: 60},
	}
	for _, e := range books {
		db.Exec(`INSERT INTO books (title, author, price) VALUES (?, ?, ?)`, e.Title, e.Author, e.Price)
	}

	// 查询价格大于 50 元的书籍，并将结果映射到 Book 结构体切片中，确保类型安全。
	var moreThan50Books []Book
	db.Select(&moreThan50Books, "select id, title, author, price from books where price > ?", 50)

	fmt.Println("大于50元的书籍列表：")
	for _, e := range moreThan50Books {
		fmt.Printf("ID: %d, Title: %s, Author: %s, Price: %v\n", e.ID, e.Title, e.Author, e.Price)
	}
}

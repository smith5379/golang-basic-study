package main

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type Employee struct {
	ID         uint64
	Name       string
	Department string
	Salary     float64
}

func main() {

	db, err := sqlx.Connect("sqlite3", "test.db")
	if err != nil {
		fmt.Printf("连接数据库失败: %v", err)
	}
	defer db.Close()

	createSql := `
			create table if not exists employees (
	        id         integer primary key autoincrement,
	        name       varchar(100),
	    	department varchar(100),
	    	salary     real )`

	_, err = db.Exec(createSql)
	if err != nil {
		fmt.Printf("无法创建表: %v", err)
	}

	employees := []Employee{
		{Name: "Alice", Department: "技术部", Salary: 8000},
		{Name: "Bob", Department: "人事部", Salary: 5000},
		{Name: "Charlie", Department: "技术部", Salary: 9000},
		{Name: "David", Department: "财务部", Salary: 6000},
	}
	for _, e := range employees {
		db.Exec(`INSERT INTO employees (name, department, salary) VALUES (?, ?, ?)`, e.Name, e.Department, e.Salary)
	}

	// 查询部门为 "技术部" 的员工
	var techEmployees []Employee
	db.Select(&techEmployees, "select id, name, department, salary from employees where department = ?", "技术部")

	fmt.Println("技术部员工列表：")
	for _, e := range techEmployees {
		fmt.Printf("ID: %d, Name: %s, Department: %s, Salary: %v\n", e.ID, e.Name, e.Department, e.Salary)
	}

	//编写Go代码，使用Sqlx查询 employees 表中工资最高的员工信息，并将结果映射到一个 Employee 结构体中。
	var maxSalaryEmployee Employee
	db.Get(&maxSalaryEmployee, "select id, name, department, salary from employees order by salary desc limit 1")
	fmt.Printf("工资最高的员工信息： ID: %d, Name: %s, Department: %s, Salary: %v\n", maxSalaryEmployee.ID, maxSalaryEmployee.Name, maxSalaryEmployee.Department, maxSalaryEmployee.Salary)

}

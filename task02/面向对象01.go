package main

import (
	"fmt"
	"math"
)

type Shape interface {
	Area() float64
	Perimeter() float64
}

// 长方形
type Rectangle struct {
	width, height float64
}

func (r Rectangle) Area() float64 {
	return r.width * r.height
}
func (r Rectangle) Perimeter() float64 {
	return (r.width + r.height) * 2
}

// 圆形
type Circle struct {
	radius float64
}

func (c Circle) Area() float64 {
	return math.Pi * c.radius * c.radius
}
func (c Circle) Perimeter() float64 {
	return 2 * math.Pi * c.radius
}

/*
*
题目 ：定义一个 Shape 接口，包含 Area() 和 Perimeter() 两个方法。然后创建 Rectangle 和 Circle 结构体，实现 Shape 接口。在主函数中，创建这两个结构体的实例，并调用它们的 Area() 和 Perimeter() 方法。
考察点 ：接口的定义与实现、面向对象编程风格。
*/
func main() {

	var r Shape = Rectangle{width: 5, height: 10}
	fmt.Println("长方形的面积: ", r.Area())
	fmt.Println("长方形的周长: ", r.Perimeter())

	var c Shape = Circle{radius: 5}
	fmt.Println("圆的面积: ", c.Area())
	fmt.Println("圆的周长: ", c.Perimeter())

}

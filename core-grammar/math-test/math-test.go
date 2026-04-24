package math_test

import (
	"fmt"
	"math"
)

func Test1() {
	//最大值
	fmt.Println(math.Max(1, 2))
	//最小值
	fmt.Println(math.Min(1, 2))
	//绝对值
	fmt.Println(math.Abs(-1))
	//余数
	fmt.Println(math.Mod(12, 10))
	//Nan 检测
	fmt.Println(math.IsNaN(math.NaN()))
	//Inf 检测
	fmt.Println(math.IsInf(1.0, 1))
	fmt.Println(math.IsInf(math.Inf(-1), -1))
	//取整
	fmt.Println(math.Trunc(1.26))
	fmt.Println(math.Trunc(-1.26))
	fmt.Println(math.Trunc(2.333))
	//向下取整
	fmt.Println(math.Floor(2.5))
	//向上取整
	fmt.Println(math.Ceil(2.5))
	//四舍五入
	fmt.Println(math.Round(2.2389))
	fmt.Println(math.Round(-2.2389))
	//求对数
	fmt.Println(math.Log(100) / math.Log(10))
	fmt.Println(math.Log(1) / math.Log(2))
	//E 的指数
	fmt.Println(math.Exp(2))
	//幂
	fmt.Println(math.Pow(2, 3))
	//平方根
	fmt.Println(math.Sqrt(4))
	//立方根
	fmt.Println(math.Cbrt(8))
	//开N方
	fmt.Println(math.Round(math.Pow(8, 1.0/3)))
	fmt.Println(math.Round(math.Pow(100, 1.0/2)))
	//sin
	fmt.Println(math.Sin(0))
	fmt.Println(math.Sin(20))
	//cos
	fmt.Println(math.Cos(0))
	fmt.Println(math.Cos(20))
	//tan
	fmt.Println(math.Tan(0))
	fmt.Println(math.Tan(20))
}

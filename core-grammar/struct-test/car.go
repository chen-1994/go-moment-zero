package struct_test

import "fmt"

type Car struct {
	id    int
	name  string
	price float64
	color string
	brand string
}

type Info struct {
	car Car
}

func Init() {
	u1 := Car{1, "卡罗拉", 22, "", "丰田"}
	fmt.Println(u1)
	u2 := Car{
		id:   2,
		name: "3系",
	}
	fmt.Println(u2)
	u2.name = "4系"
	fmt.Println(u2)
}

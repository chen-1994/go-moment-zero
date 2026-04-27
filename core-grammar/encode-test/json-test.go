package encode_test

import (
	"encoding/json"
	"fmt"
)

type User struct {
	Name string
	Age  string
	Sex  string
}

type UserNew struct { //字段重命名
	Name string `json:"n"`
	Age  string `json:"a"`
	Sex  string `json:"s"`
}

func JsonTest1() {
	u := User{
		Name: "张三",
		Age:  "18",
		Sex:  "男",
	}
	//bytes, err := json.Marshal(u) //json序列化
	bytes, err := json.MarshalIndent(u, "", "\t") //缩进
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(u)
	fmt.Println(string(bytes))

	//反序列化
	u1 := User{}
	//err1 := json.Unmarshal(bytes, &u1)
	jsonStr := "{\"Name\":\"张三1\", \"Age\":\"18\", \"Sex\":\"男\"}"
	err1 := json.Unmarshal([]byte(jsonStr), &u1)
	if err1 != nil {
		fmt.Println(err1)
		return
	}
	fmt.Printf("%+v", u1)
}

func JsonTest2() {
	u := UserNew{
		Name: "张三",
		Age:  "18",
		Sex:  "男",
	}
	bytes, err := json.Marshal(u)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(u)
	fmt.Println(string(bytes))
}

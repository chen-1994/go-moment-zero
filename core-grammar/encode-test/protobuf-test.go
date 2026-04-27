package encode_test

import (
	p "core-grammar/encode-test/proto"
	"fmt"

	"github.com/golang/protobuf/proto"
)

type Person struct {
	UserId   string `json:"id"`
	Username string `json:"name"`
	Age      int    `json:"age"`
	Address  string `json:"address"`
}

func ProtobufTest1() {
	person := p.Person{
		Name:   "wyh",
		Age:    12,
		Gender: p.Gender_FE_MAIL,
	}

	data, err := proto.Marshal(&person) //序列化
	if err != nil {
		fmt.Println(err)
		return
	}
	temp := &p.Person{}
	fmt.Println("proto buffer len: ", len(data), "bytes:", data)
	err = proto.Unmarshal(data, temp) //反序列化
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(temp)
}

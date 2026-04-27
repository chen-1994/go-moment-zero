package encode_test

import (
	"encoding/xml"
	"fmt"
	"os"
)

type UserInfo struct {
	UserId   int    `xml:"id"`
	UserName string `xml:"name"`
	Age      int    `xml:"age"`
	Address  string `xml:"address"`
}

func XmlTest1() {
	p := UserInfo{
		UserId:   1,
		UserName: "zs",
		Age:      11,
		Address:  "hx",
	}
	//bytes, err := xml.Marshal(p)
	bytes, err := xml.MarshalIndent(p, "", "\t") //格式化缩进
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(bytes))
}

func XmlTest2() {
	bytes, err := os.ReadFile("./encode-test/test.xml")
	if err != nil {
		fmt.Println(err)
		return
	}
	var p UserInfo
	err = xml.Unmarshal(bytes, &p)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(p)
}

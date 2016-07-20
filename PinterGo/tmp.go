package main

import (
	"encoding/json"
	"fmt"
)


//
type Board struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	URL  string `json:"url"`
}

type BoardsData struct {
	Data []Board `json:"data"`
}

func testJSON() {
	pindata := BoardsData{}
	jsonData := `{"data": [{"url": "https://www.pinterest.com/popvac/lab/", "id": "557461328814517506", "name": "lab"}, {"url": "https://www.pinterest.com/popvac/candy/", "id": "557461328814659076", "name": "candy"}]}`
	_ = json.Unmarshal([]byte(jsonData), &pindata)
	fmt.Println(pindata.Data[0])
	fmt.Println(pindata.Data[0].ID)
	fmt.Println(pindata.Data[1])
	fmt.Println(pindata.Data[1].ID)

	/*
	str := `{"page": 1, "fruits": ["apple", "peach"]}`
	res := Response2{}
	json.Unmarshal([]byte(str), &res)
	fmt.Println(res)
	fmt.Println(res.Fruits[0])
	*/
}

//
type Stringer interface {
	String() string
}
func testToString(){
	ToString(new (Stringer))
}

func ToString(any interface{})  {
	var t string
	if v, ok := any.(Stringer); ok {
		t = v.String()
	}
	print(t)
}
//


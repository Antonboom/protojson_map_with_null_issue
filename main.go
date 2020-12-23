package main

import (
	"fmt"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/structpb"
)

//go:generate protoc --go_out=. example.proto

var json = protojson.MarshalOptions{
	Multiline:     true,
	Indent:        "  ",
	AllowPartial:  true,
	UseProtoNames: true,
}

func main() {
	user := &User{
		Id:    "user-2",
		Email: "user2@gmail.com",
	}

	// 1. ResponseSimple
	/*
		{
		  "users":  {
			"user-1":  {},  // <- null expected!
			"user-2":  {
			  "id":  "user-2",
			  "email":  "user2@gmail.com"
			}
		  }
		}
	*/
	data, err := json.Marshal(&ResponseSimple{
		Users: map[string]*User{
			"user-1": nil,
			"user-2": user,
		},
	})
	failOnErr(err)
	fmt.Println(string(data))

	// 2. ResponseWithCustomNullable
	/*
		{
		  "users":  {
			"user-1":  {          // <- extra field!
			  "null":  null
			},
			"user-2":  {
			  "user":  {          // <- extra field!
				"id":  "user-2",
				"email":  "user2@gmail.com"
			  }
			}
		  }
		}
	*/
	data, err = json.Marshal(&ResponseWithCustomNullable{
		Users: map[string]*ResponseWithCustomNullable_NullableUser{
			"user-1": {
				Kind: &ResponseWithCustomNullable_NullableUser_Null{Null: structpb.NullValue_NULL_VALUE},
			},
			"user-2": {
				Kind: &ResponseWithCustomNullable_NullableUser_User{User: user},
			},
		},
	})
	failOnErr(err)
	fmt.Println(string(data))

	// 3. ResponseWithValue
	/*
		{
		  "users":  {
			"user-1":  null,
			"user-2":  {
			  "email":  "user2@gmail.com",
			  "id":  "user-2"
			}
		  }
		}
	*/
	uData, err := protojson.Marshal(user)
	failOnErr(err)

	userStruct := new(structpb.Struct)
	failOnErr(protojson.Unmarshal(uData, userStruct))

	data, err = json.Marshal(&ResponseWithValue{
		Users: map[string]*structpb.Value{
			"user-1": structpb.NewNullValue(),
			"user-2": structpb.NewStructValue(userStruct),
		},
	})
	failOnErr(err)
	fmt.Println(string(data))
}

func failOnErr(err error) {
	if err != nil {
		panic(err)
	}
}

package event

import (
	"encoding/json"
	"fmt"

	"github.com/Kchanit/microservice-payment-golang/internal/core/utils"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/omise/omise-go"
	"github.com/omise/omise-go/operations"
)

func Listen(message *kafka.Message) {
	// fmt.Printf("Message on %T: %s\n", *message.TopicPartition.Topic, string(message.Value))

	switch *message.TopicPartition.Topic {
	case "user_create":
		/*
				{
			   "id":"f75a53cd-3e6f-4ed6-aa96-ce3d4308e03c",
			   "last_login":null,
			   "is_superuser":false,
			   "is_staff":false,
			   "is_active":true,
			   "date_joined":"2024-01-29T21:01:07.830650+07:00",
			   "first_name":"Ratchaphon",
			   "username":"Ratchaphon1412",
			   "last_name":"Hinsui",
			   "email":"sevenknight5570@gmail.com",
			   "phone_number":"0818886119",
			   "gender":"Male",
			   "country":"Thailand",
			   "profile":"https://avatar.oxro.io/avatar.svg?name=Ratchaphon+Hinsui&background=ffd60a&color=fff&length=2",
			   "customer_omise_id":null,
			   "is_email_verified":false,
			   "groups":[
			      1
			   ],
			   "user_permissions":[
			      52,
			      54,
			      53,
			      50,
			      22,
			      55,
			      51,
			      23,
			      49,
			      56,
			      24
			   ]
			}

		*/
		fmt.Printf("User created %s\n", string(message.Value))
		// change message.Value to map[string]interface{}
		userjson := map[string]interface{}{}
		json.Unmarshal(message.Value, &userjson)
		// use facade to create customer in omise
		fmt.Printf("email is %s\n", userjson["email"])
		//call services omise for create Customer

		facade := utils.FacadeSingleton()
		client := facade.Omise
		customer := &omise.Customer{}

		err := client.Do(customer, &operations.CreateCustomer{
			Email:       userjson["email"].(string),
			Description: userjson["username"].(string),
		})
		if err != nil {
			fmt.Println("Error create customer omise")
		}
		fmt.Println(customer)
		/*
			&{{customer cust_test_5ylhp900ypl37ahckt7 false 0xc0002bc230 0001-01-01 00:00:00 +0000 UTC}  sevenknight5570@gmail.com Ratchaphon1412 0xc000270000 map[]}
		*/

		// create customer in database

		/* code here */

		// send kafka to user_customer
		err = facade.SendKafka("user_customer", map[string]interface{}{
			"email": userjson["email"],
			"token": customer.ID,
		})
		if err != nil {
			fmt.Println("Error send kafka customer Token")
		}

	default:

		fmt.Println(" Topic not found")

	}
}

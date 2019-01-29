package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"strings"
	"time"

	prompt "github.com/c-bata/go-prompt"
	"github.com/streadway/amqp"
)

var key = "doitdoitdoitdoitdoitdoitdoitdoit"

type ReqRes struct {
	Req []string
	Res string
}

var fields map[string]ReqRes

func init() {
	fields = make(map[string]ReqRes)

	fields["retry"] = ReqRes{
		Req: []string{
			"route",
			"data",
		},
		Res: "interface",
	}

	fields["ping"] = ReqRes{
		Req: []string{
			"ping(string)",
		},
		Res: "pong(string)",
	}

	fields["pong"] = ReqRes{
		Req: []string{
			"ping(string)",
		},
		Res: "errors",
	}

	fields["rpc_queue"] = ReqRes{
		Req: []string{
			"number(integer)",
		},
		Res: "fibbonaci(int)",
	}
	fields["user_register"] = ReqRes{
		Req: []string{
			"string fullname",
			"string email",
			"string password",
			"string account_type",
			"string shop_name",
			"string username",
			"string phone",
			"string dob",
			"int    school_id",
			"string invitation_code",
		},
		Res: "errors[],user_id",
	}

	fields["user_invite_familly"] = ReqRes{
		Req: []string{
			"int user_id",
		},
		Res: "errors[],URL",
	}

	fields["user_show_invitation_list"] = ReqRes{
		Req: []string{
			"int user_id",
		},
		Res: "errors[],[{user_id,invited,confirm,code,invited_status,invitation_link}]",
	}

	fields["user_default_info"] = ReqRes{
		Req: []string{
			"int user_id",
		},
		Res: "errors[],[{user_name,firstname,lastname,alias,dob,phone,gender,status,last_login}]",
	}

	fields["user_familly_member"] = ReqRes{
		Req: []string{
			"int user_id",
		},
		Res: "errors[],[{user_id,firstname,lastname,gender,dob,phone,modified,modified_by}]",
	}

	fields["create_product"] = ReqRes{
		Req: []string{
			"int64                 merchant_id",
			"int64                 school_id",
			"string                product_name",
			"string                product_detail",
			"float64               product_price",
			"int64                 product_status",
			"int                   product_calories",
			"[]ProductsNutrition   product_nutrition [NutritionPercent (int64), NutritionWeightMeasure (string), NutritionWeight (int64),NutritionName (string)]",
			"[]ProductsTopping     product_topping [ToppingName (string), ToppingPrice (float64)]",
		},

		Res: "",
	}
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func randomString(l int) string {
	bytes := make([]byte, l)
	for i := 0; i < l; i++ {
		bytes[i] = byte(randInt(65, 90))
	}
	return string(bytes)
}

func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}

func completer(d prompt.Document) []prompt.Suggest {
	s := make([]prompt.Suggest, 0)

	for key, val := range fields {
		s = append(s, prompt.Suggest{
			Text:        key,
			Description: strings.Join(val.Req, ","),
		})
	}
	return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
}

func completerNil(d prompt.Document) []prompt.Suggest {
	s := []prompt.Suggest{}
	return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
}

func send(qName string, body string) (res string, err error) {
	conn, err := amqp.Dial("amqp://root:root@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when usused
		true,  // exclusive
		false, // noWait
		nil,   // arguments
	)
	failOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	corrId := randomString(32)

	bodyEncrypt, _ := EncryptCBC([]byte(key), []byte(body))

	fmt.Println(string(bodyEncrypt))

	err = ch.Publish(
		"",    // exchange
		qName, // routing key
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType:   "text/plain",
			CorrelationId: corrId,
			ReplyTo:       q.Name,
			Body:          []byte(bodyEncrypt),
		})
	failOnError(err, "Failed to publish a message")

	for d := range msgs {
		if corrId == d.CorrelationId {

			log.Println("Response part")
			fmt.Println(string(d.Body))

			plaintext, err := DecryptCBC([]byte(key), d.Body)

			if err != nil {
				log.Println("err")
			}

			res = string(plaintext)
			fmt.Println(res)
			failOnError(err, "Failed to convert body to integer")
			break
		}
	}

	return
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	for {
		fmt.Println("Select queue name :")
		que := prompt.Input(">", completer)
		if que == "exit" {
			return
		}

		que = strings.TrimSpace(que)

		reqRes := fields[que]

		var bodyArrStr []string
		var bodyStr string

		for _, val := range reqRes.Req {
			fmt.Println("Fill the body ", val, " : ")
			body := prompt.Input(">", completerNil)
			if body == "exit" {
				return
			}

			bodyArrStr = append(bodyArrStr, body)
		}
		bodyStr = strings.Join(bodyArrStr, ",")
		log.Printf("Path : %v", que)

		log.Printf("Request : %v,%v", que, bodyStr)

		if que == "retry" {
			que = bodyArrStr[0]
			bodyStr = bodyArrStr[1]
		}

		res, err := send(que, bodyStr)

		log.Printf("Result : %v", res)

		str := fmt.Sprintf(`
## %v
- FieldsRequest : %v
- FieldsResponse : %v

### Sample : 
- Request : %v,%v
- Result: %v`, que, strings.Join(reqRes.Req, ","), reqRes.Res, que, bodyStr, res)

		err = ioutil.WriteFile(fmt.Sprintf("%v.md", que), []byte(str), 0644)

		failOnError(err, "Failed to handle RPC request")
	}

}

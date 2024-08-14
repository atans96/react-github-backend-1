package routes

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/tidwall/gjson"
	"go.mongodb.org/mongo-driver/bson"
	Mongo "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"server_go/src/service"
	"server_go/src/service/validation"
	"sigs.k8s.io/yaml"
	"strconv"
	"time"
)

type registrationData struct {
	TokenType string `mapper:"token_type" json:"token_type"`
	Token     string `mapper:"access_token" json:"access_token"`
	Scope     string `mapper:"scope" json:"scope"`
}

func shouldBeStruct(d reflect.Type) error {
	td := d.Elem()
	if td.Kind() != reflect.Struct {
		errStr := fmt.Sprintf("Input should be %v, found %v", reflect.Struct, td.Kind())
		return errors.New(errStr)
	}
	return nil
}
func mapToStruct(url url.Values) error {
	newUser := registrationData{}
	dType := reflect.TypeOf(newUser)

	if err := shouldBeStruct(dType); err != nil {
		return err
	}
	// Data Holder Value
	dhVal := reflect.ValueOf(newUser)

	// Loop over all the fields present in struct (Title, Body, JSON)
	for i := 0; i < dType.Elem().NumField(); i++ {

		// Give me ith field. Elem() is used to dereference the pointer
		field := dType.Elem().Field(i)

		// Get the value from field tag i.e in case of Title it is "title"
		key := field.Tag.Get("mapper")

		// Get the type of field
		kind := field.Type.Kind()

		// Get the value from query params with given key
		val := url.Get(key)
		//  Get reference of field value provided to input `d`
		result := dhVal.Elem().Field(i)

		// we only check for string for now so,
		if kind == reflect.String {
			// set the value to string field
			// for other kinds we need to use different functions like; SetInt, Set etc
			result.SetString(val)
		} else {
			return errors.New("only supports string")
		}

	}
	return nil
}

func postReq(tokenString chan map[string]interface{}, c *fiber.Ctx) {
	reg := new(validation.Registration)
	if err := c.BodyParser(reg); err != nil {
		panic(err)
	}
	errs := service.Validate(reg)
	if errs != nil {
		panic(errs)
	}
	out, err := json.Marshal(reg)
	if err != nil {
		panic(err)
	}
	bod := bytes.NewBuffer(out)
	makeNew, err := http.NewRequest("POST", "https://github.com/login/oauth/access_token", bod)
	if err != nil {
		panic(err)
	}
	makeNew.Header.Add("Content-Type", "application/json")

	request := service.MakeRequest{Req: makeNew}
	request.Request()
	sb := "https://example.com?" + string(request.Result)
	request.Clear()
	u, err := url.Parse(sb)
	if err != nil {
		panic(err)
	}
	q, err := url.ParseQuery(u.RawQuery)
	token := fiber.Map{"access_token": q.Get("access_token"), "token_type": q.Get("token_type")}
	tokenString <- token
	if err != nil {
		panic(err)
	}
}
func postToServer(rid string, sess_id string) {
	b := bytes.NewBuffer([]byte(rid))
	makeNew, err := http.NewRequest("POST", os.Getenv("UWEBSOCKET_HOST")+":"+os.Getenv("UWEBSOCKET_PORT")+"/server_uwebsocket/authenticate", b)
	if err != nil {
		panic(err)
	}
	makeNew.Header.Add("Origin", "https://"+os.Getenv("GOLANG_HOST")+":"+os.Getenv("GOLANG_PORT"))
	makeNew.Header.Add("Cookie", "session_id="+sess_id)
	makeNew.Header.Add("Content-Type", "application/json")
	request := service.MakeRequest{Req: makeNew}
	request.Request()
	request.Clear()
}
func setSession(c *fiber.Ctx, rid string) {
	store, err := service.SessionStore.DB.Get(c)
	if err != nil {
		return
	}
	if err != nil {
		err := c.Status(fiber.StatusForbidden).SendString("Something wrong with session")
		if err != nil {
			panic(err)
		}
	}
	newUser := fmt.Sprintf(`{"token": "%v","username": "%v"}`, gjson.Get(rid, "token").String(), gjson.Get(rid, "data.login").String())
	postToServer(rid, store.ID())
	store.Set(store.ID(), newUser)
	hour, err := strconv.ParseInt(os.Getenv("JWT_EXPIRE"), 10, 64)
	if err != nil {
		panic(err)
	}
	store.SetExpiry(time.Hour * time.Duration(hour))

	// Save session
	if err := store.Save(); err != nil {
		panic(err)
	}
}
func getReq(output chan string, subOutput chan string, tokenString chan map[string]interface{}) {
	tok := <-tokenString
	makeNew, err := http.NewRequest("GET", "https://api.github.com/user", nil)
	if err != nil {
		panic(err)
	}
	makeNew.Header.Add("Authorization", tok["token_type"].(string)+" "+tok["access_token"].(string))
	makeNew.Header.Add("User-Agent", "Chrome: 91")

	request := service.MakeRequest{Req: makeNew}
	request.Request()
	rid := fmt.Sprintf(`{"token": "%v","token_type": "%v", "data": %v}`, tok["access_token"].(string), tok["token_type"].(string), string(request.Result))
	output <- rid
	subOutput <- rid
}
func updateDB(subOutputChan chan string) {
	rid := <-subOutputChan
	if len(gjson.Get(rid, "data.message").String()) == 0 {
		txt, _ := ioutil.ReadFile("files/languages.yml")
		y, err := yaml.YAMLToJSON(txt)
		if err != nil {
			panic(err)
		}
		var objmap map[string]interface{}
		if err := json.Unmarshal(y, &objmap); err != nil {
			panic(err)
		}
		var result []map[string]interface{}
		for key := range objmap {
			res := map[string]interface{}{
				"language": key,
				"checked":  true,
			}
			result = append(result, res)
		}
		opts := options.FindOneAndUpdate().SetUpsert(true)
		query := bson.M{"userName": gjson.Get(rid, "data.login").String()}
		update := bson.D{
			{Key: "$set",
				Value: bson.D{
					{Key: "avatar", Value: gjson.Get(rid, "data.avatar_url").String()},
					{Key: "languagePreference", Value: result},
					{Key: "joinDate", Value: time.Now()},
					{Key: "token", Value: gjson.Get(rid, "token").String()},
				},
			},
		}
		err = service.Mongo.DB.FindOneAndUpdate(nil, query, update, opts).Err()
		if err != nil {
			// ErrNoDocuments means that the filter did not match any documents in the collection
			if err == Mongo.ErrNoDocuments {
				fmt.Println(err)
			}
			fmt.Println(err)
		}
	}
}
func worker(c *fiber.Ctx, outputChan chan string) {
	subOutputChan := make(chan string)
	tokenStringChan := make(chan map[string]interface{})
	go postReq(tokenStringChan, c)
	go getReq(outputChan, subOutputChan, tokenStringChan)
	go updateDB(subOutputChan)
}
func Register(c *fiber.Ctx) error {
	outputChan := make(chan string)
	go worker(c, outputChan)
	select {
	case msg2 := <-outputChan:
		setSession(c, msg2)
		return c.SendString(msg2)
	}
}

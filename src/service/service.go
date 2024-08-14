package service

import (
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"strings"
)

var (
	ErrSourceNotArray = errors.New("Source value is not an array")
	ErrReducerNil     = errors.New("Reducer function cannot be nil")
	ErrReducerNotFunc = errors.New("Reducer argument must be a function")
)

//Reduce an array of something into another thing
func Reduce(source, initialValue, reducer interface{}) (interface{}, error) {
	srcV := reflect.ValueOf(source)
	kind := srcV.Kind()
	if kind != reflect.Slice && kind != reflect.Array {
		return nil, ErrSourceNotArray
	}

	if reducer == nil {
		return nil, ErrReducerNil
	}

	rv := reflect.ValueOf(reducer)
	if rv.Kind() != reflect.Func {
		return nil, ErrReducerNotFunc
	}

	// copy initial value as accumulator, and get the reflection value
	accumulator := initialValue
	accV := reflect.ValueOf(accumulator)
	for i := 0; i < srcV.Len(); i++ {
		entry := srcV.Index(i)

		// call reducer via reflection
		reduceResults := rv.Call([]reflect.Value{
			accV,               // send accumulator value
			entry,              // send current source entry
			reflect.ValueOf(i), // send current loop index
		})

		accV = reduceResults[0]
	}

	return accV.Interface(), nil
}

func GetLocalIP() string {
	host, _ := os.Hostname()
	addrs, _ := net.LookupIP(host)
	for _, addr := range addrs {
		if ipv4 := addr.To4(); ipv4 != nil {
			fmt.Println("IPv4: ", ipv4)
			return addr.String()
		}
	}
	return ""
}
func FilePathWalkDir(root string) ([]string, error) {
	var files []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return files, err
}
func ReadFiles() {
	files, err := FilePathWalkDir("./src/https-server")
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		if strings.Contains(file, ".docx") {
			continue
		}
		b, err := ioutil.ReadFile(file)
		if err != nil {
			fmt.Print(err)
		}
		value := string(b)
		key := strings.Split(file, "\\")
		os.Setenv(key[len(key)-1], value)
	}
}
func Validate(obj interface{}) []*ErrorResponse {
	if Validator == nil {
		return nil
	}
	return Validator.ValidateStruct(obj)
}
func Contains(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}

type MakeRequest struct {
	Req         *http.Request
	Result      []byte
	ContentType string
}

func (m *MakeRequest) Request() {
	client := http.Client{}
	resp, err := client.Do(m.Req)
	if resp == nil {
		panic("Resp is Nil")
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic(err)
		}
	}(resp.Body)
	result, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	m.Result = result
	m.ContentType = resp.Header.Get("Content-Type")
}
func (m *MakeRequest) Clear() {
	p := reflect.ValueOf(m).Elem()
	p.Set(reflect.Zero(p.Type()))
}

type Server struct {
	*fiber.App
}

func (app *Server) StartServer(addr string) {
	// Create tls certificate
	cert, err := tls.X509KeyPair([]byte(os.Getenv("localhost.crt")), []byte(os.Getenv("localhost.decrypted.key")))
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
	}
	// Create custom listener
	ln, err := tls.Listen("tcp", addr, tlsConfig)
	if err != nil {
		panic(err)
	}
	log.Fatal(app.Listener(ln))
}

// export API_USER=admin
// export API_PASS=
// export API_ADDRESS=xxxx.xxxx.us-east-2.rds.amazonaws.com
// export API_SALT=api12345
package main

import (
	"crypto/aes"
	"crypto/cipher"
	crt "crypto/rand"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/nsf/termbox-go"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var (
	dbUser, dbPass, dbAddress, dbSalt string
	DBMS                              *gorm.DB
	debug                             bool
)

type member struct {
	Id       int    `gorm:"primary_key"`
	Name     string `json:name`
	Password string `json:password`
	gorm.Model
}

type responseData struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type responseList struct {
	Status  string `json:"status"`
	Members []member
}

func main() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()
	termbox.Flush()

	if len(os.Getenv("API_USER")) == 0 || len(os.Getenv("API_PASS")) == 0 || len(os.Getenv("API_ADDRESS")) == 0 || len(os.Getenv("API_SALT")) == 0 {
		fmt.Println("environment not found!")
		os.Exit(-1)
	}

	_cert := flag.String("cert", "localhost.pem", "[-cert=ssl_certificate file path (if you don't use https, haven't to use this option)]")
	_key := flag.String("key", "localhost-key.pem", "[-key=ssl_certificate_key file path (if you don't use https, haven't to use this option)]")
	_port := flag.String("port", "8080", "[-port=port number]")
	_debug := flag.Bool("debug", false, "[-debug=debug mode]")
	flag.Parse()

	debug = bool(*_debug)

	dbUser = os.Getenv("API_USER")
	dbPass = os.Getenv("API_PASS")
	dbAddress = os.Getenv("API_ADDRESS")
	dbSalt = os.Getenv("API_SALT")

	if debug == true {
		fmt.Println("API_USER: ", dbUser)
		fmt.Println("API_PASS: ", dbPass)
		fmt.Println("API_ADDRESS: ", dbAddress)
		fmt.Println("API_SALT: ", dbSalt)
	}

	fmt.Println("db conecting..")

	DBMS = GormConnect()

	http.HandleFunc("/create", ApiHandlers)
	http.HandleFunc("/raw", ApiHandlers)
	http.HandleFunc("/find", ApiHandlers)
	http.HandleFunc("/read", ApiHandlers)
	http.HandleFunc("/update", ApiHandlers)
	http.HandleFunc("/delete", ApiHandlers)
	http.HandleFunc("/auth", ApiHandlers)

	fmt.Println("api server starting..")

	go func() {
		err := http.ListenAndServeTLS(":"+string(*_port), string(*_cert), string(*_key), nil)
		if err != nil {
			log.Fatal("ListenAndServeTLS: ", err)
		}
	}()

	termbox.SetInputMode(termbox.InputEsc)

	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case 27: //Escape
				termbox.Flush()
				DBMS.Close()
				os.Exit(0)
			}
		}
	}
}

func ApiHandlers(w http.ResponseWriter, req *http.Request) {
	errorFlag := false
	respData := responseData{}
	respList := responseList{}

	if req.Method != "POST" {
		w.WriteHeader(http.StatusBadRequest)
		respData = responseData{Status: "Error", Message: "POST Error."}
		errorFlag = true
	}

	if req.Header.Get("Content-Type") != "application/json" {
		w.WriteHeader(http.StatusBadRequest)
		respData = responseData{Status: "Error", Message: "JSON Error."}
		errorFlag = true
	}

	length, err := strconv.Atoi(req.Header.Get("Content-Length"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		respData = responseData{Status: "Error", Message: "Content Length Error."}
		errorFlag = true
	}

	body := make([]byte, length)
	length, err = req.Body.Read(body)
	if err != nil && err != io.EOF {
		w.WriteHeader(http.StatusInternalServerError)
		respData = responseData{Status: "Error", Message: "Body Read Error."}
		errorFlag = true
	}

	var jsonBody map[string]interface{}
	err = json.Unmarshal(body[:length], &jsonBody)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		respData = responseData{Status: "Error", Message: "JSON Unmarshal Error."}
		errorFlag = true
	}

	command := fmt.Sprintf("%s", req.URL)

	if debug == true {
		fmt.Printf("%v\n", jsonBody)
		fmt.Println("command: ", command)
	}

	if errorFlag == false {
		w.WriteHeader(http.StatusOK)
		if command == "/find" || command == "/raw" {
			respList = DbFindOrRaw(jsonBody, command)
			outputJson, err := json.Marshal(respList)
			if err != nil {
				fmt.Println(err)
				return
			}

			if debug == true {
				fmt.Println(outputJson)
			}
			w.Write(outputJson)
		} else {
			respData = DbSwtich(jsonBody, command)
			outputJson, err := json.Marshal(respData)
			if err != nil {
				fmt.Println(err)
				return
			}

			if debug == true {
				fmt.Println(outputJson)
			}
			w.Write(outputJson)
		}
	}
}

func CallRaw(rawString string) responseList {
	listData := responseList{}

	jsonDatas := []member{}
	DBMS.Raw(rawString).Scan(&jsonDatas)
	for i := 0; i < len(jsonDatas); i++ {
		listData.Members = append(listData.Members, jsonDatas[i])
	}
	listData.Status = "Success"
	return listData
}

func CallFind(searchString string) responseList {
	listData := responseList{}

	jsonDatas := []member{}

	DBMS.Find(&jsonDatas, "name=?", searchString)
	for i := 0; i < len(jsonDatas); i++ {
		listData.Members = append(listData.Members, jsonDatas[i])
	}
	listData.Status = "Success"
	return listData
}

func DbFindOrRaw(jsonBody map[string]interface{}, command string) responseList {
	listData := responseList{}

	switch command {
	case "/raw":
		//curl -k -H "Content-Type: application/json" -X POST -d '{"raw":" SELECT * FROM test.member;"}' https://localhost:8080/raw
		rawString, okRaw := jsonBody["raw"].(string)
		if okRaw == false {
			return responseList{Status: "Error", Members: nil}
		}
		return CallRaw(rawString)
	case "/find":
		// curl -k -H "Content-Type: application/json" -X POST -d '{"search":"user3"}' https://localhost:8080/find
		searchString, okSearch := jsonBody["search"].(string)
		if okSearch == false {
			return responseList{Status: "Error", Members: nil}
		}
		return CallFind(searchString)
	}
	return listData
}

func CallCreate(tmpUser, tmpPass string) responseData {
	jsonData := member{}

	if len(tmpUser) > 8 || len(tmpPass) > 8 {
		return responseData{Status: "Error", Message: "Username or Password > 8."}
	}

	if CheckExist(tmpUser) == true {
		return responseData{Status: "Error", Message: "user already exsits."}
	}

	jsonData.Name = tmpUser
	edStr, err := Encrypt(tmpPass, []byte(AddSpace(dbSalt)))
	if err != nil {
		fmt.Println(err)
		return responseData{Status: "Error", Message: "Can't Encode."}
	}
	jsonData.Password = edStr
	DBMS.Create(&jsonData)
	return responseData{Status: "Success", Message: fmt.Sprintf("%v", jsonData)}
}

func CallUpdate(IdReq, tmpUser, tmpPass string) responseData {
	jsonData := member{}

	Ids, err := strconv.Atoi(IdReq)
	if err != nil {
		return responseData{Status: "Error", Message: "Id not found."}
	}

	if len(tmpUser) > 8 || len(tmpPass) > 8 {
		return responseData{Status: "Error", Message: "Username or Password > 8."}
	}

	if CheckExist(tmpUser) == false {
		return responseData{Status: "Error", Message: "user not found."}
	}

	jsonData.Id = Ids
	DBMS.First(&jsonData)

	updateFlag := false
	if jsonData.Name != tmpUser {
		jsonData.Name = tmpUser
		updateFlag = true
	}

	edStr, err := Encrypt(tmpPass, []byte(AddSpace(dbSalt)))
	if err != nil {
		return responseData{Status: "Error", Message: "Can't Encode."}
	}

	if jsonData.Password != edStr {
		jsonData.Password = edStr
		updateFlag = true
	}

	if updateFlag == true {
		DBMS.Save(&jsonData)
	} else {
		return responseData{Status: "Error", Message: "Username and Password same."}
	}
	return responseData{Status: "Success", Message: fmt.Sprintf("%v", jsonData)}
}

func CallAuth(tmpUser, tmpPass string) responseData {
	jsonData := member{}

	if len(tmpUser) > 8 || len(tmpPass) > 8 {
		return responseData{Status: "Error", Message: "Username or Password > 8."}
	}

	DBMS.First(&jsonData, "name=?", tmpUser)

	edStr, err := Decrypt(jsonData.Password, []byte(AddSpace(dbSalt)))
	if err != nil {
		return responseData{Status: "Error", Message: "Can't Encode."}
	}

	if edStr == tmpPass {
		return responseData{Status: "Success", Message: "auth ok."}
	}
	return responseData{Status: "Error", Message: "auth fail."}
}

func DbSwtich(jsonBody map[string]interface{}, command string) responseData {
	var resp responseData
	jsonData := member{}

	switch command {
	case "/create":
		// curl -k -H "Content-Type: application/json" -X POST -d '{"name":"user2", "password": "pass"}' https://localhost:8080/create
		tmpUser, okUser := jsonBody["name"].(string)
		tmpPass, okPass := jsonBody["password"].(string)
		if okUser == false || okPass == false {
			return responseData{Status: "Error", Message: "Username or Password not found."}
		}

		resp = CallCreate(tmpUser, tmpPass)
	case "/read":
		// curl -k -H "Content-Type: application/json" -X POST -d '{"id":"8"}' https://localhost:8080/read
		IdReq := jsonBody["id"].(string)
		Ids, err := strconv.Atoi(IdReq)

		if err != nil {
			return responseData{Status: "Error", Message: "Id not found."}
		}

		jsonData.Id = Ids
		DBMS.First(&jsonData)
	case "/update":
		// curl -k -H "Content-Type: application/json" -X POST -d '{"id": "8", "name":"user2", "password": "pass"}' https://localhost:8080/update
		IdReq := jsonBody["id"].(string)
		tmpUser, okUser := jsonBody["name"].(string)
		tmpPass, okPass := jsonBody["password"].(string)
		if okUser == false || okPass == false {
			return responseData{Status: "Error", Message: "Username or Password not found."}
		}

		resp = CallUpdate(IdReq, tmpUser, tmpPass)
	case "/delete":
		// curl -k -H "Content-Type: application/json" -X POST -d '{"id":1}' https://localhost:8080/delete
		IdReq := jsonBody["id"].(string)
		Ids, err := strconv.Atoi(IdReq)

		if err != nil {
			return responseData{Status: "Error", Message: "Id not found."}
		}

		jsonData.Id = Ids
		DBMS.First(&jsonData)
		DBMS.Delete(&jsonData)
	case "/auth":
		// curl -k -H "Content-Type: application/json" -X POST -d '{"name":"user2", "password": "pass"}' https://localhost:8080/auth
		tmpUser, okUser := jsonBody["name"].(string)
		tmpPass, okPass := jsonBody["password"].(string)

		if okUser == false || okPass == false {
			return responseData{Status: "Error", Message: "Username or Password not found."}
		}

		resp = CallAuth(tmpUser, tmpPass)
	}

	return resp
}

func CheckExist(username string) bool {
	jsonDatas := []member{}
	DBMS.Find(&jsonDatas, "name=?", username)
	if len(jsonDatas) > 0 {
		return true
	}
	return false
}

func AddSpace(strs string) string {
	for i := 0; len(strs) < 16; i++ {
		strs += "0"
	}
	return strs
}

// FYI: http://www.inanzzz.com/index.php/post/f3pe/data-encryption-and-decryption-with-a-secret-key-in-golang
// encrypt encrypts plain string with a secret key and returns encrypt string.
func Encrypt(plainData string, secret []byte) (string, error) {
	cipherBlock, err := aes.NewCipher(secret)
	if err != nil {
		return "", err
	}

	aead, err := cipher.NewGCM(cipherBlock)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, aead.NonceSize())
	if _, err = io.ReadFull(crt.Reader, nonce); err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(aead.Seal(nonce, nonce, []byte(plainData), nil)), nil
}

// decrypt decrypts encrypt string with a secret key and returns plain string.
func Decrypt(encodedData string, secret []byte) (string, error) {
	encryptData, err := base64.URLEncoding.DecodeString(encodedData)
	if err != nil {
		return "", err
	}

	cipherBlock, err := aes.NewCipher(secret)
	if err != nil {
		return "", err
	}

	aead, err := cipher.NewGCM(cipherBlock)
	if err != nil {
		return "", err
	}

	nonceSize := aead.NonceSize()
	if len(encryptData) < nonceSize {
		return "", err
	}

	nonce, cipherText := encryptData[:nonceSize], encryptData[nonceSize:]
	plainData, err := aead.Open(nil, nonce, cipherText, nil)
	if err != nil {
		return "", err
	}

	return string(plainData), nil
}

func GormConnect() *gorm.DB {
	DBTYPE := "mysql"
	CONNECT := dbUser + ":" + dbPass + "@tcp(" + dbAddress + ":3306)/test?charset=utf8&parseTime=True&loc=Local"

	DB, err := gorm.Open(DBTYPE, CONNECT)
	if err != nil {
		fmt.Println("RDS access error!")
		panic(err.Error())
	}

	DB.LogMode(true)
	DB.SingularTable(true)
	DB.AutoMigrate(&member{})

	return DB
}

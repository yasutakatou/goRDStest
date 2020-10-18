package main

import (
	"net/http"
	"reflect"
	"testing"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func dbInit() {
	dbUser = "admin"
	dbPass = ""
	dbAddress = "database-2.czhb3qjocyn9.us-east-2.rds.amazonaws.com"
	dbSalt = "api12345"
	DBMS = GormConnect()
}

func TestDecrypt(t *testing.T) {
	type args struct {
		encodedData string
		secret      []byte
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Decrypt(tt.args.encodedData, tt.args.secret)
			if (err != nil) != tt.wantErr {
				t.Errorf("Decrypt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Decrypt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEncrypt(t *testing.T) {
	type args struct {
		plainData string
		secret    []byte
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Encrypt(tt.args.plainData, tt.args.secret)
			if (err != nil) != tt.wantErr {
				t.Errorf("Encrypt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Encrypt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDbSwtich(t *testing.T) {
	type args struct {
		jsonBody map[string]interface{}
		command  string
	}
	tests := []struct {
		name string
		args args
		want responseData
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DbSwtich(tt.args.jsonBody, tt.args.command); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DbSwtich() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDbFindOrRaw(t *testing.T) {
	type args struct {
		jsonBody map[string]interface{}
		command  string
	}
	tests := []struct {
		name string
		args args
		want responseList
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DbFindOrRaw(tt.args.jsonBody, tt.args.command); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DbFindOrRaw() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestApiHandlers(t *testing.T) {
	type args struct {
		w   http.ResponseWriter
		req *http.Request
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ApiHandlers(tt.args.w, tt.args.req)
		})
	}
}

/////////////////////////////////////// Test code ok. //////////////////////////////////////////////////

func TestCallUpdate(t *testing.T) {
	dbInit()
	CallRaw("DELETE FROM member WHERE name = 'katotest';")
	CallRaw("INSERT INTO member (name, password) VALUES ('katotest' ,'password');")
	type args struct {
		IdReq   string
		tmpUser string
		tmpPass string
	}
	tests := []struct {
		name string
		args args
		want responseData
	}{
		{
			args: args{
				IdReq:   "10",
				tmpUser: "katotest",
				tmpPass: "pass",
			},
			want: responseData{Status: "Success", Message: ""},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {
					if got := CallUpdate(tt.args.IdReq, tt.args.tmpUser, tt.args.tmpPass); got.Status != "Success" {
						t.Errorf("CallUpdate() = %v, want %v", got, tt.want)
					}
				})
			}
		})
	}
	CallRaw("DELETE FROM member WHERE name = 'katotest';")
}

func TestGormConnect(t *testing.T) {
	tests := []struct {
		name string
		want *gorm.DB
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GormConnect(); got == nil {
				t.Errorf("GormConnect() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCallCreate(t *testing.T) {
	dbInit()
	CallRaw("DELETE FROM member WHERE name = 'katotest';")
	type args struct {
		tmpUser string
		tmpPass string
	}
	tests := []struct {
		name string
		args args
		want responseData
	}{
		{
			args: args{
				tmpUser: "katotest",
				tmpPass: "pass",
			},
			want: responseData{Status: "Success", Message: ""},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CallCreate(tt.args.tmpUser, tt.args.tmpPass); got.Status != "Success" {
				t.Errorf("CallCreate() = %v, want %v", got, tt.want)
			}
		})
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CallCreate(tt.args.tmpUser, tt.args.tmpPass); got.Message != "user already exsits." {
				t.Errorf("CallCreate() = %v, want %v", got, tt.want)
			}
		})
	}
	CallRaw("DELETE FROM member WHERE name = 'katotest';")
}

func TestCallFind(t *testing.T) {
	dbInit()
	type args struct {
		searchString string
	}
	tests := []struct {
		name string
		args args
		want responseList
	}{
		{
			args: args{
				searchString: "kato2",
			},
			want: responseList{Status: "Success", Members: nil},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CallFind(tt.args.searchString); got.Members[0].Name != "kato2" {
				t.Errorf("CallFind() = %v, want %v", got, "kato2")
			}
		})
	}
}

func TestCheckExist(t *testing.T) {
	dbInit()
	type args struct {
		username string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			args: args{
				username: "kato2",
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CheckExist(tt.args.username); got != tt.want {
				t.Errorf("CheckExist() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCallRaw(t *testing.T) {
	dbInit()
	type args struct {
		rawString string
	}
	tests := []struct {
		name string
		args args
		want responseList
	}{
		{
			args: args{
				rawString: "SELECT name FROM member WHERE name='kato2';",
			},
			want: responseList{Status: "Success", Members: nil},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CallRaw(tt.args.rawString); got.Members[0].Name != "kato2" {
				t.Errorf("CallRaw() = %v, want %v", got, "kato2")
			}
		})
	}
}

func TestAddSpace(t *testing.T) {
	type args struct {
		strs string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			args: args{
				strs: "test1234",
			},
			want: "test123400000000",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := AddSpace(tt.args.strs); got != tt.want {
				t.Errorf("AddSpace() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCallAuth(t *testing.T) {
	dbInit()
	type args struct {
		tmpUser string
		tmpPass string
	}
	tests := []struct {
		name string
		args args
		want responseData
	}{
		{
			args: args{
				tmpUser: "kato2",
				tmpPass: "pass",
			},
			want: responseData{Status: "Success", Message: "auth ok."},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CallAuth(tt.args.tmpUser, tt.args.tmpPass); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CallAuth() = %v, want %v", got, tt.want)
			}
		})
	}
}

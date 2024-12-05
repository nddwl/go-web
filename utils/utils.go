package utils

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"github.com/bwmarrin/snowflake"
	"github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"io"
	"os"
)

var (
	node   = &snowflake.Node{}
	key    = &rsa.PrivateKey{}
	public = ""
)

func init() {
	var err error
	node, err = snowflake.NewNode(1)
	if err != nil {
		panic(err)
	}
	loadKey()
}

func loadKey() {
	f, err := os.Open("private.pem")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	b, err := io.ReadAll(f)
	if err != nil {
		panic(err)
	}
	block, _ := pem.Decode(b)
	a, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		panic(err)
	}
	ok := false
	key, ok = a.(*rsa.PrivateKey)
	if !ok {
		panic(errors.New("unexpected private key type"))
	}
	f1, err := os.Open("public.pem")
	if err != nil {
		panic(err)
	}
	defer f1.Close()
	b1, err := io.ReadAll(f1)
	if err != nil {
		panic(err)
	}
	block1, _ := pem.Decode(b1)
	public = base64.StdEncoding.EncodeToString(block1.Bytes)
}

func MySqlErr(err error) (number uint16, message string) {
	mysqlErr := &mysql.MySQLError{}
	if errors.As(err, &mysqlErr) {
		number = mysqlErr.Number
		message = mysqlErr.Message
	}
	return
}

func GenerateUid() int64 {
	return node.Generate().Int64()
}

func GenerateToken() string {
	return uuid.New().String()
}

func Decrypt(message string) (data []byte, err error) {
	b1, err := base64.StdEncoding.DecodeString(message)
	if err != nil {
		return
	}
	data, err = rsa.DecryptPKCS1v15(rand.Reader, key, b1)
	return
}

func GetPublicKey() string {
	return public
}

func Encrypt(message []byte) (string, error) {
	b, err := rsa.EncryptPKCS1v15(rand.Reader, &key.PublicKey, message)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(b), nil
}

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash), err
}

func ComparePassword(pwdHash string, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(pwdHash), []byte(password)) == nil
}

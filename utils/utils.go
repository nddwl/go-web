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
	node = &snowflake.Node{}
	key  = &rsa.PrivateKey{}
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
	key, err = x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		panic(err)
	}
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

func generateKey() error {
	privKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return err
	}
	privBytes := x509.MarshalPKCS1PrivateKey(privKey)
	f1, err := os.Create("private.pem")
	if err != nil {
		return err
	}
	defer f1.Close()
	err = pem.Encode(f1, &pem.Block{
		Type:    "RSA PRIVATE KEY",
		Headers: nil,
		Bytes:   privBytes,
	})
	if err != nil {
		return err
	}
	publKey := &privKey.PublicKey
	publBytes := x509.MarshalPKCS1PublicKey(publKey)
	f2, err := os.Create("public.pem")
	if err != nil {
		return err
	}
	defer f2.Close()
	err = pem.Encode(f2, &pem.Block{
		Type:    "RSA PUBLIC KEY",
		Headers: nil,
		Bytes:   publBytes,
	})
	return err
}

func Decrypt(message string) (data []byte, err error) {
	b1, err := base64.StdEncoding.DecodeString(message)
	if err != nil {
		return
	}
	data, err = rsa.DecryptPKCS1v15(rand.Reader, key, b1)
	return
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

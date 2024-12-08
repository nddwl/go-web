package validate

import (
	"github.com/dlclark/regexp2"
	"github.com/go-playground/validator/v10"
	"github.com/microcosm-cc/bluemonday"
	"regexp"
)

var validate = &validator.Validate{}
var noHTML = bluemonday.NewPolicy()

func init() {
	validate = validator.New()
	validate.RegisterValidation("name", func(fl validator.FieldLevel) bool {
		name := fl.Field().String()
		return Name(name)
	})
	validate.RegisterValidation("phone", func(fl validator.FieldLevel) bool {
		phone := fl.Field().String()
		return Phone(phone)
	})
	validate.RegisterValidation("email", func(fl validator.FieldLevel) bool {
		email := fl.Field().String()
		return Email(email)
	})
	validate.RegisterValidation("uid", func(fl validator.FieldLevel) bool {
		uid := fl.Field().Int()
		return Uid(uid)
	})
	validate.RegisterValidation("password", func(fl validator.FieldLevel) bool {
		password := fl.Field().String()
		return Password(password)
	})
	validate.RegisterValidation("code", func(fl validator.FieldLevel) bool {
		code := fl.Field().String()
		return Code(code)
	})
	validate.RegisterValidation("token", func(fl validator.FieldLevel) bool {
		token := fl.Field().String()
		return Token(token)
	})
	validate.RegisterValidation("noHTML", func(fl validator.FieldLevel) bool {
		s := fl.Field().String()
		fl.Field().SetString(NoHTML(s))
		return true
	})
	validate.RegisterValidation("safeInput", func(fl validator.FieldLevel) bool {
		s := fl.Field().String()
		return IsSafeInput(s)
	})
}

func Email(email string) bool {
	// 名称允许汉字、字母、数字，域名只允许英文域名
	re := regexp.MustCompile(`^[\p{Han}a-zA-Z0-9]{3,20}@(gmail\.com|qq\.com|163\.com|126\.com|outlook\.com)`)
	return re.MatchString(email)
}

func Name(name string) bool {
	re := regexp.MustCompile(`^[\p{Han}a-zA-Z0-9_]{4,10}$`)
	return re.MatchString(name)
}

func Phone(phone string) bool {
	re := regexp.MustCompile(`^1(3[0-9]|4[01456879]|5[0-35-9]|6[2567]|7[0-8]|8[0-9]|9[0-35-9])\d{8}$`)
	return re.MatchString(phone)
}

func Struct(s interface{}) error {
	return validate.Struct(s)
}

func Uid(uid int64) bool {
	return uid > 1000000000000000000
}

func Password(password string) bool {
	re := regexp2.MustCompile(`^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[~!@#$%^&*()_+={}\[\]:;"'<>,.?\\|-]).{8,16}$`, 0)
	ok, err := re.MatchString(password)
	if err == nil && ok {
		return true
	}
	return false
}

func Avatar(avatar string) bool {
	return true
}

func Token(token string) bool {
	re := regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)
	return re.MatchString(token)
}

func Code(code string) bool {
	re := regexp.MustCompile(`^[a-z0-9]{4}$`)
	return re.MatchString(code)
}

func NoHTML(s string) string {
	return noHTML.Sanitize(s)
}

func IsSafeInput(s string) bool {
	re := regexp.MustCompile(`[\x00-\x1F\x7F\x{200B}]`)
	return !re.MatchString(s)
}

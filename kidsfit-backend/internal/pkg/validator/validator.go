package validator

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/go-playground/validator/v10"
)

// Validator 参数校验器，封装go-playground/validator提供结构体校验能力
type Validator struct {
	validate *validator.Validate
}

// NewValidator 创建校验器实例并注册自定义校验规则
func NewValidator() *Validator {
	v := &Validator{
		validate: validator.New(),
	}
	v.RegisterCustomValidations()
	return v
}

// RegisterCustomValidations 注册自定义校验规则
// 包括手机号校验、年龄校验和评分校验
func (v *Validator) RegisterCustomValidations() {
	// 注册手机号校验规则：11位数字，1开头
	_ = v.validate.RegisterValidation("phone", validatePhone)
	// 注册年龄校验规则：3-12岁
	_ = v.validate.RegisterValidation("age", validateAge)
	// 注册评分校验规则：0-100
	_ = v.validate.RegisterValidation("score", validateScore)
}

// Validate 校验结构体字段，根据struct tag中的validate规则进行校验
// s: 待校验的结构体指针
func (v *Validator) Validate(s interface{}) error {
	return v.validate.Struct(s)
}

// FormatValidationError 将校验错误格式化为字段-错误信息的映射
// 返回每个校验失败字段的中文错误描述
// err: validator返回的校验错误
func FormatValidationError(err error) map[string]string {
	result := make(map[string]string)
	if err == nil {
		return result
	}

	validationErrs, ok := err.(validator.ValidationErrors)
	if !ok {
		result["error"] = err.Error()
		return result
	}

	for _, e := range validationErrs {
		field := e.Field()
		switch e.Tag() {
		case "required":
			result[field] = fmt.Sprintf("%s为必填项", field)
		case "phone":
			result[field] = fmt.Sprintf("%s必须是11位数字且以1开头", field)
		case "age":
			result[field] = fmt.Sprintf("%s必须在3-12岁之间", field)
		case "score":
			result[field] = fmt.Sprintf("%s必须在0-100之间", field)
		case "min":
			result[field] = fmt.Sprintf("%s不能小于%s", field, e.Param())
		case "max":
			result[field] = fmt.Sprintf("%s不能大于%s", field, e.Param())
		case "len":
			result[field] = fmt.Sprintf("%s长度必须为%s", field, e.Param())
		case "email":
			result[field] = fmt.Sprintf("%s格式不正确", field)
		case "oneof":
			result[field] = fmt.Sprintf("%s必须是[%s]中的一个", field, e.Param())
		default:
			result[field] = fmt.Sprintf("%s校验失败: %s", field, e.Tag())
		}
	}
	return result
}

// validatePhone 手机号校验函数：11位数字，1开头
// fl: 校验字段信息
func validatePhone(fl validator.FieldLevel) bool {
	phone := fl.Field().String()
	matched, _ := regexp.MatchString(`^1\d{10}$`, phone)
	return matched
}

// validateAge 年龄校验函数：3-12岁
// fl: 校验字段信息
func validateAge(fl validator.FieldLevel) bool {
	ageStr := fmt.Sprintf("%v", fl.Field().Interface())
	age, err := strconv.Atoi(ageStr)
	if err != nil {
		return false
	}
	return age >= 3 && age <= 12
}

// validateScore 评分校验函数：0-100
// fl: 校验字段信息
func validateScore(fl validator.FieldLevel) bool {
	scoreStr := fmt.Sprintf("%v", fl.Field().Interface())
	score, err := strconv.Atoi(scoreStr)
	if err != nil {
		return false
	}
	return score >= 0 && score <= 100
}

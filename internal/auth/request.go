package auth

import "github.com/nayefradwi/nayef_go_common/validation"

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

func (r RefreshRequest) Validate() error {
	validator, ruleF := validation.NewValidator(), validation.NewStringValidationRuleFactory()
	validation.AddRule(validator, ruleF.IsRequired(r.RefreshToken, "refresh_token"))
	return validator.Validate()
}

type CodeRequest struct {
	Email string `json:"email"`
}

func (r CodeRequest) Validate() error {
	validator, ruleF := validation.NewValidator(), validation.NewStringValidationRuleFactory()
	validation.AddRule(validator, ruleF.IsEmail(r.Email, "email"))
	return validator.Validate()
}

type VerifyCodeRequest struct {
	Email string `json:"email"`
	Code  string `json:"code"`
}

func (r VerifyCodeRequest) Validate() error {
	validator, ruleF := validation.NewValidator(), validation.NewStringValidationRuleFactory()
	validation.AddRule(validator, ruleF.IsEmail(r.Email, "email"))
	validation.AddRule(validator, ruleF.IsRequired(r.Code, "code"))
	return validator.Validate()
}

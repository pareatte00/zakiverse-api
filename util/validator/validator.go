package validator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"slices"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

const (
	errNonStructReceived = "a non struct is received"
	errReflectError      = "reflect field error"
	errEmptyTagValue     = "value of tag '%s' in field '%s' is empty"
)

const (
	tagRequired          = "required"
	tagEmail             = "email"
	tagMin               = "min"
	tagMax               = "max"
	tagUUID              = "uuid"
	tagUUID4             = "uuid4"
	tagOneOf             = "oneof"
	tagHttpUrl           = "http_url"
	tagBase64            = "base64"
	tagStartsWith        = "startswith"
	tagEndsWith          = "endswith"
	tagJwt               = "jwt"
	tagE164              = "e164"
	tagDateTime          = "datetime"
	tagAlpha             = "alpha"
	tagIso3166Alpha2     = "iso3166_alpha2"
	tagIso3166Alpha3     = "iso3166_alpha3"
	tagIso8601Datetime   = "iso8601datetime"
	tagIso8601NoTimezone = "iso8601notz"
	tagGoUUID            = "gouuid"
	tagGoTime            = "gotime"
	tagDate              = "date"
	tagTime              = "time"
	tagAlphaSpace        = "alpha_space"
	tagAlphaSpaceTH      = "alpha_space_th"
	tagAlphaTH           = "alpha_th"
	tagUniqueField       = "uniquefield"
	tagSplitOneOf        = "splitoneof"
	tagSplitUUID         = "splituuid"
	tagAllOrNone         = "allornone"
	tagDecimal           = "decimal"
	tagUnsignedDecimal   = "udecimal"
	tagDecimalEq         = "deq"
	tagDecimalGt         = "dgt"
	tagDecimalGte        = "dgte"
	tagDecimalLt         = "dlt"
	tagDecimalLte        = "dlte"
	tagDecimalNe         = "dne"
)

type FieldError struct {
	Field    string `json:"field"`
	Value    any    `json:"value"`
	ErrorMsg string `json:"error_message"`
}

var validate *validator.Validate

func init() {
	validate = validator.New()

	validate.RegisterValidation(tagE164, vxe164)
	validate.RegisterValidation(tagGoTime, vxgotime)
	validate.RegisterValidation(tagGoUUID, vxgouuid)
	validate.RegisterValidation(tagIso8601Datetime, vxiso8601datetime)
	validate.RegisterValidation(tagIso8601NoTimezone, vxiso8601notimezone)
	validate.RegisterValidation(tagDate, vxdate)
	validate.RegisterValidation(tagTime, vxtime)
	validate.RegisterValidation(tagAlphaSpace, vxalphaspace)
	validate.RegisterValidation(tagAlphaSpaceTH, vxalphaspaceth)
	validate.RegisterValidation(tagAlphaTH, vxalphath)
	validate.RegisterValidation(tagUniqueField, vxuniquefield)
	validate.RegisterValidation(tagSplitOneOf, vxsplitoneof)
	validate.RegisterValidation(tagSplitUUID, vxsplituuid)
	validate.RegisterValidation(tagAllOrNone, vxallornone)
	validate.RegisterValidation(tagDecimal, vxdecimal)
	validate.RegisterValidation(tagUnsignedDecimal, vxudecimal)
	validate.RegisterValidation(tagDecimalEq, vxdeq)
	validate.RegisterValidation(tagDecimalGt, vxdgt)
	validate.RegisterValidation(tagDecimalGte, vxdgte)
	validate.RegisterValidation(tagDecimalLt, vxdlt)
	validate.RegisterValidation(tagDecimalLte, vxdlte)
	validate.RegisterValidation(tagDecimalNe, vxdne)
}

func getErrorMessage(errTag string, errParam any) string {
	switch errTag {
	case tagRequired:
		return "this field is required"
	case tagEmail:
		return "this field must be a valid email address"
	case tagMin:
		return fmt.Sprintf("this field must be at least %s characters", errParam)
	case tagMax:
		return fmt.Sprintf("this field must be not longer than %s characters", errParam)
	case tagUUID:
		return "this field must be a valid UUID"
	case tagUUID4:
		return "this field must be a valid UUID4"
	case tagOneOf:
		return fmt.Sprintf("this field must be one of %s", strings.Replace(errParam.(string), " ", ",", -1))
	case tagHttpUrl:
		return "this field must be a valid URL"
	case tagBase64:
		return "this field must be a valid base64 string"
	case tagStartsWith:
		return fmt.Sprintf("this field must be start with %s", errParam)
	case tagEndsWith:
		return fmt.Sprintf("this field must be end with %s", errParam)
	case tagJwt:
		return "this field must be a valid JWT"
	case tagE164:
		return "this field must be a valid E.164 phone number"
	case tagDateTime:
		return "this field must be a valid datetime"
	case tagDate:
		return "This field must be a valid date"
	case tagTime:
		return "This field must be a valid time"
	case tagGoUUID:
		return "This field must be a valid uuid.UUID"
	case tagGoTime:
		return "This field must be a valid time.Time"
	case tagIso8601Datetime:
		return "This field must be a valid ISO8601 datetime"
	case tagIso8601NoTimezone:
		return "This field must be a valid ISO8601 datetime without timezone"
	case tagIso3166Alpha2:
		return "This field must be a valid ISO3166 alpha-2 country code"
	case tagIso3166Alpha3:
		return "This field must be a valid ISO3166 alpha-3 country code"
	case tagAlpha:
		return "This field must contain letters only"
	case tagAlphaSpace:
		return "this field must contain letters and spaces only"
	case tagAlphaSpaceTH:
		return "this field must contain letters and spaces only"
	case tagAlphaTH:
		return "this field must contain letters only"
	case tagSplitOneOf:
		return fmt.Sprintf("this field must be one of %s", strings.Replace(errParam.(string), " ", ",", -1))
	case tagSplitUUID:
		return "this field must be a valid UUID"
	case tagAllOrNone:
		return fmt.Sprintf("Required all %s or none", errParam)
	case tagDecimal:
		return "this field must be a valid 2 points decimal"
	case tagUnsignedDecimal:
		return "this field must be a valid decimal"
	case tagDecimalEq:
		return fmt.Sprintf("this field must be a decimal and equal %s", errParam)
	case tagDecimalGt:
		return fmt.Sprintf("this field must be a decimal and greater than %s", errParam)
	case tagDecimalGte:
		return fmt.Sprintf("this field must be a decimal and greater than or equal %s", errParam)
	case tagDecimalLt:
		return fmt.Sprintf("this field must be a decimal and less than %s", errParam)
	case tagDecimalLte:
		return fmt.Sprintf("this field must be a decimal and less than or equal %s", errParam)
	case tagDecimalNe:
		return fmt.Sprintf("this field must be a decimal and not equal %s", errParam)
	}

	return "this field is invalid"
}

var extractIndexRegex = regexp.MustCompile(`\[(\d+)\]$`)

func extractIndex(s string) (string, string) {
	match := extractIndexRegex.FindStringSubmatch(s)

	if len(match) > 0 {
		return s[:len(s)-len(match[0])], match[0]
	}
	return s, ""
}

var defaultTags = []string{"json", "form", "uri", "header", "mapstructure"}

func findField(t reflect.Type, name string) (reflect.StructField, bool) {
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)

		if f.Name == name {
			return f, true
		}

		if f.Anonymous {
			ft := f.Type
			for ft.Kind() == reflect.Ptr {
				ft = ft.Elem()
			}

			if ft.Kind() == reflect.Struct {
				if sf, ok := findField(ft, name); ok {
					return sf, true
				}
			}
		}
	}
	return reflect.StructField{}, false
}

type GetFieldNameParam struct {
	Struct    any
	FieldName string
	TagName   []string
}

func getFieldName(config GetFieldNameParam) (string, bool, error) {
	t := reflect.TypeOf(config.Struct)
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	if t.Kind() != reflect.Struct {
		return "", false, errors.New(errReflectError)
	}

	field, ok := findField(t, config.FieldName)
	if !ok {
		return config.FieldName, false, nil
	}

	for _, tag := range config.TagName {
		if v := field.Tag.Get(tag); v != "" {
			return v, true, nil
		}
	}

	return config.FieldName, false, nil
}

func getInnerStruct(errorStruct any, fieldName string) any {
	if errorStruct == nil {
		return nil
	}

	v := reflect.ValueOf(errorStruct)
	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return nil
	}

	f := v.FieldByName(fieldName)
	if f.IsValid() {
		if f.Kind() == reflect.Slice || f.Kind() == reflect.Array {
			if f.Len() == 0 {
				return nil
			}
			f = f.Index(0)
		}

		if f.Kind() == reflect.Struct {
			return f.Interface()
		}
	}

	for i := 0; i < v.NumField(); i++ {
		sf := v.Type().Field(i)
		if !sf.Anonymous {
			continue
		}

		fv := v.Field(i)
		if fv.Kind() == reflect.Struct {
			if inner := getInnerStruct(fv.Interface(), fieldName); inner != nil {
				return inner
			}
		}
	}

	return nil
}

func Validate(vs any) (isValid bool, fieldErrors []FieldError, validatorError error) {
	defer func() {
		if r := recover(); r != nil {
			validatorError = fmt.Errorf("validator panic: %v", r)
		}
	}()

	if vs == nil {
		return false, nil, errors.New("validation input is nil")
	}

	t := reflect.TypeOf(vs)
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	if t.Kind() != reflect.Struct {
		return false, nil, errors.New("validation expects a struct")
	}

	verrors := validate.Struct(vs)
	if verrors == nil {
		return true, nil, nil
	}

	validationErrors, ok := verrors.(validator.ValidationErrors)
	if !ok {
		return false, nil, verrors
	}

	for _, ferr := range validationErrors {
		errorStruct := vs

		structInfo := strings.Split(ferr.StructNamespace(), ".")
		if len(structInfo) <= 1 {
			return false, nil, errors.New(errNonStructReceived)
		}

		var errFieldName string
		var prevOriginalFieldName string

		for fieldIndex, originalFieldName := range structInfo {
			if fieldIndex == 0 {
				continue
			}
			if fieldIndex >= 2 {
				errorStruct = getInnerStruct(errorStruct, prevOriginalFieldName)
			}

			fieldNamePattern, indexPattern := extractIndex(originalFieldName)

			fieldName, _, errGetFieldName := getFieldName(GetFieldNameParam{
				Struct:    errorStruct,
				FieldName: fieldNamePattern,
				TagName:   defaultTags,
			})
			if errGetFieldName != nil {
				return false, nil, errGetFieldName
			}

			innerErrFieldName := fieldName + indexPattern
			if fieldIndex >= 2 {
				innerErrFieldName = "." + innerErrFieldName
			}

			errFieldName += innerErrFieldName
			prevOriginalFieldName = fieldNamePattern
		}

		// --- Resolve error message ---
		customMessage, foundCustomTag, errCustom := getFieldName(GetFieldNameParam{
			Struct:    errorStruct,
			FieldName: prevOriginalFieldName,
			TagName:   []string{"validateErrorMessage"},
		})

		if errCustom != nil || !foundCustomTag {
			customMessage = getErrorMessage(ferr.Tag(), ferr.Param())
		} else {
			switch customMessage {
			case "HIDE":
				customMessage = ""
			case "DEFAULT":
				customMessage = "this field is invalid"
			}
		}

		fieldErrors = append(fieldErrors, FieldError{
			Field:    errFieldName,
			Value:    ferr.Value(),
			ErrorMsg: customMessage,
		})
	}

	return false, fieldErrors, nil
}

var vxiso8601datetimeregexstr = "^(?:[1-9]\\d{3}-(?:(?:0[1-9]|1[0-2])-(?:0[1-9]|1\\d|2[0-8])|(?:0[13-9]|1[0-2])-(?:29|30)|(?:0[13578]|1[02])-31)|(?:[1-9]\\d(?:0[48]|[2468][048]|[13579][26])|(?:[2468][048]|[13579][26])00)-02-29)T(?:[01]\\d|2[0-3]):[0-5]\\d:[0-5]\\d(?:\\.\\d{1,9})?(?:Z|[+-][01]\\d:[0-5]\\d)$"
var vxiso8601datetimeregex = regexp.MustCompile(vxiso8601datetimeregexstr)

func vxiso8601datetime(fl validator.FieldLevel) bool {
	return vxiso8601datetimeregex.MatchString(fl.Field().String())
}

var vxiso8601notimezoneregexstr = "^(?:[1-9]\\d{3}-(?:(?:0[1-9]|1[0-2])-(?:0[1-9]|1\\d|2[0-8])|(?:0[13-9]|1[0-2])-(?:29|30)|(?:0[13578]|1[02])-31)|(?:[1-9]\\d(?:0[48]|[2468][048]|[13579][26])|(?:[2468][048]|[13579][26])00)-02-29)T(?:[01]\\d|2[0-3]):[0-5]\\d:[0-5]\\d(?:\\.\\d{1,9})?$"
var vxiso8601notimezoneregex = regexp.MustCompile(vxiso8601notimezoneregexstr)

func vxiso8601notimezone(fl validator.FieldLevel) bool {
	return vxiso8601notimezoneregex.MatchString(fl.Field().String())
}

func vxgotime(fl validator.FieldLevel) bool {
	field := fl.Field()

	if field.Kind() == reflect.Ptr {
		if field.IsNil() {
			return true
		}
		field = field.Elem()
	}

	t, ok := field.Interface().(time.Time)
	if !ok {
		return false
	}

	return !t.IsZero()
}

func vxgouuid(fl validator.FieldLevel) bool {
	field := fl.Field()

	t, ok := field.Interface().(uuid.UUID)
	if !ok {
		return false
	}

	return t != uuid.Nil
}

var vxdateregexstr = "^(?:[1-9]\\d{3}-(?:(?:0[1-9]|1[0-2])-(?:0[1-9]|1\\d|2[0-8])|(?:0[13-9]|1[0-2])-(?:29|30)|(?:0[13578]|1[02])-31)|(?:[1-9]\\d(?:0[48]|[2468][048]|[13579][26])|(?:[2468][048]|[13579][26])00)-02-29)$"
var vxdateregex = regexp.MustCompile(vxdateregexstr)

func vxdate(fl validator.FieldLevel) bool {
	return vxdateregex.MatchString(fl.Field().String())
}

var vxtimeregexstr = `^(?:[01]\d|2[0-3]):[0-5]\d:[0-5]\d(?:\.\d{1,9})?$`
var vxtimeregex = regexp.MustCompile(vxtimeregexstr)

func vxtime(fl validator.FieldLevel) bool {
	return vxtimeregex.MatchString(fl.Field().String())
}

var vxalphaspaceregexstr = `^[a-zA-Z\s]+$`
var vxalphaspaceregex = regexp.MustCompile(vxalphaspaceregexstr)

func vxalphaspace(fl validator.FieldLevel) bool {
	s := fl.Field().String()
	if len(s) == 0 {
		return false
	}
	if s[0] == ' ' || s[len(s)-1] == ' ' {
		return false
	}

	return vxalphaspaceregex.MatchString(fl.Field().String())
}

var vxalphaspacethregexstr = `^[a-zA-Zก-๏\s]+$`
var vxalphaspacethregex = regexp.MustCompile(vxalphaspacethregexstr)

func vxalphaspaceth(fl validator.FieldLevel) bool {
	s := fl.Field().String()
	if len(s) == 0 {
		return false
	}
	if s[0] == ' ' || s[len(s)-1] == ' ' {
		return false
	}

	return vxalphaspacethregex.MatchString(fl.Field().String())
}

var vxalphathregexstr = `^[a-zA-Zก-๏]+$`
var vxalphathregex = regexp.MustCompile(vxalphathregexstr)

// Ref regex https://www.ninenik.com/%E0%B9%81%E0%B8%99%E0%B8%A7%E0%B8%97%E0%B8%B2%E0%B8%87%E0%B8%95%E0%B8%A3%E0%B8%A7%E0%B8%88%E0%B8%82%E0%B9%89%E0%B8%AD%E0%B8%A1%E0%B8%B9%E0%B8%A5%E0%B9%80%E0%B8%89%E0%B8%9E%E0%B8%B2%E0%B8%B0%E0%B8%A0%E0%B8%B2%E0%B8%A9%E0%B8%B2%E0%B9%84%E0%B8%97%E0%B8%A2_%E0%B8%A0%E0%B8%B2%E0%B8%A9%E0%B8%B2%E0%B8%AD%E0%B8%B1%E0%B8%87%E0%B8%81%E0%B8%A4%E0%B8%A9%E0%B8%94%E0%B9%89%E0%B8%A7%E0%B8%A2_Regular_Expression-877.html
func vxalphath(fl validator.FieldLevel) bool {
	return vxalphathregex.MatchString(fl.Field().String())
}

func vxuniquefield(fl validator.FieldLevel) bool {
	slice := fl.Field()

	if slice.Kind() != reflect.Slice && slice.Kind() != reflect.Array {
		return true
	}

	if slice.Len() == 0 {
		return true
	}

	elemType := slice.Index(0).Kind()
	if elemType != reflect.Struct {
		return false
	}

	params := strings.Split(fl.Param(), " ")
	seen := make(map[string]bool)

	for _, param := range params {
		for i := 0; i < slice.Len(); i++ {
			item := slice.Index(i)
			field := item.FieldByName(param)

			if !field.IsValid() {
				return false
			}

			value := field.String()
			if seen[value] {
				return false
			}
			seen[value] = true
		}
	}

	return true
}

func vxsplitoneof(fl validator.FieldLevel) bool {
	f := fl.Field()
	v, ok := f.Interface().(string)
	if !ok {
		return false
	}

	optionStr := fl.Param()
	options := strings.Split(optionStr, " ")
	v2 := strings.Split(v, ",")

	for _, v3 := range v2 {
		if slices.Contains(options, v3) {
			return true
		}
	}

	return false
}

func vxallornone(fl validator.FieldLevel) bool {
	p := fl.Parent()

	optionStr := fl.Param()
	options := strings.Split(optionStr, " ")

	var fieldSet []bool

	for _, option := range options {
		f := p.FieldByName(option)
		if !f.IsValid() {
			return false
		}

		if !isemptyvalue(f.Interface()) {
			fieldSet = append(fieldSet, true)
		} else {
			fieldSet = append(fieldSet, false)
		}

	}

	allSet := true
	noneSet := true
	for _, set := range fieldSet {
		if !set {
			allSet = false
		}
		if set {
			noneSet = false
		}
	}

	return allSet || noneSet
}

func isemptyvalue(value any) bool {
	v := reflect.ValueOf(value)

	switch v.Kind() {
	case reflect.String:
		return v.String() == ""
	case reflect.Slice, reflect.Array:
		return v.Len() == 0
	case reflect.Map:
		return v.Len() == 0
	case reflect.Ptr:
		return v.IsNil()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	default:
		return false
	}
}

func vxsplituuid(fl validator.FieldLevel) bool {
	f := fl.Field()
	v, ok := f.Interface().(string)
	if !ok {
		return false
	}

	v2 := strings.Split(v, ",")

	for _, v3 := range v2 {
		if _, ok := uuid.Parse(v3); ok != nil {
			return false
		}
	}

	return true
}

var vxe164regexstr = `^\+[1-9]\d{1,14}$`
var vxe164regex = regexp.MustCompile(vxe164regexstr)

var vxe164localregexstr = `^0\d{9}$`
var vxe164localregex = regexp.MustCompile(vxe164localregexstr)

func vxe164(fl validator.FieldLevel) bool {
	phone := fl.Field().String()
	if vxe164regex.MatchString(phone) {
		return true
	}
	return vxe164localregex.MatchString(phone)
}

var vxdecimalregexstr = `^-?\d+(\.\d{1,2})?$`
var vxdecimalregex = regexp.MustCompile(vxdecimalregexstr)

func vxdecimal(fl validator.FieldLevel) bool {
	f := fl.Field()

	if v, ok := f.Interface().(decimal.Decimal); ok {
		s := v.String()
		return vxdecimalregex.MatchString(s)
	}

	return false
}

var vxudecimalregexstr = `^-?\d+(\.\d*)?$`
var vxudecimalregex = regexp.MustCompile(vxudecimalregexstr)

func vxudecimal(fl validator.FieldLevel) bool {
	f := fl.Field()

	if v, ok := f.Interface().(decimal.Decimal); ok {
		return vxudecimalregex.MatchString(v.String())
	}

	return false
}

func vxdeq(fl validator.FieldLevel) bool {
	f := fl.Field()
	fparams := fl.Param()

	v2, err := decimal.NewFromString(fparams)
	if err != nil {
		return false
	}

	if v, ok := f.Interface().(decimal.Decimal); ok {
		return v.Equal(v2)
	}

	return false
}

func vxdgt(fl validator.FieldLevel) bool {
	f := fl.Field()
	fparams := fl.Param()

	v2, err := decimal.NewFromString(fparams)
	if err != nil {
		return false
	}

	if v, ok := f.Interface().(decimal.Decimal); ok {
		return v.GreaterThan(v2)
	}

	return false
}

func vxdgte(fl validator.FieldLevel) bool {
	f := fl.Field()
	fparams := fl.Param()

	v2, err := decimal.NewFromString(fparams)
	if err != nil {
		return false
	}

	if v, ok := f.Interface().(decimal.Decimal); ok {
		return v.GreaterThanOrEqual(v2)
	}

	return false
}

func vxdlt(fl validator.FieldLevel) bool {
	f := fl.Field()
	fparams := fl.Param()

	v2, err := decimal.NewFromString(fparams)
	if err != nil {
		return false
	}

	if v, ok := f.Interface().(decimal.Decimal); ok {
		return v.LessThan(v2)
	}

	return false
}

func vxdlte(fl validator.FieldLevel) bool {
	f := fl.Field()
	fparams := fl.Param()

	v2, err := decimal.NewFromString(fparams)
	if err != nil {
		return false
	}

	if v, ok := f.Interface().(decimal.Decimal); ok {
		return v.LessThanOrEqual(v2)
	}

	return false
}

func vxdne(fl validator.FieldLevel) bool {
	f := fl.Field()
	fparams := fl.Param()

	v2, err := decimal.NewFromString(fparams)
	if err != nil {
		return false
	}

	if v, ok := f.Interface().(decimal.Decimal); ok {
		return !v.Equal(v2)
	}

	return false
}

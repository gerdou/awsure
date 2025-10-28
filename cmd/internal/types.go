package internal

import (
	"crypto/sha512"
	"encoding/hex"
	"encoding/xml"
	"fmt"
	"reflect"
	"time"
)

var configFileVersion = "1.0.0"

type configuration struct {
	AzureTenantId        string `yaml:"azure_tenant_id"`
	AzureAppIdUri        string `yaml:"azure_app_id_uri"`
	AzureUsername        string `yaml:"azure_username"`
	OktaUsername         string `yaml:"okta_username"`
	RememberMe           bool   `yaml:"remember_me"`
	DefaultJumpRole      string `yaml:"default_jump_role"`
	DestinationAccountId string `yaml:"destination_account_id"`
	DestinationRoleName  string `yaml:"destination_role_name"`
	DefaultDurationHours int    `yaml:"default_duration_hours"`
	Region               string `yaml:"region"`
}

func (c *configuration) Hash() string {
	input := fmt.Sprintf("%s|%s|%s|%s",
		c.AzureTenantId,
		c.AzureAppIdUri,
		c.AzureUsername,
		c.OktaUsername)

	hash := sha512.New()
	hash.Write([]byte(input))
	hashBytes := hash.Sum(nil)

	return hex.EncodeToString(hashBytes)
}

func (c *configuration) Merge(other *configuration) {
	cValue := reflect.ValueOf(c).Elem()
	otherValue := reflect.ValueOf(other).Elem()

	for i := 0; i < cValue.NumField(); i++ {
		cField := cValue.Field(i)
		otherField := otherValue.Field(i)

		switch cField.Kind() {
		case reflect.String:
			if otherField.String() != "" {
				cField.SetString(otherField.String())
			}
		case reflect.Bool:
			if otherField.Bool() {
				cField.SetBool(otherField.Bool())
			}
		case reflect.Int:
			if otherField.Int() >= 1 && otherField.Int() <= 12 {
				cField.SetInt(otherField.Int())
			}
		case reflect.Invalid:
		case reflect.Int8:
		case reflect.Int16:
		case reflect.Int32:
		case reflect.Int64:
		case reflect.Uint:
		case reflect.Uint8:
		case reflect.Uint16:
		case reflect.Uint32:
		case reflect.Uint64:
		case reflect.Uintptr:
		case reflect.Float32:
		case reflect.Float64:
		case reflect.Complex64:
		case reflect.Complex128:
		case reflect.Array:
		case reflect.Chan:
		case reflect.Func:
		case reflect.Interface:
		case reflect.Map:
		case reflect.Pointer:
		case reflect.Slice:
		case reflect.Struct:
		case reflect.UnsafePointer:
		}
	}
}

type configurationFile struct {
	Version string                    `yaml:"version"`
	Configs map[string]*configuration `yaml:"configs"`
}

type jumpRoleCredentials struct {
	AwsAccessKeyId     string    `yaml:"aws_access_key_id"`
	AwsSecretAccessKey string    `yaml:"aws_secret_access_key"`
	AwsSessionToken    string    `yaml:"aws_session_token"`
	AwsExpiration      time.Time `yaml:"aws_expiration"`
}

type jumpRoleCredentialsFile struct {
	Version     string                          `yaml:"version"`
	Credentials map[string]*jumpRoleCredentials `yaml:"credentials"`
}

type samlResponse struct {
	XMLName   xml.Name
	Assertion samlAssertion `xml:"Assertion"`
}

type samlAssertion struct {
	XMLName            xml.Name
	AttributeStatement samlAttributeStatement
}

type samlAttributeValue struct {
	XMLName xml.Name
	Type    string `xml:"xsi:type,attr"`
	Value   string `xml:",innerxml"`
}

type samlAttribute struct {
	XMLName         xml.Name
	Name            string               `xml:",attr"`
	AttributeValues []samlAttributeValue `xml:"AttributeValue"`
}

type samlAttributeStatement struct {
	XMLName    xml.Name
	Attributes []samlAttribute `xml:"Attribute"`
}

type role struct {
	roleArn      string
	principalArn string
}

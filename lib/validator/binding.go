package validator

import (
	"github.com/gin-gonic/gin/binding"
	"reflect"
	"sync"
)

//替换gin的binding
func Init() {
	valid:= new(DefaultValidator)
	binding.Validator = valid
}


type DefaultValidator struct {
	once     sync.Once
	validate *myValidator
}

var _ binding.StructValidator = &DefaultValidator{}

func (v *DefaultValidator) ValidateStruct(obj interface{}) error {

	if kindOfData(obj) == reflect.Struct {
		v.lazyinit()
		if err := v.validate.Validate(obj); err != nil {
			return err[0]
		}
	}

	return nil
}

func (v *DefaultValidator) Engine() interface{} {
	v.lazyinit()
	return v.validate
}

func (v *DefaultValidator) lazyinit() {
	v.once.Do(func() {
		v.validate = New()
		v.validate.SetTag("binding")
		v.validate.SetValidatorSplit("|")
		// add any custom validations etc. here
	})
}

func kindOfData(data interface{}) reflect.Kind {

	value := reflect.ValueOf(data)
	valueType := value.Kind()

	if valueType == reflect.Ptr {
		valueType = value.Elem().Kind()
	}
	return valueType
}

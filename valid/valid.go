package valid

import (
	"fmt"
    "github.com/go-playground/validator/v10"	// validator验证库的依赖
	"github.com/gin-gonic/gin/binding" // gin框架的绑定，用于获取gin的validator引擎
	"github.com/go-playground/locales/en"	// 英语语言包的资源
	"github.com/go-playground/locales/zh"	// 中文语言包的资源
	ut "github.com/go-playground/universal-translator"	// 翻译器
	zhtrans "github.com/go-playground/validator/v10/translations/zh"	// 汉语翻译包
    "reflect"	// 反射相关包
	"strings"	// 字符串操作相关包
)



// 两个全局变量
var Trans ut.Translator
var Validate *validator.Validate

// 初始化 vaidator引擎：
func InitValidator() {
    // 获得gin框架的validator引擎
    Validate, _ = binding.Validator.Engine().(*validator.Validate)
    // 注册GetTag，获取参数结构体中的tag
    Validate.RegisterTagNameFunc(GetTag)
    // 初始化翻译器（固定为汉语翻译包）
    InitTrans("zh")
    
	//--------------------------------------------------------------------

    // 这里是注册翻译函数，这里随意，但是注意：
    // 1.如果结构体中binding的tag已经可以校验的规则，就不需要自定义函数了，自然也不需要自定义翻译函数
    // 2.如果结构体中tag不足以全部校验规则时，要自定义校验函数，必须 还要自定义翻译函数
    // 自定义校验函数和自定义翻译函数： “username”作为tag可以直接使用了
	Validate.RegisterValidation("username", ValidateUsername)
    Validate.RegisterTranslation("username", Trans, tranUserNameRegis, tranUserName)
	
	//--------------------------------------------------------------------
}

//获取参数结构体中的tag
func GetTag(fld reflect.StructField) string {
    name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
    if name == "-" {
        return ""
    }
    return name
}

//--------------------------------------------------------------------
//如果你有自定义函数：
// 自定义校验函数
func ValidateUsername(fl validator.FieldLevel) bool {
    if fl.Field().String() == "111" {
        return true
    }
    return false
}

// 自定义校验函数对应的注册翻译函数
func tranUserNameRegis(ut ut.Translator) error {
	// 添加特定标签和翻译的错误信息
    return ut.Add("username", "{0}非法的用户名!", true)
}

// 自定义校验函数对应的翻译函数
func tranUserName(ut ut.Translator, fe validator.FieldError) string {
	// 拼接错误信息后返回
    t, _ := ut.T("username", fe.Field())
    return t
}
// -----------------------------------------------------------------


// 初始化翻译器：根据传入的方言选择翻译方向
func InitTrans(locale string) (err error) {
    en := en.New()
    zh := zh.New()
    uni := ut.New(en, zh)
	var ok bool
    // 注意：这里一定要给全局的trans，否则空指针报错
    Trans, ok = uni.GetTranslator(locale)
    if!ok {
        return fmt.Errorf("无法获取 %s 语言对应的翻译器")
    }
    switch locale {
    case "zh":
        err = zhtrans.RegisterDefaultTranslations(Validate, Trans)
    default:
        return fmt.Errorf("不支持的语言环境: %s")
    }
    return
}

// 另外的作用（可选）：用于删除错误信息前面的结构体信息
func RemoveStructName(fields map[string]string) map[string]string {
    result := map[string]string{}
    for field, err := range fields {
        result[field[strings.Index(field, ".")+1:]] = err
    }
    return result
}

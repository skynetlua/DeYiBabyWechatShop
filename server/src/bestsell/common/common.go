package common

import (
	"crypto/md5"
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

func init() {
	//timer1 := time.NewTicker(2 * time.Second)
	//go func(t *time.Ticker) {
	//	for {
	//		<-t.C
	//		MakeResUrl("/Users/linyou/svn/plan1/server/static/api/qrcode/20200430")
	//		break
	//	}
	//}(timer1)

}

func AtoI(val string)int {
	v,err := strconv.Atoi(val)
	if err != nil {
		fmt.Println("AtoI error:",err)
		return 0
	}
	return v
}

func AtoIDefault(val string, def int)int {
	v,err := strconv.Atoi(val)
	if err != nil {
		fmt.Println("AtoI error:", err)
		return def
	}
	return v
}

func AtoF(val string)float64 {
	ret, ok := strconv.ParseFloat(val, 64)
	if ok != nil {
		ret = 0.0
	}
	return ret
}

func MakeMd5(data string)string  {
	has := md5.Sum([]byte(data))
	md5str := fmt.Sprintf("%x", has)
	return md5str
}

func MakeResUrl(rpath string)string  {
	rpath = strings.Replace(rpath, "\\", "/", -1)
	idx := strings.Index(rpath, "/static")
	if idx<0 {
		//idx = strings.Index(rpath, "\\static\\")
		//if idx<0 {
			return rpath
		//}
	}
	rpath = rpath[idx:]
	_url := StaticUrl+rpath
	return _url
}

func ConvertString(val interface{})string {
	switch val.(type) {
	case string:
		return val.(string)
	}
	return ""
}

func ConvertInt(val interface{})int{
	switch val.(type) {
	case int:
		return val.(int)
	case int32:
		return int(val.(int32))
	case int64:
		return int(val.(int64))
	case float64:
		return int(val.(float64))
	case float32:
		return int(val.(float32))
	//case string:
	//	v,err := strconv.Atoi(val.(string))
	//	if err != nil {
	//		fmt.Println("ConvertInt err:",err)
	//		return 0
	//	}
	//	return v
	}
	return 0
}

// func ConvertIntArray(val interface{})[]int{
// 	switch val.(type) {
// 		case primitive.A:
// 			var ret []int
// 			for _, v := range val.(primitive.A) {
// 				ret = append(ret, ConvertInt(v))
// 			}
// 			return ret
// 		case []interface{}:
// 			var ret []int
// 			for _, v := range val.([]interface{}) {
// 				ret = append(ret, ConvertInt(v))
// 			}
// 			return ret
// 	}
// 	return val.([]int)
// }

// func ConvertMap(val interface{})map[string]interface{}{
// 	switch val.(type) {
// 	case primitive.M:
// 		ret := make(map[string]interface{})
// 		for k, v := range val.(primitive.M) {
// 			ret[k] = v
// 		}
// 		return ret
// 	}
// 	return val.(map[string]interface{})
// }

// func ConvertArray(val interface{})[]interface{}{
// 	switch val.(type) {
// 	case primitive.A:
// 		var ret []interface{}
// 		for _, v := range val.(primitive.A) {
// 			ret = append(ret, v)
// 		}
// 		return ret
// 	case []map[string]interface{}:
// 		var ret []interface{}
// 		for _, v := range val.([]map[string]interface{}) {
// 			ret = append(ret, v)
// 		}
// 		return ret
// 	}
// 	return val.([]interface{})
// }

func ConvertFloat64(val interface{})float64{
	switch val.(type) {
	case float64:
		return val.(float64)
	case float32:
		return float64(val.(float32))
	case int32:
		return float64(val.(int32))
	case int64:
		return float64(val.(int64))
	}
	return 0.0
}

func MakeMoneyValue(val float64) float64{
	return float64(int64(val*10)) / 10
}

// func MakeMoneyValue1(val float64) float64{
// 	return float64(int64(val*10)) / 10
// }

func ForceParseJson(txtJson string)*map[string]string  {
	fmt.Println("txtJson:", txtJson)
	items := strings.Split(txtJson, ",")
	var retOjb = make(map[string]string)
	for _,item := range items {
		tmp := strings.Split(item, ":")
		key := tmp[0]
		val := tmp[1]
		if len(tmp) > 2 {
			if strings.Index(tmp[1], "{") >= 0 {
				key = tmp[1]
				val = tmp[2]
			}else{
				val = tmp[1]+":"+tmp[2]
			}
		}
		sIdx := strings.Index(key, "\"")
		key = key[sIdx+1:]
		eIdx := strings.Index(key, "\"")
		key = key[:eIdx]

		sIdx = strings.Index(val, "\"")
		if sIdx >= 0 {
			val = val[sIdx+1:]
			eIdx = strings.Index(val, "\"")
			val = val[:eIdx]
		}
		retOjb[key] = val
	}
	return &retOjb
}

func CheckHan(testStr *string) bool{
	for _, v := range *testStr {
		if unicode.Is(unicode.Han, v) {
			return true
		}
	}
	return false
}
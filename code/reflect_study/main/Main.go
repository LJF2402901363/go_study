package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"reflect"
	"strconv"
	"strings"
)

type MysqlConfig struct {
	Url      string `ini:"url"`
	Port     int    `ini:"port"`
	UserName string `ini:"userName"`
	Password string `ini:"password"`
}
type RedisConfig struct {
	Host     string `ini:"host"`
	Port     int    `ini:"port"`
	UserName string `ini:"userName"`
	Password string `ini:"password"`
	Database int    `ini:"database"`
}
type Config struct {
	MysqlConfig `ini:"mysql"`
	RedisConfig `ini:"redis"`
}

func reflectConfig(fileName string, configStruct interface{}) (err error) {
	//1.参数的校验,传入的data对象必须是指针
	dataType := reflect.TypeOf(configStruct)
	if dataType.Kind() != reflect.Ptr {
		err = errors.New("传入的data类型不是指针！")
		return err
	}
	//2.传进来的configStruct必须是结构体
	if dataType.Elem().Kind() != reflect.Struct {
		err = errors.New("传入的data类型不是结构体！")
		return err
	}

	//读取文件
	fileBytes, err := ioutil.ReadFile(fileName)
	if err != nil {
		return err
	}
	//读取文件的所有行
	fileLines := string(fileBytes)
	contentLines := strings.Split(fileLines, "\r\n")
	var structName string
	for index, line := range contentLines {
		line = strings.TrimSpace(line)
		//这行是注释
		if strings.HasPrefix(line, ";") || strings.HasPrefix(line, "#") {
			continue
		}
		if strings.HasPrefix(line, "[") {
			if strings.HasSuffix(line, "]") {
				word := line[1 : len(line)-1]
				if len(word) == 0 {
					return errors.New(fmt.Sprintf("第%d行格式错误:%s", index, line))
				}
				//根据configStruct获取对应的结构体
				//value := reflect.ValueOf(configStruct)
				for i := 0; i < dataType.Elem().NumField(); i++ {
					field := dataType.Elem().Field(i)
					if field.Tag.Get("ini") == word {
						structName = field.Name
						fmt.Printf("找到%s中的嵌套结构体%s", word, structName)
					}
				}
				fmt.Println(word)
			} else {
				return errors.New(fmt.Sprintf("第%d行格式错误:%s", index, line))
			}
		} else {
			//每一行键值对
			keyValueArr := strings.Split(line, "=")
			if len(keyValueArr) != 2 {
				return errors.New(fmt.Sprintf("第%d行格式错误:%s", index, line))
			}
			key := keyValueArr[0]
			val := keyValueArr[1]
			//获取结构体
			dataVal := reflect.ValueOf(configStruct)
			structObj := dataVal.Elem().FieldByName(structName)
			structType := structObj.Type()
			if structType.Kind() != reflect.Struct {
				err = errors.New(fmt.Sprintf("传入的%s类型不是结构体！", structName))
				return err
			}
			var fieldName string
			var fieldType reflect.StructField
			for i := 0; i < structObj.NumField(); i++ {
				field := structType.Field(i)
				fieldType = field
				if field.Tag.Get("ini") == key {
					fieldName = field.Name
					break
				}

			}
			fieldByName := structObj.FieldByName(fieldName)
			switch fieldType.Type.Kind() {
			case reflect.String:
				fieldByName.SetString(val)
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				parseInt, err := strconv.ParseInt(val, 10, 64)
				if err != nil {
					return errors.New(fmt.Sprintf("第%d行出错了，key:%s的值%s不是合法的int类型", index, key, val))
				}
				fieldByName.SetInt(parseInt)
			case reflect.Bool:
				parseInt, err := strconv.ParseBool(val)
				if err != nil {
					return errors.New(fmt.Sprintf("第%d行出错了，key:%s的值%s不是合法的bool类型", index, key, val))
				}
				fieldByName.SetBool(parseInt)
			case reflect.Float32, reflect.Float64:
				parseInt, err := strconv.ParseFloat(val, 64)
				if err != nil {
					return errors.New(fmt.Sprintf("第%d行出错了，key:%s的值%s不是合法的float类型", index, key, val))
				}
				fieldByName.SetFloat(parseInt)
			}
		}

	}
	return
}
func main() {
	var mc Config
	fileName := "reflect_study/main/config.ini"
	err := reflectConfig(fileName, &mc)
	if err != nil {
		fmt.Printf("load %s failed,%s", fileName, err)
	}
	fmt.Println(mc)

}

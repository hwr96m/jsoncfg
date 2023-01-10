package config

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

//------------ Переменные ----------------------------------------------------
var ()

//------------ Структуры ----------------------------------------------------

//------------ Функции ----------------------------------------------------

//Открытие файла конфигурации и запись данных в структуру.
//file - путь к файлу с настройками, configStruct - указатель на экземпляр структуры данных.
func ParseFile(file string, configStruct interface{}) error {
	var (
		err error
		f   *os.File
	)
	f, err = os.Open(file) //открываем файл
	if err != nil {
		return fmt.Errorf("func ParseFile: %s", err)
	}
	defer f.Close()
	err = json.NewDecoder(f).Decode(configStruct) //извлечение данных из json и запись в configStruct
	if err != nil {
		return fmt.Errorf("func ParseFile: json.NewDecoder(f).Decode: %s", err)
	}
	return nil
}

//Открытие файла конфигурации и запись данных в структуру.
//file - путь к файлу с настройками, configStruct - указатель на экземпляр структуры данных.
//Если файл не найден, то создается новый с данными из defaultFile.
func ParseFileOrDefault(file, defaultFile string, configStruct interface{}) error {
	if _, err := os.Stat(file); os.IsNotExist(err) { //если файла нет, то копируем defaultFile
		err = copyFile(file, defaultFile)
		if err != nil {
			return err
		}
	}
	return ParseFile(file, configStruct)
}

func copyFile(dst, src string) error {
	fsrc, err := os.Open(src)
	if err != nil {
		return err
	}
	defer fsrc.Close()

	fdst, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer fdst.Close()

	_, err = io.Copy(fdst, fsrc)
	if err != nil {
		return err
	}
	return nil
}

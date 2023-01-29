package jsoncfg

import (
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
)

// ------------ Переменные ----------------------------------------------------
var ()

//------------ Структуры ----------------------------------------------------

//------------ Функции ----------------------------------------------------

// Открытие файла конфигурации и запись данных в структуру.
// file - путь к файлу с настройками, configStructs - указатели на экземпляр структуры данных.
// Можно указать несколько структур
// Поддерживает файлы .jsonc
func DecodeFile(file string, configStructs ...interface{}) error {
	var (
		err error
		f   *os.File
	)
	f, err = os.OpenFile(file, os.O_RDONLY, 0440) //открываем файл
	if err != nil {
		return fmt.Errorf("func DecodeFile: ошибка открытия файла: %s", err)
	}
	defer f.Close()

	if filepath.Ext(file) == ".jsonc" { //если на вход получили .jsonc, то удаляем комментарии
		js, err := io.ReadAll(f)
		if err != nil {
			return fmt.Errorf("func DecodeFile(): io.ReadAll: %w", err)
		}
		re, err := regexp.Compile(`//.*[^\r\n]`)
		if err != nil {
			return fmt.Errorf("func DecodeFile(): regexp.Compile: %w", err)
		}
		js = re.ReplaceAll(js, []byte{})
		for _, configStruct := range configStructs {
			err = json.Unmarshal(js, configStruct) //извлечение данных из json и запись в configStruct
			if err != nil {
				return fmt.Errorf("func DecodeFile: Unmarshal: %s", err)
			}
		}
	} else {
		for _, configStruct := range configStructs {
			err = json.NewDecoder(f).Decode(configStruct) //извлечение данных из json и запись в configStruct
			if err != nil {
				return fmt.Errorf("func DecodeFile: Decode: %s", err)
			}
		}
	}
	return nil
}

// Открытие файла конфигурации и запись данных в структуру.
// file - путь к файлу с настройками, configStruct - указатель на экземпляр структуры данных.
// Если файл не найден, то создается новый со структурой данных configStructs.
func DecodeFileOrCreate(file string, configStruct interface{}) error {
	fi, err := os.Stat(file)
	switch {
	case os.IsNotExist(err): //если не существует
		err = EncodeFile(file, configStruct, 0600)
		if err != nil {
			return fmt.Errorf("func DecodeFileOrDefault(): создание файла: %+v", err)
		}
		return fmt.Errorf("файл не найден, был создан новый")
	case fi.IsDir(): //если указан не файл, а папка
		return fmt.Errorf("func DecodeFileOrDefault(): вместо файла указан путь до папки: %+v", err)
	default:
		return DecodeFile(file, configStruct)
	}
}

// Запись структуры в файл с форматированием
func EncodeFile(file string, configStruct interface{}, perm fs.FileMode) error {
	//сохраняем информацию о сессии в файл
	f, err := os.OpenFile(file, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, perm)
	if err != nil {
		return fmt.Errorf("func EncodeFile(): open file %s: %s", file, err)
	}
	defer f.Close()
	j, err := json.MarshalIndent(configStruct, "  ", "  ") //создаем json с форматированием
	if err != nil {
		return fmt.Errorf("func EncodeFile(): MarshalIndent(): %s", err)
	}
	f.Write(j)
	return nil
}

// Запись структуры в файл, минифицированная версия
func EncodeFileMinify(file string, configStruct interface{}, perm fs.FileMode) error {
	//сохраняем информацию о сессии в файл
	f, err := os.OpenFile(file, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, perm)
	if err != nil {
		return fmt.Errorf("func EncodeFileMinify(): open file %s: %s", file, err)
	}
	defer f.Close()
	err = json.NewEncoder(f).Encode(configStruct)
	if err != nil {
		return fmt.Errorf("func EncodeFileMinify(): Encode(): %s", err)
	}
	return nil
}

//------------ Локальные функции ---------------------------------------------------------

/* func copyFile(dst, src string) error {
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
}*/

package emptyfile

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func ReCreateFile(par_folder, name_ext string) *os.File {
	// Get executable path
	var file *os.File

	execPath, err := os.Executable()
	if err != nil {
		log.Fatal("Ошибка получения пути исполняемого файла:", err)
		return nil
	}

	// Get executable dir
	execDir := filepath.Dir(execPath)

	// Check exist log dir
	_, err = os.Stat(filepath.Join(execDir, par_folder))
	if os.IsNotExist(err) {
		err = os.Mkdir(filepath.Join(execDir, par_folder), 0755)
		if err != nil {
			log.Fatal(fmt.Sprintf("Не получилось создать директорию логов '%s':", filepath.Join(execDir, par_folder)), err)
			return nil
		}
	}

	// Set config path
	filename := filepath.Join(execDir, par_folder, name_ext)

	// Create file
	// Delete previous file if it exists
	// ignore log pinmsg file
	if name_ext != "WednesdayPollPinMsg.txt" {
		if _, err := os.Stat(filename); !os.IsNotExist(err) {
			err = os.Remove(filename)
			if err != nil {
				log.Fatal(err, fmt.Sprintf("Не получилось удалить файл лога '%s':", filename))
				return nil
			}
		}
	}
	file, err = os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err, fmt.Sprintf("Не получилось создать файл лога '%s':", filename))
		return nil
	}

	log.Printf((fmt.Sprintf("Файл лога '%s' создан", filename)))

	return file
}

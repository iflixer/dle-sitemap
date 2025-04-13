package helper

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

func CopyDir(src string, dst string) error {
	// Проверяем, существует ли источник
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}
	if !srcInfo.IsDir() {
		return fmt.Errorf("источник не папка: %s", src)
	}

	// Создаём целевую папку
	err = os.MkdirAll(dst, srcInfo.Mode())
	if err != nil {
		return err
	}

	// Проходим по всем файлам и поддиректориям
	return filepath.WalkDir(src, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Новый путь назначения
		relPath, _ := filepath.Rel(src, path)
		targetPath := filepath.Join(dst, relPath)

		// Копируем директории
		if d.IsDir() {
			return os.MkdirAll(targetPath, 0755)
		}

		// Копируем файл
		return copyFile(path, targetPath)
	})
}

func copyFile(srcFile, dstFile string) error {
	src, err := os.Open(srcFile)
	if err != nil {
		return err
	}
	defer src.Close()

	dst, err := os.Create(dstFile)
	if err != nil {
		return err
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)
	if err != nil {
		return err
	}

	// Копируем права
	srcInfo, err := os.Stat(srcFile)
	if err != nil {
		return err
	}
	return os.Chmod(dstFile, srcInfo.Mode())
}

package system

import (
	"archive/zip"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func CheckIsDir(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	return fileInfo.IsDir(), nil
}

func CopyDir(srcDir, dstDir string) error {
	return filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		relPath, err := filepath.Rel(srcDir, path)
		if err != nil {
			return err
		}
		dstPath := filepath.Join(dstDir, relPath)
		if info.IsDir() {
			return os.MkdirAll(dstPath, os.ModePerm)
		}
		return CopyFile(path, dstPath)
	})
}

func CopyFile(srcFile, destFile string) error {
	input, err := os.Open(srcFile)
	if err != nil {
		return err
	}
	defer input.Close()

	output, err := os.Create(destFile)
	if err != nil {
		return err
	}
	defer output.Close()

	_, err = io.Copy(output, input)
	return err
}

func ZipDir(srcDir, zipFilePath string) error {
	zipFile, err := os.Create(zipFilePath)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	archive := zip.NewWriter(zipFile)
	defer archive.Close()

	err = filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(srcDir, path)
		if err != nil {
			return err
		}

		header.Name = strings.ReplaceAll(relPath, string(os.PathSeparator), "/")

		if info.IsDir() {
			header.Name += "/"
		} else {
			header.Method = zip.Deflate
		}

		writer, err := archive.CreateHeader(header)
		if err != nil {
			return err
		}

		if !info.IsDir() {
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()
			if _, err := io.Copy(writer, file); err != nil {
				file.Close()
				return err
			}
		}
		return nil
	})
	return err
}

func UnzipDir(zipFile, destDir string) error {
	reader, err := zip.OpenReader(zipFile)
	if err != nil {
		return err
	}
	defer reader.Close()

	for _, file := range reader.File {
		path := filepath.Join(destDir, file.Name)

		if file.FileInfo().IsDir() {
			if err := os.MkdirAll(path, file.Mode()); err != nil {
				return err
			}
		} else {
			if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
				return err
			}

			outFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
			if err != nil {
				return err
			}
			defer outFile.Close()

			rc, err := file.Open()
			if err != nil {
				return err
			}
			defer rc.Close()

			if _, err := io.Copy(outFile, rc); err != nil {
				return err
			}
		}
	}
	return nil
}

func CleanAndCreateDir(dirPath string) error {
	if _, err := os.Stat(dirPath); !os.IsNotExist(err) {
		if err := os.RemoveAll(dirPath); err != nil {
			return err
		}
	}
	return os.MkdirAll(dirPath, 0755)
}

// CheckAndCreateDir 检查指定路径的文件夹是否存在，如果不存在则创建它。
func CheckAndCreateDir(path string) error {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		err := os.MkdirAll(path, 0755)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}
	return nil
}

func GetSavDir(path string) (string, error) {
	var levelFilePath string
	levelFilePath, err := GetLevelSavFilePath(path)
	if err != nil {
		return "", err
	}
	return filepath.Dir(levelFilePath), nil
}

func GetLevelSavFilePath(path string) (string, error) {
	var foundPath string
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.Name() == "Level.sav" {
			foundPath = path
			return errors.New("file found")
		}
		return nil
	})

	if foundPath == "" {
		if err != nil && !errors.Is(err, errors.New("file found")) {
			return "", err
		}
		return "", errors.New("file Level.sav not found")
	}

	return foundPath, nil
}

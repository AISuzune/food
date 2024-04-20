package file

import "os"

// CheckNotExist 检查指定文件是否存在
func CheckNotExist(src string) bool {
	_, err := os.Stat(src) // 获取文件或目录的状态信息

	return os.IsNotExist(err) // 如果错误类型是文件或目录不存在，则返回true，否则返回false
}

// MkDir 建立文件夹
func MkDir(src string) error {
	err := os.MkdirAll(src, os.ModePerm) // 创建目录，os.ModePerm表示文件权限模式
	if err != nil {
		return err // 如果创建目录失败，返回错误
	}

	return nil // 如果创建目录成功，返回nil
}

// IsNotExistMkDir 检查文件夹是否存在, 不存在则创建
func IsNotExistMkDir(src string) error {
	if notExist := CheckNotExist(src); notExist == true { // 检查目录是否不存在
		if err := MkDir(src); err != nil { // 如果目录不存在，创建目录
			return err // 如果创建目录失败，返回错误
		}
	}

	return nil // 如果目录存在或创建成功，返回nil
}

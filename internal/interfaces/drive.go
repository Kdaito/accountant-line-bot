package interfaces

import "os"

type DriveInterface interface {
	Upload(parentId string, title string, file *os.File) (string, error)
}

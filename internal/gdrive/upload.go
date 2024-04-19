package gdrive

import (
	"os"

	"google.golang.org/api/drive/v2"
)

type GDrive struct {
	Service *drive.Service
}

func (g *GDrive) Upload(parentId string, title string, file *os.File) (string, error) {
	f := &drive.File{
		Title:   title,
		Parents: []*drive.ParentReference{{Id: parentId}},
	}

	r, err := g.Service.Files.Insert(f).Media(file).Do()
	if err != nil {
		return "", err
	}

	return r.Id, nil
}

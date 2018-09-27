package models

import (
	"io"
	"mime/multipart"
	"os/exec"
	"spyc_backend/app/core"
	"strings"

	"github.com/Reti500/mgomap"
	mgo "gopkg.in/mgo.v2"
)

// File ...
type File struct {
	mgomap.DocumentBase `json:",inline" bson:",inline"`

	ACL          string     `json:"-" bson:"acl"`
	Action       string     `json:"action" bson:"action"`
	Current      bool       `json:"-" bson:"current"`
	CurrentName  string     `json:"current_name" bson:"current_name"`
	Format       string     `json:"format" bson:"format"`
	OriginalName string     `json:"original_name" bson:"original_name"`
	PATH         string     `json:"-" bson:"path"`
	Signing      string     `json:"-" bson:"signing"`
	Size         int64      `json:"size" bson:"size"`
	URL          string     `json:"url" bson:"file_url"`
	Ref          *mgo.DBRef `json:"-" bson:"user"`
}

// GetDocumentName Required Method ...
func (f *File) GetDocumentName() string {
	return "files"
}

// Upload file to Amazon S3
func (f *File) Upload(part *multipart.FileHeader) (url string, err error) {
	// S3 PATH -> /:model/:id/:action/:uuid.ext
	model := f.Ref.Collection
	id := f.Ref.Id.(string)
	action := f.Action
	uuid, err := f.GenerateUUID()
	if err != nil {
		return
	}

	nameSplit := strings.Split(strings.ToLower(part.Filename), ".")
	suffix := nameSplit[len(nameSplit)-1]

	// Open multipart file
	file, err := part.Open()
	if err != nil {
		return
	}
	defer file.Close()

	f.ACL = core.S3PublicRead
	f.Format = suffix
	f.CurrentName = uuid + "." + suffix
	f.OriginalName = part.Filename
	f.PATH = model + "/" + id + "/" + action + "/" + f.CurrentName
	f.Size, _ = file.Seek(io.SeekStart, io.SeekEnd)

	if err = core.UploadFile(file, f.PATH, f.ACL, f.Format, f.Size); err != nil {
		return
	}

	url, err = core.GetS3Object(1440, f.PATH)
	if err != nil {
		return
	}

	f.URL = url

	return f.URL, err
}

//GenerateUUID generates unique filename to save in s3
func (f *File) GenerateUUID() (string, error) {
	out, err := exec.Command("uuidgen").Output()
	if err != nil {
		return "", err
	}

	return strings.Replace(strings.Trim(string(out), "\n"), "-", "_", -1), nil
}

package core

import (
	"io"
	"mime/multipart"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/revel/revel"
)

// Const for AWS-S3 connection
const (
	AwsAccessKeyID     = "AKIAJ2MIUWROIWMSBUVQ"
	AwsSecretAccessKey = "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"
	Token              = ""
	AWSBucket          = "spyctest01"
	AWSBucketProd      = "spychattterprod"
	AWSRegion          = "us-west-2"
	AWSRegionProd      = "us-east-1"
	S3Private          = "private"
	S3PublicRead       = "public-read"
)

// UploadFile ....
func UploadFile(file multipart.File, path string, acl string, ext string, size int64) error {
	// Get s3 session with credentials
	svc, err := getS3Session()
	if err != nil {
		return err
	}

	contentType := GetFileType("file."+ext) + "/" + ext

	// Generate the svg PUT object to the defined Bucket
	params := &s3.PutObjectInput{
		Bucket:        aws.String(GetBucket()),
		Key:           aws.String(path),
		Body:          file,
		ACL:           aws.String(acl),
		ContentType:   aws.String(contentType),
		ContentLength: aws.Int64(size),
	}

	if _, err := svc.PutObject(params); err != nil {
		return err
	}

	return nil
}

// UploadOSFileReader ....
func UploadOSFileReader(file io.ReadSeeker, path string, acl string, ext string, size int64) (url string, err error) {
	// Get s3 session with credentials
	svc, err := getS3Session()
	if err != nil {
		revel.ERROR.Print(err)
		return
	}

	file.Seek(0, 0)

	contentType := GetFileType("file."+ext) + "/" + ext

	// Generate the svg PUT object to the defined Bucket
	params := &s3.PutObjectInput{
		Bucket:        aws.String(GetBucket()),
		Key:           aws.String(path),
		Body:          file,
		ACL:           aws.String(acl),
		ContentType:   aws.String(contentType),
		ContentLength: aws.Int64(size),
	}

	if _, err = svc.PutObject(params); err != nil {
		revel.ERROR.Print(err)
		return
	}

	req, _ := svc.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(GetBucket()),
		Key:    aws.String(path),
	})

	if url, err = req.Presign(DefaultSigningTime * time.Minute); err != nil {
		return
	}

	return
}

// DeleteFile remove file from S3
func DeleteFile(path string) error {
	// Get s3 session with credentials
	svc, err := getS3Session()
	if err != nil {
		revel.ERROR.Print(err)
		return err
	}

	input := &s3.DeleteObjectInput{
		Bucket: aws.String(GetBucket()),
		Key:    aws.String(path),
	}

	if _, err := svc.DeleteObject(input); err != nil {
		return err
	}

	return nil
}

// GetS3Object return url of S3 object
func GetS3Object(minutes int, path string) (string, error) {
	var url string

	// Get s3 session with credentials
	svc, err := getS3Session()
	if err != nil {
		return url, err
	}

	req, _ := svc.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(GetBucket()),
		Key:    aws.String(path),
	})

	url, err = req.Presign(time.Duration(minutes) * time.Minute)
	if err != nil {
		return url, err
	}

	return url, nil
}

// getS3Session gets the S3 current session
func getS3Session() (svc *s3.S3, err error) {
	cfg, err := getS3Creds()
	if err != nil {
		return
	}

	svc = s3.New(session.New(), cfg)

	return svc, nil
}

func getS3Creds() (cfg *aws.Config, err error) {
	creds := credentials.NewStaticCredentials(AwsAccessKeyID, AwsSecretAccessKey, Token)

	_, err = creds.Get()
	if err != nil {
		return
	}

	cfg = aws.NewConfig().WithRegion(GetRegion()).WithCredentials(creds)

	return
}

//GetFileType returns the type of the file based on the file Extension
func GetFileType(fileName string) (fileType string) {

	ext := filepath.Ext(fileName)

	switch ext {
	// AUDIO files
	case ".aif", ".cda", ".mid", ".midi", ".mp3", ".mpa", ".ogg", ".wav", ".wma", ".wpl":
		return "audio"

	// COMPRESSED files
	case ".7z", ".arj", ".deb", ".pkg", ".rar", ".rpm", ".tar.gz", ".z", ".zip":
		return "compressed"

	// MEDIAFILE files
	case ".dmg", ".iso", ".toast", ".vcd":
		return "mediafile"

	//DATABASE files
	case ".csv", ".dat", ".db", ".dbf", ".log", ".mdb", ".sav", ".sql", ".tar", ".xml":
		return "database"

	//EXECUTABLE files
	case ".apk", ".bat", ".bin", ".pl", ".com", ".exe", ".gadget", ".jar", ".wsf":
		return "exec"

	//FONT files
	case ".fnt", ".fon", ".otf", ".ttf":
		return "font"

	//IMAGE files
	case ".ai", ".bmp", ".gif", ".ico", ".jpg", ".jpeg", ".png", ".ps", ".psd", ".svg", ".tif", ".tiff":
		return "image"

	//INTERNET files
	case ".asp", ".aspx", ".cer", ".cmf", ".cgi", ".css", ".htm", ".html", ".js", ".jsp", ".part", ".php", ".rss", ".xhtml", ".cshtml":
		return "internet"

	//PRESENTATION files
	case ".key", ".odp", ".pps", ".ppt", ".pptx":
		return "presentation"

	//PROGRAMMING files
	case ".c", ".class", ".cpp", ".cs", ".java", ".py", ".h", ".sh", ".swift", ".vb", ".go", ".lisp", ".lua":
		return "programming"

	//SYSTEM files
	case ".bak", ".cab", ".cfg", ".cpl", ".cur", ".dll", ".dmp", ".drv", ".icns", ".ini", ".lnk", ".msi", ".sys", ".tmp":
		return "sys"

	// VIDEO files
	case ".3g2", ".3gp", ".avi", ".flv", ".h264", ".m4v", ".mkv", ".mov", ".mp4", ".mpg", ".mpeg", ".rm", ".swf", ".vob", ".wmv":
		return "video"

	// TEXT PROCESSOR files
	case ".doc", ".docx", ".odt", ".pdf", ".rtf", ".txt", ".tex", ".wks", ".wpd":
		return "text"

	// OTHER file format
	default:
		return "others"

	}
}

// GetBucket returns name of S3 bucket for dev and prod
func GetBucket() string {
	if revel.DevMode {
		return AWSBucket
	}

	return AWSBucketProd
}

// GetRegion returns bucket region for dev and prod
func GetRegion() string {
	if revel.DevMode {
		return AWSRegion
	}

	return AWSRegionProd
}

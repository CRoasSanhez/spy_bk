package core

import (
	"bytes"
	"crypto/tls"
	"log"
	"math/rand"
	"reflect"
	"time"

	gomail "gopkg.in/gomail.v2"
)

var (
	LetterRunes    = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ*-/&.")
	WebLetterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ().*-_%$!¡+[]")
)

// FindOnArray search over string array, if value exists on array returns at position, else return -1
func FindOnArray(array []string, value string) int {
	for k, v := range array {
		if v == value {
			return k
		}
	}

	return -1
}

// GenerateToken make string token with size (x)
func GenerateToken(dictionary []rune, size int) string {
	rand.Seed(time.Now().UTC().UnixNano())

	b := make([]rune, size)
	for i := range b {
		b[i] = dictionary[rand.Intn(len(dictionary))]
	}
	return string(b)
}

// SendEmail ...
func SendEmail(to []string, subject string, body string) error {
	m := gomail.NewMessage()

	m.SetHeader("From", "noreply@spychatter.net")
	m.SetHeader("To", to...)
	// m.SetAddressHeader("Cc", "dan@example.com", "Dan")
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)
	// m.Attach("/home/Alex/lolcat.jpg")

	d := gomail.NewDialer("smtpout.secureserver.net", 80, "noreply@spychatter.net", "dM\"K>JMtt\\em6RB")
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	return d.DialAndSend(m)
}

// ConcatArray concatenates elements from string array into single string
func ConcatArray(array []string) string {
	var buffer bytes.Buffer
	for i := 0; i < len(array); i++ {
		buffer.WriteString(array[i])
		if i <= len(array)-1 {
			buffer.WriteString(",")
		}
	}
	return buffer.String()
}

// InTimeSpan returns if a checkTime is between two dates
func InTimeSpan(start, end, check time.Time) bool {
	return check.After(start) && check.Before(end)
}

// HoursBetweenDates returns the difference in hours between two dates
func HoursBetweenDates(stime, etime time.Time) float64 {
	return stime.Sub(etime).Hours()
}

// ChangeUTCTimeToLocalZone ...
func ChangeUTCTimeToLocalZone(date time.Time) time.Time {
	_, offset := time.Now().Zone()
	return date.Local().Add(-time.Duration(offset) * time.Second)
}

// ChangeLocalTimeToUTCZone ...
func ChangeLocalTimeToUTCZone(date time.Time) time.Time {
	_, offset := time.Now().Zone()
	return date.UTC().Add(+time.Duration(offset) * time.Second)
}

// ClearFields set zero value for current field
func ClearFields(obj interface{}, fields ...string) {
	r := reflect.ValueOf(obj)

	for _, v := range fields {
		f := reflect.Indirect(r).FieldByName(v)
		if f.CanSet() {
			f.Set(reflect.Zero(f.Type()))
		} else {
			log.Println("No se puede setear ", v)
		}
	}
}

// MergeStructs set not null fields on src to dist
func MergeStructs(dist interface{}, src interface{}) {
	r := reflect.ValueOf(src)
	r2 := reflect.ValueOf(dist)

	if reflect.TypeOf(dist) != reflect.TypeOf(src) {
		return
	}

	for r.Kind() == reflect.Ptr {
		r = r.Elem()
	}

	for i := 0; i < r.NumField(); i++ {
		f := reflect.Indirect(r).Field(i).Interface()
		if f != reflect.Zero(reflect.TypeOf(f)).Interface() {
			// set field on dist interface{}
			f2 := reflect.Indirect(r2).Field(i)
			if f2.CanSet() {
				f2.Set(reflect.ValueOf(f))
			}
		}
	}

}

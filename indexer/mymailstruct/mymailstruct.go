package mymailstruct

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"strings"
)

const (
	dataBaseName           = "ENRON_MAILS"
	zincSearchBulkAPI      = "http://localhost:4080/api/_bulk"
	zincSearchUserName     = "admin"
	zincSearchUserPassword = "Complexpass#123"
)

var zincSearchAPI = fmt.Sprintf("http://localhost:4080/api/%s/_doc", dataBaseName)

// MyMailStruct contains the fields of an email.
type MyMailStruct struct {
	MessageID   string
	Date        string
	From        string
	To          string
	Subject     string
	MVersion    string
	ContentType string
	Encoding    string
	XFrom       string
	XTo         string
	Xcc         string
	Xbcc        string
	XFolder     string
	XOrigin     string
	XFileName   string
	BodyMessage string
}

// ReadMailFile This function read a mail file ("path") and returns a string with the full message.
func ReadMailFile(path string) string {
	file, err := os.ReadFile(path)
	if err != nil {
		fmt.Printf("Could not read the file due to this %s error \n", err)
	}
	return string(file)
}

// IndexMailMessage This function returns a string in json format with the essential information of the string "message".
func IndexMailMessage(message string) string {
	fieldName := [15]string{"Message-ID: ", "\r\nDate: ", "\r\nFrom: ", "\r\nTo: ", "\r\nSubject: ", "\r\nMime-Version: ", "\r\nContent-Type: ", "\r\nContent-Transfer-Encoding: ", "\r\nX-From: ", "\r\nX-To: ", "\r\nX-cc: ", "\r\nX-bcc: ", "\r\nX-Folder: ", "\r\nX-Origin: ", "\r\nX-FileName: "}
	fieldName2 := [15]string{"\r\n", "\r\n", "\r\n", "\r\nSubject: ", "\r\n", "\r\n", "\r\n", "\r\n", "\r\n", "\r\n", "\r\n", "\r\n", "\r\n", "\r\n", "\r\n"}

	var fieldData [16]string

	before := " "
	after := " "
	found := false

	for i := 0; i < 15; i++ {
		_, after, found = strings.Cut(message, fieldName[i])
		if found {
			before, _, _ = strings.Cut(after, fieldName2[i])
			fieldData[i] = before
		} else {
			fieldData[i] = " "
		}
	}

	_, after, found = strings.Cut(message, fieldName[14])
	if found {
		_, after, _ = strings.Cut(after, "\r\n")
		fieldData[15] = after
	} else {
		fieldData[15] = " "
	}

	data := MyMailStruct{
		MessageID:   fieldData[0],
		Date:        fieldData[1],
		From:        fieldData[2],
		To:          fieldData[3],
		Subject:     fieldData[4],
		MVersion:    fieldData[5],
		ContentType: fieldData[6],
		Encoding:    fieldData[7],
		XFrom:       fieldData[8],
		XTo:         fieldData[9],
		Xcc:         fieldData[10],
		Xbcc:        fieldData[11],
		XFolder:     fieldData[12],
		XOrigin:     fieldData[13],
		XFileName:   fieldData[14],
		BodyMessage: fieldData[15],
	}

	b, err := json.Marshal(data)
	if err != nil {
		fmt.Println("error:", err)
	}

	return string(b)
}

// IndexInZincSearc This function indexes (in ZincSearch) all email contained in root.
func IndexInZincSearc(root string) {
	fileSystem := os.DirFS(root)
	fs.WalkDir(fileSystem, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			log.Fatal(err)
		}

		if !d.IsDir() {
			file, err := os.ReadFile(root + "/" + path)
			if err != nil {
				fmt.Printf("Could not read the file due to this %s error \n", err)
			}

			data := IndexMailMessage(string(file))

			req, err := http.NewRequest(http.MethodPost, zincSearchAPI, strings.NewReader(data))
			if err != nil {
				log.Fatal(err)
			}

			req.SetBasicAuth(zincSearchUserName, zincSearchUserPassword)
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/81.0.4044.138 Safari/537.36")

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				log.Fatal(err)
			}

			defer resp.Body.Close()
		}

		return nil
	})
}

// IndexData sends a batch of mails to ZincSearch.
func IndexData(buffer []string) error {
	req, err := http.NewRequest(http.MethodPost, zincSearchBulkAPI, strings.NewReader(strings.Join(buffer, "\n")))
	if err != nil {
		log.Println("Error: IndexData: NewRequest: ", err)
		return err
	}

	req.SetBasicAuth(zincSearchUserName, zincSearchUserPassword)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/81.0.4044.138 Safari/537.36")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("Error: Requesting: ", err)
		return err
	}

	defer resp.Body.Close()

	return nil
}

// BulkIndexInZincSearc. This function indexes (in ZincSearch) all email contained in root.
func BulkIndexInZincSearc(root string) {
	N := 10000
	buffCounter := 0

	bufferStr := make([]string, N)

	fileSystem := os.DirFS(root)
	fs.WalkDir(fileSystem, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			log.Fatal(err)
		}

		if !d.IsDir() {
			file, err := os.ReadFile(root + "/" + path)
			if err != nil {
				fmt.Printf("Could not read the file due to this %s error \n", err)
			}

			if buffCounter == N-1 {
				err = IndexData(bufferStr)
				if err != nil {
					fmt.Printf("Error: IndexDate: %s \n", err)
				}

				buffCounter = 0
			}

			data := IndexMailMessage(string(file))
			bufferStr[buffCounter] = `{"index" : {"_index" : "games4"}}` + "\n" + data
			buffCounter++
		}

		return nil
	})

	err := IndexData(bufferStr)
	if err != nil {
		fmt.Printf("Error: IndexDate: %s \n", err)
	}
}

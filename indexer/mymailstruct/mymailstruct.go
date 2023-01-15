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

// This function read a mail file ("path") and returns a string with the full message.
func ReadMailFile(path string) string {
	file, err := os.ReadFile(path)
	if err != nil {
		fmt.Printf("Could not read the file due to this %s error \n", err)
	}
	return string(file)
}

// This function returns a string in json format with the essential information of the string "message".
func IndexMailMessage(message string) string {
	allfound := true
	before, after, found := strings.Cut(message, "Message-ID: ")
	allfound = allfound && found
	before, after, found = strings.Cut(after, "\r\nDate: ")
	id := before
	allfound = allfound && found
	before, after, found = strings.Cut(after, "\r\nFrom: ")
	date := before
	allfound = allfound && found
	before, after, found = strings.Cut(after, "\r\nTo: ")
	from := before
	allfound = allfound && found
	before, after, found = strings.Cut(after, "\r\nSubject: ")
	to := before
	allfound = allfound && found
	before, after, found = strings.Cut(after, "\r\nMime-Version: ")
	subject := before
	allfound = allfound && found
	before, after, found = strings.Cut(after, "\r\nContent-Type: ")
	version := before
	allfound = allfound && found
	before, after, found = strings.Cut(after, "\r\nContent-Transfer-Encoding: ")
	ctype := before
	allfound = allfound && found
	before, after, found = strings.Cut(after, "\r\nX-From: ")
	ctencoding := before
	allfound = allfound && found
	before, after, found = strings.Cut(after, "\r\nX-To: ")
	xfrom := before
	allfound = allfound && found
	before, after, found = strings.Cut(after, "\r\nX-cc: ")
	xto := before
	allfound = allfound && found
	before, after, found = strings.Cut(after, "\r\nX-bcc: ")
	xcc := before
	allfound = allfound && found
	before, after, found = strings.Cut(after, "\r\nX-Folder: ")
	xbcc := before
	allfound = allfound && found
	before, after, found = strings.Cut(after, "\r\nX-Origin: ")
	xfolder := before
	allfound = allfound && found
	before, after, found = strings.Cut(after, "\r\nX-FileName: ")
	xorigin := before
	allfound = allfound && found
	before, after, found = strings.Cut(after, "\r\n")
	xfilename := before
	body := after
	allfound = allfound && found

	if !allfound {
		fmt.Println("error: unexpected file format", message)
		return ""
	}

	data := MyMailStruct{
		MessageID:   id,
		Date:        date,
		From:        from,
		To:          to,
		Subject:     subject,
		MVersion:    version,
		ContentType: ctype,
		Encoding:    ctencoding,
		XFrom:       xfrom,
		XTo:         xto,
		Xcc:         xcc,
		Xbcc:        xbcc,
		XFolder:     xfolder,
		XOrigin:     xorigin,
		XFileName:   xfilename,
		BodyMessage: body,
	}

	b, err := json.Marshal(data)
	if err != nil {
		fmt.Println("error:", err)
	}

	return string(b)
}

// This function indexes (in ZincSearch) all email contained in root.
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

			// usar hilos
			// wait groups (concurrencia)
			// variables de entorno
			req, err := http.NewRequest("POST", "http://localhost:4080/api/games3/_doc", strings.NewReader(data))
			if err != nil {
				log.Fatal(err)
			}
			req.SetBasicAuth("admin", "Complexpass#123")
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

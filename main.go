package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func main() {
	servedFolderName := "myfsfolder"
	// show directory
	http.Handle("/"+servedFolderName+"/", http.StripPrefix("/"+servedFolderName+"/", http.FileServer(http.Dir("./"+servedFolderName))))
	// html
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/html")
		http.ServeFile(w, r, "uploadform.html")
	})
	// POST /upload
	http.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
		var MaxUploadSize int64
		MaxUploadSize = 1024 * 1024 * 1024
		if r.Method != "POST" {
			http.Error(w, "HTTP method must be POST. Check if you are using correct http method.", http.StatusMethodNotAllowed)
			return
		}

		r.Body = http.MaxBytesReader(w, r.Body, MaxUploadSize)
		if err := r.ParseMultipartForm(MaxUploadSize); err != nil {
			http.Error(w, "1GB以下のファイルを選択してください。", http.StatusBadRequest)
			return
		}

		// Read uploaded file
		file, fileHeader, err := r.FormFile("file")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer file.Close()

		// Create the directory to save the uploaded files if not exists
		err = os.MkdirAll("./"+servedFolderName, os.ModePerm)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// default filename -> ./foldername/current_time_as_UNIX_time.xxx
		filenameUploaded := fmt.Sprintf("./%s/%d%s", servedFolderName, time.Now().UnixNano(), filepath.Ext(fileHeader.Filename))

		// Get filename posted from the form
		customFilename := r.FormValue("filename")
		if customFilename != "" {
			// if custom filename is posted, use it as the filename
			fmt.Println("custom filename:", customFilename)
			filenameUploaded = fmt.Sprintf("./%s/%s%s", servedFolderName, customFilename, filepath.Ext(fileHeader.Filename))
		}

		// Create a file in fileserver
		dst, err := os.Create(filenameUploaded)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer dst.Close()

		// Copy the content of file uploaded to the newly created file in fileserver
		_, err = io.Copy(dst, file)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Respond with success page
		w.Header().Add("Content-Type", "text/html")
		http.ServeFile(w, r, "uploadsuccess.html")
	})
	// Respond pong!
	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "pong!")
	})
	http.ListenAndServe(":8080", nil)
}

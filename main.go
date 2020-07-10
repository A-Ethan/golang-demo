package main

import (
	"image/png"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/disintegration/imaging"
)

func main() {
	http.HandleFunc("/", homePage)
	http.HandleFunc("/upload", uploadHandle)
	fs := http.FileServer(http.Dir("./uploaded"))
	http.Handle("/uploaded/", http.StripPrefix("/uploaded", fs))

	log.Print("Server started on localhost:80, use /upload for uploading files and /files/{fileName} for downloading files.")
	log.Fatal(http.ListenAndServe(":80", nil))
}

func homePage(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	tmpl, err := os.OpenFile("./upload.gtpl", os.O_RDONLY, 0666)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	io.Copy(w, tmpl)
	return
}

func uploadHandle(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	req.ParseForm()
	if req.Method != "POST" {
		w.Write([]byte("method not support"))
		return
	} else {
		// 接收图片
		uploadFile, handle, err := req.FormFile("uploadFile")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// 检查图片后缀
		ext := strings.ToLower(path.Ext(handle.Filename))
		if strings.ToLower(ext) != ".png" {
			http.Error(w, "只支持png图片上传", http.StatusForbidden)
			return
		}
		// file to image
		image, err := png.Decode(uploadFile)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// 修改图片大小
		dstImage128 := imaging.Resize(image, 512, 512, imaging.Lanczos)
		// 保存图片
		os.Mkdir("./uploaded/", 0777)
		saveFile, err := os.OpenFile("./uploaded/"+handle.Filename, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		png.Encode(saveFile, dstImage128)
		defer uploadFile.Close()
		defer saveFile.Close()
		// 上传图片成功
		w.Write([]byte(`<meta charset="utf-8"> 查看上传图片: <a target='_blank' href='/uploaded/` + handle.Filename + "'>" + handle.Filename + "</a>"))
		// 下载图片
		// fs := http.FileServer(http.Dir(uploadPath))
		// http.Handle("/download/", http.StripPrefix("/files", fs))

	}
}

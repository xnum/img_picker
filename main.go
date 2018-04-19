package main

import (
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"text/template"
	"time"
)

const (
	WWWROOT = "public"
	DATA    = "data"
)

var store ImageFaceInfoStore = nil

type Message struct {
	Code int
	Msg  string
	Pos  [][]int
}

func writeResp(msg *Message, w http.ResponseWriter) {
	b, err := json.Marshal(msg)
	if err != nil {
		log.Fatal(err)
		return
	}

	w.Write(b)
}

func dataHandler(w http.ResponseWriter, req *http.Request) {
	store, _ = StoreNew(DATA)
	store.Check(WWWROOT)
	str := store.Jsonify()
	io.WriteString(w, str)
}

func uploadHandler(w http.ResponseWriter, req *http.Request) {

	if req.Method == "GET" {
		crutime := time.Now().Unix()
		h := md5.New()
		io.WriteString(h, strconv.FormatInt(crutime, 10))
		token := fmt.Sprintf("%x", h.Sum(nil))

		t, err := template.ParseFiles("upload.gtpl")
		if err != nil {
			log.Println(err)
			return
		}
		t.Execute(w, token)

		return
	}

	if req.Method != "POST" {
		io.WriteString(w, "POST Only!\n")
		return
	}

	req.ParseMultipartForm(32 << 20)
	file, handler, err := req.FormFile("uploadfile")
	if err != nil {
		fmt.Println(err)
		return
	}

	defer file.Close()

	resp := &Message{200, "", nil}
	defer writeResp(resp, w)

	filename, err := input_check(handler.Filename)
	if err != nil {
		resp = &Message{500, err.Error(), nil}
		return
	}

	www_filename := WWWROOT + "/" + filename

	f, err := os.OpenFile(www_filename,
		os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		resp = &Message{500, err.Error(), nil}
		return
	}

	defer f.Close()
	io.Copy(f, file)

	rows, err := face_detect(www_filename)
	if err != nil {
		resp = &Message{500, err.Error(), nil}
		return
	}

	index := strings.LastIndexByte(www_filename, '.')
	if index < 0 {
		log.Fatal("filename.")
		return
	}

	resp.Pos = make([][]int, len(rows)-1)
	// If we wanna split image
	for k, v := range rows {
		pos := strings.Fields(v)
		if len(pos) != 4 {
			break
		}

		arr := make([]int, 4)
		for i := 0; i < 4; i++ {
			arr[i], _ = strconv.Atoi(pos[i])
		}
		resp.Pos[k] = arr

		target_name := fmt.Sprintf("%s-%d%s", www_filename[:index], k, www_filename[index:])
		if err := split_image(www_filename, target_name, pos); err != nil {
			log.Println(err)
			continue
		}
	}

	store.Add(filename, resp.Pos)
	store.Save(DATA)
}

func input_check(filename string) (string, error) {
	if len(filename) < 3 {
		return "", errors.New("Filename too short")
	}

	any := func(val string, arr []string, f func(string, string) bool) bool {
		for _, v := range arr {
			if f(val, v) {
				return true
			}
		}

		return false
	}

	if !any(filename, []string{".jpg", ".png"}, strings.HasSuffix) {
		return "", errors.New("unsupport file format")
	}

	if m, _ := regexp.MatchString("^[a-zA-Z0-9_-]+[.][a-zA-Z]+$", filename); !m {
		return "", errors.New("bad filename")
	}

	return filename, nil
}

// Returns split image filenames
func face_detect(filename string) ([]string, error) {
	out, err := exec.Command("facedetect", filename).Output()
	if err != nil {
		log.Fatal(err)
		return nil, nil
	}

	rows := strings.Split(string(out), "\n")

	return rows, err
}

// Return single split image
func split_image(filename, target_name string, pos []string) error {
	// convert -crop 70x70+46+139 /tmp/upload.jpg /home/num/upload-00.jpg
	pos_str := fmt.Sprintf("%vx%v+%v+%v", pos[2], pos[3], pos[0], pos[1])

	err := exec.Command("convert", "-crop", pos_str, filename, target_name).Run()
	if err != nil {
		return err
	}

	return nil
}

func main() {
	var err error
	store, err = StoreNew(DATA)
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/upload", uploadHandler)
	http.HandleFunc("/data", dataHandler)
	http.Handle("/", http.FileServer(http.Dir(WWWROOT)))

	fmt.Println("Listening...")
	log.Fatal(http.ListenAndServe(":4000", nil))
}

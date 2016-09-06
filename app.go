package main

import (
  "net/http"
  "io/ioutil"
  "log"
  "encoding/json"
  "os"
  "errors"
  "regexp"
)

type Package struct {
  Name string
  Size int64
}

const rootDirectory = "/resources"

func getPackages(path string)([]Package, error) {
  files, err := ioutil.ReadDir(path)

  if err != nil {
    log.Println(err)
    return nil, errors.New("Directory does not exist")
  }

  packages := make([]Package, 0, len(files))

  for _, file := range files {
    packages = append(packages, Package{Name: file.Name(), Size: file.Size()})
  }

  return packages, nil
} 

func ensurePath(path string){
  re, err := regexp.Compile(`^(.*)\/(.*)$`)
  if err != nil{
    log.Println(err)
  }

  result := re.FindStringSubmatch(path)

  if result[1] != ""{
    rw := int(0777)
    err = os.MkdirAll(rootDirectory + result[1], os.FileMode(rw))
    if err != nil {
      log.Println(err.Error())
    }
  }
}

func HandlePost(writer http.ResponseWriter, request *http.Request) {
  buf, _ := ioutil.ReadAll(request.Body)
  ro := int(0444)
  ensurePath(request.URL.Path)
  err := ioutil.WriteFile(rootDirectory + request.URL.Path, buf, os.FileMode(ro))

  if err != nil {
    http.Error(writer, err.Error(), 500)
  } else {
    writer.Write([]byte("File succesfuly uploaded"))
  }
}


func HandleGet(writer http.ResponseWriter, request *http.Request) {
  path := rootDirectory + request.URL.Path
  fileOrDir, err := os.Open(path) 

  if err != nil {
    http.Error(writer, err.Error(), 404)
    return
  }

  defer fileOrDir.Close()

  fd, err := fileOrDir.Stat()

  if err != nil {
    http.Error(writer, err.Error(), 500)
    return
  }

  if fd.IsDir() {
    packages, err := getPackages(path)

    if err != nil {
      http.Error(writer, err.Error(), 500)
    }else{
      payload, _ := json.Marshal(packages)
      writer.Header().Add("Content-Type", "application/json")
      writer.Write(payload)
    }

  } else {
    content, err := ioutil.ReadFile(path)

    if err != nil {
      http.Error(writer, err.Error(), 500)
    }else{
      writer.Header().Add("Content-Type", "application/octet-stream")
      writer.Write(content)
    }
  }


}

func rootHandler(writer http.ResponseWriter, request *http.Request) {

  log.Println(request.Method + " " + request.URL.Host + request.URL.Path)

  if request.Method == "GET" {
    HandleGet(writer, request)
  } else if request.Method == "POST" {
    HandlePost(writer, request)
  } else{
    http.Error(writer, "Method not allowed", 403)
  }
}

func main() {
  port := "3535"
  http.HandleFunc("/", rootHandler);
  log.Print("Starting on port " + port)
  http.ListenAndServe(":" + port, nil)
}

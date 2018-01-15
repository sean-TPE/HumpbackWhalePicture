/*****************************************************************************
 * Step 3. Serving static files (images) which are stored in database.       *
 *                                                                           *
 * To visit:                                                                 *
 *       http://localhost:3000/public/image/jpg/humpback_whale.jpg           *
 *       http://localhost:3000/public/image/png/humpback_whale.png           *
 *       http://localhost:3000/public/image/gif/humpback_whale.gif           *
 *****************************************************************************/
package main

import (
    "fmt"
    "net/http"
    "io/ioutil"
    "database/sql"
    "strings"
    "path/filepath"
    "encoding/base64"
    "strconv"

    _ "github.com/mattn/go-sqlite3"
)

const homepage = "public/html/index.html"

var db *sql.DB
const dbName = "imgDB"
const tbName = "images"

func main() {
    dbname := dbName + ".db"
    var err error
    db, err = sql.Open("sqlite3", dbname)
    if err != nil {
        fmt.Printf("Open %v Error: %v\n", dbname, err)
        return
    }
    defer db.Close()

    http.HandleFunc("/", homeHandleFunc)
    http.HandleFunc("/public/", imgHandleFunc)
    http.ListenAndServe(":3000", nil)
}

func homeHandleFunc(w http.ResponseWriter, r *http.Request) {
    data, err := ioutil.ReadFile(homepage)
    if err != nil {
        fmt.Printf("ReadFile %v Error: %v\n", homepage, err)
        http.NotFound(w, r)
    }
    html := string(data)

    w.Header().Set("Content-Type", "text/html; charset=utf-8")
    fmt.Fprint(w, html)
}

func imgHandleFunc(w http.ResponseWriter, r *http.Request) {
    file := strings.TrimLeft(r.URL.String(), "/")
    dir := filepath.Dir(file)
    file = filepath.Base(file)

    stmt := fmt.Sprintf("SELECT * FROM %v WHERE name=? AND path=?", tbName)
    rows, err := db.Query(stmt, file, dir)
    if err != nil {
        fmt.Printf("Query (%v) Error: %v\n", filepath.Join(dir, file), err)
        http.NotFound(w, r)
        return
    }

    var id int
    var name string
    var path string
    var data string
    for rows.Next() {
        err = rows.Scan(&id, &name, &path, &data)
        if err != nil {
            fmt.Println("Scan Error: ", err)
            http.Redirect(w, r, "/", http.StatusSeeOther) // TODO: modify redirect route
        } else {
            break
        }
    }

    err = rows.Err()
    if err != nil {
        fmt.Println("rows.Err(): ", err)
        http.Redirect(w, r, "/", http.StatusSeeOther) // TODO: modify redirect route
        return
    }

    if data == "" {
        fmt.Println("Oops! Empty data...")
        http.Redirect(w, r, "/", http.StatusSeeOther) // TODO: modify redirect route
        return
    }

    decodeBytes, err := base64.StdEncoding.DecodeString(data)
    if err != nil {
        fmt.Println("Base64 decode Error: ", err)
        http.Redirect(w, r, "/", http.StatusSeeOther) // TODO: modify redirect route
        return
    }

    ext := strings.TrimLeft(filepath.Ext(name), ".")
    w.Header().Set("Content-Type", fmt.Sprintf("image/%v", ext))
    w.Header().Set("Content-Length", strconv.Itoa(len(decodeBytes)))

    if _, err := w.Write(decodeBytes); err != nil {
        fmt.Println("Unable to wriet image. Error: ", err)
        http.Redirect(w, r, "/", http.StatusSeeOther) // TODO: modify redirect route
        return
    }
}
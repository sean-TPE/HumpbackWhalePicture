/*****************************************************************************
 * Step 2. Store images in SQLite database.                                  *
 *****************************************************************************/
package main

import (
    "fmt"
    "database/sql"
    "encoding/base64"
    "path/filepath"
    "io/ioutil"
    "os"

    _ "github.com/mattn/go-sqlite3"
    . "pathmgr"
)

var db *sql.DB
const dbName = "imgDB"
const tbName = "images"
const outputFolder = "output"
const sourcePath = "public"

func main() {
    dbname := dbName + ".db"
    var err error
    db, err = sql.Open("sqlite3", dbname)
    if err != nil {
        fmt.Printf("Open %v Error: %v\n", dbname, err)
        return
    }
    defer db.Close()

    stmt := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %v (id INTEGER PRIMARY KEY AUTOINCREMENT, name TEXT, path TEXT, data BLOB)", tbName)
    sqlcmd, err := db.Prepare(stmt)
    if err != nil {
        fmt.Println("Prepare 'CREATE TABLE' command Error: ", err)
        return
    }

    _, err = sqlcmd.Exec()
    if err != nil {
        fmt.Println("Execute 'CREATE TABLE' command Error: ", err)
        return
    }

    filePaths, err := getPicturePathsIn(sourcePath)
    if err != nil {
        fmt.Println("Get picture in ", sourcePath, " Error: ", err)
        return
    }

    storeToDB(filePaths)

    storeImageToOutputFolder(outputFolder)
}

func getPicturePathsIn(path string) ([]string, error) {
    return GetListByOption(path, Recursive | HideFolder | ToSlash, "gif,jpg,png")
}

func storeToDB(filePaths []string) {
    for _, file := range filePaths {
        var cnt int
        stmt := fmt.Sprintf("SELECT COUNT(*) FROM %v WHERE name=? AND path=?", tbName)
        _ = db.QueryRow(stmt, filepath.Base(file), filepath.Dir(file)).Scan(&cnt)
        if cnt != 0 {
            continue
        }

        data, err := ioutil.ReadFile(file)
        if err != nil {
            fmt.Printf("ReadFile %v Error: %v\n", file, err)
            continue
        }

        encodeString := base64.StdEncoding.EncodeToString(data)

        stmt = fmt.Sprintf("INSERT INTO %v (name, path, data) VALUES (?, ?, ?)", tbName)
        sqlcmd, err := db.Prepare(stmt)
        if err != nil {
            fmt.Println("Prepare 'INSERT' command Error: ", err)
            continue
        }

        _, err = sqlcmd.Exec(filepath.Base(file), filepath.Dir(file), encodeString)
        if err != nil {
            fmt.Println("Execute 'INSERT' command Error: ", err)
            continue
        }
    }
}

func storeImageToOutputFolder(outputPath string) {
    stmt := fmt.Sprintf("SELECT * FROM %v", tbName)
    rows, err := db.Query(stmt)
    if err != nil {
        fmt.Println("Query Error: ", err)
        return
    }

    for rows.Next() {
        var id int
        var name string
        var path string
        var data string
        err = rows.Scan(&id, &name, &path, &data)
        if err != nil {
            fmt.Printf("Scan %d row Error: %v\n", id, err)
            continue
        }

        bytes, err := base64.StdEncoding.DecodeString(data)
        if err != nil {
            fmt.Println("Base64 decode Error: ", err)
            continue
        }

        dirpath := filepath.Join(outputPath, path)
        err = os.MkdirAll(dirpath, os.ModePerm)
        if err != nil {
            fmt.Printf("Mkdir (%v) Error: %v\n", dirpath, err)
            continue
        }

        err = ioutil.WriteFile(filepath.Join(dirpath, name), bytes, 0644)
        if err != nil {
            fmt.Printf("WriteFile (path: %v, name: %v) Error: %v\n", path, name, err)
            continue
        }
    }
}
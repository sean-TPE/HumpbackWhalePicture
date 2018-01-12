package pathmgr

import (
    "fmt"
    "path/filepath"
    "strings"
    "os"
)

// pathMgr is a struct that with some methods
// which allows you list all files in a directory.
type pathMgr struct {
    fileList []string
}

// WalkOption is a enumeration (enum)
// which act as options for the method of pathMgr
type WalkOption uint16

// Default, Recursive, etc. are options for the methods of pathMgr
const (
    Default    WalkOption = 0
    Recursive  WalkOption = 1 << iota
    HideFolder
    HideRoot
    ToSlash
    FileBase
)

func wtob(number WalkOption) bool {
    if number != 0 {
        return true
    } else {
        return false
    }
}

const separator = " "

// GetListByOption list files in the given path.
// It recursively list all files or just list all files in current folder
// depend on the option and constraints.
func GetListByOption(root string, option WalkOption, constraints ...string) ([]string, error) {
    var files pathMgr
    var args string
    for _, v := range constraints {
        args += separator
        args += strings.TrimSpace(v)
    }
    args = strings.TrimLeft(args, separator)

    err := files.getListByOption(root, option, args)
    return files.fileList, err
}

// getListByOption list files in the given path.
// It recursively list all files or just list all files in current folder
// depend on the option and constraints.
func (f *pathMgr) getListByOption(root string, option WalkOption, constraints ...string) error {
    f.fileList = nil
    var err error = nil
    var recursive  = wtob(option & Recursive)
    var hidefolder = wtob(option & HideFolder)
    var hideroot   = wtob(option & HideRoot)

    var args string
    for _, v := range constraints {
        args += separator
        args += strings.TrimSpace(v)
    }
    args = strings.TrimLeft(args, separator)

    err = filepath.Walk(root, func (path string, info os.FileInfo, err error) error {
        if err != nil {
            fmt.Println("Error occur when walk through ", path)
            fmt.Println(err)
            return nil
        }

        if info.IsDir() {

            if !hidefolder {
                if path != root {
                    appendWithOption(&(f.fileList), path, option, args)
                } else {
                    if !hideroot {
                        appendWithOption(&(f.fileList), path, option, args)
                    }
                }
            }
            if !recursive && path != root {
                return filepath.SkipDir
            }

        } else {
            appendWithOption(&(f.fileList), path, option, args)
        }
        return nil
    })

    return err
}

func appendWithOption(slice *[]string, path string, option WalkOption, constraints ...string) {
    var toslash  = wtob(option & ToSlash)
    var filebase = wtob(option & FileBase)

    var args string
    for _, v := range constraints {
        args += separator
        args += strings.TrimSpace(v)
    }
    args = strings.TrimLeft(args, separator)

    if toslash {
        path = filepath.ToSlash(path)
    }
    if filebase {
        path = filepath.Base(path)
    }

    if len(args) != 0 {
        args = strings.Replace(args, ",", separator, -1)
        args = strings.Replace(args, ";", separator, -1)
        types := strings.Split(args, separator)

        for _, v := range types {
            if strings.TrimSpace(v) == "" {
                continue
            }
            if strings.TrimLeft(filepath.Ext(path), ".") == strings.TrimLeft(v, ".") {
                *slice = append(*slice, path)
                break
            }
        }
    } else {
        *slice = append(*slice, path)
    }
}
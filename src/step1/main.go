/*****************************************************************************
 * Step 1. Recursively list all files in a directory and its subdirectories. *
 *****************************************************************************/
package main

import (
    "fmt"
    "path/filepath"
    "os"
)

// PathMgr is a struct that with some methods
// which allows you list all files in a directory.
type PathMgr struct {
    FileList []string
}

// WalkOption is a enumeration (enum)
// which act as options for the method of PathMgr
type WalkOption uint8

// Default, Recursive, etc. are options for the methods of PathMgr
const (
    Default    WalkOption = 0
    Recursive  WalkOption = 1 << iota
    HideFolder
    HideRoot
    ToSlash
)

func wtob(number WalkOption) bool {
    if number != 0 {
        return true
    } else {
        return false
    }
}

// GetListByOption list files in the given path.
// It recursively list all files or just list all files in current folder
// depend on the option.
func (f *PathMgr) GetListByOption(root string, option WalkOption) error {
    f.FileList = nil
    var err error = nil
    var recursive  = wtob(option & Recursive)
    var hidefolder = wtob(option & HideFolder)
    var hideroot   = wtob(option & HideRoot)
    var toslash    = wtob(option & ToSlash)

    err = filepath.Walk(root, func (path string, info os.FileInfo, err error) error {
        if err != nil {
            fmt.Println("Error occur when walk through ", path)
            fmt.Println(err)
            return nil
        }

        rtpath := root
        if toslash {
            rtpath = filepath.ToSlash(rtpath)
            path = filepath.ToSlash(path)
        }

        if info.IsDir() {

            if !hidefolder {
                if path != rtpath {
                    f.FileList = append(f.FileList, path)
                } else {
                    if !hideroot {
                        f.FileList = append(f.FileList, path)
                    }
                }
            }
            if !recursive && path != rtpath {
                return filepath.SkipDir
            }

        } else {
            f.FileList = append(f.FileList, path)
        }
        return nil
    })

    return err
}

func main() {
    var files PathMgr
    files.GetListByOption("public", Recursive | HideRoot | ToSlash)

    for _, v := range files.FileList {
        fmt.Println(v)
    }
}
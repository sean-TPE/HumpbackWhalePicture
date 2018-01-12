/*****************************************************************************
 *    This is a simple example which shows how to use package 'pathmgr'.     *
 *****************************************************************************/
package main

import (
    "fmt"

    . "pathmgr"
)

func main() {
    fmt.Println("1. Recursively list files and ignore folders. <without constraints>")
    files, err := GetListByOption("public", Recursive | HideFolder)

    if err != nil {
        fmt.Println("Error: ", err)
        return
    }

    for _, v := range files {
        fmt.Println(v)
    }

    fmt.Println("====================================================================")
    fmt.Println("2. Recursively list all image files(jpg,png,gif). <with constraints>")
    files, err = GetListByOption("public", Recursive | ToSlash, "jpg;png", "gif")

    if err != nil {
        fmt.Println("Error: ", err)
        return
    }

    for _, v := range files {
        fmt.Println(v)
    }
}
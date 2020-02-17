# wordscalc
Для тестового задания

## 

    package main

    import (
        "fmt"
        "local/wordscalc/wordscalc"
        "log"
    )

    func main() {
        // Подсчет символов из файлов в папке
        words, err := wordscalc.CountInFolderFiles("testsfiles")
        if err != nil {
            log.Fatalln(err)
        }
        for runeN, count := range words {
            fmt.Printf("%#U %d\n", runeN, count)
        }
        fmt.Println()
        // подсчет символов в одном файле
        wordsInFile, err := wordscalc.CountInFile("testsfiles/цувцувцувцу.txt")
        if err != nil {
            log.Fatalln(err)
        }
        for runeN, count := range wordsInFile {
            fmt.Printf("%#U %d\n", runeN, count)
        }
    }


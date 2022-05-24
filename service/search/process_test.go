package search

import (
    "fmt"
    "hegosearch/data/doc"
    "hegosearch/data/index"
    "hegosearch/data/storage"
    "hegosearch/data/tokenize"
    "testing"
    "time"
)

// please test after import the data
func TestSearchText(t *testing.T) {
    docDB := doc.NewDocDriver(
        storage.NewLevelDBStorage("../../data/db/doc"),
    )
    indexDB := index.NewIndexDriver(
        storage.NewLevelDBStorage("../../data/db/index"),
    )
    token := tokenize.NewToken()
    defer indexDB.CloseIndexDB()
    defer docDB.CloseDocDB()
    newSearch := NewSearch(indexDB, docDB, token)
    start := time.Now()
    results := SearchText("开个", newSearch)
    duration := time.Since(start)
    fmt.Println("time: ", duration.Milliseconds(), "ms")
    for _, result := range results {
        fmt.Println("DocId: ", result.DocId)
        fmt.Println("Score: ", result.Score)
    }
}

// please test after import the data
func TestSearchResult(t *testing.T) {
    docDB := doc.NewDocDriver(
        storage.NewLevelDBStorage("../../data/db/doc"),
    )
    indexDB := index.NewIndexDriver(
        storage.NewLevelDBStorage("../../data/db/index"),
    )
    token := tokenize.NewToken()
    defer indexDB.CloseIndexDB()
    defer docDB.CloseDocDB()
    newSearch := NewSearch(indexDB, docDB, token)
    start := time.Now()
    results := SearchResult("开个玩笑", newSearch)
    duration := time.Since(start)
    fmt.Println("time: ", duration.Milliseconds(), "ms")
    for _, result := range results {
        fmt.Println("DocId: ", result.DocId)
        fmt.Println("Score: ", result.Score)
        fmt.Println("Url: ", result.Url)
        fmt.Println("Text: ", result.Text)
    }
}

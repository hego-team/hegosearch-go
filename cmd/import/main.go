package main

import (
    "bufio"
    "encoding/csv"
    "fmt"
    "hegosearch/data"
    "hegosearch/data/model"
    "io"
    "log"
    "os"
    "time"
)

// import the csv data from dataset

var (
    filenames []string
)

func Init() {
    filenameInfo, err := os.Open("data/file")
    if err != nil {
        panic(err)
    }
    reader := bufio.NewReader(filenameInfo)
    for {
        line, _, err := reader.ReadLine()
        if err != nil {
            if err == io.EOF {
                break
            }
        }
        filenames = append(filenames, string(line))
    }
}

func main() {
    Init()
    data.JiebaInit()
    csvfile, err := os.Open("data/dataset/wukong_100m_1.csv")
    if err != nil {
        panic(err)
    }
    defer csvfile.Close()

    csvReader := csv.NewReader(csvfile)
    docDB := data.DocDataInit("data/db/doc")
    indexDB := data.IndexDataInit("data/db/index")
    defer docDB.DocDB.Close()
    defer indexDB.IndexDB.Close()
    _, err = csvReader.Read()
    if err == io.EOF {
        log.Fatalf("read first error")
    }
    start := time.Now()
    count := 0

    wordMap := make(map[string][]uint64)
    for {
        count++
        record, err := csvReader.Read()
        if err == io.EOF {
            break
        }
        if err != nil {
            log.Fatal(err)
        }
        // fmt.Printf("image: %s,text: %s \n", record[0], record[1])
        if len(record[0]) == 0 || len(record[1]) == 0 {
            continue
        }
        doc := model.Document{
            Url:  record[0],
            Text: record[1],
        }
        id, err := docDB.InsertIntoDocDB(&doc)
        if err != nil {
            panic(err)
        }
        words := data.PartWord(doc.Text)

        for i := range words {
            if ids, ok := wordMap[words[i]]; ok {
                wordMap[words[i]] = append(ids, id)
            } else {
                wordMap[words[i]] = []uint64{id}
            }
        }
        if count%10000 == 0 {
            cost := time.Since(start)
            fmt.Println("import cost:", cost.Seconds(), "s")
            start = time.Now()
        }
    }
    for k, v := range wordMap {
        err := indexDB.InsertIdsIntoIndexDB(v, k)
        if err != nil {
            panic(err)
        }
    }
    fmt.Println("import count:", count)
}

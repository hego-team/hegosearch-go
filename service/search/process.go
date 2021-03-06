package search

import (
    "hegosearch/data/model"
    "log"
    "sort"
    "time"
)

type DocResult struct {
    Score     float64
    WordCount uint64
}

type SearchContent struct {
    DocId uint64
    Score float64
    Url   string
    Text  string
}

type ProcessParams struct {
    Limit    uint64
    Text     string
    StopWord string
}

type SearchResp struct {
    Content []*SearchContent
    Time    int64
    Count   uint64
}

func SearchText(req *model.SearchReq, engine *Search) []*ProcessResult {
    words := engine.Tokenize.PartWord(req.Text)
    allScoreMap := make(map[uint64]*DocResult)
    for i := range words {
        scoreMap, err := engine.SearchKey(words[i])
        if err != nil {
            log.Println("when search word ", words[i], "error")
            continue
        }
        for key := range scoreMap {
            if _, ok := allScoreMap[key]; ok {
                allScoreMap[key].WordCount++
                allScoreMap[key].Score = allScoreMap[key].Score + scoreMap[key]
            } else {
                allScoreMap[key] = &DocResult{
                    Score:     scoreMap[key],
                    WordCount: 1,
                }
            }
        }
        for key := range allScoreMap {
            allScoreMap[key].Score = allScoreMap[key].Score * float64(allScoreMap[key].WordCount) / float64(len(words))
        }
    }
    processResults := make([]*ProcessResult, len(allScoreMap))
    index := 0
    for key := range allScoreMap {
        processResults[index] = &ProcessResult{
            Score: allScoreMap[key].Score,
            DocId: key,
        }
        index++
    }
    // sort by the score desc
    sort.Sort(ScoreSlice(processResults))
    return processResults
}

func SearchResult(req *model.SearchReq, engine *Search) *SearchResp {
    start := time.Now()
    proRes := SearchText(req, engine)
    length := len(proRes)
    if uint64(length) > req.Limit {
        length = int(req.Limit)
    }
    res := make([]*SearchContent, length)
    count := uint64(0)
    for i, result := range proRes {
        docRes, err := engine.docDB.FindFromDocDB(result.DocId)
        if err != nil {
            log.Println("Search error , when search the id from the docDB", err)
            break
        }
        res[i] = &SearchContent{
            DocId: result.DocId,
            Score: result.Score,
            Url:   docRes.Url,
            Text:  docRes.Text,
        }
        count++
        if count >= req.Limit {
            break
        }
    }
    since := time.Since(start).Milliseconds()
    return &SearchResp{
        Content: res,
        Time:    since,
        Count:   uint64(len(proRes)),
    }
}

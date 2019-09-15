package anime

import (
    "fmt"
    "strings"
    "time"
    "math/rand"
    "github.com/tealeg/xlsx"
    "douban_spider/collector"
)

const MAX_SIZE_OF_ANIME_LIST = 10000

type Anime struct {
    Name string          `json:anime_name`
    DoubanIds []string   `json:douban_ids`
}

type AnimeList struct {
    Animes []*Anime       `json:anime_list`
}

func (al *AnimeList) ReadXlsx(file_path string) error {
    xlsxFileHandler, err := xlsx.OpenFile(file_path)
    if err != nil {
        fmt.Println("Open file %s fail, error : %s", file_path, err.Error())
        return err
    }
    al.Animes = make([]*Anime, 0, MAX_SIZE_OF_ANIME_LIST)

    for _, sheet := range xlsxFileHandler.Sheets {
        fmt.Println("Sheet name :", sheet.Name)
        for _, row := range sheet.Rows {
            cell_num := len(row.Cells)
            if cell_num != 2 {
                continue
            }
            animeName := row.Cells[0].String()
            if animeName == "" {
                //empty anime name
                continue
            }
            ids := row.Cells[1].String()
            if ids == "none" || ids == "" {
                fmt.Println("Find", animeName, " has no douban ids, it will be ignored")
                continue
            }
            idList := strings.Split(ids, ",")
            var anime = &Anime{}
            anime.Name = animeName
            anime.DoubanIds = idList
            al.Animes = append(al.Animes, anime)
        }
    }
    return nil
}

func (ani *Anime) CollectAnime(useAccount bool, userName string, password string) error {
    if len(ani.DoubanIds) == 0 {
        return nil
    }
    fmt.Println("Collecting ", ani.Name, ani.DoubanIds)

    var collector = &collector.DoubanMovieCommentCollector{}

    if useAccount == true {
        collector.UseDoubanAccount(userName, password)
    } else {
        collector.CancelUseAccount()
    }

    for _, did := range ani.DoubanIds {
        collector.MovieID = did
        collector.MovieName = ani.Name
        collector.FetchMovieComment()
        time.Sleep(time.Duration(rand.Intn(5)) * time.Second)
    }
    return nil
}

func (al *AnimeList) CollectAnimeList(useAccount bool, userName string, password string) error {
    if len(al.Animes) == 0 {
        return nil
    }
    total := len(al.Animes)
    for i, anime := range al.Animes {
        fmt.Println("Now handling anime: anime.Name ", i + 1, "/", total)
        anime.CollectAnime(useAccount, userName, password)
        time.Sleep(time.Duration(rand.Intn(4)) * time.Second)
    }
    return nil
}

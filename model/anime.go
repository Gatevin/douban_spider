package model

import (
    "fmt"
    "strings"
    "time"
    //"math/rand"
    "github.com/tealeg/xlsx"
    "douban_spider/collector"
    . "douban_spider/utils"
)

const MAX_SIZE_OF_ANIME_LIST = 10000

type Anime struct {
    Name string          `json:"anime_name"`
    DoubanIds []string   `json:"douban_ids"`
}

type AnimeList struct {
    Animes []*Anime       `json:"anime_list"`
}

func (al *AnimeList) ReadXlsx(file_path string) error {
    xlsxFileHandler, err := xlsx.OpenFile(file_path)
    if err != nil {
        fmt.Println("Open file ", file_path, " failed with error : ", err.Error())
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

func (al *AnimeList) CollectAnimeList(useAccount bool, userName string, password string) error {
    if len(al.Animes) == 0 {
        return nil
    }
    Ipmgr.Prepare()
    Ipmgr.FetchIpList()
    Ipmgr.PrintPoolInfo()
    animeIDCount := 0
    refreshIp := 5
    total := len(al.Animes)
    for i, ani := range al.Animes {
        fmt.Println("Now handling anime: anime.Name ", i + 1, "/", total)
        if len(ani.DoubanIds) == 0 {
            return nil
        }
        fmt.Println("Collecting ", ani.Name, ani.DoubanIds)
        for _, did := range ani.DoubanIds {
            //拿个高匿名ip
            ip := Ipmgr.GetAnonymousIpWithIndex(animeIDCount)
            animeIDCount += 1
            for {
                if ip != nil {
                    break
                }
                fmt.Println("ip == nil in GetAnonymousIpWithIndex occured, animeIDCount = ", animeIDCount)
                ip = Ipmgr.GetAnonymousIpWithIndex(animeIDCount)
                animeIDCount += 1
            }
            fmt.Println(ip.Address)
            var collector = &collector.DoubanMovieCommentCollector{}
            if useAccount == true {
                collector.UseDoubanAccount(userName, password)
            } else {
                collector.CancelUseAccount()
            }
            collector.UseAnonymousIp(ip)
            
            collector.MovieID = did
            collector.MovieName = ani.Name
            
            fmt.Println("collector ", did, ani.Name, " ip address: ", ip)
            collector.ConfigCollyRule()
            collector.FetchMovieComment()

            // ip用完了返还回去
            Ipmgr.ReturnAnonymousIp(ip)
            //time.Sleep(time.Duration(rand.Intn(10)) * time.Second)

            //if animeIDCount >= 20 {
            //    return nil
            //}
            
            if animeIDCount >= refreshIp {
                for {
                    refreshIp += 20
                    if animeIDCount < refreshIp {
                        break
                    }
                }
                fmt.Println("ip pool is updating")
                Ipmgr.FetchIpList()
                Ipmgr.PrintPoolInfo()
            }
            
            time.Sleep(time.Duration(10*time.Second))
        }
        time.Sleep(time.Duration(20 * time.Second))
    }
    return nil
}

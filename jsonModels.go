package main

type Summoner struct {
	Id            int    `json:"id"`
	Name          string `json:"name"`
	ProfileIconId int    `json:"profileIconId"`
	SummonerLevel int    `json:"summonerLevel"`
	RevisionDate  int64  `json:"revisionDate"`
}

package models

type SongQuery struct {
	Title  string
	Artist string
}

type Artist struct {
	ID     int
	Name   string
	IsMain bool
}

type SongNode struct {
	ID      int
	Title   string
	Artists []Artist
	Genres  []string
}

type SearchResult struct {
	SourceSong  SongNode
	MatchedSong SongNode
	Distance    int
	Path        []SongNode
}

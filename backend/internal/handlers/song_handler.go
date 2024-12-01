package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/Emeruem-Kennedy1/ghopper/internal/models"
	"github.com/Emeruem-Kennedy1/ghopper/internal/repository"
	"github.com/Emeruem-Kennedy1/ghopper/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/zmb3/spotify"
)

type SongSearchRequest struct {
	Songs    []models.SongQuery `json:"songs"`
	Genre    string             `json:"genre"`
	MaxDepth int                `json:"maxDepth"`
}

type TopTracksAnalysisRequest struct {
	Genre string `json:"genre"`
}

type TopTrackResponseSong struct {
	Title    string `json:"title"`
	Artist   string `json:"artist"`
	URL      string `json:"url"`
	ImageURL string `json:"imageURL"`
}
type TopTracksAnalysisResponse struct {
	Songs    []TopTrackResponseSong `json:"songs"`
	Playlist string                 `json:"playlist"`
}

type GraphResponse struct {
	// Map of song ID to its connections
	AdjacencyList map[string]map[string]interface{} `json:"adjacencyList"`
	// Map of song ID to song details
	Nodes map[string]SongNode `json:"nodes"`
	// Original path information
	Paths []PathInfo `json:"paths"`
}

type ArtistInfo struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	IsMain bool   `json:"isMain"`
}

type SongNode struct {
	ID      int          `json:"id"`
	Title   string       `json:"title"`
	Artists []ArtistInfo `json:"artists"`
	Genres  []string     `json:"genres"`
}

type PathInfo struct {
	Start     string   `json:"start"`     // Starting song ID
	End       string   `json:"end"`       // Ending song ID
	PathNodes []string `json:"pathNodes"` // List of song IDs in path
	Distance  int      `json:"distance"`
}

func transformToArtistInfo(artists []models.Artist) []ArtistInfo {
	result := make([]ArtistInfo, len(artists))
	for i, artist := range artists {
		result[i] = ArtistInfo{
			ID:     artist.ID,
			Name:   artist.Name,
			IsMain: artist.IsMain,
		}
	}
	return result
}

func AnalyzeSongsGivenGenre(songRepo *repository.SongRepository, clientManager *services.ClientManager, spotifyService *services.SpotifyService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userID, exists := ctx.Get("userID")
		if !exists {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}
		client, clientExists := clientManager.GetClient(userID.(string))
		if !clientExists {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		var req TopTracksAnalysisRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request format"})
			return
		}

		if req.Genre == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "genre is required"})
			return
		}

		limit := 30
		timeRange := "long"

		tracks, err := client.CurrentUsersTopTracksOpt(&spotify.Options{Limit: &limit, Timerange: &timeRange})
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get user's top tracks"})
			return
		}

		// store the tracks as SongQuery
		songs := make([]models.SongQuery, len(tracks.Tracks))
		for _, track := range tracks.Tracks {
			songs = append(songs, models.SongQuery{
				Title:  track.Name,
				Artist: track.Artists[0].Name,
			})
		}

		analysisResults, err := songRepo.FindSongsByGenreBFS(songs, req.Genre, 2)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to analyze songs"})
			return
		}

		songResults := make([]models.SongQuery, 0, len(analysisResults))
		for _, result := range analysisResults {
			song := models.SongQuery{
				Title:  result.MatchedSong.Title,
				Artist: result.MatchedSong.Artists[0].Name,
			}
			songResults = append(songResults, song)
		}

		// find the urls of the songs
		var response TopTracksAnalysisResponse
		topTrackSongs := make([]TopTrackResponseSong, 0, len(songResults))
		var songIDs []spotify.ID

		for _, song := range songResults {
			url, err := spotifyService.GetSongURL(userID.(string), song.Title, song.Artist)

			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get song url"})
				return
			}

			if url != "" {
				// get the song ID
				songID := strings.Split(url, "/")[len(strings.Split(url, "/"))-1]
				songIDs = append(songIDs, spotify.ID(strings.TrimSpace(songID)))

				topTrackSongs = append(topTrackSongs, TopTrackResponseSong{
					Title:  song.Title,
					Artist: song.Artist,
					URL:    url,
				})
			}
		}

		// add the songs to a playlist
		playlistName := fmt.Sprintf("Explore %s songs", req.Genre)
		playlistDescription := fmt.Sprintf("Playlist of songs in the genre %s", req.Genre)
		playlistURL, err := spotifyService.CreatePlaylistFromSongs(userID.(string), songIDs, playlistName, playlistDescription)

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create playlist"})
			return
		}

		response = TopTracksAnalysisResponse{
			Songs:    topTrackSongs,
			Playlist: playlistURL,
		}

		ctx.JSON(http.StatusOK, response)
	}
}

func SearchSongByGenre(songRepo *repository.SongRepository) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req SongSearchRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request format"})
			return
		}

		if req.Genre == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "genre is required"})
			return
		}

		if req.MaxDepth <= 0 {
			req.MaxDepth = 5 // Default max depth
		}

		results, err := songRepo.FindSongsByGenreBFS(req.Songs, req.Genre, req.MaxDepth)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to search songs"})
			return
		}

		graphResponse := GraphResponse{
			AdjacencyList: make(map[string]map[string]interface{}),
			Nodes:         make(map[string]SongNode),
			Paths:         make([]PathInfo, 0),
		}

		for _, result := range results {
			// Add source and matched songs to nodes
			sourceID := fmt.Sprintf("%d", result.SourceSong.ID)
			matchedID := fmt.Sprintf("%d", result.MatchedSong.ID)

			// Add nodes
			if _, exists := graphResponse.Nodes[sourceID]; !exists {
				graphResponse.Nodes[sourceID] = SongNode{
					ID:      result.SourceSong.ID,
					Title:   result.SourceSong.Title,
					Artists: transformToArtistInfo(result.SourceSong.Artists),
					Genres:  result.SourceSong.Genres,
				}
			}

			if _, exists := graphResponse.Nodes[matchedID]; !exists {
				graphResponse.Nodes[matchedID] = SongNode{
					ID:      result.MatchedSong.ID,
					Title:   result.MatchedSong.Title,
					Artists: transformToArtistInfo(result.MatchedSong.Artists),
					Genres:  result.MatchedSong.Genres,
				}
			}

			// Build adjacency list from path
			for i := 0; i < len(result.Path)-1; i++ {
				currentID := fmt.Sprintf("%d", result.Path[i].ID)
				nextID := fmt.Sprintf("%d", result.Path[i+1].ID)

				// Add nodes from path
				if _, exists := graphResponse.Nodes[currentID]; !exists {
					graphResponse.Nodes[currentID] = SongNode{
						ID:      result.Path[i].ID,
						Title:   result.Path[i].Title,
						Artists: transformToArtistInfo(result.Path[i].Artists),
						Genres:  result.Path[i].Genres,
					}
				}

				// Add to adjacency list (bidirectional)
				if graphResponse.AdjacencyList[currentID] == nil {
					graphResponse.AdjacencyList[currentID] = make(map[string]interface{})
				}
				if graphResponse.AdjacencyList[nextID] == nil {
					graphResponse.AdjacencyList[nextID] = make(map[string]interface{})
				}
				graphResponse.AdjacencyList[currentID][nextID] = struct{}{}
				graphResponse.AdjacencyList[nextID][currentID] = struct{}{}
			}

			// Add path info
			pathNodes := make([]string, len(result.Path))
			for i, node := range result.Path {
				pathNodes[i] = fmt.Sprintf("%d", node.ID)
			}

			graphResponse.Paths = append(graphResponse.Paths, PathInfo{
				Start:     sourceID,
				End:       matchedID,
				PathNodes: pathNodes,
				Distance:  result.Distance,
			})
		}

		ctx.JSON(http.StatusOK, graphResponse)
	}
}

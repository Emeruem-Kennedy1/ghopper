package handlers

import (
	"fmt"
	"math/rand"
	"net/http"
	"strings"

	"github.com/Emeruem-Kennedy1/ghopper/internal/models"
	"github.com/Emeruem-Kennedy1/ghopper/internal/repository"
	"github.com/Emeruem-Kennedy1/ghopper/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/zmb3/spotify"
	"go.uber.org/zap"
)

var genreGroups = map[string]string{
	"hip-hop":    "Hip-Hop / Rap / R&B",
	"rap":        "Hip-Hop / Rap / R&B",
	"r&b":        "Hip-Hop / Rap / R&B",
	"electronic": "Electronic / Dance",
	"dance":      "Electronic / Dance",
	"rock":       "Rock / Pop",
	"pop":        "Rock / Pop",
	"soul":       "Soul / Funk / Disco",
	"funk":       "Soul / Funk / Disco",
	"disco":      "Soul / Funk / Disco",
	"jazz":       "Jazz / Blues",
	"blues":      "Jazz / Blues",
	"reggae":     "Reggae / Dub",
	"dub":        "Reggae / Dub",
	"country":    "Country / Folk",
	"folk":       "Country / Folk",
	"world":      "World / Latin",
	"latin":      "World / Latin",
	"soundtrack": "Soundtrack / Library",
	"library":    "Soundtrack / Library",
	"classical":  "Classical",
}

// Representative search genre for each group (for database queries)
var searchGenres = map[string]string{
	"Hip-Hop / Rap / R&B":  "hip-hop", // We can use any of hip-hop, rap, or r&b
	"Electronic / Dance":   "electronic",
	"Rock / Pop":           "rock",
	"Soul / Funk / Disco":  "soul",
	"Jazz / Blues":         "jazz",
	"Reggae / Dub":         "reggae",
	"Country / Folk":       "country",
	"World / Latin":        "world",
	"Soundtrack / Library": "soundtrack",
	"Classical":            "classical",
}

var defaultGenrePlaylists = map[string][]string{
	"hip-hop": {
		"https://open.spotify.com/playlist/37i9dQZF1DXbkfWVLd8wE3?si=zqZ10XC9S2a095CXo8vW6Q",
		"https://open.spotify.com/playlist/0h9Gaqt2sNJ8M5aMV3h9BO?si=eB4jwKD0RFKXFz9WUoqXkA",
		"https://open.spotify.com/playlist/37i9dQZF1DX04mASjTsvf0?si=r6S-aGh0Q7a2gb3nCFLemQ",
	},
	"electronic": {
		"https://open.spotify.com/playlist/37i9dQZF1DWZBCPUIUs2iR?si=0REkYg79Sp6ghVwZnO0D6Q",
		"https://open.spotify.com/playlist/44dFP8mNyCi3UcBlyaRICH?si=_bljU3wzSH6h7tfYIGQ0cw",
	},
	"rock": {
		"https://open.spotify.com/playlist/37i9dQZF1DWXRqgorJj26U?si=nYX-RaKmTOqRcq73fyGEnQ",
		"https://open.spotify.com/playlist/37i9dQZF1EIctsc1CJao2L?si=6MJxt2cpQeyYtAMYg0RKJw",
	},
	"soul": {
		"https://open.spotify.com/playlist/73sIU7MIIIrSh664eygyjm?si=uKn9DEovQSqxb9P4l5u7RQ",
		"https://open.spotify.com/playlist/37i9dQZF1DWWvhKV4FBciw?si=1lzYUcgJR-6DBpoluxX2OA",
		"https://open.spotify.com/playlist/37i9dQZF1DX1MUPbVKMgJE?si=0nXq73dyRYGzdrrzEpd2Gw",
	},
	"jazz": {
		"https://open.spotify.com/playlist/4pIwPQAiZk4JGiWRzAaxwK?si=gxXRMeGQS_SeqeTLhTIGbQ",
		"https://open.spotify.com/playlist/0A1IHcqjyImN9uoHRsVtBn?si=f1z51sRCRx6gd7dPEfg5_g",
	},
	"reggae": {
		"https://open.spotify.com/playlist/37i9dQZF1EQpjs4F0vUZ1x?si=03B58nZISEeLoLDvbZ7jbw",
		"https://open.spotify.com/playlist/7AI62FuDUugLcg1IyVgMwU?si=GiAjkEvDT2mdloDn-zsKpQ",
	},
	"country": {
		"https://open.spotify.com/playlist/0QFaFgDQQiKBob7VIZIilG?si=kjNuzW4iS86D9Udo0VSeUA",
		"https://open.spotify.com/playlist/37i9dQZF1DWVmps5U8gHNv?si=ElAZI6eRS9OWr0EMMlh2Ow",
	},
	"world": {
		"https://open.spotify.com/playlist/37i9dQZF1DXcIme26eJxid?si=F382j4bBTjSa8oS2NC3C_w",
		"https://open.spotify.com/playlist/37i9dQZF1DX6ThddIjWuGT?si=B0r6U0sNQlS8DD1xC4g9ng",
	},
	"soundtrack": {
		"https://open.spotify.com/playlist/3vDe8D64ytZRKXt0AsJT0B?si=4HFGEkdPRwSPajsaU23T2A",
	},
	"classical": {
		"https://open.spotify.com/playlist/2AIyLES2xJfPa6EOxmKySl?si=qpeStpIiQO--nr__oaFQCg",
	},
}

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

type DeletePlaylistRequest struct {
	PlaylistID string `json:"playlistID"`
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

// Helper function to normalize genre for playlist creation and display
func normalizeGenre(genre string) string {
	normalizedGenre := strings.ToLower(strings.TrimSpace(genre))
	if groupName, exists := genreGroups[normalizedGenre]; exists {
		return groupName
	}
	return normalizedGenre
}

// Helper function to get the database search genre
func getSearchGenre(groupedGenre string) string {
	if searchGenre, exists := searchGenres[groupedGenre]; exists {
		return searchGenre
	}
	return strings.ToLower(strings.TrimSpace(groupedGenre))
}

func getRandomPlaylist(genre string) string {
	if playlists, exists := defaultGenrePlaylists[genre]; exists && len(playlists) > 0 {
		return playlists[rand.Intn(len(playlists))]
	}
	return "" // Return empty string if no playlist found
}

func AnalyzeSongsGivenGenre(songRepo repository.SongRepositoryInterface, clientManager services.ClientManagerInterface, spotifyService services.SpotifyServiceInterface) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userID, exists := ctx.Get("userID")
		if !exists {
			zap.L().Warn("Unauthorized attempt to analyze songs")
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}
		client, clientExists := clientManager.GetClient(userID.(string))
		if !clientExists {
			zap.L().Warn("No Spotify client found for user",
				zap.String("userID", userID.(string)))
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		var req TopTracksAnalysisRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			zap.L().Error("Invalid request format",
				zap.Error(err),
				zap.String("userID", userID.(string)),
			)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request format"})
			return
		}

		if req.Genre == "" {
			zap.L().Error("Genre not provided",
				zap.String("userID", userID.(string)),
			)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "genre is required"})
			return
		}

		limit := 50
		timeRange := "short"

		tracks, err := client.CurrentUsersTopTracksOpt(&spotify.Options{Limit: &limit, Timerange: &timeRange})
		if err != nil {
			zap.L().Error("Failed to fetch top tracks from Spotify",
				zap.String("userID", userID.(string)),
				zap.Error(err))

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

		if len(songs) == 0 {
			zap.L().Warn("No songs found for user",
				zap.String("userID", userID.(string)))

			ctx.JSON(http.StatusNotFound, gin.H{"error": "no songs found for user"})
			return
		}

		// Normalize genre for playlist creation and UI display
		normalizedGenre := normalizeGenre(req.Genre)
		// Get the single search genre for database lookup
		searchGenre := getSearchGenre(normalizedGenre)

		analysisResults, err := songRepo.FindSongsByGenreBFS(songs, searchGenre, 2)
		if err != nil {
			zap.L().Error("Failed to analyze songs",
				zap.String("userID", userID.(string)),
				zap.Error(err))

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
				zap.L().Error("Failed to get song URL",
					zap.String("userID", userID.(string)),
					zap.Error(err))

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

		if len(songIDs) == 0 {
			response = TopTracksAnalysisResponse{
				Songs:    topTrackSongs,
				Playlist: getRandomPlaylist(searchGenre),
			}
			zap.L().Info("No songs found for genre",
				zap.String("userID", userID.(string)),
				zap.String("genre", req.Genre),
			)
			ctx.JSON(http.StatusOK, response)
			return
		}

		// add the songs to a playlist
		playlistName := fmt.Sprintf("Explore %s songs", normalizedGenre)
		playlistDescription := fmt.Sprintf("Playlist of songs in the genre %s", normalizedGenre)
		playlistURL, err := spotifyService.CreatePlaylistFromSongs(userID.(string), songIDs, playlistName, playlistDescription)

		if err != nil {
			zap.L().Error("Failed to create playlist",
				zap.String("userID", userID.(string)),
				zap.Error(err))

			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create playlist"})
			return
		}

		response = TopTracksAnalysisResponse{
			Songs:    topTrackSongs,
			Playlist: playlistURL,
		}

		zap.L().Info("Successfully analyzed songs",
			zap.String("userID", userID.(string)),
			zap.Any("songs", songResults),
			zap.String("genre", req.Genre),
		)
		ctx.JSON(http.StatusOK, response)
	}
}

func SearchSongByGenre(songRepo repository.SongRepositoryInterface) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req SongSearchRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			zap.L().Error("Invalid request format",
				zap.Error(err))

			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request format"})
			return
		}

		if req.Genre == "" {
			zap.L().Error("Genre not provided")
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "genre is required"})
			return
		}

		if req.MaxDepth <= 0 {
			req.MaxDepth = 5 // Default max depth
		}

		results, err := songRepo.FindSongsByGenreBFS(req.Songs, req.Genre, req.MaxDepth)
		if err != nil {
			zap.L().Error("Failed to search songs",
				zap.Error(err),
			)
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

		zap.L().Info("Successfully searched songs",
			zap.Any("songs", req.Songs),
			zap.String("genre", req.Genre),
		)
		ctx.JSON(http.StatusOK, graphResponse)
	}
}

func DeletePlaylist(spotifyService services.SpotifyServiceInterface, spotifySongRepo repository.SpotifySongRepositoryInterface) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userID, exists := ctx.Get("userID")
		if !exists {
			zap.L().Warn("Unauthorized attempt to delete playlist")
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		playlistID := ctx.Param("playlistID")

		if playlistID == "" {
			zap.L().Error("Playlist ID not provided")
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "playlistID is required"})
			return
		}

		playlist, err := spotifySongRepo.FindPlaylistByIDAndUser(playlistID, userID.(string))
		if err != nil {
			zap.L().Error("Failed to find playlist",
				zap.String("userID", userID.(string)),
				zap.String("playlistID", playlistID),
				zap.Error(err))

			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to find playlist"})
			return
		}

		if playlist == nil {
			zap.L().Warn("Playlist not found",
				zap.String("userID", userID.(string)),
				zap.String("playlistID", playlistID))

			ctx.JSON(http.StatusNotFound, gin.H{"error": "playlist not found"})
			return
		}

		err = spotifySongRepo.DeletePlaylist(playlist)
		if err != nil {
			fmt.Printf("Failed to delete playlist: %v\n", err)
			zap.L().Error("Failed to delete playlist",
				zap.String("userID", userID.(string)),
				zap.String("playlistID", playlistID),
				zap.Error(err))

			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete playlist"})
			return
		}

		err = spotifyService.DeletePlaylist(userID.(string), playlist.ID)
		if err != nil {
			zap.L().Error("Failed to delete playlist from Spotify",
				zap.String("userID", userID.(string)),
				zap.String("playlistID", playlistID),
				zap.Error(err))

			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete playlist from Spotify"})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"message": "playlist deleted"})
	}
}

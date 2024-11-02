package repository

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	"github.com/Emeruem-Kennedy1/ghopper/internal/models"
)

type SongRepository struct {
	db *sql.DB
}

func NewSongRepository(db *sql.DB) *SongRepository {
	return &SongRepository{db: db}
}

func (r *SongRepository) GetSongIDsByTitleAndArtist(title, artist string) ([]int, error) {
	query := `
        SELECT DISTINCT s.id
        FROM Song s
        JOIN SongArtist sa ON s.id = sa.songId
        JOIN Artist a ON sa.artistId = a.id
        WHERE s.title = ? AND a.name = ?
    `

	rows, err := r.db.Query(query, title, artist)

	if err != nil {
		return nil, fmt.Errorf("error getting song ids by title and artist: %v", err)
	}

	defer rows.Close()

	var songIDs []int

	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("error scanning song id: %v", err)
		}
		songIDs = append(songIDs, id)
	}

	return songIDs, nil
}

func (r *SongRepository) GetSongWithDetails(SongID int) (*models.SongNode, error) {

	query := `
		SELECT 
			s.id,
			s.title,
			a.name as artist_name,
			a.id as artist_id,
			GROUP_CONCAT(DISTINCT g.name) as genres
		FROM Song s
		JOIN SongArtist sa ON s.id = sa.songId
		JOIN Artist a ON sa.artistId = a.id
		LEFT JOIN _SongToGenre sg ON sg.B = s.id
		LEFT JOIN Genre g ON g.id = sg.A
		WHERE s.id = ?
		GROUP BY s.id
	`

	var genres sql.NullString
	var artistName string
	var artistID int
	song := &models.SongNode{
		ID: SongID,
	}

	err := r.db.QueryRow(query, SongID).Scan(
		&song.ID,
		&song.Title,
		&artistName,
		&artistID,
		&genres,
	)

	if err != nil {
		return nil, fmt.Errorf("error getting song: %v", err)
	}

	if genres.Valid {
		song.Genres = strings.Split(genres.String, ",")
	}

	artistQuery := `
        SELECT 
            a.id,
            a.name,
            sa.isMainArtist
        FROM SongArtist sa
        JOIN Artist a ON sa.artistId = a.id
        WHERE sa.songId = ?
        ORDER BY sa.isMainArtist DESC
    `

	artistRows, err := r.db.Query(artistQuery, SongID)
	if err != nil {
		return nil, fmt.Errorf("error getting artists: %v", err)
	}
	defer artistRows.Close()

	for artistRows.Next() {
		var artist models.Artist
		if err := artistRows.Scan(&artist.ID, &artist.Name, &artist.IsMain); err != nil {
			return nil, fmt.Errorf("error scanning artist: %v", err)
		}
		song.Artists = append(song.Artists, artist)
	}

	return song, nil
}

func (r *SongRepository) GetAllSampledSongs(songID int) ([]int, error) {
	query := `
        WITH RECURSIVE SameSongs AS (
            -- Base case: Find the original song and its duplicates
            SELECT id
            FROM Song s1
            WHERE EXISTS (
                SELECT 1 FROM Song s2
                WHERE s2.id = ?
                AND s1.title = s2.title
				AND ((s1.releaseYear IS NULL AND s2.releaseYear IS NULL) OR s1.releaseYear = s2.releaseYear)
            )
        )
        SELECT DISTINCT sampled_song_id
        FROM (
            -- Songs that sample our song
            SELECT sampled_in_song_id as sampled_song_id
            FROM SameSongs ss
            JOIN Sample s ON ss.id = s.original_song_id
            
            UNION
            
            -- Songs that our song samples
            SELECT original_song_id as sampled_song_id
            FROM SameSongs ss
            JOIN Sample s ON ss.id = s.sampled_in_song_id
        ) all_samples
    `

	rows, err := r.db.Query(query, songID)
	if err != nil {
		return nil, fmt.Errorf("error getting sampled songs: %v", err)
	}
	defer rows.Close()

	var sampledSongs []int
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		sampledSongs = append(sampledSongs, id)
	}

	return sampledSongs, nil
}

// func (r *SongRepository) FindSongsByGenreBFS(songQueries []models.SongQuery, targetGenre string, maxDepth int) ([]models.SearchResult, error) {
// 	var results []models.SearchResult
// 	visited := make(map[string]bool)

// 	type QueueItem struct {
// 		SongID   int
// 		Distance int
// 		Path     []models.SongNode
// 	}

// 	getMainArtist := func(artists []models.Artist) string {
// 		for _, artist := range artists {
// 			if artist.IsMain {
// 				return artist.Name
// 			}
// 		}
// 		return ""
// 	}

// 	for _, query := range songQueries {
// 		startIDs, err := r.GetSongIDsByTitleAndArtist(query.Title, query.Artist)
// 		if err != nil {
// 			return nil, err
// 		}

// 		for _, startID := range startIDs {
// 			startSong, err := r.GetSongWithDetails(startID)
// 			if err != nil {
// 				return nil, err
// 			}

// 			mainArtist := getMainArtist(startSong.Artists)
// 			visitKey := fmt.Sprintf("%s-%s", startSong.Title, mainArtist)

// 			if visited[visitKey] {
// 				continue
// 			}

// 			visited[visitKey] = true

// 			queue := []QueueItem{
// 				{
// 					SongID:   startID,
// 					Distance: 0,
// 					Path:     []models.SongNode{*startSong},
// 				},
// 			}

// 			// BFS implementation
// 			// Just ensure we're using the new visitKey format consistently
// 			for len(queue) > 0 {
// 				current := queue[0]
// 				queue = queue[1:]

// 				if current.Distance > maxDepth {
// 					continue
// 				}

// 				currentSong, err := r.GetSongWithDetails(current.SongID)
// 				if err != nil {
// 					continue
// 				}

// 				containsGenre := func(genres []string, targetGenre string) bool {
// 					normalizedTarget := strings.TrimSpace(strings.ToLower(targetGenre))
// 					for _, genre := range genres {
// 						if strings.TrimSpace(strings.ToLower(genre)) == normalizedTarget {
// 							return true
// 						}
// 					}
// 					return false
// 				}

// 				// Then in your BFS function:
// 				if containsGenre(currentSong.Genres, targetGenre) {
// 					results = append(results, models.SearchResult{
// 						SourceSong:  *startSong,
// 						MatchedSong: *currentSong,
// 						Distance:    current.Distance,
// 						Path:        current.Path,
// 					})
// 					// continue
// 				}

// 				// all sampled songs
// 				sampledSongs, err := r.GetAllSampledSongs(current.SongID)
// 				if err != nil {
// 					return nil, err
// 				}

// 				for _, sampledSongID := range sampledSongs {
// 					sampledSong, err := r.GetSongWithDetails(sampledSongID)
// 					if err != nil {
// 						return nil, err
// 					}

// 					sampleMainArtist := getMainArtist(sampledSong.Artists)
// 					visitKey := fmt.Sprintf("%s-%s", sampledSong.Title, sampleMainArtist)

// 					if !visited[visitKey] {
// 						visited[visitKey] = true
// 						newPath := append([]models.SongNode{}, current.Path...)
// 						newPath = append(newPath, *sampledSong)

// 						queue = append(queue, QueueItem{
// 							SongID:   sampledSongID,
// 							Distance: current.Distance + 1,
// 							Path:     newPath,
// 						})
// 					}
// 				}

// 			}

// 		}
// 	}

// 	return results, nil
// }

func (r *SongRepository) FindSongsByGenreBFS(songQueries []models.SongQuery, targetGenre string, maxDepth int) ([]models.SearchResult, error) {
	var startConditions []string
	var params []interface{}

	for _, query := range songQueries {
		startConditions = append(startConditions, "(s.title = ? AND a.name = ?)")
		params = append(params, query.Title, query.Artist)
	}

	query := `
        WITH RECURSIVE SongPath AS (
            -- Base case: start with input songs
            SELECT 
                s.id,
                s.id as source_id,
                0 as distance,
                CAST(CONCAT('[', s.id, ']') AS CHAR(1000)) as path
            FROM Song s
            JOIN SongArtist sa ON s.id = sa.songId
            JOIN Artist a ON sa.artistId = a.id
            WHERE ` + strings.Join(startConditions, " OR ") + `
            GROUP BY s.id

            UNION ALL

            -- Recursive case: follow sampling relationships
            SELECT 
                CASE 
                    WHEN sam.sampled_in_song_id = sp.id THEN sam.original_song_id
                    ELSE sam.sampled_in_song_id
                END,
                sp.source_id,
                sp.distance + 1,
                CONCAT(sp.path, ',', CASE 
                    WHEN sam.sampled_in_song_id = sp.id THEN sam.original_song_id
                    ELSE sam.sampled_in_song_id
                END)
            FROM SongPath sp
            JOIN Sample sam ON sp.id = sam.original_song_id OR sp.id = sam.sampled_in_song_id
            WHERE sp.distance < ?
        )
        SELECT DISTINCT
            sp.id as song_id,
            sp.source_id,
            sp.distance,
            sp.path,
            s.title,
            GROUP_CONCAT(DISTINCT g.name) as genres
        FROM SongPath sp
        JOIN Song s ON s.id = sp.id
        LEFT JOIN _SongToGenre sg ON sg.B = s.id
        LEFT JOIN Genre g ON g.id = sg.A
        WHERE EXISTS (
            SELECT 1
            FROM _SongToGenre sg2
            JOIN Genre g2 ON g2.id = sg2.A
            WHERE sg2.B = sp.id AND g2.name = ?
        )
        GROUP BY sp.id, sp.source_id, sp.distance, sp.path, s.title
        ORDER BY sp.distance;
    `

	params = append(params, maxDepth, targetGenre)

	rows, err := r.db.Query(query, params...)
	if err != nil {
		return nil, fmt.Errorf("error executing search query: %v", err)
	}
	defer rows.Close()

	var results []models.SearchResult
	for rows.Next() {
		var (
			songID, sourceID int
			title, genreStr  string
			distance         int
			pathStr          string
		)

		err := rows.Scan(
			&songID,
			&sourceID,
			&distance,
			&pathStr,
			&title,
			&genreStr,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning results: %v", err)
		}

		// Build source and matched song nodes
		sourceSong, err := r.GetSongWithDetails(sourceID)
		if err != nil {
			continue
		}

		matchedSong, err := r.GetSongWithDetails(songID)
		if err != nil {
			continue
		}

		// Convert path string to array of SongNodes
		pathStr = strings.Trim(pathStr, "[]")
		pathIDs := strings.Split(pathStr, ",")
		var path []models.SongNode
		for _, idStr := range pathIDs {
			id, _ := strconv.Atoi(idStr)
			songDetails, err := r.GetSongWithDetails(id)
			if err != nil {
				continue
			}
			path = append(path, *songDetails)
		}

		results = append(results, models.SearchResult{
			SourceSong:  *sourceSong,
			MatchedSong: *matchedSong,
			Distance:    distance,
			Path:        path,
		})
	}

	return results, nil
}

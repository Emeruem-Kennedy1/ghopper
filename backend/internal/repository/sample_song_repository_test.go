package repository

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Emeruem-Kennedy1/ghopper/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupSongTestDB(t *testing.T) (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err, "Failed to create mock database")
	return db, mock
}

func TestSongRepository_GetSongIDsByTitleAndArtist(t *testing.T) {
	db, mock := setupSongTestDB(t)
	defer db.Close()
	repo := NewSongRepository(db)

	t.Run("Found_Songs", func(t *testing.T) {
		// Arrange
		title := "Test Song"
		artist := "Test Artist"
		expectedRows := sqlmock.NewRows([]string{"id"}).
			AddRow(1).
			AddRow(2)

		mock.ExpectQuery("SELECT DISTINCT s.id FROM Song").
			WithArgs(title, artist).
			WillReturnRows(expectedRows)

		// Act
		songIDs, err := repo.GetSongIDsByTitleAndArtist(title, artist)

		// Assert
		require.NoError(t, err, "Should not return error when songs are found")
		assert.Equal(t, []int{1, 2}, songIDs, "Should return correct song IDs")
		assert.NoError(t, mock.ExpectationsWereMet(), "All expectations should be met")
	})

	t.Run("No_Songs_Found", func(t *testing.T) {
		// Arrange
		title := "Unknown Song"
		artist := "Unknown Artist"
		expectedRows := sqlmock.NewRows([]string{"id"})

		mock.ExpectQuery("SELECT DISTINCT s.id FROM Song").
			WithArgs(title, artist).
			WillReturnRows(expectedRows)

		// Act
		songIDs, err := repo.GetSongIDsByTitleAndArtist(title, artist)

		// Assert
		require.NoError(t, err, "Should not return error when no songs are found")
		assert.Empty(t, songIDs, "Should return empty slice when no songs are found")
		assert.NoError(t, mock.ExpectationsWereMet(), "All expectations should be met")
	})

	t.Run("Database_Error", func(t *testing.T) {
		// Arrange
		title := "Error Song"
		artist := "Error Artist"

		mock.ExpectQuery("SELECT DISTINCT s.id FROM Song").
			WithArgs(title, artist).
			WillReturnError(sql.ErrConnDone)

		// Act
		songIDs, err := repo.GetSongIDsByTitleAndArtist(title, artist)

		// Assert
		assert.Error(t, err, "Should return error when database fails")
		assert.Nil(t, songIDs, "Should return nil when database error occurs")
		assert.NoError(t, mock.ExpectationsWereMet(), "All expectations should be met")
	})
}

func TestSongRepository_GetSongWithDetails(t *testing.T) {
	db, mock := setupSongTestDB(t)
	defer db.Close()
	repo := NewSongRepository(db)

	t.Run("Get_Song_With_Details", func(t *testing.T) {
		// Arrange
		songID := 1
		songTitle := "Test Song"
		artistName := "Test Artist"
		artistID := 101
		genres := "Rock,Pop"

		// Mock for main song query
		songRows := sqlmock.NewRows([]string{"id", "title", "artist_name", "artist_id", "genres"}).
			AddRow(songID, songTitle, artistName, artistID, genres)

		mock.ExpectQuery("SELECT s.id, s.title, a.name as artist_name, a.id as artist_id, GROUP_CONCAT").
			WithArgs(songID).
			WillReturnRows(songRows)

		// Mock for artists query
		artistRows := sqlmock.NewRows([]string{"id", "name", "isMainArtist"}).
			AddRow(101, "Test Artist", true).
			AddRow(102, "Featured Artist", false)

		mock.ExpectQuery("SELECT a.id, a.name, sa.isMainArtist FROM SongArtist sa").
			WithArgs(songID).
			WillReturnRows(artistRows)

		// Act
		song, err := repo.GetSongWithDetails(songID)

		// Assert
		require.NoError(t, err, "Should not return error when song is found")
		assert.NotNil(t, song, "Should return a song")
		assert.Equal(t, songID, song.ID, "Song ID should match")
		assert.Equal(t, songTitle, song.Title, "Song title should match")
		assert.Equal(t, []string{"Rock", "Pop"}, song.Genres, "Genres should match")
		assert.Len(t, song.Artists, 2, "Should have two artists")
		assert.NoError(t, mock.ExpectationsWereMet(), "All expectations should be met")
	})

	t.Run("Song_Not_Found", func(t *testing.T) {
		// Arrange
		songID := 999

		mock.ExpectQuery("SELECT s.id, s.title, a.name as artist_name, a.id as artist_id, GROUP_CONCAT").
			WithArgs(songID).
			WillReturnError(sql.ErrNoRows)

		// Act
		song, err := repo.GetSongWithDetails(songID)

		// Assert
		assert.Error(t, err, "Should return error when song is not found")
		assert.Nil(t, song, "Should return nil when song is not found")
		assert.NoError(t, mock.ExpectationsWereMet(), "All expectations should be met")
	})
}

func TestSongRepository_GetAllSampledSongs(t *testing.T) {
	db, mock := setupSongTestDB(t)
	defer db.Close()
	repo := NewSongRepository(db)

	t.Run("Get_Sampled_Songs", func(t *testing.T) {
		// Arrange
		songID := 1
		expectedSampledSongs := []int{2, 3, 4}

		rows := sqlmock.NewRows([]string{"sampled_song_id"}).
			AddRow(2).
			AddRow(3).
			AddRow(4)

		mock.ExpectQuery("WITH RECURSIVE SameSongs AS").
			WithArgs(songID).
			WillReturnRows(rows)

		// Act
		sampledSongs, err := repo.GetAllSampledSongs(songID)

		// Assert
		require.NoError(t, err, "Should not return error when sampled songs are found")
		assert.Equal(t, expectedSampledSongs, sampledSongs, "Should return correct sampled song IDs")
		assert.NoError(t, mock.ExpectationsWereMet(), "All expectations should be met")
	})

	t.Run("No_Sampled_Songs", func(t *testing.T) {
		// Arrange
		songID := 5
		rows := sqlmock.NewRows([]string{"sampled_song_id"})

		mock.ExpectQuery("WITH RECURSIVE SameSongs AS").
			WithArgs(songID).
			WillReturnRows(rows)

		// Act
		sampledSongs, err := repo.GetAllSampledSongs(songID)

		// Assert
		require.NoError(t, err, "Should not return error when no sampled songs are found")
		assert.Empty(t, sampledSongs, "Should return empty slice when no sampled songs are found")
		assert.NoError(t, mock.ExpectationsWereMet(), "All expectations should be met")
	})

	t.Run("Database_Error", func(t *testing.T) {
		// Arrange
		songID := 6

		mock.ExpectQuery("WITH RECURSIVE SameSongs AS").
			WithArgs(songID).
			WillReturnError(sql.ErrConnDone)

		// Act
		sampledSongs, err := repo.GetAllSampledSongs(songID)

		// Assert
		assert.Error(t, err, "Should return error when database fails")
		assert.Nil(t, sampledSongs, "Should return nil when database error occurs")
		assert.NoError(t, mock.ExpectationsWereMet(), "All expectations should be met")
	})
}

func TestSongRepository_FindSongsByGenreBFS(t *testing.T) {
	db, mock := setupSongTestDB(t)
	defer db.Close()
	repo := NewSongRepository(db)

	t.Run("Find_Songs_By_Genre", func(t *testing.T) {
		// Arrange
		songQueries := []models.SongQuery{
			{Title: "Song 1", Artist: "Artist 1"},
			{Title: "Song 2", Artist: "Artist 2"},
		}
		targetGenre := "Rock"
		maxDepth := 3

		// Mock for the recursive query results
		searchRows := sqlmock.NewRows([]string{"song_id", "source_id", "distance", "path", "title", "genres"}).
			AddRow(101, 1, 1, "[1,101]", "Found Song 1", "Rock,Alternative").
			AddRow(102, 2, 2, "[2,103,102]", "Found Song 2", "Rock,Pop")

		// We need to use ExpectQuery with a regex pattern because the query has multiple placeholders
		mock.ExpectQuery("WITH RECURSIVE SongPath AS").
			WillReturnRows(searchRows)

		// For each result, we'll need to mock GetSongWithDetails calls for source and matched songs
		// First result source song (ID: 1)
		sourceSong1Rows := sqlmock.NewRows([]string{"id", "title", "artist_name", "artist_id", "genres"}).
			AddRow(1, "Song 1", "Artist 1", 201, "Pop")
		mock.ExpectQuery("SELECT s.id, s.title, a.name as artist_name, a.id as artist_id, GROUP_CONCAT").
			WithArgs(1).
			WillReturnRows(sourceSong1Rows)

		sourceSong1ArtistRows := sqlmock.NewRows([]string{"id", "name", "isMainArtist"}).
			AddRow(201, "Artist 1", true)
		mock.ExpectQuery("SELECT a.id, a.name, sa.isMainArtist FROM SongArtist sa").
			WithArgs(1).
			WillReturnRows(sourceSong1ArtistRows)

		// First result matched song (ID: 101)
		matchedSong1Rows := sqlmock.NewRows([]string{"id", "title", "artist_name", "artist_id", "genres"}).
			AddRow(101, "Found Song 1", "Rock Artist", 301, "Rock,Alternative")
		mock.ExpectQuery("SELECT s.id, s.title, a.name as artist_name, a.id as artist_id, GROUP_CONCAT").
			WithArgs(101).
			WillReturnRows(matchedSong1Rows)

		matchedSong1ArtistRows := sqlmock.NewRows([]string{"id", "name", "isMainArtist"}).
			AddRow(301, "Rock Artist", true)
		mock.ExpectQuery("SELECT a.id, a.name, sa.isMainArtist FROM SongArtist sa").
			WithArgs(101).
			WillReturnRows(matchedSong1ArtistRows)

		// Path songs for result 1
		// Path node 1 (id: 1) - already mocked above for source song

		// Second result source song (ID: 2)
		sourceSong2Rows := sqlmock.NewRows([]string{"id", "title", "artist_name", "artist_id", "genres"}).
			AddRow(2, "Song 2", "Artist 2", 202, "Electronic")
		mock.ExpectQuery("SELECT s.id, s.title, a.name as artist_name, a.id as artist_id, GROUP_CONCAT").
			WithArgs(2).
			WillReturnRows(sourceSong2Rows)

		sourceSong2ArtistRows := sqlmock.NewRows([]string{"id", "name", "isMainArtist"}).
			AddRow(202, "Artist 2", true)
		mock.ExpectQuery("SELECT a.id, a.name, sa.isMainArtist FROM SongArtist sa").
			WithArgs(2).
			WillReturnRows(sourceSong2ArtistRows)

		// Second result matched song (ID: 102)
		matchedSong2Rows := sqlmock.NewRows([]string{"id", "title", "artist_name", "artist_id", "genres"}).
			AddRow(102, "Found Song 2", "Rock Pop Artist", 302, "Rock,Pop")
		mock.ExpectQuery("SELECT s.id, s.title, a.name as artist_name, a.id as artist_id, GROUP_CONCAT").
			WithArgs(102).
			WillReturnRows(matchedSong2Rows)

		matchedSong2ArtistRows := sqlmock.NewRows([]string{"id", "name", "isMainArtist"}).
			AddRow(302, "Rock Pop Artist", true)
		mock.ExpectQuery("SELECT a.id, a.name, sa.isMainArtist FROM SongArtist sa").
			WithArgs(102).
			WillReturnRows(matchedSong2ArtistRows)

		// Path songs for result 2
		// Path node 1 (id: 2) - already mocked above for source song

		// Path node 2 (id: 103)
		pathNode2Rows := sqlmock.NewRows([]string{"id", "title", "artist_name", "artist_id", "genres"}).
			AddRow(103, "Intermediate Song", "Intermediate Artist", 203, "Electronic,Rock")
		mock.ExpectQuery("SELECT s.id, s.title, a.name as artist_name, a.id as artist_id, GROUP_CONCAT").
			WithArgs(103).
			WillReturnRows(pathNode2Rows)

		pathNode2ArtistRows := sqlmock.NewRows([]string{"id", "name", "isMainArtist"}).
			AddRow(203, "Intermediate Artist", true)
		mock.ExpectQuery("SELECT a.id, a.name, sa.isMainArtist FROM SongArtist sa").
			WithArgs(103).
			WillReturnRows(pathNode2ArtistRows)

		// Act
		results, err := repo.FindSongsByGenreBFS(songQueries, targetGenre, maxDepth)

		// Assert
		require.NoError(t, err, "Should not return error when songs are found")
		assert.Len(t, results, 2, "Should return 2 results")
		assert.Equal(t, 1, results[0].Distance, "First result should have distance 1")
		assert.Equal(t, 2, results[1].Distance, "Second result should have distance 2")
		assert.Equal(t, "Found Song 1", results[0].MatchedSong.Title, "First result should match correct song")
		assert.Equal(t, "Found Song 2", results[1].MatchedSong.Title, "Second result should match correct song")
		assert.NoError(t, mock.ExpectationsWereMet(), "All expectations should be met")
	})

	t.Run("No_Matching_Songs", func(t *testing.T) {
		// Arrange
		songQueries := []models.SongQuery{
			{Title: "Rare Song", Artist: "Rare Artist"},
		}
		targetGenre := "Experimental Jazz"
		maxDepth := 2

		// Mock empty result set
		mock.ExpectQuery("WITH RECURSIVE SongPath AS").
			WillReturnRows(sqlmock.NewRows([]string{"song_id", "source_id", "distance", "path", "title", "genres"}))

		// Act
		results, err := repo.FindSongsByGenreBFS(songQueries, targetGenre, maxDepth)

		// Assert
		require.NoError(t, err, "Should not return error when no songs are found")
		assert.Empty(t, results, "Should return empty results when no songs match")
		assert.NoError(t, mock.ExpectationsWereMet(), "All expectations should be met")
	})

	t.Run("Database_Error", func(t *testing.T) {
		// Arrange
		songQueries := []models.SongQuery{
			{Title: "Error Song", Artist: "Error Artist"},
		}
		targetGenre := "Rock"
		maxDepth := 3

		mock.ExpectQuery("WITH RECURSIVE SongPath AS").
			WillReturnError(sql.ErrConnDone)

		// Act
		results, err := repo.FindSongsByGenreBFS(songQueries, targetGenre, maxDepth)

		// Assert
		assert.Error(t, err, "Should return error when database fails")
		assert.Nil(t, results, "Should return nil when database error occurs")
		assert.NoError(t, mock.ExpectationsWereMet(), "All expectations should be met")
	})
}

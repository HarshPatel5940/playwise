package datastructures

import (
	"fmt"
	"src/internal/models"
	"strings"
	"testing"
	"time"
)

func TestNewPlaylistTreeNode(t *testing.T) {
	parent := NewPlaylistTreeNode("Parent", GenreNode, nil)
	child := NewPlaylistTreeNode("Child", SubgenreNode, parent)

	if child == nil {
		t.Fatal("Expected non-nil node")
	}
	if child.Name != "Child" {
		t.Errorf("Expected name 'Child', got %s", child.Name)
	}
	if child.NodeType != SubgenreNode {
		t.Errorf("Expected SubgenreNode type, got %v", child.NodeType)
	}
	if child.Parent != parent {
		t.Error("Parent not set correctly")
	}
	if child.Children == nil {
		t.Error("Children map should be initialized")
	}
	if child.Songs == nil {
		t.Error("Songs slice should be initialized")
	}
}

func TestAddChild(t *testing.T) {
	parent := NewPlaylistTreeNode("Parent", GenreNode, nil)

	// Add first child
	child1 := parent.AddChild("Child1", SubgenreNode)
	if child1 == nil {
		t.Fatal("Expected non-nil child")
	}
	if child1.Name != "Child1" {
		t.Errorf("Expected name 'Child1', got %s", child1.Name)
	}
	if child1.Parent != parent {
		t.Error("Parent not set correctly")
	}

	// Add second child with same name (should return existing)
	child2 := parent.AddChild("Child1", MoodNode)
	if child2 != child1 {
		t.Error("Adding child with same name should return existing child")
	}
	if child2.NodeType != SubgenreNode {
		t.Error("Node type should not change when adding existing child")
	}

	// Add different child
	child3 := parent.AddChild("Child2", MoodNode)
	if child3 == child1 {
		t.Error("Different children should be different objects")
	}
	if len(parent.Children) != 2 {
		t.Errorf("Expected 2 children, got %d", len(parent.Children))
	}
}

func TestGetChild(t *testing.T) {
	parent := NewPlaylistTreeNode("Parent", GenreNode, nil)
	child := parent.AddChild("TestChild", SubgenreNode)

	// Test getting existing child
	retrieved := parent.GetChild("TestChild")
	if retrieved != child {
		t.Error("GetChild should return the correct child")
	}

	// Test getting non-existent child
	nonExistent := parent.GetChild("NonExistent")
	if nonExistent != nil {
		t.Error("GetChild should return nil for non-existent child")
	}
}

func TestHasChildren(t *testing.T) {
	node := NewPlaylistTreeNode("Test", GenreNode, nil)

	// Initially should have no children
	if node.HasChildren() {
		t.Error("New node should not have children")
	}

	// After adding child
	node.AddChild("Child", SubgenreNode)
	if !node.HasChildren() {
		t.Error("Node should have children after adding one")
	}
}

func TestGetChildrenNames(t *testing.T) {
	parent := NewPlaylistTreeNode("Parent", GenreNode, nil)

	// Test empty children
	names := parent.GetChildrenNames()
	if len(names) != 0 {
		t.Errorf("Expected 0 children names, got %d", len(names))
	}

	// Add children
	parent.AddChild("Child1", SubgenreNode)
	parent.AddChild("Child2", SubgenreNode)
	parent.AddChild("Child3", SubgenreNode)

	names = parent.GetChildrenNames()
	if len(names) != 3 {
		t.Errorf("Expected 3 children names, got %d", len(names))
	}

	// Check all names are present (order doesn't matter for map iteration)
	expectedNames := map[string]bool{"Child1": false, "Child2": false, "Child3": false}
	for _, name := range names {
		if _, exists := expectedNames[name]; !exists {
			t.Errorf("Unexpected child name: %s", name)
		} else {
			expectedNames[name] = true
		}
	}

	// Verify all expected names were found
	for name, found := range expectedNames {
		if !found {
			t.Errorf("Expected child name %s not found", name)
		}
	}
}

func TestAddSongAndGetSongs(t *testing.T) {
	artistNode := NewPlaylistTreeNode("Artist", ArtistNode, nil)
	genreNode := NewPlaylistTreeNode("Genre", GenreNode, nil)

	song1 := &models.Song{ID: "1", Title: "Song 1", Artist: "Artist"}
	song2 := &models.Song{ID: "2", Title: "Song 2", Artist: "Artist"}

	// Test adding songs to artist node
	artistNode.AddSong(song1)
	artistNode.AddSong(song2)

	songs := artistNode.GetSongs()
	if len(songs) != 2 {
		t.Errorf("Expected 2 songs, got %d", len(songs))
	}
	if songs[0] != song1 || songs[1] != song2 {
		t.Error("Songs not added correctly")
	}

	// Test adding song to non-artist node (should be ignored)
	genreNode.AddSong(song1)
	genreSongs := genreNode.GetSongs()
	if len(genreSongs) != 0 {
		t.Error("Non-artist nodes should not store songs")
	}
}

func TestGetPath(t *testing.T) {
	// Create hierarchy: Root -> Genre -> Subgenre -> Mood -> Artist
	root := NewPlaylistTreeNode("Root", GenreNode, nil)
	genre := root.AddChild("Rock", GenreNode)
	subgenre := genre.AddChild("Alternative", SubgenreNode)
	mood := subgenre.AddChild("Energetic", MoodNode)
	artist := mood.AddChild("Nirvana", ArtistNode)

	path := artist.GetPath()
	expectedPath := []string{"Rock", "Alternative", "Energetic", "Nirvana"}

	if len(path) != len(expectedPath) {
		t.Errorf("Expected path length %d, got %d", len(expectedPath), len(path))
	}

	for i, segment := range expectedPath {
		if path[i] != segment {
			t.Errorf("Path segment %d: expected %s, got %s", i, segment, path[i])
		}
	}

	// Test root path
	rootPath := root.GetPath()
	if len(rootPath) != 0 {
		t.Errorf("Root path should be empty, got %v", rootPath)
	}
}

func TestNewPlaylistExplorerTree(t *testing.T) {
	tree := NewPlaylistExplorerTree()

	if tree == nil {
		t.Fatal("Expected non-nil tree")
	}
	if tree.Root == nil {
		t.Fatal("Expected non-nil root")
	}
	if tree.Root.Name != "Root" {
		t.Errorf("Expected root name 'Root', got %s", tree.Root.Name)
	}
	if tree.TotalSongs != 0 {
		t.Errorf("Expected 0 total songs, got %d", tree.TotalSongs)
	}
	if len(tree.Stats) == 0 {
		t.Error("Stats should be initialized")
	}
}

func createPlaylistTestSong(id, title, artist, genre, subgenre, mood string) *models.Song {
	return &models.Song{
		ID:        id,
		Title:     title,
		Artist:    artist,
		Genre:     genre,
		SubGenre:  subgenre,
		Mood:      mood,
		Duration:  180,
		Rating:    4,
		PlayCount: 10,
		AddedAt:   time.Now(),
	}
}

func TestAddSong(t *testing.T) {
	tree := NewPlaylistExplorerTree()

	song1 := createPlaylistTestSong("1", "Smells Like Teen Spirit", "Nirvana", "Rock", "Alternative", "Energetic")
	song2 := createPlaylistTestSong("2", "Come As You Are", "Nirvana", "Rock", "Alternative", "Melancholic")
	song3 := createPlaylistTestSong("3", "Bohemian Rhapsody", "Queen", "Rock", "Classic Rock", "Epic")

	// Test adding songs
	tree.AddSong(song1)
	tree.AddSong(song2)
	tree.AddSong(song3)

	if tree.TotalSongs != 3 {
		t.Errorf("Expected 3 total songs, got %d", tree.TotalSongs)
	}

	// Verify stats
	stats := tree.GetStats()
	if stats["genres"].(int) != 1 {
		t.Errorf("Expected 1 genre, got %v", stats["genres"])
	}
	if stats["subgenres"].(int) != 2 {
		t.Errorf("Expected 2 subgenres, got %v", stats["subgenres"])
	}
	if stats["moods"].(int) != 3 {
		t.Errorf("Expected 3 moods, got %v", stats["moods"])
	}
	if stats["artists"].(int) != 2 {
		t.Errorf("Expected 2 artists, got %v", stats["artists"])
	}
}

func TestAddSongWithEmptyFields(t *testing.T) {
	tree := NewPlaylistExplorerTree()

	// Song with empty fields
	song := &models.Song{
		ID:    "1",
		Title: "Test Song",
		// All other fields empty
	}

	tree.AddSong(song)

	if tree.TotalSongs != 1 {
		t.Errorf("Expected 1 song, got %d", tree.TotalSongs)
	}

	// Should create "Unknown" categories
	genres := tree.GetGenres()
	if len(genres) != 1 || genres[0] != "Unknown Genre" {
		t.Errorf("Expected Unknown Genre, got %v", genres)
	}

	subgenres := tree.GetSubgenres("Unknown Genre")
	if len(subgenres) != 1 || subgenres[0] != "Unknown Subgenre" {
		t.Errorf("Expected Unknown Subgenre, got %v", subgenres)
	}
}

func TestAddNilSong(t *testing.T) {
	tree := NewPlaylistExplorerTree()

	// Should not panic or add anything
	tree.AddSong(nil)

	if tree.TotalSongs != 0 {
		t.Errorf("Expected 0 songs after adding nil, got %d", tree.TotalSongs)
	}
}

func TestGetGenres(t *testing.T) {
	tree := NewPlaylistExplorerTree()

	// Initially empty
	genres := tree.GetGenres()
	if len(genres) != 0 {
		t.Errorf("Expected 0 genres, got %d", len(genres))
	}

	// Add songs from different genres
	tree.AddSong(createPlaylistTestSong("1", "Song 1", "Artist 1", "Rock", "Alternative", "Energetic"))
	tree.AddSong(createPlaylistTestSong("2", "Song 2", "Artist 2", "Pop", "Mainstream", "Happy"))
	tree.AddSong(createPlaylistTestSong("3", "Song 3", "Artist 3", "Jazz", "Smooth", "Relaxed"))

	genres = tree.GetGenres()
	if len(genres) != 3 {
		t.Errorf("Expected 3 genres, got %d", len(genres))
	}

	expectedGenres := map[string]bool{"Rock": false, "Pop": false, "Jazz": false}
	for _, genre := range genres {
		if _, exists := expectedGenres[genre]; !exists {
			t.Errorf("Unexpected genre: %s", genre)
		} else {
			expectedGenres[genre] = true
		}
	}
}

func TestGetSubgenres(t *testing.T) {
	tree := NewPlaylistExplorerTree()
	tree.AddSong(createPlaylistTestSong("1", "Song 1", "Artist 1", "Rock", "Alternative", "Energetic"))
	tree.AddSong(createPlaylistTestSong("2", "Song 2", "Artist 2", "Rock", "Classic Rock", "Epic"))

	subgenres := tree.GetSubgenres("Rock")
	if len(subgenres) != 2 {
		t.Errorf("Expected 2 subgenres for Rock, got %d", len(subgenres))
	}

	// Test non-existent genre
	emptySubgenres := tree.GetSubgenres("NonExistent")
	if len(emptySubgenres) != 0 {
		t.Errorf("Expected 0 subgenres for non-existent genre, got %d", len(emptySubgenres))
	}
}

func TestGetMoods(t *testing.T) {
	tree := NewPlaylistExplorerTree()
	tree.AddSong(createPlaylistTestSong("1", "Song 1", "Artist 1", "Rock", "Alternative", "Energetic"))
	tree.AddSong(createPlaylistTestSong("2", "Song 2", "Artist 2", "Rock", "Alternative", "Melancholic"))

	moods := tree.GetMoods("Rock", "Alternative")
	if len(moods) != 2 {
		t.Errorf("Expected 2 moods, got %d", len(moods))
	}

	// Test non-existent path
	emptyMoods := tree.GetMoods("NonExistent", "Genre")
	if len(emptyMoods) != 0 {
		t.Errorf("Expected 0 moods for non-existent path, got %d", len(emptyMoods))
	}

	emptyMoods = tree.GetMoods("Rock", "NonExistent")
	if len(emptyMoods) != 0 {
		t.Errorf("Expected 0 moods for non-existent subgenre, got %d", len(emptyMoods))
	}
}

func TestGetArtists(t *testing.T) {
	tree := NewPlaylistExplorerTree()
	tree.AddSong(createPlaylistTestSong("1", "Song 1", "Nirvana", "Rock", "Alternative", "Energetic"))
	tree.AddSong(createPlaylistTestSong("2", "Song 2", "Pearl Jam", "Rock", "Alternative", "Energetic"))

	artists := tree.GetArtists("Rock", "Alternative", "Energetic")
	if len(artists) != 2 {
		t.Errorf("Expected 2 artists, got %d", len(artists))
	}

	// Test non-existent paths
	emptyArtists := tree.GetArtists("NonExistent", "Alternative", "Energetic")
	if len(emptyArtists) != 0 {
		t.Errorf("Expected 0 artists for non-existent genre, got %d", len(emptyArtists))
	}
}

func TestGetSongs(t *testing.T) {
	tree := NewPlaylistExplorerTree()
	song1 := createPlaylistTestSong("1", "Smells Like Teen Spirit", "Nirvana", "Rock", "Alternative", "Energetic")
	song2 := createPlaylistTestSong("2", "Come As You Are", "Nirvana", "Rock", "Alternative", "Energetic")

	tree.AddSong(song1)
	tree.AddSong(song2)

	songs := tree.GetSongs("Rock", "Alternative", "Energetic", "Nirvana")
	if len(songs) != 2 {
		t.Errorf("Expected 2 songs, got %d", len(songs))
	}

	// Test non-existent paths
	emptySongs := tree.GetSongs("NonExistent", "Alternative", "Energetic", "Nirvana")
	if len(emptySongs) != 0 {
		t.Errorf("Expected 0 songs for non-existent path, got %d", len(emptySongs))
	}
}

func TestGetAllSongsInGenre(t *testing.T) {
	tree := NewPlaylistExplorerTree()

	// Add multiple songs in Rock genre
	tree.AddSong(createPlaylistTestSong("1", "Song 1", "Artist 1", "Rock", "Alternative", "Energetic"))
	tree.AddSong(createPlaylistTestSong("2", "Song 2", "Artist 2", "Rock", "Classic Rock", "Epic"))
	tree.AddSong(createPlaylistTestSong("3", "Song 3", "Artist 3", "Pop", "Mainstream", "Happy"))

	rockSongs := tree.GetAllSongsInGenre("Rock")
	if len(rockSongs) != 2 {
		t.Errorf("Expected 2 Rock songs, got %d", len(rockSongs))
	}

	// Test non-existent genre
	emptySongs := tree.GetAllSongsInGenre("NonExistent")
	if len(emptySongs) != 0 {
		t.Errorf("Expected 0 songs for non-existent genre, got %d", len(emptySongs))
	}
}

func TestGetAllSongsInMood(t *testing.T) {
	tree := NewPlaylistExplorerTree()

	// Add songs with same mood across different genres
	tree.AddSong(createPlaylistTestSong("1", "Song 1", "Artist 1", "Rock", "Alternative", "Energetic"))
	tree.AddSong(createPlaylistTestSong("2", "Song 2", "Artist 2", "Pop", "Dance", "Energetic"))
	tree.AddSong(createPlaylistTestSong("3", "Song 3", "Artist 3", "Jazz", "Fusion", "Relaxed"))

	energeticSongs := tree.GetAllSongsInMood("Energetic")
	if len(energeticSongs) != 2 {
		t.Errorf("Expected 2 Energetic songs, got %d", len(energeticSongs))
	}

	relaxedSongs := tree.GetAllSongsInMood("Relaxed")
	if len(relaxedSongs) != 1 {
		t.Errorf("Expected 1 Relaxed song, got %d", len(relaxedSongs))
	}

	// Test non-existent mood
	emptySongs := tree.GetAllSongsInMood("NonExistent")
	if len(emptySongs) != 0 {
		t.Errorf("Expected 0 songs for non-existent mood, got %d", len(emptySongs))
	}
}

func TestDepthFirstSearch(t *testing.T) {
	tree := NewPlaylistExplorerTree()
	tree.AddSong(createPlaylistTestSong("1", "Song 1", "Artist 1", "Rock", "Alternative", "Energetic"))
	tree.AddSong(createPlaylistTestSong("2", "Song 2", "Artist 2", "Pop", "Mainstream", "Happy"))

	visitedNodes := []string{}
	tree.DepthFirstSearch(func(node *PlaylistTreeNode) {
		if node.Name != "Root" {
			visitedNodes = append(visitedNodes, node.Name)
		}
	})

	// Should visit all nodes
	expectedMinimumNodes := 6 // 2 genres + 2 subgenres + 2 moods + 2 artists
	if len(visitedNodes) < expectedMinimumNodes {
		t.Errorf("Expected at least %d visited nodes, got %d", expectedMinimumNodes, len(visitedNodes))
	}
}

func TestBreadthFirstSearch(t *testing.T) {
	tree := NewPlaylistExplorerTree()
	tree.AddSong(createPlaylistTestSong("1", "Song 1", "Artist 1", "Rock", "Alternative", "Energetic"))

	visitedNodes := []string{}
	tree.BreadthFirstSearch(func(node *PlaylistTreeNode) {
		visitedNodes = append(visitedNodes, node.Name)
	})

	// Should start with Root
	if len(visitedNodes) == 0 || visitedNodes[0] != "Root" {
		t.Error("BFS should start with Root node")
	}

	// Should visit all nodes
	if len(visitedNodes) < 5 { // Root + Genre + Subgenre + Mood + Artist
		t.Errorf("Expected at least 5 visited nodes, got %d", len(visitedNodes))
	}
}

func TestFindSongPath(t *testing.T) {
	tree := NewPlaylistExplorerTree()
	song := createPlaylistTestSong("test123", "Test Song", "Test Artist", "Rock", "Alternative", "Energetic")
	tree.AddSong(song)

	// Test finding existing song
	path, err := tree.FindSongPath("test123")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	expectedPath := []string{"Rock", "Alternative", "Energetic", "Test Artist"}
	if len(path) != len(expectedPath) {
		t.Errorf("Expected path length %d, got %d", len(expectedPath), len(path))
	}

	for i, segment := range expectedPath {
		if path[i] != segment {
			t.Errorf("Path segment %d: expected %s, got %s", i, segment, path[i])
		}
	}

	// Test finding non-existent song
	_, err = tree.FindSongPath("nonexistent")
	if err == nil {
		t.Error("Expected error for non-existent song")
	}
}

func TestRemoveSong(t *testing.T) {
	tree := NewPlaylistExplorerTree()
	song1 := createPlaylistTestSong("1", "Song 1", "Artist", "Rock", "Alternative", "Energetic")
	song2 := createPlaylistTestSong("2", "Song 2", "Artist", "Rock", "Alternative", "Energetic")

	tree.AddSong(song1)
	tree.AddSong(song2)

	if tree.TotalSongs != 2 {
		t.Errorf("Expected 2 songs before removal, got %d", tree.TotalSongs)
	}

	// Remove existing song
	err := tree.RemoveSong("1")
	if err != nil {
		t.Errorf("Expected no error removing existing song, got %v", err)
	}

	if tree.TotalSongs != 1 {
		t.Errorf("Expected 1 song after removal, got %d", tree.TotalSongs)
	}

	// Try to remove non-existent song
	err = tree.RemoveSong("nonexistent")
	if err == nil {
		t.Error("Expected error removing non-existent song")
	}

	// Verify remaining song is still accessible
	songs := tree.GetSongs("Rock", "Alternative", "Energetic", "Artist")
	if len(songs) != 1 || songs[0].ID != "2" {
		t.Error("Remaining song should still be accessible")
	}
}

func TestGetTreeStructure(t *testing.T) {
	tree := NewPlaylistExplorerTree()
	tree.AddSong(createPlaylistTestSong("1", "Song 1", "Artist 1", "Rock", "Alternative", "Energetic"))
	tree.AddSong(createPlaylistTestSong("2", "Song 2", "Artist 2", "Rock", "Alternative", "Energetic"))
	tree.AddSong(createPlaylistTestSong("3", "Song 3", "Artist 1", "Pop", "Mainstream", "Happy"))

	structure := tree.GetTreeStructure()

	// Check Rock genre structure
	rockGenre, exists := structure["Rock"]
	if !exists {
		t.Error("Rock genre should exist in structure")
	}

	rockMap, ok := rockGenre.(map[string]interface{})
	if !ok {
		t.Error("Rock genre should be a map")
	}

	alternativeSubgenre, exists := rockMap["Alternative"]
	if !exists {
		t.Error("Alternative subgenre should exist under Rock")
	}

	alternativeMap, ok := alternativeSubgenre.(map[string]interface{})
	if !ok {
		t.Error("Alternative subgenre should be a map")
	}

	energeticMood, exists := alternativeMap["Energetic"]
	if !exists {
		t.Error("Energetic mood should exist under Alternative")
	}

	energeticMap, ok := energeticMood.(map[string]interface{})
	if !ok {
		t.Error("Energetic mood should be a map")
	}

	// Should have 2 artists under Rock->Alternative->Energetic
	if len(energeticMap) != 2 {
		t.Errorf("Expected 2 artists under Energetic, got %d", len(energeticMap))
	}
}

func TestString(t *testing.T) {
	tree := NewPlaylistExplorerTree()

	// Test empty tree
	emptyStr := tree.String()
	if emptyStr != "Empty Playlist Explorer Tree" {
		t.Errorf("Expected empty tree string, got %s", emptyStr)
	}

	// Test with songs
	tree.AddSong(createPlaylistTestSong("1", "Song 1", "Artist 1", "Rock", "Alternative", "Energetic"))
	tree.AddSong(createPlaylistTestSong("2", "Song 2", "Artist 2", "Pop", "Mainstream", "Happy"))

	str := tree.String()
	if str == "Empty Playlist Explorer Tree" {
		t.Error("Tree with songs should not return empty string")
	}

	// Should contain total songs count
	if !strings.Contains(str, "Total Songs: 2") {
		t.Error("String should contain total songs count")
	}

	// Should contain stats
	if !strings.Contains(str, "genres") {
		t.Error("String should contain genre count")
	}
}

func TestNormalization(t *testing.T) {
	tree := NewPlaylistExplorerTree()

	// Test case normalization
	song1 := createPlaylistTestSong("1", "Song 1", "artist name", "ROCK", "alternative", "Energetic")
	song2 := createPlaylistTestSong("2", "Song 2", "Artist Name", "rock", "ALTERNATIVE", "energetic")

	tree.AddSong(song1)
	tree.AddSong(song2)

	// Should be normalized to title case
	genres := tree.GetGenres()
	if len(genres) != 1 || genres[0] != "Rock" {
		t.Errorf("Expected normalized genre 'Rock', got %v", genres)
	}

	subgenres := tree.GetSubgenres("Rock")
	if len(subgenres) != 1 || subgenres[0] != "Alternative" {
		t.Errorf("Expected normalized subgenre 'Alternative', got %v", subgenres)
	}

	artists := tree.GetArtists("Rock", "Alternative", "Energetic")
	if len(artists) != 1 || artists[0] != "Artist Name" {
		t.Errorf("Expected normalized artist 'Artist Name', got %v", artists)
	}

	// Should have 2 songs under same normalized path
	songs := tree.GetSongs("Rock", "Alternative", "Energetic", "Artist Name")
	if len(songs) != 2 {
		t.Errorf("Expected 2 songs under normalized path, got %d", len(songs))
	}
}

func TestWhitespaceHandling(t *testing.T) {
	tree := NewPlaylistExplorerTree()

	// Test with leading/trailing whitespace
	song := createPlaylistTestSong("1", "Song", "  Artist  ", " Rock ", "  Alternative  ", " Energetic ")
	tree.AddSong(song)

	genres := tree.GetGenres()
	if len(genres) != 1 || genres[0] != "Rock" {
		t.Errorf("Expected trimmed genre 'Rock', got %v", genres)
	}

	artists := tree.GetArtists("Rock", "Alternative", "Energetic")
	if len(artists) != 1 || artists[0] != "Artist" {
		t.Errorf("Expected trimmed artist 'Artist', got %v", artists)
	}
}

// Benchmark tests
func BenchmarkAddSong(b *testing.B) {
	tree := NewPlaylistExplorerTree()
	songs := make([]*models.Song, b.N)

	for i := 0; i < b.N; i++ {
		songs[i] = createPlaylistTestSong(
			fmt.Sprintf("song_%d", i),
			fmt.Sprintf("Title %d", i),
			fmt.Sprintf("Artist %d", i%100),  // 100 different artists
			fmt.Sprintf("Genre %d", i%10),    // 10 different genres
			fmt.Sprintf("Subgenre %d", i%50), // 50 different subgenres
			fmt.Sprintf("Mood %d", i%20),     // 20 different moods
		)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tree.AddSong(songs[i])
	}
}

func BenchmarkDepthFirstSearch(b *testing.B) {
	tree := NewPlaylistExplorerTree()

	// Add many songs to create a large tree
	for i := 0; i < 1000; i++ {
		song := createPlaylistTestSong(
			fmt.Sprintf("song_%d", i),
			fmt.Sprintf("Title %d", i),
			fmt.Sprintf("Artist %d", i%100),
			fmt.Sprintf("Genre %d", i%10),
			fmt.Sprintf("Subgenre %d", i%50),
			fmt.Sprintf("Mood %d", i%20),
		)
		tree.AddSong(song)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tree.DepthFirstSearch(func(node *PlaylistTreeNode) {
			// Do minimal work
			_ = node.Name
		})
	}
}

func BenchmarkFindSongPath(b *testing.B) {
	tree := NewPlaylistExplorerTree()

	// Add many songs
	songIDs := make([]string, 1000)
	for i := 0; i < 1000; i++ {
		songID := fmt.Sprintf("song_%d", i)
		songIDs[i] = songID
		song := createPlaylistTestSong(
			songID,
			fmt.Sprintf("Title %d", i),
			fmt.Sprintf("Artist %d", i%100),
			fmt.Sprintf("Genre %d", i%10),
			fmt.Sprintf("Subgenre %d", i%50),
			fmt.Sprintf("Mood %d", i%20),
		)
		tree.AddSong(song)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tree.FindSongPath(songIDs[i%len(songIDs)])
	}
}

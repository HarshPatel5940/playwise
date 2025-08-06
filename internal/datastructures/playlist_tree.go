package datastructures

import (
	"fmt"
	"src/internal/models"
	"strings"
)

// PlaylistTreeNodeType defines the type of node in the playlist tree
type PlaylistTreeNodeType int

const (
	GenreNode PlaylistTreeNodeType = iota
	SubgenreNode
	MoodNode
	ArtistNode
)

// PlaylistTreeNode represents a node in the playlist explorer tree
// Each node can have multiple children and stores songs at artist level
// Time Complexity: O(1) for field access
// Space Complexity: O(k) where k is the number of children
type PlaylistTreeNode struct {
	Name     string
	NodeType PlaylistTreeNodeType
	Children map[string]*PlaylistTreeNode
	Songs    []*models.Song // Only populated for artist nodes
	Parent   *PlaylistTreeNode
}

// NewPlaylistTreeNode creates a new playlist tree node
// Time Complexity: O(1)
// Space Complexity: O(1)
func NewPlaylistTreeNode(name string, nodeType PlaylistTreeNodeType, parent *PlaylistTreeNode) *PlaylistTreeNode {
	return &PlaylistTreeNode{
		Name:     name,
		NodeType: nodeType,
		Children: make(map[string]*PlaylistTreeNode),
		Songs:    make([]*models.Song, 0),
		Parent:   parent,
	}
}

// AddChild adds a child node to the current node
// Time Complexity: O(1)
// Space Complexity: O(1)
func (node *PlaylistTreeNode) AddChild(childName string, childType PlaylistTreeNodeType) *PlaylistTreeNode {
	if _, exists := node.Children[childName]; !exists {
		node.Children[childName] = NewPlaylistTreeNode(childName, childType, node)
	}
	return node.Children[childName]
}

// GetChild retrieves a child node by name
// Time Complexity: O(1) average
// Space Complexity: O(1)
func (node *PlaylistTreeNode) GetChild(childName string) *PlaylistTreeNode {
	return node.Children[childName]
}

// HasChildren checks if the node has any children
// Time Complexity: O(1)
// Space Complexity: O(1)
func (node *PlaylistTreeNode) HasChildren() bool {
	return len(node.Children) > 0
}

// GetChildrenNames returns all child names
// Time Complexity: O(k) where k is the number of children
// Space Complexity: O(k)
func (node *PlaylistTreeNode) GetChildrenNames() []string {
	names := make([]string, 0, len(node.Children))
	for name := range node.Children {
		names = append(names, name)
	}
	return names
}

// AddSong adds a song to an artist node
// Time Complexity: O(1)
// Space Complexity: O(1)
func (node *PlaylistTreeNode) AddSong(song *models.Song) {
	if node.NodeType == ArtistNode {
		node.Songs = append(node.Songs, song)
	}
}

// GetSongs returns all songs in an artist node
// Time Complexity: O(1)
// Space Complexity: O(1)
func (node *PlaylistTreeNode) GetSongs() []*models.Song {
	return node.Songs
}

// GetPath returns the full path from root to current node
// Time Complexity: O(d) where d is the depth
// Space Complexity: O(d)
func (node *PlaylistTreeNode) GetPath() []string {
	path := make([]string, 0)
	current := node

	for current != nil && current.Name != "Root" {
		path = append([]string{current.Name}, path...)
		current = current.Parent
	}

	return path
}

// PlaylistExplorerTree represents the hierarchical song organization
// Structure: Genre → Subgenre → Mood → Artist → Songs
// Time Complexity: O(1) for root access, O(d) for traversal where d is depth
// Space Complexity: O(n) where n is the total number of unique categories + songs
type PlaylistExplorerTree struct {
	Root       *PlaylistTreeNode
	TotalSongs int
	Stats      map[string]int // Statistics for each level
}

// NewPlaylistExplorerTree creates a new playlist explorer tree
// Time Complexity: O(1)
// Space Complexity: O(1)
func NewPlaylistExplorerTree() *PlaylistExplorerTree {
	return &PlaylistExplorerTree{
		Root:       NewPlaylistTreeNode("Root", GenreNode, nil),
		TotalSongs: 0,
		Stats: map[string]int{
			"genres":    0,
			"subgenres": 0,
			"moods":     0,
			"artists":   0,
		},
	}
}

// AddSong adds a song to the tree, creating the hierarchy as needed
// Time Complexity: O(1) average for hash map operations
// Space Complexity: O(1) for the song, O(d) for path creation if needed
func (pet *PlaylistExplorerTree) AddSong(song *models.Song) {
	if song == nil {
		return
	}

	// Normalize the category names
	genre := strings.Title(strings.ToLower(strings.TrimSpace(song.Genre)))
	subgenre := strings.Title(strings.ToLower(strings.TrimSpace(song.SubGenre)))
	mood := strings.Title(strings.ToLower(strings.TrimSpace(song.Mood)))
	artist := strings.Title(strings.ToLower(strings.TrimSpace(song.Artist)))

	// Handle empty categories
	if genre == "" {
		genre = "Unknown Genre"
	}
	if subgenre == "" {
		subgenre = "Unknown Subgenre"
	}
	if mood == "" {
		mood = "Unknown Mood"
	}
	if artist == "" {
		artist = "Unknown Artist"
	}

	// Navigate/create the hierarchy: Root -> Genre -> Subgenre -> Mood -> Artist
	genreNode := pet.Root.GetChild(genre)
	if genreNode == nil {
		genreNode = pet.Root.AddChild(genre, GenreNode)
		pet.Stats["genres"]++
	}

	subgenreNode := genreNode.GetChild(subgenre)
	if subgenreNode == nil {
		subgenreNode = genreNode.AddChild(subgenre, SubgenreNode)
		pet.Stats["subgenres"]++
	}

	moodNode := subgenreNode.GetChild(mood)
	if moodNode == nil {
		moodNode = subgenreNode.AddChild(mood, MoodNode)
		pet.Stats["moods"]++
	}

	artistNode := moodNode.GetChild(artist)
	if artistNode == nil {
		artistNode = moodNode.AddChild(artist, ArtistNode)
		pet.Stats["artists"]++
	}

	// Add the song to the artist node
	artistNode.AddSong(song)
	pet.TotalSongs++
}

// GetGenres returns all available genres
// Time Complexity: O(g) where g is the number of genres
// Space Complexity: O(g)
func (pet *PlaylistExplorerTree) GetGenres() []string {
	return pet.Root.GetChildrenNames()
}

// GetSubgenres returns all subgenres for a given genre
// Time Complexity: O(1) for genre lookup + O(s) for subgenres where s is number of subgenres
// Space Complexity: O(s)
func (pet *PlaylistExplorerTree) GetSubgenres(genre string) []string {
	genreNode := pet.Root.GetChild(genre)
	if genreNode == nil {
		return []string{}
	}
	return genreNode.GetChildrenNames()
}

// GetMoods returns all moods for a given genre and subgenre
// Time Complexity: O(1) for navigation + O(m) for moods where m is number of moods
// Space Complexity: O(m)
func (pet *PlaylistExplorerTree) GetMoods(genre, subgenre string) []string {
	genreNode := pet.Root.GetChild(genre)
	if genreNode == nil {
		return []string{}
	}

	subgenreNode := genreNode.GetChild(subgenre)
	if subgenreNode == nil {
		return []string{}
	}

	return subgenreNode.GetChildrenNames()
}

// GetArtists returns all artists for a given genre, subgenre, and mood
// Time Complexity: O(1) for navigation + O(a) for artists where a is number of artists
// Space Complexity: O(a)
func (pet *PlaylistExplorerTree) GetArtists(genre, subgenre, mood string) []string {
	genreNode := pet.Root.GetChild(genre)
	if genreNode == nil {
		return []string{}
	}

	subgenreNode := genreNode.GetChild(subgenre)
	if subgenreNode == nil {
		return []string{}
	}

	moodNode := subgenreNode.GetChild(mood)
	if moodNode == nil {
		return []string{}
	}

	return moodNode.GetChildrenNames()
}

// GetSongs returns all songs for a specific artist in a given category hierarchy
// Time Complexity: O(1) for navigation
// Space Complexity: O(1)
func (pet *PlaylistExplorerTree) GetSongs(genre, subgenre, mood, artist string) []*models.Song {
	genreNode := pet.Root.GetChild(genre)
	if genreNode == nil {
		return []*models.Song{}
	}

	subgenreNode := genreNode.GetChild(subgenre)
	if subgenreNode == nil {
		return []*models.Song{}
	}

	moodNode := subgenreNode.GetChild(mood)
	if moodNode == nil {
		return []*models.Song{}
	}

	artistNode := moodNode.GetChild(artist)
	if artistNode == nil {
		return []*models.Song{}
	}

	return artistNode.GetSongs()
}

// GetAllSongsInGenre returns all songs in a specific genre
// Time Complexity: O(n) where n is the number of songs in the genre
// Space Complexity: O(n)
func (pet *PlaylistExplorerTree) GetAllSongsInGenre(genre string) []*models.Song {
	genreNode := pet.Root.GetChild(genre)
	if genreNode == nil {
		return []*models.Song{}
	}

	songs := make([]*models.Song, 0)
	pet.collectAllSongs(genreNode, &songs)
	return songs
}

// GetAllSongsInMood returns all songs with a specific mood across all genres
// Time Complexity: O(n) where n is the total number of songs
// Space Complexity: O(k) where k is the number of matching songs
func (pet *PlaylistExplorerTree) GetAllSongsInMood(mood string) []*models.Song {
	songs := make([]*models.Song, 0)
	pet.searchByMood(pet.Root, mood, &songs)
	return songs
}

// searchByMood recursively searches for songs with a specific mood
// Time Complexity: O(n) where n is the total number of nodes
// Space Complexity: O(d) for recursion stack where d is depth
func (pet *PlaylistExplorerTree) searchByMood(node *PlaylistTreeNode, mood string, songs *[]*models.Song) {
	if node.NodeType == MoodNode && node.Name == mood {
		// Found a mood node, collect all songs from its artist children
		pet.collectAllSongs(node, songs)
		return
	}

	// Recursively search in children
	for _, child := range node.Children {
		pet.searchByMood(child, mood, songs)
	}
}

// collectAllSongs recursively collects all songs from a subtree
// Time Complexity: O(n) where n is the number of nodes in subtree
// Space Complexity: O(d) for recursion stack where d is depth
func (pet *PlaylistExplorerTree) collectAllSongs(node *PlaylistTreeNode, songs *[]*models.Song) {
	if node.NodeType == ArtistNode {
		*songs = append(*songs, node.Songs...)
		return
	}

	// Recursively collect from children
	for _, child := range node.Children {
		pet.collectAllSongs(child, songs)
	}
}

// DepthFirstSearch performs DFS traversal and applies a function to each node
// Time Complexity: O(n) where n is the total number of nodes
// Space Complexity: O(d) for recursion stack where d is depth
func (pet *PlaylistExplorerTree) DepthFirstSearch(visitFunc func(*PlaylistTreeNode)) {
	pet.dfsHelper(pet.Root, visitFunc)
}

// dfsHelper is the recursive helper for DFS
func (pet *PlaylistExplorerTree) dfsHelper(node *PlaylistTreeNode, visitFunc func(*PlaylistTreeNode)) {
	if node == nil {
		return
	}

	visitFunc(node)

	for _, child := range node.Children {
		pet.dfsHelper(child, visitFunc)
	}
}

// BreadthFirstSearch performs BFS traversal and applies a function to each node
// Time Complexity: O(n) where n is the total number of nodes
// Space Complexity: O(w) where w is the maximum width of the tree
func (pet *PlaylistExplorerTree) BreadthFirstSearch(visitFunc func(*PlaylistTreeNode)) {
	if pet.Root == nil {
		return
	}

	queue := []*PlaylistTreeNode{pet.Root}

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		visitFunc(current)

		for _, child := range current.Children {
			queue = append(queue, child)
		}
	}
}

// GetStats returns statistics about the tree
// Time Complexity: O(1)
// Space Complexity: O(1)
func (pet *PlaylistExplorerTree) GetStats() map[string]interface{} {
	return map[string]interface{}{
		"total_songs": pet.TotalSongs,
		"genres":      pet.Stats["genres"],
		"subgenres":   pet.Stats["subgenres"],
		"moods":       pet.Stats["moods"],
		"artists":     pet.Stats["artists"],
	}
}

// FindSongPath finds the hierarchical path for a song
// Time Complexity: O(n) worst case
// Space Complexity: O(d) where d is depth
func (pet *PlaylistExplorerTree) FindSongPath(songID string) ([]string, error) {
	var foundPath []string
	var foundSong *models.Song

	// Search for the song using DFS
	pet.DepthFirstSearch(func(node *PlaylistTreeNode) {
		if node.NodeType == ArtistNode && foundSong == nil {
			for _, song := range node.Songs {
				if song.ID == songID {
					foundSong = song
					foundPath = node.GetPath()
					return
				}
			}
		}
	})

	if foundSong == nil {
		return nil, fmt.Errorf("song with ID %s not found", songID)
	}

	return foundPath, nil
}

// RemoveSong removes a song from the tree
// Time Complexity: O(n) worst case to find the song
// Space Complexity: O(d) for recursion stack
func (pet *PlaylistExplorerTree) RemoveSong(songID string) error {
	var removed bool

	pet.DepthFirstSearch(func(node *PlaylistTreeNode) {
		if node.NodeType == ArtistNode && !removed {
			for i, song := range node.Songs {
				if song.ID == songID {
					// Remove song from slice
					node.Songs = append(node.Songs[:i], node.Songs[i+1:]...)
					pet.TotalSongs--
					removed = true

					// If artist has no more songs, consider removing the artist node
					// (Implementation could be extended to clean up empty branches)
					return
				}
			}
		}
	})

	if !removed {
		return fmt.Errorf("song with ID %s not found", songID)
	}

	return nil
}

// GetTreeStructure returns a structured representation of the tree
// Time Complexity: O(n) where n is the total number of nodes
// Space Complexity: O(n)
func (pet *PlaylistExplorerTree) GetTreeStructure() map[string]interface{} {
	structure := make(map[string]interface{})

	for genreName, genreNode := range pet.Root.Children {
		genreMap := make(map[string]interface{})

		for subgenreName, subgenreNode := range genreNode.Children {
			subgenreMap := make(map[string]interface{})

			for moodName, moodNode := range subgenreNode.Children {
				moodMap := make(map[string]interface{})

				for artistName, artistNode := range moodNode.Children {
					moodMap[artistName] = len(artistNode.Songs)
				}

				subgenreMap[moodName] = moodMap
			}

			genreMap[subgenreName] = subgenreMap
		}

		structure[genreName] = genreMap
	}

	return structure
}

// String returns a string representation of the tree
// Time Complexity: O(n)
// Space Complexity: O(n)
func (pet *PlaylistExplorerTree) String() string {
	if pet.TotalSongs == 0 {
		return "Empty Playlist Explorer Tree"
	}

	result := fmt.Sprintf("Playlist Explorer Tree (Total Songs: %d)\n", pet.TotalSongs)
	result += fmt.Sprintf("Stats: %d genres, %d subgenres, %d moods, %d artists\n\n",
		pet.Stats["genres"], pet.Stats["subgenres"], pet.Stats["moods"], pet.Stats["artists"])

	pet.printTreeHelper(pet.Root, "", &result)
	return result
}

// printTreeHelper recursively builds the string representation
func (pet *PlaylistExplorerTree) printTreeHelper(node *PlaylistTreeNode, prefix string, result *string) {
	if node.Name == "Root" {
		for _, child := range node.Children {
			pet.printTreeHelper(child, "", result)
		}
		return
	}

	indent := prefix
	if node.NodeType == ArtistNode {
		*result += fmt.Sprintf("%s└── %s (%d songs)\n", indent, node.Name, len(node.Songs))
	} else {
		*result += fmt.Sprintf("%s├── %s\n", indent, node.Name)
		newPrefix := prefix + "│   "
		for _, child := range node.Children {
			pet.printTreeHelper(child, newPrefix, result)
		}
	}
}

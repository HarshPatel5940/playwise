package datastructures

import (
	"src/internal/models"
	"testing"
	"time"
)

func TestNewPlaylistSorter(t *testing.T) {
	sorter := NewPlaylistSorter(SortByTitle)
	if sorter == nil {
		t.Fatal("Expected non-nil sorter")
	}
	if sorter.GetCriteria() != SortByTitle {
		t.Errorf("Expected SortByTitle, got %v", sorter.GetCriteria())
	}
}

func TestSetAndGetCriteria(t *testing.T) {
	sorter := NewPlaylistSorter(SortByTitle)

	criteria := []SortCriteria{
		SortByTitle,
		SortByArtist,
		SortByDurationAsc,
		SortByDurationDesc,
		SortByRecentlyAdded,
		SortByOldestAdded,
		SortByRating,
		SortByPlayCount,
	}

	for _, criterion := range criteria {
		sorter.SetCriteria(criterion)
		if sorter.GetCriteria() != criterion {
			t.Errorf("Expected %v, got %v", criterion, sorter.GetCriteria())
		}
	}
}

func TestGetSortCriteriaString(t *testing.T) {
	sorter := NewPlaylistSorter(SortByTitle)

	tests := []struct {
		criteria SortCriteria
		expected string
	}{
		{SortByTitle, "Title (A-Z)"},
		{SortByArtist, "Artist (A-Z)"},
		{SortByDurationAsc, "Duration (Shortest First)"},
		{SortByDurationDesc, "Duration (Longest First)"},
		{SortByRecentlyAdded, "Recently Added"},
		{SortByOldestAdded, "Oldest Added"},
		{SortByRating, "Rating (Highest First)"},
		{SortByPlayCount, "Play Count (Most Played First)"},
	}

	for _, test := range tests {
		sorter.SetCriteria(test.criteria)
		result := sorter.GetSortCriteriaString()
		if result != test.expected {
			t.Errorf("For criteria %v, expected %q, got %q", test.criteria, test.expected, result)
		}
	}
}

func createTestSongs() []*models.Song {
	baseTime := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)

	return []*models.Song{
		{
			ID:        "1",
			Title:     "Zebra Song",
			Artist:    "Artist B",
			Duration:  240,
			Rating:    4,
			PlayCount: 10,
			AddedAt:   baseTime,
		},
		{
			ID:        "2",
			Title:     "Alpha Track",
			Artist:    "Artist A",
			Duration:  180,
			Rating:    5,
			PlayCount: 25,
			AddedAt:   baseTime.Add(24 * time.Hour),
		},
		{
			ID:        "3",
			Title:     "Beta Beat",
			Artist:    "Artist C",
			Duration:  300,
			Rating:    3,
			PlayCount: 5,
			AddedAt:   baseTime.Add(48 * time.Hour),
		},
		{
			ID:        "4",
			Title:     "Charlie Song",
			Artist:    "Artist A",
			Duration:  200,
			Rating:    4,
			PlayCount: 15,
			AddedAt:   baseTime.Add(12 * time.Hour),
		},
	}
}

func TestMergeSort(t *testing.T) {
	sorter := NewPlaylistSorter(SortByTitle)
	songs := createTestSongs()

	// Test sorting by title
	sorted := sorter.MergeSort(songs)

	if len(sorted) != len(songs) {
		t.Errorf("Expected %d songs, got %d", len(songs), len(sorted))
	}

	expectedTitles := []string{"Alpha Track", "Beta Beat", "Charlie Song", "Zebra Song"}
	for i, song := range sorted {
		if song.Title != expectedTitles[i] {
			t.Errorf("Position %d: expected %q, got %q", i, expectedTitles[i], song.Title)
		}
	}

	// Verify original slice is not modified
	if songs[0].Title != "Zebra Song" {
		t.Error("Original slice was modified")
	}
}

func TestQuickSort(t *testing.T) {
	sorter := NewPlaylistSorter(SortByArtist)
	songs := createTestSongs()

	sorted := sorter.QuickSort(songs)

	if len(sorted) != len(songs) {
		t.Errorf("Expected %d songs, got %d", len(songs), len(sorted))
	}

	// Should be sorted by artist, then by title within same artist
	expectedOrder := []string{"Alpha Track", "Charlie Song", "Zebra Song", "Beta Beat"}
	for i, song := range sorted {
		if song.Title != expectedOrder[i] {
			t.Errorf("Position %d: expected %q, got %q", i, expectedOrder[i], song.Title)
		}
	}
}

func TestHeapSort(t *testing.T) {
	sorter := NewPlaylistSorter(SortByDurationAsc)
	songs := createTestSongs()

	sorted := sorter.HeapSort(songs)

	if len(sorted) != len(songs) {
		t.Errorf("Expected %d songs, got %d", len(songs), len(sorted))
	}

	// Should be sorted by duration ascending
	expectedDurations := []int{180, 200, 240, 300}
	for i, song := range sorted {
		if song.Duration != expectedDurations[i] {
			t.Errorf("Position %d: expected duration %d, got %d", i, expectedDurations[i], song.Duration)
		}
	}
}

func TestSortingWithEmptySlice(t *testing.T) {
	sorter := NewPlaylistSorter(SortByTitle)
	songs := []*models.Song{}

	// Test all algorithms with empty slice
	mergeSorted := sorter.MergeSort(songs)
	quickSorted := sorter.QuickSort(songs)
	heapSorted := sorter.HeapSort(songs)

	if len(mergeSorted) != 0 {
		t.Error("MergeSort should return empty slice for empty input")
	}
	if len(quickSorted) != 0 {
		t.Error("QuickSort should return empty slice for empty input")
	}
	if len(heapSorted) != 0 {
		t.Error("HeapSort should return empty slice for empty input")
	}
}

func TestSortingWithSingleElement(t *testing.T) {
	sorter := NewPlaylistSorter(SortByTitle)
	songs := []*models.Song{
		{
			ID:     "1",
			Title:  "Single Song",
			Artist: "Single Artist",
		},
	}

	mergeSorted := sorter.MergeSort(songs)
	quickSorted := sorter.QuickSort(songs)
	heapSorted := sorter.HeapSort(songs)

	if len(mergeSorted) != 1 || mergeSorted[0].Title != "Single Song" {
		t.Error("MergeSort failed with single element")
	}
	if len(quickSorted) != 1 || quickSorted[0].Title != "Single Song" {
		t.Error("QuickSort failed with single element")
	}
	if len(heapSorted) != 1 || heapSorted[0].Title != "Single Song" {
		t.Error("HeapSort failed with single element")
	}
}

func TestSortByCriteria(t *testing.T) {
	songs := createTestSongs()

	tests := []struct {
		criteria  SortCriteria
		checkFunc func([]*models.Song) bool
		name      string
	}{
		{
			criteria: SortByRating,
			name:     "by rating (highest first)",
			checkFunc: func(sorted []*models.Song) bool {
				return sorted[0].Rating == 5 && sorted[len(sorted)-1].Rating == 3
			},
		},
		{
			criteria: SortByPlayCount,
			name:     "by play count (highest first)",
			checkFunc: func(sorted []*models.Song) bool {
				return sorted[0].PlayCount == 25 && sorted[len(sorted)-1].PlayCount == 5
			},
		},
		{
			criteria: SortByDurationDesc,
			name:     "by duration (longest first)",
			checkFunc: func(sorted []*models.Song) bool {
				return sorted[0].Duration == 300 && sorted[len(sorted)-1].Duration == 180
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			sorter := NewPlaylistSorter(test.criteria)
			sorted := sorter.MergeSort(songs)

			if !test.checkFunc(sorted) {
				t.Errorf("Sorting %s failed", test.name)
			}
		})
	}
}

func TestSortByTime(t *testing.T) {
	sorter := NewPlaylistSorter(SortByRecentlyAdded)
	songs := createTestSongs()

	sorted := sorter.MergeSort(songs)

	// Most recently added should be first
	if sorted[0].Title != "Beta Beat" {
		t.Errorf("Expected Beta Beat first (most recent), got %s", sorted[0].Title)
	}

	// Test oldest added
	sorter.SetCriteria(SortByOldestAdded)
	sorted = sorter.MergeSort(songs)

	// Oldest added should be first
	if sorted[0].Title != "Zebra Song" {
		t.Errorf("Expected Zebra Song first (oldest), got %s", sorted[0].Title)
	}
}

func TestMultiCriteriaSort(t *testing.T) {
	sorter := NewPlaylistSorter(SortByTitle)
	songs := createTestSongs()

	// Sort by artist first, then by title
	criteria := []SortCriteria{SortByArtist, SortByTitle}
	sorted := sorter.MultiCriteriaSort(songs, criteria)

	// First two songs should be from Artist A, sorted by title
	if sorted[0].Artist != "Artist A" || sorted[0].Title != "Alpha Track" {
		t.Errorf("Multi-criteria sort failed: expected Alpha Track first, got %s", sorted[0].Title)
	}
	if sorted[1].Artist != "Artist A" || sorted[1].Title != "Charlie Song" {
		t.Errorf("Multi-criteria sort failed: expected Charlie Song second, got %s", sorted[1].Title)
	}
}

func TestMultiCriteriaSortEdgeCases(t *testing.T) {
	sorter := NewPlaylistSorter(SortByTitle)
	songs := createTestSongs()

	// Test with empty criteria
	sorted := sorter.MultiCriteriaSort(songs, []SortCriteria{})
	if len(sorted) != len(songs) {
		t.Error("MultiCriteriaSort with empty criteria should return copy of original")
	}

	// Test with single criterion
	sorted = sorter.MultiCriteriaSort(songs, []SortCriteria{SortByTitle})
	if sorted[0].Title != "Alpha Track" {
		t.Error("MultiCriteriaSort with single criterion failed")
	}
}

func TestIsStableSorted(t *testing.T) {
	sorter := NewPlaylistSorter(SortByTitle)
	songs := createTestSongs()

	// Test unsorted
	if sorter.IsStableSorted(songs) {
		t.Error("IsStableSorted should return false for unsorted songs")
	}

	// Test sorted
	sorted := sorter.MergeSort(songs)
	if !sorter.IsStableSorted(sorted) {
		t.Error("IsStableSorted should return true for sorted songs")
	}

	// Test empty slice
	if !sorter.IsStableSorted([]*models.Song{}) {
		t.Error("IsStableSorted should return true for empty slice")
	}

	// Test single element
	if !sorter.IsStableSorted([]*models.Song{songs[0]}) {
		t.Error("IsStableSorted should return true for single element")
	}
}

func TestSortPlaylist(t *testing.T) {
	playlist := NewDoublyLinkedList()
	songs := createTestSongs()

	// Add songs to playlist
	for _, song := range songs {
		playlist.AddSong(song)
	}

	sorter := NewPlaylistSorter(SortByTitle)

	// Test different algorithms
	algorithms := []string{"merge", "quick", "heap", "unknown"}

	for _, algorithm := range algorithms {
		// Reset playlist
		playlist.Clear()
		for _, song := range songs {
			playlist.AddSong(song)
		}

		sorter.SortPlaylist(playlist, algorithm)

		// Check if sorted correctly
		current := playlist.Head
		prevTitle := ""
		for current != nil {
			if prevTitle != "" && current.Song.Title < prevTitle {
				t.Errorf("Playlist not sorted correctly with algorithm %s", algorithm)
				break
			}
			prevTitle = current.Song.Title
			current = current.Next
		}
	}
}

func TestSortPlaylistEmpty(t *testing.T) {
	playlist := NewDoublyLinkedList()
	sorter := NewPlaylistSorter(SortByTitle)

	// Should not panic on empty playlist
	sorter.SortPlaylist(playlist, "merge")

	if !playlist.IsEmpty() {
		t.Error("Empty playlist should remain empty after sort")
	}
}

func TestBenchmarkSort(t *testing.T) {
	sorter := NewPlaylistSorter(SortByTitle)
	songs := createTestSongs()

	// Test with normal dataset
	benchmarks := sorter.BenchmarkSort(songs)

	expectedAlgorithms := []string{"merge_sort", "quick_sort", "heap_sort"}
	for _, algorithm := range expectedAlgorithms {
		if _, exists := benchmarks[algorithm]; !exists {
			t.Errorf("Benchmark missing for %s", algorithm)
		}
		if benchmarks[algorithm] < 0 {
			t.Errorf("Benchmark time cannot be negative for %s", algorithm)
		}
	}

	// Test with empty dataset
	emptyBenchmarks := sorter.BenchmarkSort([]*models.Song{})
	if len(emptyBenchmarks) != 0 {
		t.Error("Empty dataset should return empty benchmarks")
	}
}

func TestCompareFunction(t *testing.T) {
	baseTime := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)

	song1 := &models.Song{
		Title:     "Alpha",
		Artist:    "Artist A",
		Duration:  180,
		Rating:    4,
		PlayCount: 10,
		AddedAt:   baseTime,
	}

	song2 := &models.Song{
		Title:     "Beta",
		Artist:    "Artist B",
		Duration:  240,
		Rating:    5,
		PlayCount: 15,
		AddedAt:   baseTime.Add(time.Hour),
	}

	tests := []struct {
		criteria SortCriteria
		expected int // -1 if song1 < song2, 0 if equal, 1 if song1 > song2
		name     string
	}{
		{SortByTitle, -1, "title comparison"},
		{SortByArtist, -1, "artist comparison"},
		{SortByDurationAsc, -1, "duration ascending"},
		{SortByDurationDesc, 1, "duration descending"},
		{SortByRating, 1, "rating (higher first)"},
		{SortByPlayCount, 1, "play count (higher first)"},
		{SortByRecentlyAdded, 1, "recently added"},
		{SortByOldestAdded, -1, "oldest added"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			sorter := NewPlaylistSorter(test.criteria)
			result := sorter.compare(song1, song2)

			if (result < 0 && test.expected >= 0) ||
				(result == 0 && test.expected != 0) ||
				(result > 0 && test.expected <= 0) {
				t.Errorf("Compare failed for %s: expected %d, got %d", test.name, test.expected, result)
			}
		})
	}
}

func TestCompareEdgeCases(t *testing.T) {
	baseTime := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)

	// Same title songs
	song1 := &models.Song{Title: "Same", Artist: "Artist A", AddedAt: baseTime}
	song2 := &models.Song{Title: "Same", Artist: "Artist B", AddedAt: baseTime}

	sorter := NewPlaylistSorter(SortByTitle)
	result := sorter.compare(song1, song2)
	if result != 0 {
		t.Error("Songs with same title should compare equal")
	}

	// Same artist, different titles
	song1.Artist = "Same Artist"
	song2.Artist = "Same Artist"
	song1.Title = "Alpha"
	song2.Title = "Beta"

	sorter.SetCriteria(SortByArtist)
	result = sorter.compare(song1, song2)
	if result >= 0 {
		t.Error("When artists are same, should sort by title")
	}

	// Same rating, should fall back to title
	song1.Rating = 4
	song2.Rating = 4
	song1.Title = "Alpha"
	song2.Title = "Beta"

	sorter.SetCriteria(SortByRating)
	result = sorter.compare(song1, song2)
	if result >= 0 {
		t.Error("When ratings are same, should sort by title")
	}
}

// Benchmark tests
func BenchmarkMergeSort(b *testing.B) {
	sorter := NewPlaylistSorter(SortByTitle)
	songs := createLargeSongDataset(1000)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sorter.MergeSort(songs)
	}
}

func BenchmarkQuickSort(b *testing.B) {
	sorter := NewPlaylistSorter(SortByTitle)
	songs := createLargeSongDataset(1000)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sorter.QuickSort(songs)
	}
}

func BenchmarkHeapSort(b *testing.B) {
	sorter := NewPlaylistSorter(SortByTitle)
	songs := createLargeSongDataset(1000)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sorter.HeapSort(songs)
	}
}

// Helper function to create large dataset for benchmarking
func createLargeSongDataset(size int) []*models.Song {
	songs := make([]*models.Song, size)
	baseTime := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)

	for i := 0; i < size; i++ {
		songs[i] = &models.Song{
			ID:        string(rune('A' + i%26)),
			Title:     generateRandomTitle(i),
			Artist:    generateRandomArtist(i),
			Duration:  120 + (i*13)%300, // Random-ish duration
			Rating:    1 + i%5,          // Rating 1-5
			PlayCount: i % 100,          // Random play count
			AddedAt:   baseTime.Add(time.Duration(i) * time.Hour),
		}
	}

	return songs
}

func generateRandomTitle(seed int) string {
	titles := []string{
		"Amazing Song", "Beautiful Track", "Cool Beat", "Dynamic Music",
		"Epic Melody", "Fantastic Tune", "Great Rhythm", "Harmonious Song",
		"Incredible Track", "Joyful Music", "Killer Beat", "Lovely Song",
		"Magnificent Tune", "Nice Track", "Outstanding Music", "Perfect Song",
	}
	return titles[seed%len(titles)]
}

func generateRandomArtist(seed int) string {
	artists := []string{
		"Artist Alpha", "Band Beta", "Creator Charlie", "DJ Delta",
		"Echo Band", "Fire Group", "Golden Voice", "Harmony Makers",
		"Indie Artists", "Jazz Masters", "Key Players", "Live Band",
		"Music Makers", "New Sound", "Orchestra Plus", "Pop Stars",
	}
	return artists[seed%len(artists)]
}

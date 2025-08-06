package datastructures

import (
	"src/internal/models"
	"strings"
	"time"
)

// SortCriteria defines the sorting criteria options
type SortCriteria int

const (
	SortByTitle SortCriteria = iota
	SortByArtist
	SortByDurationAsc
	SortByDurationDesc
	SortByRecentlyAdded
	SortByOldestAdded
	SortByRating
	SortByPlayCount
)

// PlaylistSorter provides various sorting algorithms for playlists
// Time Complexity varies by algorithm: Merge Sort O(n log n), Quick Sort O(n log n) average
// Space Complexity: Merge Sort O(n), Quick Sort O(log n) average
type PlaylistSorter struct {
	criteria SortCriteria
}

// NewPlaylistSorter creates a new playlist sorter with specified criteria
// Time Complexity: O(1)
// Space Complexity: O(1)
func NewPlaylistSorter(criteria SortCriteria) *PlaylistSorter {
	return &PlaylistSorter{
		criteria: criteria,
	}
}

// MergeSort sorts the playlist using merge sort algorithm
// Time Complexity: O(n log n) - guaranteed
// Space Complexity: O(n) - requires additional space for merging
func (ps *PlaylistSorter) MergeSort(songs []*models.Song) []*models.Song {
	if len(songs) <= 1 {
		return songs
	}

	// Create a copy to avoid modifying original slice
	result := make([]*models.Song, len(songs))
	copy(result, songs)

	ps.mergeSortHelper(result, 0, len(result)-1)
	return result
}

// mergeSortHelper is the recursive helper for merge sort
// Time Complexity: O(n log n)
// Space Complexity: O(n) due to temporary arrays and recursion stack
func (ps *PlaylistSorter) mergeSortHelper(songs []*models.Song, left, right int) {
	if left < right {
		mid := left + (right-left)/2

		// Recursively sort left and right halves
		ps.mergeSortHelper(songs, left, mid)
		ps.mergeSortHelper(songs, mid+1, right)

		// Merge the sorted halves
		ps.merge(songs, left, mid, right)
	}
}

// merge combines two sorted subarrays
// Time Complexity: O(n) where n is the size of the subarray being merged
// Space Complexity: O(n) for temporary arrays
func (ps *PlaylistSorter) merge(songs []*models.Song, left, mid, right int) {
	// Create temporary arrays for left and right subarrays
	leftSize := mid - left + 1
	rightSize := right - mid

	leftArray := make([]*models.Song, leftSize)
	rightArray := make([]*models.Song, rightSize)

	// Copy data to temporary arrays
	for i := 0; i < leftSize; i++ {
		leftArray[i] = songs[left+i]
	}
	for j := 0; j < rightSize; j++ {
		rightArray[j] = songs[mid+1+j]
	}

	// Merge the temporary arrays back into songs[left..right]
	i, j, k := 0, 0, left

	for i < leftSize && j < rightSize {
		if ps.compare(leftArray[i], rightArray[j]) <= 0 {
			songs[k] = leftArray[i]
			i++
		} else {
			songs[k] = rightArray[j]
			j++
		}
		k++
	}

	// Copy remaining elements
	for i < leftSize {
		songs[k] = leftArray[i]
		i++
		k++
	}

	for j < rightSize {
		songs[k] = rightArray[j]
		j++
		k++
	}
}

// QuickSort sorts the playlist using quick sort algorithm
// Time Complexity: O(n log n) average, O(n²) worst case
// Space Complexity: O(log n) average due to recursion stack
func (ps *PlaylistSorter) QuickSort(songs []*models.Song) []*models.Song {
	if len(songs) <= 1 {
		return songs
	}

	// Create a copy to avoid modifying original slice
	result := make([]*models.Song, len(songs))
	copy(result, songs)

	ps.quickSortHelper(result, 0, len(result)-1)
	return result
}

// quickSortHelper is the recursive helper for quick sort
// Time Complexity: O(n log n) average, O(n²) worst case
// Space Complexity: O(log n) average due to recursion stack
func (ps *PlaylistSorter) quickSortHelper(songs []*models.Song, low, high int) {
	if low < high {
		// Partition the array and get pivot index
		pivotIndex := ps.partition(songs, low, high)

		// Recursively sort elements before and after partition
		ps.quickSortHelper(songs, low, pivotIndex-1)
		ps.quickSortHelper(songs, pivotIndex+1, high)
	}
}

// partition rearranges the array around a pivot element
// Time Complexity: O(n) where n is the size of the subarray
// Space Complexity: O(1)
func (ps *PlaylistSorter) partition(songs []*models.Song, low, high int) int {
	// Choose the rightmost element as pivot
	pivot := songs[high]
	i := low - 1 // Index of smaller element

	for j := low; j < high; j++ {
		// If current element is smaller than or equal to pivot
		if ps.compare(songs[j], pivot) <= 0 {
			i++
			songs[i], songs[j] = songs[j], songs[i]
		}
	}

	// Place pivot in correct position
	songs[i+1], songs[high] = songs[high], songs[i+1]
	return i + 1
}

// HeapSort sorts the playlist using heap sort algorithm
// Time Complexity: O(n log n) - guaranteed
// Space Complexity: O(1) - in-place sorting
func (ps *PlaylistSorter) HeapSort(songs []*models.Song) []*models.Song {
	if len(songs) <= 1 {
		return songs
	}

	// Create a copy to avoid modifying original slice
	result := make([]*models.Song, len(songs))
	copy(result, songs)

	n := len(result)

	// Build max heap
	for i := n/2 - 1; i >= 0; i-- {
		ps.heapify(result, n, i)
	}

	// Extract elements from heap one by one
	for i := n - 1; i > 0; i-- {
		// Move current root to end
		result[0], result[i] = result[i], result[0]

		// Call heapify on the reduced heap
		ps.heapify(result, i, 0)
	}

	return result
}

// heapify maintains the heap property for a subtree rooted at index i
// Time Complexity: O(log n)
// Space Complexity: O(1)
func (ps *PlaylistSorter) heapify(songs []*models.Song, n, i int) {
	largest := i
	left := 2*i + 1
	right := 2*i + 2

	// If left child is larger than root
	if left < n && ps.compare(songs[left], songs[largest]) > 0 {
		largest = left
	}

	// If right child is larger than largest so far
	if right < n && ps.compare(songs[right], songs[largest]) > 0 {
		largest = right
	}

	// If largest is not root
	if largest != i {
		songs[i], songs[largest] = songs[largest], songs[i]
		ps.heapify(songs, n, largest)
	}
}

// compare compares two songs based on the current sorting criteria
// Returns: < 0 if song1 < song2, 0 if song1 == song2, > 0 if song1 > song2
// Time Complexity: O(1) for most criteria, O(k) for string comparisons
// Space Complexity: O(1)
func (ps *PlaylistSorter) compare(song1, song2 *models.Song) int {
	switch ps.criteria {
	case SortByTitle:
		return strings.Compare(strings.ToLower(song1.Title), strings.ToLower(song2.Title))

	case SortByArtist:
		artistCmp := strings.Compare(strings.ToLower(song1.Artist), strings.ToLower(song2.Artist))
		if artistCmp == 0 {
			// If same artist, sort by title
			return strings.Compare(strings.ToLower(song1.Title), strings.ToLower(song2.Title))
		}
		return artistCmp

	case SortByDurationAsc:
		return song1.Duration - song2.Duration

	case SortByDurationDesc:
		return song2.Duration - song1.Duration

	case SortByRecentlyAdded:
		if song1.AddedAt.After(song2.AddedAt) {
			return -1 // song1 is more recent
		} else if song1.AddedAt.Before(song2.AddedAt) {
			return 1 // song2 is more recent
		}
		return 0

	case SortByOldestAdded:
		if song1.AddedAt.Before(song2.AddedAt) {
			return -1 // song1 is older
		} else if song1.AddedAt.After(song2.AddedAt) {
			return 1 // song2 is older
		}
		return 0

	case SortByRating:
		ratingDiff := song2.Rating - song1.Rating // Higher ratings first
		if ratingDiff == 0 {
			// If same rating, sort by title
			return strings.Compare(strings.ToLower(song1.Title), strings.ToLower(song2.Title))
		}
		return ratingDiff

	case SortByPlayCount:
		playCountDiff := song2.PlayCount - song1.PlayCount // Higher play counts first
		if playCountDiff == 0 {
			// If same play count, sort by title
			return strings.Compare(strings.ToLower(song1.Title), strings.ToLower(song2.Title))
		}
		return playCountDiff

	default:
		return strings.Compare(strings.ToLower(song1.Title), strings.ToLower(song2.Title))
	}
}

// SetCriteria updates the sorting criteria
// Time Complexity: O(1)
// Space Complexity: O(1)
func (ps *PlaylistSorter) SetCriteria(criteria SortCriteria) {
	ps.criteria = criteria
}

// GetCriteria returns the current sorting criteria
// Time Complexity: O(1)
// Space Complexity: O(1)
func (ps *PlaylistSorter) GetCriteria() SortCriteria {
	return ps.criteria
}

// SortPlaylist sorts a doubly linked list playlist using the specified algorithm
// Time Complexity: O(n) to convert + O(n log n) to sort + O(n) to reconstruct
// Space Complexity: O(n)
func (ps *PlaylistSorter) SortPlaylist(playlist *DoublyLinkedList, algorithm string) {
	if playlist.IsEmpty() {
		return
	}

	// Convert playlist to slice
	songs := playlist.ToSlice()

	// Sort using specified algorithm
	var sortedSongs []*models.Song
	switch algorithm {
	case "merge":
		sortedSongs = ps.MergeSort(songs)
	case "quick":
		sortedSongs = ps.QuickSort(songs)
	case "heap":
		sortedSongs = ps.HeapSort(songs)
	default:
		sortedSongs = ps.MergeSort(songs) // Default to merge sort
	}

	// Reconstruct the playlist with sorted songs
	playlist.Clear()
	for _, song := range sortedSongs {
		playlist.AddSong(song)
	}
}

// MultiCriteriaSort sorts songs using multiple criteria with priority
// Time Complexity: O(n log n)
// Space Complexity: O(n)
func (ps *PlaylistSorter) MultiCriteriaSort(songs []*models.Song, criteria []SortCriteria) []*models.Song {
	if len(songs) <= 1 || len(criteria) == 0 {
		return songs
	}

	result := make([]*models.Song, len(songs))
	copy(result, songs)

	// Sort by each criterion in reverse order (last criterion first)
	for i := len(criteria) - 1; i >= 0; i-- {
		ps.criteria = criteria[i]
		result = ps.MergeSort(result) // Use stable sort for multi-criteria
	}

	return result
}

// BenchmarkSort compares performance of different sorting algorithms
// Time Complexity: Depends on algorithm and input size
// Space Complexity: O(n) for copies
func (ps *PlaylistSorter) BenchmarkSort(songs []*models.Song) map[string]time.Duration {
	if len(songs) == 0 {
		return map[string]time.Duration{}
	}

	benchmarks := make(map[string]time.Duration)

	// Benchmark Merge Sort
	start := time.Now()
	ps.MergeSort(songs)
	benchmarks["merge_sort"] = time.Since(start)

	// Benchmark Quick Sort
	start = time.Now()
	ps.QuickSort(songs)
	benchmarks["quick_sort"] = time.Since(start)

	// Benchmark Heap Sort
	start = time.Now()
	ps.HeapSort(songs)
	benchmarks["heap_sort"] = time.Since(start)

	return benchmarks
}

// IsStableSorted checks if the songs are sorted according to current criteria
// Time Complexity: O(n)
// Space Complexity: O(1)
func (ps *PlaylistSorter) IsStableSorted(songs []*models.Song) bool {
	if len(songs) <= 1 {
		return true
	}

	for i := 1; i < len(songs); i++ {
		if ps.compare(songs[i-1], songs[i]) > 0 {
			return false
		}
	}

	return true
}

// GetSortCriteriaString returns a human-readable string for the sorting criteria
// Time Complexity: O(1)
// Space Complexity: O(1)
func (ps *PlaylistSorter) GetSortCriteriaString() string {
	switch ps.criteria {
	case SortByTitle:
		return "Title (A-Z)"
	case SortByArtist:
		return "Artist (A-Z)"
	case SortByDurationAsc:
		return "Duration (Shortest First)"
	case SortByDurationDesc:
		return "Duration (Longest First)"
	case SortByRecentlyAdded:
		return "Recently Added"
	case SortByOldestAdded:
		return "Oldest Added"
	case SortByRating:
		return "Rating (Highest First)"
	case SortByPlayCount:
		return "Play Count (Most Played First)"
	default:
		return "Unknown"
	}
}

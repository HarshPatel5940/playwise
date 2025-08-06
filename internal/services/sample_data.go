package services

import (
	"src/internal/models"
)

// SampleDataLoader provides sample songs for demonstration
type SampleDataLoader struct {
	songs []*models.Song
}

// NewSampleDataLoader creates a new sample data loader
func NewSampleDataLoader() *SampleDataLoader {
	return &SampleDataLoader{
		songs: generateSampleSongs(),
	}
}

// LoadSampleData loads sample songs into the playlist engine
func (sdl *SampleDataLoader) LoadSampleData(engine *PlaylistEngine) error {
	for _, song := range sdl.songs {
		err := engine.AddSong(
			song.Title, song.Artist, song.Album,
			song.Genre, song.SubGenre, song.Mood,
			song.Duration, song.BPM,
		)
		if err != nil {
			// Continue loading other songs even if one fails
			continue
		}

		// Set rating if provided
		if song.Rating > 0 {
			engine.RateSong(song.ID, song.Rating)
		}
	}
	return nil
}

// GetSampleSongs returns all sample songs
func (sdl *SampleDataLoader) GetSampleSongs() []*models.Song {
	return sdl.songs
}

// generateSampleSongs creates a comprehensive set of sample songs
func generateSampleSongs() []*models.Song {
	sampleData := []struct {
		title    string
		artist   string
		album    string
		genre    string
		subgenre string
		mood     string
		duration int
		bpm      int
		rating   int
	}{
		// Rock Songs
		{"Bohemian Rhapsody", "Queen", "A Night at the Opera", "Rock", "Progressive Rock", "Dramatic", 355, 72, 5},
		{"Stairway to Heaven", "Led Zeppelin", "Led Zeppelin IV", "Rock", "Hard Rock", "Epic", 482, 82, 5},
		{"Hotel California", "Eagles", "Hotel California", "Rock", "Soft Rock", "Mysterious", 391, 75, 5},
		{"Sweet Child O' Mine", "Guns N' Roses", "Appetite for Destruction", "Rock", "Hard Rock", "Energetic", 356, 125, 4},
		{"Smells Like Teen Spirit", "Nirvana", "Nevermind", "Rock", "Alternative Rock", "Aggressive", 301, 117, 4},
		{"Wonderwall", "Oasis", "What's the Story Morning Glory?", "Rock", "Alternative Rock", "Nostalgic", 258, 87, 4},
		{"Creep", "Radiohead", "Pablo Honey", "Rock", "Alternative Rock", "Melancholic", 238, 92, 4},
		{"Black", "Pearl Jam", "Ten", "Rock", "Grunge", "Emotional", 341, 69, 4},
		{"Paranoid Android", "Radiohead", "OK Computer", "Rock", "Alternative Rock", "Complex", 383, 64, 5},
		{"Jeremy", "Pearl Jam", "Ten", "Rock", "Grunge", "Dark", 318, 86, 4},

		// Pop Songs
		{"Shape of You", "Ed Sheeran", "÷", "Pop", "Pop Rock", "Happy", 233, 96, 4},
		{"Blinding Lights", "The Weeknd", "After Hours", "Pop", "Synthpop", "Energetic", 200, 171, 5},
		{"Bad Guy", "Billie Eilish", "When We All Fall Asleep Where Do We Go?", "Pop", "Electropop", "Dark", 194, 135, 4},
		{"Levitating", "Dua Lipa", "Future Nostalgia", "Pop", "Dance Pop", "Upbeat", 203, 103, 4},
		{"Anti-Hero", "Taylor Swift", "Midnights", "Pop", "Indie Pop", "Introspective", 200, 97, 4},
		{"As It Was", "Harry Styles", "Harry's House", "Pop", "Pop Rock", "Nostalgic", 167, 173, 4},
		{"Good 4 U", "Olivia Rodrigo", "Sour", "Pop", "Pop Punk", "Angry", 178, 166, 4},
		{"Stay", "The Kid LAROI & Justin Bieber", "F*ck Love 3", "Pop", "Pop Rap", "Romantic", 141, 169, 3},
		{"Watermelon Sugar", "Harry Styles", "Fine Line", "Pop", "Pop Rock", "Happy", 174, 95, 4},
		{"Don't Start Now", "Dua Lipa", "Future Nostalgia", "Pop", "Dance Pop", "Confident", 183, 124, 4},

		// Hip Hop Songs
		{"HUMBLE.", "Kendrick Lamar", "DAMN.", "Hip Hop", "Conscious Rap", "Aggressive", 177, 150, 5},
		{"God's Plan", "Drake", "Scorpion", "Hip Hop", "Pop Rap", "Confident", 198, 77, 4},
		{"Sicko Mode", "Travis Scott", "Astroworld", "Hip Hop", "Trap", "Dark", 312, 155, 4},
		{"Old Town Road", "Lil Nas X", "7 EP", "Hip Hop", "Country Rap", "Fun", 113, 136, 3},
		{"Lose Yourself", "Eminem", "8 Mile Soundtrack", "Hip Hop", "Hardcore Hip Hop", "Motivational", 326, 86, 5},
		{"Alright", "Kendrick Lamar", "To Pimp a Butterfly", "Hip Hop", "Conscious Rap", "Hopeful", 219, 100, 5},
		{"Money Trees", "Kendrick Lamar", "Good Kid M.A.A.D City", "Hip Hop", "West Coast Hip Hop", "Reflective", 384, 80, 4},
		{"INDUSTRY BABY", "Lil Nas X & Jack Harlow", "Montero", "Hip Hop", "Pop Rap", "Confident", 212, 149, 3},
		{"Life Is Good", "Future & Drake", "High Off Life", "Hip Hop", "Trap", "Boastful", 243, 81, 3},
		{"Rockstar", "Post Malone & 21 Savage", "Beerbongs & Bentleys", "Hip Hop", "Pop Rap", "Braggadocious", 218, 160, 4},

		// Electronic Songs
		{"Levels", "Avicii", "Original Mix", "Electronic", "Progressive House", "Euphoric", 203, 126, 4},
		{"Titanium", "David Guetta ft. Sia", "Nothing But The Beat", "Electronic", "Electro House", "Empowering", 245, 126, 4},
		{"Clarity", "Zedd ft. Foxes", "Clarity", "Electronic", "Progressive House", "Emotional", 271, 128, 4},
		{"Animals", "Martin Garrix", "Single", "Electronic", "Big Room House", "Aggressive", 302, 128, 3},
		{"Strobe", "Deadmau5", "For Lack of a Better Name", "Electronic", "Progressive House", "Atmospheric", 645, 128, 5},
		{"One More Time", "Daft Punk", "Discovery", "Electronic", "French House", "Joyful", 320, 123, 5},
		{"Midnight City", "M83", "Hurry Up We're Dreaming", "Electronic", "Synthwave", "Dreamy", 244, 104, 4},
		{"Breathe Me", "Sia", "Colour The Small One", "Electronic", "Electropop", "Vulnerable", 268, 75, 4},
		{"Scary Monsters and Nice Sprites", "Skrillex", "Scary Monsters and Nice Sprites", "Electronic", "Dubstep", "Chaotic", 225, 140, 3},
		{"Ghosts 'n' Stuff", "Deadmau5", "For Lack of a Better Name", "Electronic", "Electro House", "Dark", 335, 128, 4},

		// Jazz Songs
		{"Take Five", "Dave Brubeck Quartet", "Time Out", "Jazz", "Cool Jazz", "Sophisticated", 324, 175, 5},
		{"Kind of Blue", "Miles Davis", "Kind of Blue", "Jazz", "Modal Jazz", "Contemplative", 567, 120, 5},
		{"A Love Supreme", "John Coltrane", "A Love Supreme", "Jazz", "Spiritual Jazz", "Transcendent", 487, 80, 5},
		{"So What", "Miles Davis", "Kind of Blue", "Jazz", "Modal Jazz", "Cool", 563, 132, 5},
		{"Giant Steps", "John Coltrane", "Giant Steps", "Jazz", "Hard Bop", "Complex", 287, 290, 4},
		{"Blue in Green", "Miles Davis", "Kind of Blue", "Jazz", "Modal Jazz", "Melancholic", 337, 66, 4},
		{"Autumn Leaves", "Bill Evans Trio", "Sunday at the Village Vanguard", "Jazz", "Post Bop", "Nostalgic", 472, 108, 4},
		{"Maiden Voyage", "Herbie Hancock", "Maiden Voyage", "Jazz", "Post Bop", "Adventurous", 503, 120, 4},
		{"Summertime", "Ella Fitzgerald", "Porgy and Bess", "Jazz", "Vocal Jazz", "Dreamy", 253, 72, 4},
		{"Round Midnight", "Thelonious Monk", "Genius of Modern Music", "Jazz", "Bebop", "Mysterious", 311, 55, 4},

		// Classical Songs
		{"Symphony No. 9", "Ludwig van Beethoven", "Symphony No. 9", "Classical", "Romantic", "Triumphant", 4200, 120, 5},
		{"The Four Seasons - Spring", "Antonio Vivaldi", "The Four Seasons", "Classical", "Baroque", "Joyful", 600, 100, 5},
		{"Canon in D", "Johann Pachelbel", "Canon and Gigue", "Classical", "Baroque", "Peaceful", 360, 54, 4},
		{"Für Elise", "Ludwig van Beethoven", "Bagatelle No. 25", "Classical", "Classical", "Gentle", 195, 120, 4},
		{"Ave Maria", "Franz Schubert", "Ellens Gesang III", "Classical", "Romantic", "Sacred", 390, 72, 4},
		{"Moonlight Sonata", "Ludwig van Beethoven", "Piano Sonata No. 14", "Classical", "Classical", "Melancholic", 900, 27, 5},
		{"Eine kleine Nachtmusik", "Wolfgang Amadeus Mozart", "Serenade No. 13", "Classical", "Classical", "Elegant", 1800, 120, 4},
		{"Clair de Lune", "Claude Debussy", "Suite Bergamasque", "Classical", "Impressionist", "Dreamy", 300, 50, 5},
		{"The Blue Danube", "Johann Strauss II", "The Blue Danube", "Classical", "Romantic", "Graceful", 720, 180, 4},
		{"Ride of the Valkyries", "Richard Wagner", "Die Walküre", "Classical", "Romantic", "Epic", 500, 138, 4},

		// Country Songs
		{"Friends in Low Places", "Garth Brooks", "No Fences", "Country", "Country Pop", "Nostalgic", 259, 120, 4},
		{"Sweet Caroline", "Neil Diamond", "Brother Love's Travelling Salvation Show", "Country", "Country Pop", "Happy", 201, 125, 4},
		{"Wagon Wheel", "Darius Rucker", "True Believers", "Country", "Country Rock", "Uplifting", 191, 150, 3},
		{"Cruise", "Florida Georgia Line", "Here's to the Good Times", "Country", "Country Pop", "Fun", 200, 120, 3},
		{"Need You Now", "Lady Antebellum", "Need You Now", "Country", "Country Pop", "Longing", 236, 120, 4},
		{"Before He Cheats", "Carrie Underwood", "Some Hearts", "Country", "Country Pop", "Vengeful", 199, 120, 4},
		{"Body Like a Back Road", "Sam Hunt", "Montevallo", "Country", "Country Pop", "Romantic", 157, 98, 3},
		{"Chicken Fried", "Zac Brown Band", "The Foundation", "Country", "Country Rock", "Carefree", 239, 120, 4},
		{"Live Like You Were Dying", "Tim McGraw", "Live Like You Were Dying", "Country", "Country Pop", "Inspirational", 289, 76, 4},
		{"Man! I Feel Like a Woman!", "Shania Twain", "Come On Over", "Country", "Country Pop", "Empowering", 298, 135, 4},

		// R&B Songs
		{"Superstition", "Stevie Wonder", "Talking Book", "R&B", "Funk", "Groovy", 245, 100, 5},
		{"What's Going On", "Marvin Gaye", "What's Going On", "R&B", "Soul", "Conscious", 231, 74, 5},
		{"Respect", "Aretha Franklin", "I Never Loved a Man", "R&B", "Soul", "Empowering", 147, 115, 5},
		{"I Want You Back", "The Jackson 5", "Diana Ross Presents The Jackson 5", "R&B", "Motown", "Joyful", 179, 100, 4},
		{"Let's Stay Together", "Al Green", "Let's Stay Together", "R&B", "Southern Soul", "Romantic", 199, 96, 4},
		{"I Heard It Through the Grapevine", "Marvin Gaye", "In the Groove", "R&B", "Motown", "Dramatic", 195, 82, 4},
		{"My Girl", "The Temptations", "The Temptations Sing Smokey", "R&B", "Motown", "Loving", 175, 120, 4},
		{"Stand By Me", "Ben E. King", "Don't Play That Song!", "R&B", "Doo-wop", "Comforting", 181, 118, 4},
		{"I Got You (I Feel Good)", "James Brown", "Papa's Got a Brand New Bag", "R&B", "Funk", "Energetic", 158, 144, 4},
		{"Sexual Healing", "Marvin Gaye", "Midnight Love", "R&B", "Contemporary R&B", "Sensual", 241, 103, 4},
	}

	songs := make([]*models.Song, 0, len(sampleData))

	for _, data := range sampleData {
		song := models.NewSong(
			"", // ID will be generated
			data.title,
			data.artist,
			data.album,
			data.genre,
			data.subgenre,
			data.mood,
			data.duration,
			data.bpm,
		)
		song.Rating = data.rating
		songs = append(songs, song)
	}

	return songs
}

// GetSamplePlaylistsByGenre returns sample playlists organized by genre
func (sdl *SampleDataLoader) GetSamplePlaylistsByGenre() map[string][]*models.Song {
	playlists := make(map[string][]*models.Song)

	for _, song := range sdl.songs {
		genre := song.Genre
		if playlists[genre] == nil {
			playlists[genre] = make([]*models.Song, 0)
		}
		playlists[genre] = append(playlists[genre], song)
	}

	return playlists
}

// GetTopRatedSongs returns songs with rating >= minRating
func (sdl *SampleDataLoader) GetTopRatedSongs(minRating int) []*models.Song {
	topSongs := make([]*models.Song, 0)

	for _, song := range sdl.songs {
		if song.Rating >= minRating {
			topSongs = append(topSongs, song)
		}
	}

	return topSongs
}

// GetSongsByMood returns songs with specified mood
func (sdl *SampleDataLoader) GetSongsByMood(mood string) []*models.Song {
	moodSongs := make([]*models.Song, 0)

	for _, song := range sdl.songs {
		if song.Mood == mood {
			moodSongs = append(moodSongs, song)
		}
	}

	return moodSongs
}

// GetGenreStatistics returns statistics about genres in sample data
func (sdl *SampleDataLoader) GetGenreStatistics() map[string]interface{} {
	stats := make(map[string]interface{})
	genreCounts := make(map[string]int)
	totalDuration := 0
	totalSongs := len(sdl.songs)

	for _, song := range sdl.songs {
		genreCounts[song.Genre]++
		totalDuration += song.Duration
	}

	stats["total_songs"] = totalSongs
	stats["total_genres"] = len(genreCounts)
	stats["total_duration"] = totalDuration
	stats["average_duration"] = totalDuration / totalSongs
	stats["genre_distribution"] = genreCounts

	return stats
}

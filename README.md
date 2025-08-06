# 🎵 Playwise
###  Advanced Music Playlist Engine

A comprehensive music playlist management system implementing advanced data structures and algorithms for efficient music organization, intelligent recommendations, and real-time analytics.

**Assignment**: Common Core Implementation: Modules & Scenarios + Specialized Use Cases
**Student**: Harsh N Patel
**Roll No**: RA2211028010127
**Email**: hp8823@srmist.edu.in
**Framework**: Go + Echo + Templ + Tailwind CSS

## 🌟 Features
### Core Data Structures Implemented

1. **🔗 Playlist Engine using Doubly Linked Lists**
   - Add, delete, reorder, and reverse songs
   - Bidirectional traversal optimization
   - O(1) insertion/deletion at ends

2. **📚 Playback History using Stack**
   - LIFO structure for undo functionality
   - Bounded stack to prevent memory bloat
   - Recent play tracking with statistics

3. **🌳 Song Rating Tree using Binary Search Tree**
   - Organize songs by 1-5 star ratings
   - Rating buckets for multiple songs per rating
   - Efficient range queries and sorted retrieval

4. **⚡ Instant Song Lookup using HashMap**
   - O(1) average lookup by song ID or title
   - Custom hash function with collision resolution
   - Automatic resizing with load factor monitoring

5. **🔄 Time-based Sorting using Multiple Algorithms**
   - Merge Sort (stable, guaranteed O(n log n))
   - Quick Sort (fast average case)
   - Heap Sort (in-place, guaranteed performance)
   - Multiple sorting criteria support

6. **📊 System Snapshot Module**
   - Live dashboard with real-time statistics
   - Performance metrics and analytics
   - Top songs, rating distribution, and usage stats

### Specialized Features

7. **🌲 Playlist Explorer Tree (N-ary Tree)**
   - Hierarchical navigation: Genre → Subgenre → Mood → Artist
   - DFS/BFS traversal capabilities
   - Dynamic category organization

8. **🎯 Smart Recommendations**
   - AI-powered song suggestions
   - Based on listening history and similarity
   - Filters recently played songs

## 🚀 Getting Started

### Prerequisites

- Go 1.21 or higher
- Git

### Installation

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd playwise
   ```

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Install Templ (if not already installed)**
   ```bash
   go install github.com/a-h/templ/cmd/templ@latest
   ```

4. **Generate templates**
   ```bash
   ~/go/bin/templ generate
   ```

5. **Build the application**
   ```bash
   go build -o main ./cmd/api
   ```

6. **Run the application**
   ```bash
   ./main
   ```

The application will start on `http://localhost:8080`

## 🔧 API Endpoints

### Playlist Management
```http
GET    /api/playlist                    # Get current playlist
POST   /api/playlist/songs             # Add new song
DELETE /api/playlist/songs/:index      # Delete song by index
PUT    /api/playlist/songs/:from/move/:to # Move song
POST   /api/playlist/reverse           # Reverse playlist
POST   /api/playlist/sample-data       # Load sample data
```

### Playback Operations
```http
POST   /api/playlist/songs/:index/play # Play song
POST   /api/playlist/undo              # Undo last play
GET    /api/playlist/history           # Get playback history
```

### Search & Sorting
```http
GET    /api/playlist/search            # Search songs (by ID/title)
POST   /api/playlist/sort              # Sort playlist
GET    /api/playlist/benchmark         # Benchmark sorting algorithms
```

### Rating System
```http
POST   /api/playlist/songs/:id/rate    # Rate a song (1-5 stars)
GET    /api/playlist/rating/:rating    # Get songs by rating
```

### Music Explorer
```http
GET    /api/explorer/genres                    # Get all genres
GET    /api/explorer/genres/:genre/subgenres   # Get subgenres
GET    /api/explorer/songs                     # Get songs by path
```

### Analytics
```http
GET    /api/playlist/recommendations   # Smart recommendations
GET    /api/playlist/stats             # Playlist statistics
GET    /api/dashboard                  # Live dashboard snapshot
```

## 🏗️ Architecture

### Project Structure
```
playwise/
├── cmd/
│   ├── api/                    # Main application entry
│   └── web/                    # Web templates and handlers
├── internal/
│   ├── datastructures/         # Core data structure implementations
│   │   ├── doubly_linked_list.go
│   │   ├── stack.go
│   │   ├── bst.go
│   │   ├── hashmap.go
│   │   ├── sorting.go
│   │   └── playlist_tree.go
│   ├── models/                 # Data models
│   │   └── song.go
│   ├── services/               # Business logic layer
│   │   ├── playlist_engine.go
│   │   └── sample_data.go
│   └── server/                 # HTTP handlers and routing
│       ├── server.go
│       ├── routes.go
│       └── playlist_handlers.go
└── TECHNICAL_DESIGN.md         # Comprehensive technical documentation
```

### Data Flow
```
Web Interface → HTTP Handlers → Service Layer → Data Structures → Models
```

## 📊 Performance Analysis

### Time Complexity Summary
| Operation | Data Structure | Average | Worst Case | Space |
|-----------|---------------|---------|------------|-------|
| Add Song | Doubly Linked List | O(1) | O(1) | O(1) |
| Search Song | HashMap | O(1) | O(n) | O(1) |
| Sort Playlist | Merge/Quick Sort | O(n log n) | O(n²)* | O(n) |
| Rate Song | BST | O(log n) | O(n) | O(1) |
| Tree Navigation | N-ary Tree | O(1) | O(1) | O(1) |

*Quick Sort worst case

### Memory Usage
- **Small Playlist (100 songs)**: ~45KB total
- **Large Playlist (10,000 songs)**: ~4.5MB total
- **Bounded History Stack**: Max 100 entries (~25KB)

## 🧪 Testing

### Run Tests
```bash
go test ./...
```

### Performance Benchmarks
```bash
go test -bench=. -benchmem ./internal/datastructures/
```

### API Testing
Use the built-in benchmark endpoint:
```bash
curl http://localhost:8080/api/playlist/benchmark
```

## 🎯 Key Algorithms Implemented

### 1. Doubly Linked List Operations
- **Add Song**: O(1) tail insertion
- **Move Song**: O(n) with bidirectional optimization
- **Reverse**: O(n) pointer manipulation

### 2. Hash Map with Collision Resolution
- **DJB2 Hash Function**: Excellent distribution
- **Separate Chaining**: Handles collisions gracefully
- **Dynamic Resizing**: Maintains optimal load factor

### 3. BST with Rating Buckets
- **Insertion**: O(log n) average case
- **Range Queries**: Efficient rating-based searches
- **Balanced Operations**: Minimizes tree depth

### 4. Sorting Algorithm Comparison
- **Merge Sort**: Stable, predictable performance
- **Quick Sort**: Fast average case with pivot optimization
- **Heap Sort**: In-place with guaranteed performance

### 5. N-ary Tree Traversal
- **DFS**: Complete tree exploration
- **BFS**: Level-order processing
- **Path Navigation**: O(1) hierarchical lookup

## 🔍 Design Decisions & Trade-offs

### 1. Doubly vs Singly Linked List
**✅ Chose Doubly**: Enables efficient bidirectional traversal and O(n) reverse operation
**❌ Trade-off**: Higher memory overhead (extra pointer per node)

### 2. Custom HashMap vs Built-in Map
**✅ Chose Custom**: Educational value, performance monitoring, custom optimizations
**❌ Trade-off**: More development time, potential for bugs

### 3. BST vs Array for Ratings
**✅ Chose BST**: Efficient range queries, ordered traversal, dynamic operations
**❌ Trade-off**: More complex than simple array buckets

### 4. Bounded vs Unbounded History
**✅ Chose Bounded**: Prevents memory bloat, predictable usage
**❌ Trade-off**: Loses oldest history entries

## 📈 Live Dashboard Features

- **Real-time Statistics**: Song counts, duration, ratings
- **Top 5 Longest Songs**: Dynamic ranking
- **Rating Distribution**: Visual representation of user preferences
- **Performance Metrics**: Hash map load factors, operation counts
- **Genre Statistics**: Hierarchical data breakdown

## 🎵 Sample Data

The system includes 80+ carefully curated sample songs across genres:
- **Rock**: Classic and alternative tracks
- **Pop**: Modern hits and classics
- **Hip Hop**: Conscious rap and mainstream hits
- **Electronic**: House, dubstep, and synthwave
- **Jazz**: Traditional and modern jazz standards
- **Classical**: Orchestral masterpieces
- **Country**: Traditional and contemporary country
- **R&B**: Soul, funk, and contemporary R&B

## 🚀 Future Enhancements

### Planned Features
- **Persistent Storage**: Database integration
- **User Authentication**: Multi-user support
- **Real-time Updates**: WebSocket integration
- **Advanced Search**: Full-text search with fuzzy matching
- **Machine Learning**: Enhanced recommendation algorithms

### Scalability Improvements
- **Microservices Architecture**: Service decomposition
- **Distributed Storage**: Sharded playlist storage
- **Caching Layer**: Redis integration for performance
- **API Rate Limiting**: Protection against abuse

## 📚 Learning Outcomes

This project demonstrates:
- **Data Structure Implementation**: From scratch implementations with complexity analysis
- **Algorithm Optimization**: Performance-conscious design decisions
- **System Design**: Scalable architecture with clean separation of concerns
- **Web Development**: Modern Go web development practices
- **Performance Analysis**: Benchmarking and optimization techniques

## 📄 License
This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

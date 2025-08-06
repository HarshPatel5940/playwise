#!/bin/bash

# Playwise Demo Script
# This script demonstrates the key features of the Playwise music playlist engine

echo "🎵 Welcome to Playwise - Advanced Music Playlist Engine Demo"
echo "=========================================================="
echo ""

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Base URL
BASE_URL="http://localhost:8080"

# Function to make API calls and display results
api_call() {
    local method=$1
    local endpoint=$2
    local data=$3
    local description=$4

    echo -e "${BLUE}📡 $description${NC}"
    echo "   $method $endpoint"

    if [ -n "$data" ]; then
        response=$(curl -s -X $method "$BASE_URL$endpoint" -H "Content-Type: application/json" -d "$data")
    else
        response=$(curl -s -X $method "$BASE_URL$endpoint")
    fi

    echo "   Response: $(echo $response | jq -r '.message // .error // "Success"')"
    echo ""
}

# Function to check if server is running
check_server() {
    echo -e "${YELLOW}🔍 Checking if Playwise server is running...${NC}"
    if curl -s "$BASE_URL/health" > /dev/null; then
        echo -e "${GREEN}✅ Server is running!${NC}"
    else
        echo -e "${RED}❌ Server is not running. Please start the server first:${NC}"
        echo "   ./main"
        exit 1
    fi
    echo ""
}

# Function to wait for user input
wait_for_user() {
    echo -e "${YELLOW}Press Enter to continue...${NC}"
    read
}

# Main demo flow
main() {
    check_server

    echo -e "${GREEN}🚀 Starting Playwise Demo${NC}"
    echo ""

    # Demo 1: Load Sample Data
    echo -e "${BLUE}📊 Demo 1: Loading Sample Data${NC}"
    echo "Loading 80+ curated songs across multiple genres..."
    api_call "POST" "/api/playlist/sample-data" "" "Load sample music library"
    wait_for_user

    # Demo 2: Basic Playlist Operations
    echo -e "${BLUE}📝 Demo 2: Basic Playlist Operations${NC}"

    # Get current playlist
    api_call "GET" "/api/playlist" "" "Get current playlist"

    # Add a custom song
    custom_song='{"title":"Demo Song","artist":"Test Artist","album":"Demo Album","genre":"Pop","subgenre":"Indie Pop","mood":"Happy","duration":210,"bpm":120}'
    api_call "POST" "/api/playlist/songs" "$custom_song" "Add custom song to playlist"

    # Move song (move last song to position 5)
    api_call "PUT" "/api/playlist/songs/80/move/5" "" "Move song from position 80 to position 5"

    wait_for_user

    # Demo 3: Playback and History
    echo -e "${BLUE}🎵 Demo 3: Playback and History Management${NC}"

    # Play some songs
    api_call "POST" "/api/playlist/songs/0/play" "" "Play first song"
    api_call "POST" "/api/playlist/songs/10/play" "" "Play song at index 10"
    api_call "POST" "/api/playlist/songs/25/play" "" "Play song at index 25"

    # Check history
    api_call "GET" "/api/playlist/history?count=5" "" "Get recent playback history"

    # Undo last play
    api_call "POST" "/api/playlist/undo" "" "Undo last played song"

    wait_for_user

    # Demo 4: Rating System
    echo -e "${BLUE}⭐ Demo 4: Song Rating System (BST)${NC}"

    # Rate some songs (we need to get song IDs first)
    echo "Getting song IDs for rating demo..."
    response=$(curl -s "$BASE_URL/api/playlist")
    song_id_1=$(echo $response | jq -r '.data.songs[0].id')
    song_id_2=$(echo $response | jq -r '.data.songs[1].id')
    song_id_3=$(echo $response | jq -r '.data.songs[2].id')

    # Rate songs
    api_call "POST" "/api/playlist/songs/$song_id_1/rate" '{"rating":5}' "Rate first song 5 stars"
    api_call "POST" "/api/playlist/songs/$song_id_2/rate" '{"rating":4}' "Rate second song 4 stars"
    api_call "POST" "/api/playlist/songs/$song_id_3/rate" '{"rating":5}' "Rate third song 5 stars"

    # Get songs by rating
    api_call "GET" "/api/playlist/rating/5" "" "Get all 5-star rated songs"

    wait_for_user

    # Demo 5: Search Functionality
    echo -e "${BLUE}🔍 Demo 5: Hash Map Search (O(1) Lookup)${NC}"

    # Search by title
    api_call "GET" "/api/playlist/search?type=title&q=Bohemian Rhapsody" "" "Search for 'Bohemian Rhapsody' by title"

    # Search by ID
    api_call "GET" "/api/playlist/search?type=id&q=$song_id_1" "" "Search for song by ID"

    wait_for_user

    # Demo 6: Sorting Algorithms
    echo -e "${BLUE}🔄 Demo 6: Sorting Algorithms Comparison${NC}"

    # Benchmark different sorting algorithms
    api_call "GET" "/api/playlist/benchmark" "" "Benchmark sorting algorithm performance"

    # Sort by different criteria
    api_call "POST" "/api/playlist/sort" '{"criteria":"title","algorithm":"merge"}' "Sort playlist by title (Merge Sort)"
    api_call "POST" "/api/playlist/sort" '{"criteria":"duration_desc","algorithm":"quick"}' "Sort playlist by duration descending (Quick Sort)"
    api_call "POST" "/api/playlist/sort" '{"criteria":"rating","algorithm":"heap"}' "Sort playlist by rating (Heap Sort)"

    wait_for_user

    # Demo 7: Playlist Explorer Tree
    echo -e "${BLUE}🌳 Demo 7: Hierarchical Music Explorer (N-ary Tree)${NC}"

    # Navigate the music hierarchy
    api_call "GET" "/api/explorer/genres" "" "Get all music genres"
    api_call "GET" "/api/explorer/genres/Rock/subgenres" "" "Get Rock subgenres"
    api_call "GET" "/api/explorer/genres/Rock/subgenres/Alternative Rock/moods" "" "Get Alternative Rock moods"
    api_call "GET" "/api/explorer/genres/Rock/subgenres/Alternative Rock/moods/Energetic/artists" "" "Get Energetic Alternative Rock artists"
    api_call "GET" "/api/explorer/songs?genre=Rock&subgenre=Alternative Rock&mood=Energetic&artist=Foo Fighters" "" "Get Foo Fighters energetic alternative rock songs"

    wait_for_user

    # Demo 8: Smart Recommendations
    echo -e "${BLUE}🎯 Demo 8: AI-Powered Smart Recommendations${NC}"

    # Get recommendations based on listening history
    api_call "GET" "/api/playlist/recommendations?count=10" "" "Get 10 smart song recommendations"

    wait_for_user

    # Demo 9: Live Dashboard and Analytics
    echo -e "${BLUE}📈 Demo 9: Live Dashboard and Analytics${NC}"

    # Get comprehensive statistics
    api_call "GET" "/api/playlist/stats" "" "Get detailed playlist statistics"
    api_call "GET" "/api/dashboard" "" "Get live dashboard snapshot"

    wait_for_user

    # Demo 10: Advanced Operations
    echo -e "${BLUE}⚡ Demo 10: Advanced Playlist Operations${NC}"

    # Reverse the entire playlist
    api_call "POST" "/api/playlist/reverse" "" "Reverse entire playlist order"

    # Update playlist name
    api_call "PUT" "/api/playlist/name" '{"name":"My Awesome Demo Playlist"}' "Update playlist name"

    echo ""
    echo -e "${GREEN}🎉 Demo Complete!${NC}"
    echo ""
    echo -e "${YELLOW}🌐 Want to explore more?${NC}"
    echo "   • Open your browser: $BASE_URL/playlist"
    echo "   • Try the interactive web interface"
    echo "   • Explore the technical documentation: TECHNICAL_DESIGN.md"
    echo ""
    echo -e "${BLUE}📊 Key Features Demonstrated:${NC}"
    echo "   ✅ Doubly Linked List (Playlist Management)"
    echo "   ✅ Stack (Playback History with Undo)"
    echo "   ✅ Binary Search Tree (Song Rating System)"
    echo "   ✅ Hash Map (O(1) Song Lookup)"
    echo "   ✅ Sorting Algorithms (Merge, Quick, Heap Sort)"
    echo "   ✅ N-ary Tree (Hierarchical Music Explorer)"
    echo "   ✅ Smart Recommendations (AI-powered)"
    echo "   ✅ Live Dashboard (Real-time Analytics)"
    echo ""
    echo -e "${GREEN}🎵 Thank you for trying Playwise!${NC}"
}

# Check if jq is installed
if ! command -v jq &> /dev/null; then
    echo -e "${RED}❌ jq is required for this demo. Please install it:${NC}"
    echo "   • Ubuntu/Debian: sudo apt-get install jq"
    echo "   • macOS: brew install jq"
    echo "   • Windows: Download from https://stedolan.github.io/jq/"
    exit 1
fi

# Run the main demo
main

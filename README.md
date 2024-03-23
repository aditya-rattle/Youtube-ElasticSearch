# YouTube Video Search and Indexing Service

This project is a YouTube video search and indexing service built in Go. It continuously queries the YouTube API for the latest videos matching a predefined search query, fetches relevant video data, and stores it in Elasticsearch for further analysis or retrieval.

## Features

- **Continuous YouTube API Querying**: The service continuously queries the YouTube API at a predefined interval to fetch the latest videos matching a predefined search query.

- **Background Processing**: YouTube API querying and data retrieval are performed asynchronously in the background to ensure minimal impact on service responsiveness.

- **Elasticsearch Integration**: Fetched video data is stored in Elasticsearch for efficient indexing and searching.

- **Search and Retrieval**: The service provides endpoints to search for videos in Elasticsearch based on various criteria and retrieve video data as needed.

## Setup

1. **Clone the Repository**: Clone this repository to your local machine.

2. **Install Dependencies**: Make sure you have Go installed on your machine. Install the project dependencies using `go mod tidy`.

3. **Configure Elasticsearch**: Ensure that Elasticsearch is installed and running locally.

4. **Run the Service**: Start the service by running `go run main.go`.

5. **Access the API**: Once the service is running, you can access the API endpoints to search for videos and retrieve video data.



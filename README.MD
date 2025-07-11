# 📊 TTV-Statistics API

A lightweight Go API that fetches and aggregates Twitch streamer video statistics via the Helix API.

---

## 📚 Contents

* [📌 Endpoints](#endpoints)
* [💻 Running the Application Locally](#running-the-application-locally)
* [🐳 Running the Application Using Docker](#-running-the-application-using-docker)
* [📈 Get Streamer Video Statistics](#get-streamer-video-statistics)

---

## 📌 Endpoints

* [`GET /streamer/{username}/statistics`](#get-streamer-video-statistics)

---

## 💻 Running the Application Locally

To run the application in a local environment for development purposes, you can use the following command:

```bash
go run . \
  --host=:<PORT_NUMBER> \
  --client-id=<YOUR_CLIENT_ID> \
  --client-secret=<YOUR_CLIENT_SECRET> \
  --helix-host=https://api.twitch.tv/helix
```

---

## 🐳 Running the Application Using Docker

First, create a .env file at the root of the repository, supplying your credentials following the [example .env file](.env.example)

To build and run the application using Docker Compose:

```bash
docker compose up --build
```

This command uses the provided `docker-compose.yaml` to:

* 🛠️ Build the application image
* 🚀 Start the container and run the server

> ✅ **Note:** The server runs in the foreground, and logs will be visible directly in your terminal.

### 🔻 Graceful Shutdown

To properly shut down the application and observe the server’s graceful shutdown logs:

1. 💾 **Keep the first terminal open** to view logs and runtime behavior.
2. ➕ **Open a second terminal** and run:

```bash
docker compose down
```

This command stops the container and triggers a clean application shutdown.

---

## 📈 Get Streamer Video Statistics

Endpoint:
`GET /streamer/{username}/statistics?N={number_of_videos}`

Query Parameters:

* `N`: (Required) Number of most recent videos to include in the statistics

Response:

```json
{
  "video_lengths_sum": 3600000000000,
  "most_viewed_video_title": "Title: Sample Video 1. View Count: 150",
  "view_count_sum": 300,
  "view_count_avg": 100,
  "avg_view_per_minute": 5
}
```

Error cases handled include:

* Missing or invalid `N` param
* No user data found
* Twitch API errors

---

package main

import (
    "context"
    "net/http"
    "os"
    "os/signal"
    "strconv"
    "sync"
    "time"

    "github.com/gin-gonic/gin"
)

// album represents data about a record album.
// Price is represented in integer cents to avoid floating-point issues.
type album struct {
    ID         string `json:"id"`
    Title      string `json:"title" binding:"required"`
    Artist     string `json:"artist" binding:"required"`
    PriceCents int64  `json:"price_cents" binding:"required,gte=0"`
}

// createAlbumRequest is the input payload for creating albums (no client-supplied ID).
type createAlbumRequest struct {
    Title      string `json:"title" binding:"required"`
    Artist     string `json:"artist" binding:"required"`
    PriceCents int64  `json:"price_cents" binding:"required,gte=0"`
}

// albumStore is a simple in-memory, concurrency-safe repository.
type albumStore struct {
    mu     sync.RWMutex
    albums []album
    nextID int64 // monotonically increasing numeric ID used as string
}

func newAlbumStore(seed []album) *albumStore {
    s := &albumStore{}
    var maxID int64
    for _, a := range seed {
        // Determine max existing numeric ID; if non-numeric, ignore.
        if n, err := strconv.ParseInt(a.ID, 10, 64); err == nil && n > maxID {
            maxID = n
        }
        s.albums = append(s.albums, a)
    }
    s.nextID = maxID
    return s
}

func (s *albumStore) List() []album {
    s.mu.RLock()
    defer s.mu.RUnlock()
    out := make([]album, len(s.albums))
    copy(out, s.albums)
    return out
}

func (s *albumStore) GetByID(id string) (album, bool) {
    s.mu.RLock()
    defer s.mu.RUnlock()
    for _, a := range s.albums {
        if a.ID == id {
            return a, true
        }
    }
    return album{}, false
}

func (s *albumStore) Create(in createAlbumRequest) album {
    s.mu.Lock()
    defer s.mu.Unlock()
    s.nextID++
    a := album{
        ID:         strconv.FormatInt(s.nextID, 10),
        Title:      in.Title,
        Artist:     in.Artist,
        PriceCents: in.PriceCents,
    }
    s.albums = append(s.albums, a)
    return a
}

// seed data using cents
var seedAlbums = []album{
    {ID: "1", Title: "Blue Train", Artist: "John Coltrane", PriceCents: 5699},
    {ID: "2", Title: "Jeru", Artist: "Gerry Mulligan", PriceCents: 1799},
    {ID: "3", Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", PriceCents: 3999},
}

var store = newAlbumStore(seedAlbums)

// getAlbums responds with the list of all albums as JSON.
func getAlbums(c *gin.Context) {
    c.JSON(http.StatusOK, store.List())
}

// getAlbumByID responds with a single album by ID.
func getAlbumByID(c *gin.Context) {
    id := c.Param("id")
    if a, ok := store.GetByID(id); ok {
        c.JSON(http.StatusOK, a)
        return
    }
    c.JSON(http.StatusNotFound, gin.H{"error": "album not found"})
}

// postAlbums adds an album from JSON received in the request body.
func postAlbums(c *gin.Context) {
    var req createAlbumRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    created := store.Create(req)
    c.JSON(http.StatusCreated, created)
}

// healthz is a simple liveness probe.
func healthz(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"status": "ok"}) }

// limitBodyBytes limits request body size for the handler chain.
func limitBodyBytes(n int64) gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, n)
        c.Next()
    }
}

func main() {
    router := gin.Default()

    // Routes
    router.GET("/healthz", healthz)
    router.GET("/albums", getAlbums)
    router.GET("/albums/:id", getAlbumByID)
    router.POST("/albums", limitBodyBytes(1<<20), postAlbums) // 1 MiB limit

    // Server with graceful shutdown
    addr := ":8080"
    if p := os.Getenv("PORT"); p != "" {
        if p[0] == ':' {
            addr = p
        } else {
            addr = ":" + p
        }
    }

    srv := &http.Server{Addr: addr, Handler: router}

    go func() {
        // Start server
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            // Gin's default logger already logs; in a real app, log this error.
        }
    }()

    // Wait for interrupt signal to gracefully shutdown the server with a timeout.
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, os.Interrupt)
    <-quit

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    _ = srv.Shutdown(ctx)
}

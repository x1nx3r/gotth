# Case Study: IMPHISE

*A collaborative whiteboard built from this boilerplate.*

**Live: [canvas.x1nx3r.dev](https://canvas.x1nx3r.dev)** · **Source: [github.com/x1nx3r/ingin-menjadi-programmer-handal-namun-enggan-subscribe-excalidraw](https://github.com/x1nx3r/ingin-menjadi-programmer-handal-namun-enggan-subscribe-excalidraw)**

So you read the README. You saw the landing page with the chart and the wisdom quote. And you're thinking: *"Okay, cool demo. But where's the real stuff? Auth? Database? Real-time? How do I build an actual app with this?"*

Fair question. Here's the answer.

IMPHISE is a real-time collaborative whiteboard. Google sign-in. SQLite persistence. WebSocket broadcasting when someone else watches you draw. Admin dashboard. ~20 routes, ~15 components, ~8 MB of someone else's React code for the canvas. All running on a $5 VPS behind Caddy and Cloudflare. ~16 MB RAM at rest. One binary to deploy.

It started from this boilerplate. Here's how each piece works.

## Auth

Firebase Admin SDK verifies Google ID tokens server-side. The session is an http-only cookie set for 14 days with `SameSite=Strict`. No JWT on the client. No localStorage tokens that XSS can steal.

```go
func LoginHandler(w http.ResponseWriter, r *http.Request) {
    idToken := r.FormValue("id_token")
    token, err := auth.VerifyIDToken(idToken)
    if err != nil {
        http.Error(w, "invalid token", http.StatusUnauthorized)
        return
    }
    session := generateSessionToken()
    storeSession(session, token.UID)
    http.SetCookie(w, &http.Cookie{
        Name:     "session",
        Value:    session,
        HttpOnly: true,
        SameSite: http.SameSiteStrictMode,
        Secure:   r.TLS != nil,
        MaxAge:   3600 * 24 * 14,
    })
    w.Header().Set("HX-Redirect", "/drawings")
}
```

A middleware checks the cookie on every protected route, queries the user from SQLite, and attaches it to the request context. ~200 lines.

## Persistence

Started with Firebase Firestore. Migrated to SQLite when we needed joins and transactions.

```go
db, err := sql.Open("sqlite3", "./canvas.db")
db.SetMaxOpenConns(1)
db.Exec("PRAGMA journal_mode=WAL")
db.Exec("PRAGMA busy_timeout=5000")
```

The WAL file grows forever unless checkpointed. Auto-checkpoint threshold is 1000 pages. By the time you notice, your VPS disk is full. A background goroutine handles it:

```go
func StartWALCheckpoint(ctx context.Context, db *sql.DB) {
    ticker := time.NewTicker(5 * time.Minute)
    for {
        select {
        case <-ctx.Done():
            return
        case <-ticker.C:
            db.Exec("PRAGMA wal_checkpoint(PASSIVE)")
        }
    }
}
```

The database lives outside the release tree via a symlink. Deployments never touch it.

## Saving (Solo Mode)

Most users draw alone. They don't need a WebSocket. The browser uses a dirty-bit loop:

```javascript
let dirty = false
canvas.on('onChange', () => { dirty = true })
setInterval(() => {
    if (!dirty) return
    dirty = false
    fetch('/api/draw/' + id + '/save', { method: 'POST', body: JSON.stringify(scene) })
}, 3000)
window.addEventListener('beforeunload', () => {
    navigator.sendBeacon('/api/draw/' + id + '/save', blob)
})
```

1,500 concurrent users in the load test. Zero errors. 15 MB Go heap. The bottleneck was Linux deciding to OOM us at 128 MB of RAM, not anything in our code.

## Real-Time Collaboration

When someone opens a shared link, the server pushes an SSE event to the owner's page. The owner reluctantly opens a WebSocket. Both sides exchange full scene snapshots. When the last guest leaves, the WebSocket closes and everyone goes back to HTTP.

The hub is straightforward:

```go
type Room struct {
    clients map[*Client]bool
    mu      sync.RWMutex
}

func (r *Room) Broadcast(msg []byte) {
    r.mu.RLock()
    defer r.mu.RUnlock()
    for client := range r.clients {
        select {
        case client.send <- msg:
        default:
            delete(r.clients, client)
            close(client.send)
        }
    }
}
```

No message broker. No CRDT library. No event bus. ~400 lines.

### Wire Protocol

| Type | Direction | What it does |
|---|---|---|
| `SCENE_INIT` | Server → You | Here's the current scene |
| `SCENE_UPDATE` | Both | I changed something |
| `MOUSE_LOCATION` | Both | I'm over here |
| `COLLAB_ENDED` | Server → You | You're alone again |

## Admin Dashboard

Same Templ + HTMX pattern as everything else. A middleware checks the email against `SUPER_ADMIN_EMAIL`:

```go
func RequireSuperAdmin(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        user := GetUser(r.Context())
        if user.Email != os.Getenv("SUPER_ADMIN_EMAIL") {
            http.Error(w, "forbidden", http.StatusForbidden)
            return
        }
        next.ServeHTTP(w, r)
    })
}
```

Shows active WebSocket sessions, database size, user list. No JavaScript.

## Project Structure

The boilerplate's file layout grew but didn't change:

```
main.go                        # ~20 routes, still one file
app/
  lib/
    auth.go                    # Firebase token verification
    auth_handlers.go           # Login/logout
    db.go                      # SQLite + WAL checkpoint
    middleware.go              # Session, rate limit, admin gate
  api/
    draw.go                    # CRUD handlers
    shared.go                  # Public drawing access
    hub.go                     # WebSocket rooms
    ws.go                      # WS upgrade + pumps
  canvas/
    page.go + page.templ       # Canvas with Excalidraw
    shared.go + shared.templ   # Read-only shared view
  dashboard/
    page.go + page.templ       # Drawing grid
  admin/
    page.go + page.templ       # Admin panel
  components/                  # Reusable Templ components
  assets/
    excalidraw/                # esbuild entry point
    public/                    # Static files
    assets.go                  # //go:embed
```

Everything under `app/` follows the same patterns as the boilerplate. More files, same rules.

## What to Watch For

**Tailwind v4 doesn't scan `.templ` files.** Responsive classes silently disappear from production CSS. The fix is a Go tool that scans `.templ` files and injects classes into an `@source` inline declaration for Tailwind.

**CGO is required for SQLite.** The boilerplate compiles without it. You'll need `CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build`.

**Cloudflare caches everything.** Content-hashed URLs with `immutable` cache headers. Purge on deploy. We learned this the hard way when CSS was stuck for 3 days.

**SCENE_DELTA was a mistake.** We built a delta protocol to save bandwidth. It caused phantom diffs and race conditions. Deleted it. Full SCENE_UPDATE works fine because the lazy socket means 99% of users never open one.

## Deployment

```bash
make css && make templ
CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o bin/server main.go
rsync bin/server user@host:/srv/app/releases/$(date +%s)/
ssh user@host "ln -nfs current/ releases/$(date +%s)/ && systemctl restart app"
```

Binary is ~61 MB (8 MB of that is Excalidraw). Cold start ~500ms. Memory at rest: ~16 MB. Under 50 WebSocket users: ~96 MB. Death at ~128 MB (systemd `MemoryMax`). It's happened twice. Nobody noticed.

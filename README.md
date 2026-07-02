<p align="center">
  <img src="public/goth_lady_logo.jpg" width="120" height="120" alt="GOTTH Logo" style="border: 2px solid #1a1a24;" />
</p>

<h1 align="center">GOTTH Stack Boilerplate</h1>

This is a [Go](https://golang.org/), [Templ](https://templ.guide/), [Tailwind CSS v4](https://tailwindcss.com/), and [HTMX](https://htmx.org/) project bootstrapped with zero `node_modules` and zero respect for the modern JavaScript build tax.

We call it the **GOTTH** stack because it stands for Go, Templ, Tailwind, and HTMX, and because Next.js is a beautiful, bloated lie.

## Getting Started

First, run the development server:

```bash
npm run dev
# or
make dev
```

*(Note: `package.json` exists solely as a compatibility proxy so your developer muscle memory doesn't trigger a panic attack when you can't run `npm run dev`. It routes commands to a clean Linux Makefile, preventing you from downloading a 400MB headless Chrome binary just to watch a single text string hot-reload).*

Open [http://localhost:3000](http://localhost:3000) with your browser to see the result. 

The proxy server running on port `3000` intercepts your HTML responses and injects a low-overhead WebSocket reload script. The actual Go application executes natively on port `3001`.

---

## Architectural Rationale (Born of Pure Developer Trauma)

Every technical design choice in the GOTTH stack was forged in the fires of operational PTSD. Below is the documentation of our choices and the trauma that inspired them.

### 1. In-Place Compilation & Watcher Loop Prevention
* **The Trauma**: Infinite hot-reload rebuild loops. When `templ generate` runs, it outputs `_templ.go` files. If a file watcher is set to monitor all `.go` files, compiling a template modifies a Go file, which triggers the watcher, which then triggers another compilation loop—maxing out the CPU. To bypass this, developers often use convoluted directory syncing (copying files to separate build folders), which breaks serverless hosting (like Vercel) because their Go compiler only compiles at the root of the repository.
* **The Solution**: We compile templates directly in-place at the root. We prevent the infinite watcher loop by adding a regex pattern to `.air.toml` (`exclude_regex = ["_test\\.go$", ".*_templ\\.go$"]`). Because files are generated in-place, the project remains a standard Go module, allowing zero-configuration cloud builds.

### 2. In-Memory Assets & Docs (Go `//go:embed`)
* **The Trauma**: You deploy a Go binary to a VPS or Vercel. You configure the server, launch it, and test it. Everything works—until you hit `/docs/getting-started` or load the CSS stylesheet, and the server crashes with a `panic: open docs/getting-started.md: no such file or directory` because the relative directory structure in the container was mapped differently.
* **The Solution**: We eliminated the file system entirely. Using Go's standard library `embed` package, we package `assets.CSS` (the compiled Tailwind output) and `assets.Docs` (the documentation markdown files) directly into the compiled Go binary. 
  - The binary is **100% self-contained**.
  - Serves assets directly from RAM (zero disk I/O, zero file descriptor leaks, zero path traversal vulnerabilities).
  - Works out-of-the-box in stateless, read-only serverless lambdas.

### 3. Standalone Tailwind v4 Binary Fallbacks
* **The Trauma**: You add Tailwind to a Node project. Next week, you run `npm install` and your build breaks because `postcss` had a minor version bump, `autoprefixer` is throwing deprecation warnings, and `node-gyp` is failing to compile C++ bindings on your machine's CPU architecture.
* **The Solution**: We use Tailwind's official standalone binary (`./bin/tailwindcss`) downloaded directly to the workspace. No Node.js runtime required. To survive restricted CI/CD environments (like Vercel builders where curl downloads might be blocked or fail architecture checks), the `Makefile` automatically falls back to `npx @tailwindcss/cli` if the local binary is missing or non-executable.

### 4. Zero-Dependency Router Precedence
* **The Trauma**: Choosing a Go web framework is a cycle of despair. You start with one, only for the maintainers to rewrite the API in v2, deprecate your middleware libraries, and force you to rewrite your entire routing file just to upgrade Go versions.
* **The Solution**: We use Go 1.22's native `http.ServeMux`. Thanks to recent standard library upgrades, it supports HTTP method prefixing (`GET /docs`) and dynamic slug matching (`/docs/{slug}`) natively. We route static files with a single wildcard handler (`GET /{file}`), letting Go's standard routing precedence handle explicit pages first, and fall back to public files automatically. No external router dependencies.

### 5. HTMX View Transitions (Instead of JS State Engines)
* **The Trauma**: You want to animate a simple page switch. In JavaScript, this requires installing `framer-motion`, wrapping your layout in an `AnimatePresence` provider, setting up a React context store to track page state, and resolving console warnings about "hydration mismatches" because the server and client clocks differed by 1 millisecond.
* **The Summary**: We let the browser do the work. We use HTMX's `hx-boost="true"` to intercept standard links, and `hx-swap="outerHTML transition:true"` to trigger swaps. By declaring `view-transition-name` in `globals.css` for our containers, the browser natively morphs the box dimensions and cross-fades content on the GPU. The theme switch is wrapped in `document.startViewTransition()`, causing the entire viewport to cross-fade between light and dark modes in a single, hardware-accelerated pass.

---

## Deploy on Vercel

The easiest way to deploy your GOTTH app is to use the **Vercel Platform**.

Because Vercel supports native Go runtimes, you can connect this repository to your Vercel dashboard and deploy instantly. We have pre-configured `/vercel.json` to handle routing rewrites:

```json
{
  "version": 2,
  "routes": [
    { "src": "/(.*)", "dest": "/main.go" }
  ]
}
```

### Serverless Gotchas
* **Cold Starts**: Serverless functions scale to zero. When a request comes in after a period of silence, a container must boot up. Thankfully, Go starts in under 5ms, meaning Vercel's infrastructure startup latency is the only thing slowing your users down.
* **Locked Dependencies**: Go's compiler will throw errors if packages are imported but not used. Conversely, `go mod tidy` will aggressively delete dependencies it thinks are unused. Because `templ` templates generate code outside the root folder, running `go mod tidy` at the root will delete the `templ` dependency. We prevent this by anonymously importing it in `main.go` (`_ "github.com/a-h/templ"`), locking the dependency in place.

## Learn More

To learn how to escape the Node dependency trap:
- [GOTTH Documentation](http://localhost:3000/docs) — read the grimoire, run the commands, learn the runes.
- [HTMX Documentation](https://htmx.org/docs/) — understand how to build single-page app interactivity without loading a 100KB React runtime.
- [Templ Guide](https://templ.guide/) — type-safe markup compiler for Go.

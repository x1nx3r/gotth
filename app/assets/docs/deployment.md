# Summoning & Embeds

To keep the GOTTH stack serverless-compatible and dead-simple to deploy, all static assets and documentation markdown files are embedded directly into the compiled Go binary.

## Virtual Filesystem (Go Embed)

The assets are managed in the `app/assets` subpackage:
- `assets.CSS` (compiled Tailwind stylesheet)
- `assets.Docs` (read-only virtual filesystem of all documentation markdown files)

Because these assets are loaded directly from memory, the binary has zero file-system query overhead at runtime, making it highly secure and invulnerable to read-only container constraints.

## Deploying to Vercel

Vercel natively supports zero-configuration Go deployments. To deploy:

1. Run `make build` locally to generate your `*_templ.go` templates and compiled `app/assets/globals.css.output` stylesheet.
2. Commit all generated files to Git and push them.
3. Connect your Git repository to the Vercel dashboard (leave the "Build Command" blank).
4. Vercel automatically detects the root `main.go` and Go compiler version, compiles the self-contained binary in 1 second, and deploys it immediately.

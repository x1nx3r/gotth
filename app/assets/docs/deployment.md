# Summoning & Embeds

To keep the GOTTH stack serverless-compatible and dead-simple to deploy, all static assets and documentation markdown files are embedded directly into the compiled Go binary.

## Virtual Filesystem (Go Embed)

The assets are managed in the `app/assets` subpackage:
- `assets.CSS` (compiled Tailwind stylesheet)
- `assets.Docs` (read-only virtual filesystem of all documentation markdown files)

Because these assets are loaded directly from memory, the binary has zero file-system query overhead at runtime, making it highly secure and invulnerable to read-only container constraints.

## Deploying to Vercel

Vercel natively supports zero-configuration Go deployments. To deploy:

1. Connect your Git repository to the Vercel dashboard.
2. Vercel automatically detects the root `main.go` and `go.mod` file.
3. During build compilation, Go embeds the assets virtual directory into the binary.
4. Scale-to-zero serverless handlers serve your pages immediately.

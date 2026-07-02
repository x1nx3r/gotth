# The Incantation

The GOTTH stack demonstrates that you do not need Node.js or a heavy client-side Javascript framework runtime to build lightning fast, highly interactive web applications. 

The core ingredients of the stack are:
- `Go` — Core execution runtime
- `Templ` — Server-side template compiler
- `Tailwind CSS v4` — Standalone styling compiler
- `HTMX 2` — Declarative server interactivity

## Cloning the Grimoire

To summon the boilerplate onto your local machine, clone the repository and navigate into the project directory:

```bash
git clone https://github.com/x1nx3r/gotth.git
cd gotth
make setup
```

This downloads the standalone Tailwind v4 compiler executable, tidies Go modules, compiles initial Templ files, and gets your environment ready.

## Summoning the Dev Server

To launch the local developer environment, run:
```bash
make dev
```
This runs the Tailwind v4 build watch loop and starts the `Air` proxy server. Whenever you save a `.templ` or `.go` file, the workspace compiles in milliseconds.

## Compiling for Production

To build a standalone, zero-dependency production binary:
```bash
make build
```
The resulting compiled executable is stored in `bin/server`. Deploying is as simple as copying this binary to a virtual private server, binding your ports, and running it.

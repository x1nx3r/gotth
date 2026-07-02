# Runes of Rendering

Every component in GOTTH is a Server Component. Markup layout elements are declared using `Templ` and executed purely on the server.

## Building Components

Components are declared inside `.templ` files and compiled to type-safe Go functions during the build stage:
```templ
templ Page(title string, body string) {
	@Layout(title, "/docs") {
		<main>
			<h1>{ title }</h1>
			<div>@templ.Raw(body)</div>
		</main>
	}
}
```

## Hydration is a Lie

Because pages render completely on the server before transmitting to the client, there is no expensive hydration phase, no virtual DOM to maintain, and zero client-side framework runtime footprint. The browser receives optimized HTML.

# OOB Swaps

Out-of-band updates allow the server to target and update multiple disparate components in a single HTTP request, keeping the client interactive without single-page app overhead.

## Cast a Sigil (HTMX OOB Trigger)

By returning elements with `hx-swap-oob="true"` inside your response payload, HTMX will intercept the markup and swap it directly into matching DOM targets across the document, regardless of which element triggered the request.

For example, when a custom cart update is processed, the server returns the updated cart side-drawer as the main response, but appends a matching badge snippet:
```html
<span id="cart-badge" hx-swap-oob="true">
	3 items
</span>
```

HTMX will automatically swap the `#cart-badge` element, letting you orchestrate complex client interactions on the server.

package app

import (
	"math/rand"
	"net/http"
	"sync"
	"gotth/app/components"
)

type Quote struct {
	Text   string
	Author string
}

var (
	quotes = []Quote{
		{"Simplicity is the ultimate sophistication.", "Leonardo da Vinci"},
		{"The best code is no code at all.", "Jeff Atwood"},
		{"Make it work, make it right, make it fast.", "Kent Beck"},
		{"One of my most productive days was throwing away 1000 lines of code.", "Ken Thompson"},
		{"Complexity is the enemy of execution.", "Tony Robbins"},
		{"Measuring programming progress by lines of code is like measuring aircraft building progress by weight.", "Bill Gates"},
		{"Simple things should be simple, complex things should be possible.", "Alan Kay"},
		{"Programs must be written for people to read, and only incidentally for machines to execute.", "Harold Abelson"},
		{"Premature optimization is the root of all evil.", "Donald Knuth"},
		{"First, solve the problem. Then, write the code.", "John Johnson"},
	}
	lastQuoteIdx = -1
	quoteMutex   sync.Mutex
)

func PageHandler(w http.ResponseWriter, r *http.Request) {
	quoteMutex.Lock()
	idx := rand.Intn(len(quotes))
	lastQuoteIdx = idx
	q := quotes[idx]
	quoteMutex.Unlock()

	Page(q.Text, q.Author).Render(r.Context(), w)
}

func ChartHandler(w http.ResponseWriter, r *http.Request) {
	metricType := r.URL.Query().Get("type")
	components.Chart(metricType).Render(r.Context(), w)
}

func WisdomHandler(w http.ResponseWriter, r *http.Request) {
	quoteMutex.Lock()
	idx := rand.Intn(len(quotes))
	for idx == lastQuoteIdx && len(quotes) > 1 {
		idx = rand.Intn(len(quotes))
	}
	lastQuoteIdx = idx
	q := quotes[idx]
	quoteMutex.Unlock()

	components.WisdomQuote(q.Text, q.Author).Render(r.Context(), w)
}

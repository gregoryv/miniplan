package miniplan

import (
	"regexp"
	"time"

	"github.com/google/uuid"
)

type Entry struct {
	UUID        uuid.UUID
	Title       string
	Description string
	Priority    int
	Done        bool `json:",omitempty"`

	JustCreated bool      `json:",omitempty"`
	RemovedOn   time.Time `json:",omitempty"`
}

func (me *Entry) Ref() string {
	v := me.UUID.String()
	return v[len(v)-5:]
}

// Tags returns all tags, words starting with # in the title and
// description.
func (e *Entry) Tags() []string {
	return append(
		retags.FindAllString(e.Title, -1),
		retags.FindAllString(e.Description, -1)...,
	)
}

var retags = regexp.MustCompile(`#\w+`)

type ByPriority []*Entry

func (b ByPriority) Len() int           { return len(b) }
func (b ByPriority) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }
func (b ByPriority) Less(i, j int) bool { return b[i].Priority > b[j].Priority }

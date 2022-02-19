package discordapi

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/kensh1ro/willie/config"
)

const POLL_INTERVAL time.Duration = config.POLL_INTERVAL * time.Second

var FileAttachment Attachment

type PollQueue struct {
	Q Queue
}

func New() *PollQueue {
	pq := &PollQueue{}
	return pq
}

func (pq *PollQueue) Run() {
	var last_id string
	for {
		r := rand.New(rand.NewSource(10))
		c := time.NewTicker(POLL_INTERVAL)
		select {
		case <-c.C:
			jitter := time.Duration(r.Int31n(5000)) * time.Millisecond
			time.Sleep(jitter)
			m := GetMessage(config.CHANNEL_ID)
			if m.Author.Bot || last_id == m.ID {
				continue
			}
			if len(m.Files) != 0 {
				FileAttachment = m.Files[0]
			}
			if m.Reference.ID != "" {
				FileAttachment = m.Reference.Files[0]
			}
			last_id = m.ID
			fmt.Printf("PUSH: %s\n", m.Content)
			pq.Q.Push(m.Content)
		}
	}
}

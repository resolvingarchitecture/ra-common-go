package util

import "container/list"

// Nonce tracks recently-seen ids and rejects duplicates. Ports ra.common.Nonce.
type Nonce struct {
	maxSize      int
	prunePercent int
	seen         map[string]*list.Element
	order        *list.List
}

func NewNonce(maxSize int, prunePercent int) *Nonce {
	if prunePercent < 0 {
		prunePercent = 0
	}
	if prunePercent > 100 {
		prunePercent = 100
	}
	return &Nonce{
		maxSize:      maxSize,
		prunePercent: prunePercent,
		seen:         make(map[string]*list.Element),
		order:        list.New(),
	}
}

func DefaultNonce() *Nonce { return NewNonce(1_000_000, 10) }

// ContinueOn registers id. Returns true if new, false if a replay.
func (n *Nonce) ContinueOn(id string) bool {
	n.prune()
	if _, exists := n.seen[id]; exists {
		return false
	}
	n.seen[id] = n.order.PushBack(id)
	return true
}

func (n *Nonce) Size() int { return n.order.Len() }

func (n *Nonce) prune() {
	if n.order.Len() <= n.maxSize {
		return
	}
	if n.prunePercent == 100 {
		n.seen = make(map[string]*list.Element)
		n.order.Init()
		return
	}
	toPrune := n.maxSize * n.prunePercent / 100
	for i := 0; i < toPrune && n.order.Len() > 0; i++ {
		front := n.order.Front()
		delete(n.seen, front.Value.(string))
		n.order.Remove(front)
	}
}

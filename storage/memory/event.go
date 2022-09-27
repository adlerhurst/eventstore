package memory

import (
	"github.com/adlerhurst/eventstore"
	"github.com/adlerhurst/eventstore/storage"
)

func (m *Memory) saveEvents(cmds ...eventstore.Command) (events []*eventstore.Event, err error) {
	m.seqMu.Lock()
	defer m.seqMu.Unlock()

	seq := m.currentSequence

	events = make([]*eventstore.Event, len(cmds))
	for i, cmd := range cmds {
		seq++
		events[i], err = cmdToEvent(cmd, seq)
		if err != nil {
			return nil, err
		}
	}

	m.currentSequence = seq
	m.events = append(m.events, events...)

	return events, nil
}

func cmdToEvent(cmd eventstore.Command, seq uint64) (*eventstore.Event, error) {
	payload, err := storage.PayloadToBytes(cmd.Payload())
	if err != nil {
		return nil, err
	}
	return &eventstore.Event{
		EditorService: cmd.EditorService(),
		EditorUser:    cmd.EditorUser(),
		Payload:       payload,
		Sequence:      seq,
		Subjects:      cmd.Subjects(),
	}, nil
}

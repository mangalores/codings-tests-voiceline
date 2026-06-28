package message

type TranscribeCommand struct {
	ID string
}

type ExtractCommand struct {
	ID string
}

type ExportCommand struct {
	ID string
}

type ChannelCommandPublisher[T any] struct {
	commands chan<- T
}

func NewChannelCommandPublisher[T any](commands chan<- T) *ChannelCommandPublisher[T] {
	return &ChannelCommandPublisher[T]{commands: commands}
}

func (p *ChannelCommandPublisher[T]) Publish(command T) error {
	if p == nil || p.commands == nil {
		return nil
	}

	p.commands <- command

	return nil
}

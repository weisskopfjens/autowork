package connection

type Communicator interface {
	Begin() error
	End() error
	Write(string) error
	Read() (string, error)
	IsConnected() bool
}

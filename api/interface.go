package api

type Push interface {
	Send(m *Message) (result bool, err error)
}

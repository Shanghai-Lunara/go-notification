package api

type Push interface {
	Send(workerId, pid, infoType int, token string) (result bool, err error)
}

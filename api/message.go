package api

type Message struct {
	WorkerId int
	Pid      int
	InfoType int
	Token    string
	Title    string
	Body     string
}

func NewMessage(workerId, pid, infoType int, Token string) *Message {
	titleMap := map[int]string{
		1: "t1",
		2: "t2",
		3: "t3",
		4: "t4",
		5: "t5",
		6: "t6",
	}
	bodyMap := map[int]string{
		1: "b1",
		2: "b2",
		3: "b3",
		4: "b4",
		5: "b5",
		6: "b6",
	}
	var (
		title, body string
	)
	if t, ok := titleMap[infoType]; ok {
		title = t
	} else {
		title = "default title"
	}
	if t, ok := bodyMap[infoType]; ok {
		body = t
	} else {
		body = "default body"
	}
	m := &Message{
		WorkerId: workerId,
		Pid:      pid,
		InfoType: infoType,
		Token:    Token,
		Title:    title,
		Body:     body,
	}
	return m
}

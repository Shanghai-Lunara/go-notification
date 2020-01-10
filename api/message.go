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
		1: "スパニクル",
		2: "スパニクル",
		3: "スパニクル",
		4: "スパニクル",
		5: "スパニクル",
		6: "スパニクル",
	}
	bodyMap := map[int]string{
		1: "体力が満タンになりました。さっそく冒険を続けましょう！",
		2: "これ以上冒険サイコロを持てません。いますぐパイッキと一緒に宝探しに行きましょう！",
		3: "隊長、国家探索隊が任務を達成しましたよ。次のご指示をお願いします。",
		4: "螺旋研究が完了しました。引き続き螺旋の力の奥義を探求しませんか？",
		5: "体力が満タンになりました。さっそく冒険を続けましょう！",
		6: "神秘なる予言契約が再び解放されました！あなたの運命を占ってみませんか？",
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

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
		1: "螺旋英雄譚",
		2: "螺旋英雄譚",
		3: "螺旋英雄譚",
		4: "螺旋英雄譚",
		5: "螺旋英雄譚",
		6: "螺旋英雄譚",
	}
	bodyMap := map[int]string{
		1: "體力恢復滿:您的體力回復滿嘍，快來加入我們的冒險吧！",
		2: "骰子恢復滿:裝不下更多冒險骰子啦，趕緊和派琪一起去尋寶吧！",
		3: "國家探索完成:國家探索隊已經完成使命，隊長請給予指示吧！",
		4: "螺旋研究完成:螺旋研究完成了，你還想繼續探索螺旋的奧秘嗎？",
		5: "限時活動開啟:XXXX活動已經限時開啟，隊長大顯身手的時候到了！",
		6: "預言契約刷新:神秘的預言契約再次開啟，試著占卜你的命運吧！",
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

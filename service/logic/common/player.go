package common

func (w *Worker) PullPlayerOne(pid int) (meet map[int]int, min int, err error) {
	var (
		info, del []string
	)
	meet = make(map[int]int, 0)
	if info, err = w.dao.GetSinglePlayerList(pid); err != nil {
		return meet, min, err
	}
	del, meet, min = w.Combine(info, pid)
	if len(del) > 0 {
		if err = w.dao.UpdateSinglePlayerList(pid, del); err != nil {
			return meet, min, err
		}
	}
	return meet, min, nil
}

func (w *Worker) UpdatePlayerListNodes(pid, min int) {
	w.mu.Lock()
	defer w.mu.Unlock()
	if p, ok := w.listNodes.Players[pid]; ok {
		if p.Player.Value != min {
			p.Player.Value = min
			w.listNodes.AppendOrModify(p.Player)
		}
	} else {
		p := &Player{
			Pid:   pid,
			Value: min,
		}
		w.listNodes.AppendOrModify(p)
	}
}

func (w *Worker) RefreshOne(pid int) (err error) {
	if _, min, err := w.PullPlayerOne(pid); err != nil {
		return err
	} else {
		w.UpdatePlayerListNodes(pid, min)
		return nil
	}
}

func (w *Worker) CheckOne(pid int) (err error) {
	if meet, _, err := w.PullPlayerOne(pid); err != nil {
		return err
	} else {
		if len(meet) > 0 {
			if cid, close, err := w.dao.GetPlayerSettings(pid); err != nil {
				return err
			} else {
				if cid == "" {
					return
				}
				for k, v := range meet {
					_ = v
					if err = w.ApiPost(pid, k, cid); err != nil {
						return nil
					}
				}
			}

		}
	}
}

func (w *Worker) ApiPost(pid, infoType int, cid string) (err error) {
	return nil
}

package pipe

import (
	"fmt"
	"log"
	"slices"
	"sort"
	"sync"
)

func execCmd(wg *sync.WaitGroup, c cmd, in, out chan interface{}) {
	defer wg.Done()
	c(in, out)
	if out != nil {
		close(out)
	}
}

func RunPipeline(cmds ...cmd) {
	wg := &sync.WaitGroup{}
	defer wg.Wait()

	pipes := make([]chan interface{}, 0, len(cmds)+1)
	pipes = append(pipes, make(chan interface{}))

	for cmdId, cmd := range cmds {
		wg.Add(1)
		pipes = append(pipes, make(chan interface{}))
		go execCmd(wg, cmd, pipes[cmdId], pipes[cmdId+1])
	}
}

func getUserIteration(wg *sync.WaitGroup, checkedEmailsMu *sync.Mutex, checkedEmails map[string]struct{}, Iemail interface{}, out chan interface{}) {
	defer wg.Done()

	email, casted := Iemail.(string)
	if !casted {
		log.Println("getUserIteration() - email cast: bad cast")
		return
	}

	user := GetUser(email)
	checkedEmailsMu.Lock()
	_, found := checkedEmails[user.Email]
	if !found {
		checkedEmails[user.Email] = struct{}{}
		out <- user
	}
	checkedEmailsMu.Unlock()
}

func SelectUsers(in, out chan interface{}) {
	wg := &sync.WaitGroup{}
	checkedEmails := map[string]struct{}{}
	checkedEmailsMu := &sync.Mutex{}
	defer wg.Wait()

	for Iemail := range in {
		wg.Add(1)
		go getUserIteration(wg, checkedEmailsMu, checkedEmails, Iemail, out)
	}
}

func prepareUsersBatch(Iusers []interface{}) ([]User, error) {
	usersBatch := make([]User, 0, cap(Iusers))
	for _, Iuser := range Iusers {
		user, casted := Iuser.(User)
		if casted {
			usersBatch = append(usersBatch, user)
		}
	}

	if len(usersBatch) == 0 {
		return nil, fmt.Errorf("prepareUserBatch() - user cast: bad cast for all data")
	}
	return usersBatch, nil
}

func getMessagesByUsers(wg *sync.WaitGroup, IusersBatch []interface{}, out chan interface{}) {
	defer wg.Done()

	users, err := prepareUsersBatch(IusersBatch)
	if err != nil {
		log.Printf("getMessagesByUsers() - %v", err)
		return
	}

	msgIDs, err := GetMessages(users...)
	if err != nil {
		log.Printf("getMessagesByUsers() - %v\n", err)
		return
	}

	for _, msgID := range msgIDs {
		out <- msgID
	}
}

func SelectMessages(in, out chan interface{}) {
	wg := &sync.WaitGroup{}
	IusersBatch := make([]interface{}, 0, GetMessagesMaxUsersBatch)
	defer wg.Wait()

	for Iuser := range in {
		IusersBatch = append(IusersBatch, Iuser)
		if len(IusersBatch) == GetMessagesMaxUsersBatch {
			wg.Add(1)
			go getMessagesByUsers(wg, slices.Clone(IusersBatch), out)
			IusersBatch = make([]interface{}, 0, GetMessagesMaxUsersBatch)
		}
	}

	if len(IusersBatch) > 0 {
		wg.Add(1)
		go getMessagesByUsers(wg, slices.Clone(IusersBatch), out)
	}
}

func checkSpamIteration(wg *sync.WaitGroup, rateLimitChan chan struct{}, ImsgID interface{}, out chan interface{}) {
	defer wg.Done()

	msgId, casted := ImsgID.(MsgID)
	if !casted {
		log.Println("checkSpamIteration() - message cast: bad cast")
		return
	}

	rateLimitChan <- struct{}{}
	isSpam, err := HasSpam(msgId)
	<-rateLimitChan
	if err != nil {
		log.Printf("checkSpamIteration() - %v\n", err)
		return
	}

	out <- MsgData{msgId, isSpam}
}

func CheckSpam(in, out chan interface{}) {
	wg := &sync.WaitGroup{}
	defer wg.Wait()

	rateLimitChan := make(chan struct{}, HasSpamMaxAsyncRequests)
	for ImsgID := range in {
		wg.Add(1)
		go checkSpamIteration(wg, rateLimitChan, ImsgID, out)
	}
}

func CombineResults(in, out chan interface{}) {
	stat := make([]MsgData, 0)
	for ImsgData := range in {
		msgData, casted := ImsgData.(MsgData)
		if !casted {
			log.Println("CombineResults() - message cast: bad cast")
			continue
		}
		stat = append(stat, msgData)
	}

	sort.Slice(stat, func(i, j int) bool {
		if stat[i].HasSpam == stat[j].HasSpam {
			return stat[i].ID < stat[j].ID
		}
		return stat[i].HasSpam
	})

	for _, msgData := range stat {
		out <- fmt.Sprint(msgData.HasSpam, msgData.ID)
	}
}

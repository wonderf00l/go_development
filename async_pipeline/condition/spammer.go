package main

// func to read root in, read data from io.Reader

func RunPipeline(cmds ...cmd) {

}

func SelectUsers(in, out chan interface{}) {
	// 	in - string
	// 	out - User
	for userEmail := range in {
		userEmail := (userEmail).(string)
		user := GetUser(userEmail)
		out <- user // проверка уникальности
	}

}

func SelectMessages(in, out chan interface{}) {
	// 	in - User
	// 	out - MsgID
	
}

func CheckSpam(in, out chan interface{}) {
	// in - MsgID
	// out - MsgData
}

func CombineResults(in, out chan interface{}) {
	// in - MsgData
	// out - string
}

// func to write to io.Writer
// will add this job as cmd in runPipeline

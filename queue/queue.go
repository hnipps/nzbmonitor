package queue

var queue chan string

func GetQueue() chan string {
	return queue
}

func init() {
	queue = make(chan string, 100)
}

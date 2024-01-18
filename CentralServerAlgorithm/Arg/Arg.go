package Arg

type Arg struct {
	WaitIp        string
	WaitPort      string
	RequireToken  bool
	RealeaseToken bool
}
type Reply struct {
	GetToken     bool
	ReleaseToken bool
}

type Queue struct {
	Content []string
	tail    int
}

func NewQueue(size int) *Queue {
	return &Queue{
		Content: make([]string, size),
		tail:    0,
	}
}

func (q *Queue) Push(s string) int {
	q.Content = append(q.Content, s)
	q.tail++
	return q.tail
}

func (q *Queue) Pop() (s string) {
	i := 1
	s = q.Content[0]
	for i < q.tail {
		q.Content[i-1] = q.Content[i]
		i++
	}
	q.tail--
	q.Content = q.Content[:q.tail]
	return
}

func (q Queue) IsEmpty() bool {
	if q.tail == 0 {
		return true
	} else {
		return false
	}
}

package RingBasedAlgorithmElection

type Arg struct {
	Version  int
	Election bool
} //要大写，要不然不能夸包传递

type Reply struct {
	Replymsg int
}

package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

type State struct {
	CPUH     bool
	MemH     bool
	RpsH     bool
	RtimeH   bool
	R5xxH    bool
	Replicas int
}

type Action struct {
	Plus  float64
	Stay  float64
	Minus float64
}

type QLearn struct {
	Gamma   float32 //0.4
	Epsilon float32
	Rnd     *rand.Rand
	States  map[State]Action
	//	Action int -1,0,+1
	//StateList
	CurrentState State
}

// Choose and Fill Reward
// esp = epislon -- explore prob
func (q QLearn) ChooseAction() int {
	// Rand prob
	// If prob < eps then Explore
	// Else then choose maxRewardAction

	// Don,t allow to cause app to be 0 replicas
	action := 1
	return action
}

func (q QLearn) ValidAction(action int) bool {
	if q.CurrentState.Replicas+action <= 0 {
		return false
	}
	return true
}

// TODO: This fn may not need
func (q *QLearn) GoNextState(action int) {
	q.CurrentState.Replicas += action
	// TODO: What about cpu,mem,etc..
}

func (q QLearn) MaximumOp() float32 {
	// Find Best action then return it Q-Matrix

	action := q.States[q.CurrentState]

	return float32(math.Max(action.Plus, math.Max(action.Stay, action.Minus)))
}

func (q QLearn) Reward() {
	// R(current,action)+gamma*MaximumOp
}

func (q *QLearn) Init() {
	seed := rand.NewSource(time.Now().UnixNano())
	q.Rnd = rand.New(seed)

	q.CurrentState.Replicas = 1
	q.States = make(map[State]Action)
	q.States[q.CurrentState] = Action{Stay: 1}
}

// Change current state according to real cpu,mem,etc....
func (q *QLearn) SetCurrentState(cpu, mem, rps, rtime float32, r5xx, replicas int) {
	// TODO: Creat Fn to set all this value
	if cpu > 50 {
		q.CurrentState.CPUH = true
	}

	if mem > 50 {
		q.CurrentState.MemH = true
	}
	if rps > 150 {
		q.CurrentState.RpsH = true
	}
	if rtime > 5 {
		q.CurrentState.RtimeH = true
	}
	if r5xx > 10 {
		q.CurrentState.R5xxH = true
	}
	q.CurrentState.Replicas = replicas
}

func main() {
	agent := QLearn{Gamma: 0.4, Epsilon: 0.6}
	agent.Init()
	agent.States[State{Replicas: 2}] = Action{Plus: 1, Stay: 2, Minus: -1}
	fmt.Println(agent)
	agent.CurrentState = State{Replicas: 2}
	fmt.Println("MAX", agent.MaximumOp())
}

package learn

import (
	//	"bufio"
	//	"bytes"
	"encoding/json"
	"fmt"
	//"io"
	"io/ioutil"
	"math"
	"math/rand"
	"os"
	"time"
)

type State struct {
	CPUH     bool `json:"cpu"`
	MemH     bool `json:"mem"`
	RpsH     bool `json:"rps"`
	RtimeH   bool `json:"rtime"`
	R5xxH    bool `json:"r5xx"`
	Replicas int  `json:"replicas"`
}

type Action struct {
	Plus  float64 `json:"plus"`
	Stay  float64 `json:"stay"`
	Minus float64 `json:"minus"`
}

type QLearn struct {
	Gamma   float64 `json:"gamma"` //0.4
	Epsilon float64 `json:"epsilon"`
	rnd     *rand.Rand
	States  map[string]Action `json:"states"`
	//	Action int -1,0,+1
	//StateList
	CurrentState State
}

// If prob < eps then Explore
// Else then choose maxRewardAction
func (q QLearn) ChooseAction() int {
	var nextAction int
	prob := q.rnd.Float64()
	fmt.Println("Prob ", prob, "? Epsilon", q.Epsilon)
	if prob < q.Epsilon {
		fmt.Println("Go Explore")
		nextAction = q.rnd.Intn(3) - 1 // -1,0,1
	} else {
		fmt.Println("Go Best")
		action := q.States[toKey(q.CurrentState)]
		maxR := math.Max(action.Plus, math.Max(action.Stay, action.Minus))
		if action.Minus == maxR {
			nextAction = -1
		} else if action.Plus == maxR {
			nextAction = +1
		} else { // Stay
			nextAction = 0
		}
	}
	fmt.Println("Action :", nextAction)
	return nextAction
}

/*
func (q QLearn) ValidAction(action int) bool {
	// TODO: May be 0 is good when no one using it
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
*/
func (q QLearn) MaximumOp(state State) float64 {
	// Find Best action then return it Q-Matrix
	action := q.States[toKey(state)]
	max := math.Max(action.Plus, math.Max(action.Stay, action.Minus))
	fmt.Println("Max ", action)
	fmt.Println("Max=", max)
	return max
}

func (q *QLearn) Reward(state State, action int, nowStatus map[string]float64) float64 {
	reward := 0.0
	// ACTION
	if action == 1 {
		reward -= 10
	} else if action == -1 {
		reward += 10
	} else {
		reward += 0
	}
	// Replicas - More replicas more penalty
	reward += 1 - nowStatus["replicas"]
	// Replicas 1-1=0
	if nowStatus["replicas"] == 1 {
		reward -= 100
	}

	// RTime
	reward += 5 - nowStatus["rtime"]

	// 5XX
	reward -= nowStatus["r5xx"]

	// R(current,action)+gamma*MaximumOp
	reward += q.Gamma * q.MaximumOp(state)
	fmt.Println("Reward ", reward)
	// Update Q-Matrix
	qAction := q.States[toKey(state)]
	if action == 0 {
		qAction.Plus += reward
		q.States[toKey(state)] = qAction
	} else if action == -1 {
		qAction.Minus += reward
		q.States[toKey(state)] = qAction
	} else {
		qAction.Stay += reward
		q.States[toKey(state)] = qAction
	}
	return reward
}

func (q *QLearn) Init() {
	seed := rand.NewSource(time.Now().UnixNano())
	q.rnd = rand.New(seed)

	//	q.CurrentState.Replicas = 1
	q.States = make(map[string]Action)
	//q.States[q.CurrentState] = Action{Stay: 1}
}

// Change current state according to real cpu,mem,etc....
func (q *QLearn) SetCurrentState(cpu, mem, rps, rtime float64, r5xx, replicas int) {
	// TODO: Create Fn to set all this value
	if cpu > 30 {
		q.CurrentState.CPUH = true
	} else {
		q.CurrentState.CPUH = false
	}

	if mem > 50 {
		q.CurrentState.MemH = true
	} else {
		q.CurrentState.MemH = false
	}

	if rps > 50 {
		q.CurrentState.RpsH = true
	} else {
		q.CurrentState.RpsH = false
	}

	if rtime > 5 {
		q.CurrentState.RtimeH = true
	} else {
		q.CurrentState.RtimeH = false
	}

	if r5xx > 10 {
		q.CurrentState.R5xxH = true
	} else {
		q.CurrentState.R5xxH = false
	}
	q.CurrentState.Replicas = replicas
	if _, have := q.States[toKey(q.CurrentState)]; have == false {
		q.States[toKey(q.CurrentState)] = Action{}
	}
}

func (q QLearn) Save(path string) {
	jsonData, _ := json.Marshal(q)
	file, err := os.Create(path)
	if err != nil {
		//		return err
	}
	defer file.Close()
	file.WriteString(string(jsonData))
	file.Sync()
}

func (q *QLearn) Load(path string) error {

	seed := rand.NewSource(time.Now().UnixNano())
	q.rnd = rand.New(seed)

	dat, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	//fmt.Print(string(dat))
	err = json.Unmarshal(dat, q)
	if err != nil {
		return err
	}
	fmt.Println("Load success")
	return nil

}

func toKey(c State) string {
	return fmt.Sprint(c.CPUH, c.MemH, c.RpsH, c.RtimeH, c.R5xxH, c.Replicas)
}

/*
func main() {
	agent := QLearn{Gamma: 0.4, Epsilon: 0.6}
	agent.Init()
	agent.States[State{Replicas: 2}] = Action{Plus: 1, Stay: 2, Minus: -1}
	fmt.Println(agent)
	agent.CurrentState = State{Replicas: 2}
	fmt.Println("MAX", agent.MaximumOp())
}
*/

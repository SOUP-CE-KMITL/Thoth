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
	Gamma   float32 `json:"gamma"` //0.4
	Epsilon float32 `json:"epsilon"`
	rnd     *rand.Rand
	States  map[string]Action `json:"states"`
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

	action := q.States[toKey(q.CurrentState)]

	return float32(math.Max(action.Plus, math.Max(action.Stay, action.Minus)))
}

func (q QLearn) Reward() {
	// R(current,action)+gamma*MaximumOp
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

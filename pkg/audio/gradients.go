package audio

import (
	"encoding/json"
	"sort"
	"sync"
)

var empty = Gradient{}

type Gradients struct {
	m    sync.Mutex
	data map[string]Gradient
}

func NewGradients() *Gradients {
	return &Gradients{
		m:    sync.Mutex{},
		data: make(map[string]Gradient),
	}
}

func (s *Gradients) Add(name string, g Gradient) {
	s.m.Lock()
	defer s.m.Unlock()

	s.data[name] = g
}

func (s *Gradients) Has(name string) bool {
	s.m.Lock()
	defer s.m.Unlock()

	return len(s.data[name]) != 0
}

func (s *Gradients) Get(name string) Gradient {
	s.m.Lock()
	defer s.m.Unlock()

	return s.data[name]
}

func (s *Gradients) Clear() {
	s.m.Lock()
	defer s.m.Unlock()

	s.data = make(map[string]Gradient)
}

func (s *Gradients) IsEmpty() bool {
	return s.Size() == 0
}

func (s *Gradients) List() []string {
	list := make([]string, 0, len(s.data))

	for name := range s.data {
		list = append(list, name)
	}
	sort.Strings(list)

	return list
}

func (s *Gradients) Delete(name string) {
	delete(s.data, name)
}

func (s *Gradients) Size() int {
	return len(s.data)
}

func (s *Gradients) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.data)
}

func (s *Gradients) UnmarshalJSON(data []byte) error {
	var gradients map[string]Gradient
	if err := json.Unmarshal(data, &gradients); err != nil {
		return err
	}
	s.data = gradients
	return nil
}

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package dag

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuild(t *testing.T) {
	tcs := []struct {
		name     string
		tasks    []Task
		deps     map[string][]string
		expected bool
	}{
		{
			name:     "none-task",
			tasks:    []Task{},
			deps:     map[string][]string{},
			expected: true,
		},
		{
			name:     "new-task1",
			tasks:    []Task{{"task1"}},
			deps:     map[string][]string{},
			expected: true,
		},
		{
			name:     "duplicate-task1",
			tasks:    []Task{{"task1"}},
			deps:     map[string][]string{},
			expected: false,
		},
		{
			name:     "new-task2",
			tasks:    []Task{{"task2"}},
			deps:     map[string][]string{},
			expected: true,
		},
		{
			name:     "new-task3",
			tasks:    []Task{{"task3"}},
			deps:     map[string][]string{},
			expected: true,
		},
		{
			name:     "cycle-task3",
			tasks:    []Task{},
			deps:     map[string][]string{"task2": {"task1", "task3"}, "task3": {"task1", "task2"}},
			expected: false,
		},
	}

	d := &dag{
		cfg:   DefaultConfig(),
		graph: &Graph{Nodes: map[string]*Node{}},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			err := d.Build(tc.tasks, tc.deps)
			if err == nil {
				assert.Equal(t, tc.expected, true)
			} else {
				assert.Equal(t, tc.expected, false)
			}
		})
	}
}

func TestGet(t *testing.T) {
	_ = &dag{
		cfg:   DefaultConfig(),
		graph: &Graph{Nodes: map[string]*Node{}},
	}

	// TODO

	assert.Equal(t, nil, nil)
}

func TestAddTask(t *testing.T) {
	d := &dag{
		cfg:   DefaultConfig(),
		graph: &Graph{Nodes: map[string]*Node{}},
	}

	_, err := d.addTask(Task{"task"})
	assert.Equal(t, nil, err)

	_, err = d.addTask(Task{"task"})
	assert.NotEqual(t, nil, err)
}

func TestAddLink(t *testing.T) {
	d := &dag{
		cfg:   DefaultConfig(),
		graph: &Graph{Nodes: map[string]*Node{}},
	}

	_, err := d.addTask(Task{"task1"})
	assert.Equal(t, nil, err)

	_, err = d.addTask(Task{"task2"})
	assert.Equal(t, nil, err)

	err = d.addLink("task2", "task1")
	assert.Equal(t, nil, err)
}

func TestLookForNode(t *testing.T) {
	node2 := &Node{
		Task: Task{"task2"},
	}

	node3 := &Node{
		Task: Task{"task3"},
		Prev: []*Node{node2},
	}

	d := &dag{
		cfg:   DefaultConfig(),
		graph: &Graph{Nodes: map[string]*Node{}},
	}

	err := d.lookForNode("task1", []*Node{node2})
	assert.Equal(t, nil, err)

	err = d.lookForNode("task2", []*Node{node3})
	assert.NotEqual(t, nil, err)
}

func TestFindSchedulable(t *testing.T) {
	_ = &dag{
		cfg:   DefaultConfig(),
		graph: &Graph{Nodes: map[string]*Node{}},
	}

	// TODO

	assert.Equal(t, nil, nil)
}

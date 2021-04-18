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
	"context"
	"fmt"

	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/util/sets"

	"github.com/craftslab/pipeflow/pipe/list"
)

type Dag interface {
	Build([]Task, map[string][]string) error
	Get(...string) (sets.String, error)
}

type dag struct {
	cfg   *Config
	graph *Graph
}

type Config struct {
}

type Graph struct {
	Nodes map[string]*Node
}

type Node struct {
	Task Task
	Prev []*Node
	Next []*Node
}

type Task struct {
	HashKey string
}

func New(_ context.Context, cfg *Config) Dag {
	return &dag{
		cfg:   cfg,
		graph: &Graph{Nodes: map[string]*Node{}},
	}
}

func DefaultConfig() *Config {
	return &Config{}
}

func (d *dag) Build(tasks []Task, deps map[string][]string) error {
	for _, t := range tasks {
		if _, err := d.addTask(t); err != nil {
			return fmt.Errorf("failed to add task %s: %w", t.HashKey, err)
		}
	}

	for task, dep := range deps {
		for _, prev := range dep {
			if err := d.addLink(task, prev, d.graph.Nodes); err != nil {
				return fmt.Errorf("failed to add link between %s and %s: %w", task, prev, err)
			}
		}
	}

	return nil
}

func (d *dag) Get(task ...string) (sets.String, error) {
	helper := func() []*Node {
		var node []*Node
		for _, n := range d.graph.Nodes {
			if len(n.Prev) == 0 {
				node = append(node, n)
			}
		}
		return node
	}

	roots := helper()

	buf := sets.NewString()
	visited := sets.NewString()

	ts := sets.NewString(task...)
	for _, r := range roots {
		for _, t := range d.findSchedulable(r, visited, ts) {
			buf.Insert(t.HashKey)
		}
	}

	var names []string
	for v := range visited {
		names = append(names, v)
	}

	if len(list.DiffLeft(task, names)) > 0 {
		return nil, fmt.Errorf("invalid list of done tasks")
	}

	return buf, nil
}

func (d *dag) addTask(t Task) (*Node, error) {
	if _, ok := d.graph.Nodes[t.HashKey]; ok {
		return nil, errors.New("task duplicate")
	}

	n := &Node{
		Task: t,
	}

	d.graph.Nodes[t.HashKey] = n

	return n, nil
}

func (d *dag) addLink(task, prev string, nodes map[string]*Node) error {
	helper := func(task, prev *Node) error {
		if task.Task.HashKey == prev.Task.HashKey {
			return fmt.Errorf("cycle detected: %q", task.Task.HashKey)
		}
		if err := d.lookForNode(task.Task.HashKey, prev.Prev); err != nil {
			return fmt.Errorf("cycle detected: %w", err)
		}
		task.Prev = append(task.Prev, prev)
		prev.Next = append(prev.Next, task)
		return nil
	}

	p, ok := nodes[prev]
	if !ok {
		return fmt.Errorf("%s not present in nodes", prev)
	}

	t := nodes[task]
	if err := helper(t, p); err != nil {
		return fmt.Errorf("failed to create link from %s to %s: %w", p.Task.HashKey, t.Task.HashKey, err)
	}

	return nil
}

func (d *dag) lookForNode(task string, prevs []*Node) error {
	for _, p := range prevs {
		if p.Task.HashKey == task {
			return errors.New("task found in prevs")
		}

		if err := d.lookForNode(task, p.Prev); err != nil {
			return err
		}
	}

	return nil
}

func (d *dag) findSchedulable(node *Node, visited, task sets.String) []Task {
	helper := func(task sets.String, prevs []*Node) bool {
		if len(prevs) == 0 {
			return true
		}
		var collected []string
		for _, n := range prevs {
			if task.Has(n.Task.HashKey) {
				collected = append(collected, n.Task.HashKey)
			}
		}
		return len(collected) == len(prevs)
	}

	if visited.Has(node.Task.HashKey) {
		return []Task{}
	}

	visited.Insert(node.Task.HashKey)

	if task.Has(node.Task.HashKey) {
		var buf []Task
		for _, next := range node.Next {
			if _, ok := visited[next.Task.HashKey]; !ok {
				buf = append(buf, d.findSchedulable(next, visited, task)...)
			}
		}
		return buf
	}

	if helper(task, node.Prev) {
		return []Task{node.Task}
	}

	return []Task{}
}

/*
	Copyright 2022 Loophole Labs

	Licensed under the Apache License, Version 2.0 (the "License");
	you may not use this file except in compliance with the License.
	You may obtain a copy of the License at

		   http://www.apache.org/licenses/LICENSE-2.0

	Unless required by applicable law or agreed to in writing, software
	distributed under the License is distributed on an "AS IS" BASIS,
	WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
	See the License for the specific language governing permissions and
	limitations under the License.
*/

package interval

type node struct {
	min   uint64
	max   uint64
	left  *node
	right *node
}

type Interval struct {
	root *node
}

func New() *Interval {
	return new(Interval)
}

func (i *Interval) Insert(x uint64, y uint64) {
	x, y = sort(x, y)
	i.root = insert(x, y, i.root)
}

func (i *Interval) Contains(x uint64, y uint64) bool {
	x, y = sort(x, y)
	return intersection(x, y, i.root) == y-x+1
}

func sort(x, y uint64) (uint64, uint64) {
	if y > x {
		return x, y
	}
	return y, x
}

func insert(x uint64, y uint64, d *node) *node {
	if d == nil {
		return &node{
			min:   x,
			max:   y,
			left:  nil,
			right: nil}
	}
	switch {
	case x >= d.min && y <= d.max:
		return d
	case y < d.min:
		if y+1 == d.min {
			return left(x, d.max, d.left, d.right)
		}
		return &node{
			min:   d.min,
			max:   d.max,
			left:  insert(x, y, d.left),
			right: d.right,
		}

	case x > d.max:
		if x == d.max+1 {
			return right(d.min, y, d.left, d.right)
		}
		return &node{
			min:   d.min,
			max:   d.max,
			left:  d.left,
			right: insert(x, y, d.right),
		}
	case x < d.min && y <= d.max:
		return left(x, d.max, d.left, d.right)
	case x >= d.min && y > d.max:
		return right(d.min, y, d.left, d.right)
	case x < d.min && y > d.max:
		newNode := left(x, d.max, d.left, d.right)
		return right(newNode.min, y, newNode.left, newNode.right)
	}
	return d
}

func left(min, max uint64, left, right *node) *node {
	if left != nil {
		newX, newY, newLeft := maxNode(left.min, left.max, left.left, left.right)
		if newY+1 == min {
			return &node{newX, max, newLeft, right}
		}
	}
	return &node{min, max, left, right}
}

func right(min, max uint64, left, right *node) *node {
	if right != nil {
		newX, newY, newRight := minNode(right.min, right.max, right.left, right.right)
		if max+1 == newX {
			return &node{min, newY, left, newRight}
		}
	}
	return &node{min, max, left, right}
}

func maxNode(min uint64, max uint64, left *node, right *node) (uint64, uint64, *node) {
	if right == nil {
		return min, max, left
	}
	u, v, newRight := maxNode(right.min, right.max, right.left, right.right)
	return u, v, &node{min, max, left, newRight}
}

func minNode(min uint64, max uint64, left *node, right *node) (uint64, uint64, *node) {
	if left == nil {
		return min, max, right
	}
	u, v, newLeft := minNode(left.min, left.max, left.left, left.right)
	return u, v, &node{min, max, newLeft, right}
}

func intersection(left uint64, right uint64, node *node) uint64 {
	if node == nil {
		return 0
	}
	if left > node.max {
		if node.right == nil {
			return 0
		}
		return intersection(left, right, node.right)
	}
	if right < node.min {
		if node.left == nil {
			return 0
		}
		return intersection(left, right, node.left)
	}
	if left >= node.min {
		if right <= node.max {
			return right - left + 1
		}
		intersect := node.max - left + 1
		if node.right != nil {
			intersect += intersection(node.max+1, right, node.right)
		}
		return intersect
	}
	if right <= node.max {
		intersect := right - node.min + 1
		if node.left != nil {
			intersect += intersection(left, node.min-1, node.left)
		}
		return intersect
	}
	if left <= node.min && right >= node.max {
		intersect := node.max - node.min + 1
		if node.left != nil {
			intersect += intersection(left, node.min-1, node.left)
		}
		if node.right != nil {
			intersect += intersection(node.max+1, right, node.right)
		}
		return intersect
	}
	return 0
}

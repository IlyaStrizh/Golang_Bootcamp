package ex00

import "sync"

type TreeNode struct {
	HasToy bool
	Left   *TreeNode
	Right  *TreeNode
}

func NewTreeNode(toy bool) *TreeNode {
	return &TreeNode{
		HasToy: toy,
		Left:   nil,
		Right:  nil,
	}
}

func AreToysBalanced(root *TreeNode) bool {
	if root == nil {
		return false
	}

	var left, right uint
	wg := new(sync.WaitGroup)
	wg.Add(2)
	go func() {
		defer wg.Done()
		left = countToys(root.Left)
	}()
	go func() {
		defer wg.Done()
		right = countToys(root.Right)
	}()
	wg.Wait()

	return left == right
}

func countToys(node *TreeNode) uint {
	var result uint = 0

	if node != nil {
		if node.HasToy {
			result++
		}
		result += countToys(node.Left) + countToys(node.Right)
	}

	return result
}

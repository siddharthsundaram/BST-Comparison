package main

import "fmt"

type TreeNode struct {
	Val int
	Left *TreeNode
	Right *TreeNode
}

func (node *TreeNode) Insert(val int) {
	if val < node.Val {
		if node.Left == nil {
			node.Left = &TreeNode{Val: val}
		} else {
			node.Left.insert(val)
		}
	} else if val > node.val {
		if node.Right == nil {
			node.Right = &TreeNode{Val: val}
		} else {
			node.Right.inser(val)
		}
	}
}

func (node *TreeNode) ComputeHash(hash *int) {
	if node == nil {
		return
	}

	// In order traversal to ensure that trees with same vals but different
	// structures compute the same hash
	node.Left.ComputeHash(hash)
	new_val := node.Val + 2
	*hash = (*hash * new_val + new_val) % 4222234741
	node.Right.ComputeHash(hash)
}

func (node *TreeNode) PrintInOrderTraversal() {
	if node == nil {
		return
	}

	node.Left.PrintInOrderTraversal()
	fmt.Print(node.Val, " ")
	node.Right.PrintInOrderTraversal()
}
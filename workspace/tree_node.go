package main

import "fmt"

type TreeNode struct {
	Val int
	Left *TreeNode
	Right *TreeNode
}

func NewTreeNode(val int) *TreeNode {
	return &TreeNode{Val: val}
}

func (node *TreeNode) Insert(val int) {
	if val < node.Val {
		if node.Left == nil {
			node.Left = &TreeNode{Val: val}
		} else {
			node.Left.Insert(val)
		}
	} else if val > node.Val {
		if node.Right == nil {
			node.Right = &TreeNode{Val: val}
		} else {
			node.Right.Insert(val)
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

func CompareTrees(idx1 int, idx2 int) {
	node := trees[idx1]
	other := trees[idx2]
	this_slice := []int{}
	other_slice := []int{}
	InOrderTraversal(node, &this_slice)
	InOrderTraversal(other, &other_slice)

	if len(this_slice) != len(other_slice) {
		return
	}

	for i := range this_slice {
		if this_slice[i] != other_slice[i] {
			return
		}
	}

	identical_trees[idx1][idx2] = true
	identical_trees[idx2][idx1] = true
}

func InOrderTraversal(node *TreeNode, res *[] int) {
	if node != nil {
		InOrderTraversal(node.Left, res)
		*res = append(*res, node.Val)
		InOrderTraversal(node.Right, res)
	}
}

func (node *TreeNode) PrintInOrderTraversal() {
	if node == nil {
		return
	}

	node.Left.PrintInOrderTraversal()
	fmt.Print(node.Val, " ")
	node.Right.PrintInOrderTraversal()
}
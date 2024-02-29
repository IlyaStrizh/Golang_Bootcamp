package tests

import (
	tree "Day05/ex00"
	decorating "Day05/ex01"
	heap "Day05/ex02"
	knapsack "Day05/ex03"
	"reflect"

	"testing"

	"github.com/fatih/color"
)

func TestTree(t *testing.T) {
	/*
			  0
		     / \
		    0   1
		   / \
		  0   1

	*/
	t.Run("AreToysBalanced_1", func(t *testing.T) {
		root := tree.NewTreeNode(false)
		root.Left = tree.NewTreeNode(false)
		root.Left.Left = tree.NewTreeNode(false)
		root.Left.Right = tree.NewTreeNode(true)
		root.Right = tree.NewTreeNode(true)
		res := tree.AreToysBalanced(root)

		if true != res {
			errorMsg := color.MagentaString(`AreToysBalanced: "%v", want: "%v"`, res, true)
			t.Error(errorMsg)
		}
	})
	/*
		    1
		   /  \
		  1     0
		 / \   / \
		1   0 1   1
	*/
	t.Run("AreToysBalanced_2", func(t *testing.T) {
		root := tree.NewTreeNode(true)
		root.Left = tree.NewTreeNode(true)
		root.Left.Left = tree.NewTreeNode(true)
		root.Left.Right = tree.NewTreeNode(false)
		root.Right = tree.NewTreeNode(false)
		root.Right.Left = tree.NewTreeNode(true)
		root.Right.Right = tree.NewTreeNode(true)
		res := tree.AreToysBalanced(root)

		if true != res {
			errorMsg := color.MagentaString(`AreToysBalanced: "%v", want: "%v"`, res, true)
			t.Error(errorMsg)
		}
	})
	/*
	   	 1
	   	/ \
	   1   0
	*/
	t.Run("AreToysBalanced_3", func(t *testing.T) {
		root := tree.NewTreeNode(true)
		root.Left = tree.NewTreeNode(true)
		root.Right = tree.NewTreeNode(false)
		res := tree.AreToysBalanced(root)

		if false != res {
			errorMsg := color.MagentaString(`AreToysBalanced: "%v", want: "%v"`, res, false)
			t.Error(errorMsg)
		}
	})
	/*
		  0
		 / \
		1   0
		 \   \
		  1   1
	*/
	t.Run("AreToysBalanced_4", func(t *testing.T) {
		root := tree.NewTreeNode(false)
		root.Left = tree.NewTreeNode(true)
		root.Left.Right = tree.NewTreeNode(true)
		root.Right = tree.NewTreeNode(false)
		root.Right.Right = tree.NewTreeNode(true)
		res := tree.AreToysBalanced(root)

		if false != res {
			errorMsg := color.MagentaString(`AreToysBalanced: "%v", want: "%v"`, res, false)
			t.Error(errorMsg)
		}
	})
}

func TestUnrollGarland(t *testing.T) {
	/*
	       1
	      /  \
	     1     0
	    / \   / \
	   1   0 1   1

	*/
	t.Run("UnrollGarland_1", func(t *testing.T) {
		root := tree.NewTreeNode(true)
		root.Left = tree.NewTreeNode(true)
		root.Right = tree.NewTreeNode(false)
		root.Left.Left = tree.NewTreeNode(true)
		root.Left.Right = tree.NewTreeNode(false)
		root.Right.Left = tree.NewTreeNode(true)
		root.Right.Right = tree.NewTreeNode(true)
		res := decorating.UnrollGarland(root)
		Res := []bool{true, true, false, true, true, false, true}

		if !reflect.DeepEqual(Res, res) {
			errorMsg := color.MagentaString(`UnrollGarland: "%v", want: "%v"`, res, Res)
			t.Error(errorMsg)
		}
	})
}

func TestGetNCoolestPresents(t *testing.T) {
	/*
		(5, 1)
		(4, 5)
		(3, 1)
		(5, 2)
	*/
	presents := make([]heap.Present, 4)
	presents[0].Value = 5
	presents[0].Size = 1
	presents[1].Value = 4
	presents[1].Size = 5
	presents[2].Value = 3
	presents[2].Size = 1
	presents[3].Value = 5
	presents[3].Size = 2

	var res []heap.Present
	var err error

	t.Run("GetNCoolestPresents_1", func(t *testing.T) {
		res, err = heap.GetNCoolestPresents(presents, 3)
		Res := []heap.Present{
			{Value: 5, Size: 1},
			{Value: 5, Size: 2},
			{Value: 4, Size: 5},
		}

		if err != nil {
			t.Error(err)
		} else if !reflect.DeepEqual(res, Res) {
			errorMsg := color.MagentaString(`GetNCoolestPresents: "%v", want: "%v"`, res, Res)
			t.Error(errorMsg)
		}
	})

	t.Run("GetNCoolestPresents_2", func(t *testing.T) {
		res, err = heap.GetNCoolestPresents(presents, 4)
		Res := []heap.Present{
			{Value: 5, Size: 1},
			{Value: 5, Size: 2},
			{Value: 4, Size: 5},
			{Value: 3, Size: 1},
		}

		if err != nil {
			t.Error(err)
		} else if !reflect.DeepEqual(res, Res) {
			errorMsg := color.MagentaString(`GetNCoolestPresents: "%v", want: "%v"`, res, Res)
			t.Error(errorMsg)
		}
	})

}

func TestGrabPresents(t *testing.T) {
	/*
		(5, 1)
		(4, 5)
		(3, 1)
		(5, 2)
	*/
	presents := make([]heap.Present, 4)
	presents[0].Value = 5
	presents[0].Size = 1
	presents[1].Value = 4
	presents[1].Size = 5
	presents[2].Value = 3
	presents[2].Size = 1
	presents[3].Value = 5
	presents[3].Size = 2

	var res []heap.Present

	t.Run("GrabPresents_1", func(t *testing.T) {
		m := 3
		res = knapsack.GrabPresents(presents, m)
		Res := []heap.Present{
			{Value: 5, Size: 1},
			{Value: 5, Size: 2},
		}

		if !reflect.DeepEqual(res, Res) {
			errorMsg := color.MagentaString(`GrabPresents: "%v", want: "%v"`, res, Res)
			t.Error(errorMsg)
		}
	})

	t.Run("GrabPresents_2", func(t *testing.T) {
		m := 8
		res = knapsack.GrabPresents(presents, m)
		Res := []heap.Present{
			{Value: 5, Size: 1},
			{Value: 4, Size: 5},
			{Value: 5, Size: 2},
		}

		if !reflect.DeepEqual(res, Res) {
			errorMsg := color.MagentaString(`GrabPresents: "%v", want: "%v"`, res, Res)
			t.Error(errorMsg)
		}
	})
}

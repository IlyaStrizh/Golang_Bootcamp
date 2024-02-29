package ex01

import tree "Day05/ex00"

func UnrollGarland(root *tree.TreeNode) []bool {
	if root == nil {
		return []bool{}
	}

	var result []bool
	level := []*tree.TreeNode{root}
	isEvenLevel := true

	for len(level) > 0 {
		var nextLevel []*tree.TreeNode
		var layerResult []bool

		levelSize := len(level)
		for i := 0; i < levelSize; i++ {
			node := level[i]

			// Добавляем значение узла в текущий слой
			layerResult = append(layerResult, node.HasToy)

			if node.Left != nil {
				nextLevel = append(nextLevel, node.Left)
			}
			if node.Right != nil {
				nextLevel = append(nextLevel, node.Right)
			}
		}
		// Добавляем слой к результату, в зависимости от направления обхода слоев
		if isEvenLevel {
			// Добавляем слой в обратном порядке
			size := len(layerResult)
			for i := size - 1; i >= 0; i-- {
				result = append(result, layerResult[i])
			}
		} else {
			result = append(result, layerResult...)
		}

		level = nextLevel
		isEvenLevel = !isEvenLevel
	}

	return result
}

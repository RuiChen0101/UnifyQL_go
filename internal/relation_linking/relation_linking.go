package relation_linking

import (
	"errors"
	"fmt"

	"github.com/RuiChen0101/unfiyql/internal/expression_tree"
	"github.com/RuiChen0101/unfiyql/internal/relation_chain"
)

type RelationLinker struct {
	relationChain  *relation_chain.RelationChain
	resultTable    string
	expressionTree expression_tree.ExpressionTreeNode
}

func NewRelationLinker(relationChain *relation_chain.RelationChain, expressionTree expression_tree.ExpressionTreeNode) RelationLinker {
	return RelationLinker{
		relationChain:  relationChain,
		expressionTree: expressionTree,
	}
}

func (rl *RelationLinker) GetExpressionTree() expression_tree.ExpressionTreeNode {
	return rl.expressionTree
}

func (rl *RelationLinker) Link() (string, error) {
	resultTable, err := rl.linking()
	if err != nil {
		return "", err
	}
	rl.resultTable = resultTable
	if node, ok := rl.expressionTree.(*expression_tree.BinaryOperatorNode); ok {
		node.OutputTarget = resultTable
	}
	return rl.resultTable, nil
}

func (rl *RelationLinker) finalize(target string) error {
	if target == rl.resultTable {
		return nil
	}

	if rl.relationChain.IsParentOf(target, rl.resultTable) {
		path := rl.relationChain.FindRelationPath(rl.resultTable, target)
		for _, p := range *path {
			relationNode := expression_tree.RelationNode{
				FromTable: p.FromTable,
				FromField: p.FromField,
				ToTable:   p.ToTable,
				ToField:   p.ToField,
			}
			relationNode.SetLeftNode(rl.expressionTree)
			rl.expressionTree = &relationNode
		}
		return nil
	}
	return fmt.Errorf("RelationLinker: %s is not parent of %s", target, rl.resultTable)
}

func (rl *RelationLinker) linking() (string, error) {
	rootNode := rl.expressionTree
	if node, ok := rootNode.(*expression_tree.ConditionNode); ok {
		return node.OutputTarget, nil
	}
	if node, ok := rootNode.(*expression_tree.OutputTargetNode); ok {
		if node.GetLeftNode() == nil {
			return node.OutputTarget, nil
		}
		linker := NewRelationLinker(rl.relationChain, *node.GetLeftNode())
		linker.Link()
		if err := linker.finalize(node.OutputTarget); err != nil {
			return "", err
		}
		node.SetLeftNode(linker.GetExpressionTree())
		rl.expressionTree = node
		return node.OutputTarget, nil
	}

	if rootNode.GetLeftNode() == nil || rootNode.GetRightNode() == nil {
		return "", errors.New("RelationLinker: invalid expression tree")
	}

	leftLinker := NewRelationLinker(rl.relationChain, *rootNode.GetLeftNode())
	leftResultTable, err := leftLinker.Link()
	if err != nil {
		return "", err
	}
	rootNode.SetLeftNode(leftLinker.GetExpressionTree())

	rightLinker := NewRelationLinker(rl.relationChain, *rootNode.GetRightNode())
	rightResultTable, err := rightLinker.Link()
	if err != nil {
		return "", err
	}
	rootNode.SetRightNode(rightLinker.GetExpressionTree())

	if rightResultTable == leftResultTable {
		return leftResultTable, nil
	}

	commonParent, err := rl.relationChain.FindLowestCommonParent(leftResultTable, rightResultTable)
	if err != nil {
		return "", err
	}

	if leftResultTable != commonParent {
		path := rl.relationChain.FindRelationPath(leftResultTable, commonParent)
		for _, p := range *path {
			relationNode := expression_tree.RelationNode{
				FromTable: p.FromTable,
				FromField: p.FromField,
				ToTable:   p.ToTable,
				ToField:   p.ToField,
			}
			relationNode.SetLeftNode(*rootNode.GetLeftNode())
			rootNode.SetLeftNode(&relationNode)
		}
	}

	if rightResultTable != commonParent {
		path := rl.relationChain.FindRelationPath(rightResultTable, commonParent)
		for _, p := range *path {
			relationNode := expression_tree.RelationNode{
				FromTable: p.FromTable,
				FromField: p.FromField,
				ToTable:   p.ToTable,
				ToField:   p.ToField,
			}
			relationNode.SetLeftNode(*rootNode.GetRightNode())
			rootNode.SetRightNode(&relationNode)
		}
	}

	return commonParent, nil
}

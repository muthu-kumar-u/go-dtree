package dtree

import (
	"encoding/json"
	"errors"
	"fmt"
)

var (
	ErrorInvalidTreeData = errors.New("invalid tree data")
	ErrorInvalidBranchFormat = errors.New("invalid branch format")
	ErrorInvalidBranchData = errors.New("invalid branches data")
	ErrorPredictData =  errors.New("invalid predict value")
	ErrorMissingCondition = errors.New("missing condition data")
	ErrorConditionData = errors.New("invalid condition data")
	ErrorInvalidNodeName = errors.New("invalid node name")
)

type Tree struct {
	Root *Node `json:"root"`
}

type Condition struct {
    Type  string      `json:"type"`
    Value interface{} `json:"value"` 
}

type Node struct {
	Id int `json:"node_id"`
	Predict  string    `json:"predict"`
	Branches []*Branch `json:"branches"`
}

type Outcome struct {
	Result interface{} `json:"result"`
	NextNode *Node `json:"nextNode"`
}

type Branch struct {
	Condition *Condition `json:"condition"`
	Outcome   *Outcome   `json:"outcome"`
}

func CreateTreeFromjson(jsonByteData []byte) (*Tree, error) {
	tree := &Tree{}
	err := json.Unmarshal(jsonByteData, &tree)
	if err != nil {
		return nil, err
	}

	err = tree.InitializeTree()
	if err != nil {
		return nil, err
	}

	return tree, nil
}

func CreateNewTree() *Tree {
	return &Tree{}
}

func (t *Tree) InitializeTree() error {
    if t.Root == nil || t.Root.Branches == nil {
        return ErrorInvalidTreeData
    }

    nodes := []*Node{t.Root}

    for len(nodes) > 0 {
        currNode := nodes[0]
        nodes = nodes[1:]

        for _, branch := range currNode.Branches {
            if branch.Condition == nil || branch.Outcome == nil {
                return ErrorInvalidBranchFormat
            }

            if branch.Condition.Type == "" || branch.Condition.Value == nil {
                return ErrorConditionData
            }

            if branch.Outcome.NextNode != nil {
                nodes = append(nodes, branch.Outcome.NextNode)
            }
        }
    }
    return nil
}

func (t *Tree) Decide(data map[string]interface{}) (*Node, error) {
    if t.Root == nil {
        return nil, ErrorInvalidTreeData
    }

    currentNode := t.Root
    for currentNode != nil {
		if currentNode.Predict == data["predict"].(string) {
			nextNode, err := currentNode.Traverse(data)
			if err != nil {
				return nil, err
			}
	
			if nextNode == nil {
				return nil, nil
			} else {
				return nextNode, nil 
			}
		}
		fmt.Println("104")
    }

    return nil, ErrorInvalidNodeName
}

func (n *Node) Traverse(input map[string]interface{}) (*Node, error) {
    for _, branch := range n.Branches {
        conditionValue, ok := branch.Condition.Value.(bool)
        inputValue, inputOk := input["value"].(bool)

        if ok && inputOk && conditionValue == inputValue { 
            return branch.Outcome.NextNode, nil
        }
    }
    return nil, ErrorInvalidNodeName
}

func (n *Node) AddNewNodeToTree(node *Node) (*Node, error) {
	var branches []*Branch
	for _, branchData := range node.Branches {
		condition := &Condition{
			Type:  branchData.Condition.Type,
			Value: branchData.Condition.Value,
		}

		outcome := &Outcome{NextNode: nil} 
		if branchData.Outcome != nil {
			if branchData.Outcome.NextNode != nil {
				nextNode, err := n.AddOrCreateNode(branchData.Outcome.NextNode)
				if err != nil {
					return nil, err
				}
				outcome.NextNode = nextNode
			}
		}

		branch := &Branch{
			Condition: condition,
			Outcome:   outcome,
		}
		branches = append(branches, branch)
	}

	return &Node{
		Predict:  node.Predict,
		Branches: branches,
	}, nil
}

func (n *Node) AddOrCreateNode(node *Node) (*Node, error) {
	var branches []*Branch

	for _, branch := range node.Branches {
		if branch.Condition != nil && branch.Outcome != nil {
			branches = append(branches, branch)
		}
	}

	return &Node{
		Predict: node.Predict,
		Branches: branches,
	}, nil
}

func (b *Branch) AddOrCreateBranch(branch *Branch) (*Branch, error) {
	newBranch := &Branch{}

	if branch.Condition != nil {
		newBranch.Condition = branch.Condition 
	}

	if branch.Outcome != nil {
		newBranch.Outcome = branch.Outcome
	}

	return newBranch, nil
}

func (c *Condition) AddOrCreateCondition(condition *Condition) (*Condition, error){
	switch condition.Type {
	case "bool":
		return &Condition{Type: condition.Type, Value: condition.Value}, nil
	case "comparison":
		return &Condition{Type: condition.Type, Value: condition.Value}, nil	
	default:
		return nil, ErrorConditionData
	}
}

func (o *Outcome) AddOrCreateOutcome(outcome *Outcome) (*Outcome, error) {
	if outcome.NextNode != nil {
		return &Outcome{NextNode: outcome.NextNode}, nil
	}

	return nil, nil
}

func (t *Tree) ConvertTreeToJson() ([]byte, error) {
	return json.Marshal(t)
}

  
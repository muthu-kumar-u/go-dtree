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
        value, exists := data[currentNode.Predict]
        if !exists {
            return nil, fmt.Errorf("missing decision input for: %s", currentNode.Predict)
        }

        nextNode, err := currentNode.Traverse(value)
        if err != nil {
            return nil, err
        }

        if nextNode == nil || len(nextNode.Branches) == 0 {
            return nextNode, nil
        }

        currentNode = nextNode
    }

    return nil, ErrorInvalidNodeName
}

func (n *Node) Traverse(input interface{}) (*Node, error) {
    for _, branch := range n.Branches {
        conditionValue, ok := branch.Condition.Value.(bool)
        inputValue, inputOk := input.(bool)

        if ok && inputOk && conditionValue == inputValue { 
            return branch.Outcome.NextNode, nil
        }
    }
    return nil, ErrorInvalidNodeName
}

func (n *Tree) AddNewNodeToTree(node *Node) (*Node, error) {
	var branches []*Branch
	for _, branchData := range node.Branches {
		condition := &Condition{
			Type:  branchData.Condition.Type,
			Value: branchData.Condition.Value,
		}

		outcome := &Outcome{NextNode: branchData.Outcome.NextNode} 
		if branchData.Outcome != nil {
			if branchData.Outcome.NextNode != nil {
				nextNode, err := AddOrCreateNode(branchData.Outcome.NextNode)
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

func AddOrCreateNode(node *Node) (*Node, error) {
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

func AddOrCreateBranch(branch []*Branch) ([]*Branch, error) {
	newBranch := []*Branch{}
	for i, branch := range branch {
	
		if branch.Condition != nil {
			newBranch[i].Condition = branch.Condition 
		}
	
		if branch.Outcome != nil {
			newBranch[i].Outcome = branch.Outcome
		}
	}

	return newBranch, nil
}

func AddOrCreateCondition(condition []*Condition) ([]*Condition, error){
 	newCondition := []*Condition{}
	for _, con := range condition {
		switch con.Type {
		case "bool":
			newCondition = append(newCondition, &Condition{Type: con.Type, Value: con.Value})
			continue
		case "comparison":
			newCondition = append(newCondition, &Condition{Type: con.Type, Value: con.Value})
			continue
		default:
			 return nil, ErrorConditionData
		}
	}

	return newCondition, nil
}

func AddOrCreateOutcome(outcome []*Outcome) ([]*Outcome, error) {
	newOutCome := []*Outcome{}
	for i, out := range outcome {
		if outcome[i].NextNode != nil {
			newOutCome = append(newOutCome, &Outcome{NextNode: out.NextNode})
		}
	}

	return newOutCome, nil
}

func (t *Tree) ConvertTreeToJson() ([]byte, error) {
	return json.Marshal(t)
}

  
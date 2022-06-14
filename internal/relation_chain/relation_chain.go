package relation_chain

import (
	"fmt"
	"reflect"
)

type RelationChainNode struct {
	FromTable string
	FromField string
	ToTable   string
	ToField   string
}

type relationChainMap map[string]map[string]RelationChainNode

type RelationChain struct {
	forwardRelationMap  relationChainMap
	backwardRelationMap relationChainMap
}

func NewRelationChain(forwardMap relationChainMap, backwardMap relationChainMap) RelationChain {
	return RelationChain{
		forwardRelationMap:  forwardMap,
		backwardRelationMap: backwardMap,
	}
}

func (rc *RelationChain) GetForwardRelationMap() relationChainMap {
	return rc.forwardRelationMap
}

func (rc *RelationChain) GetBackwardRelationMap() relationChainMap {
	return rc.backwardRelationMap
}

func (rc *RelationChain) FindLowestCommonParent(t1 string, t2 string) (string, error) {
	visitedTable := map[string]bool{
		t1: true,
	}
	for current, ok := rc.backwardRelationMap[t1]; ok; {
		parent := reflect.ValueOf(current).MapKeys()[0].String()
		visitedTable[parent] = true
		current, ok = rc.backwardRelationMap[parent]
	}
	if _, ok := visitedTable[t2]; ok {
		return t2, nil
	}
	for current, ok := rc.backwardRelationMap[t2]; ok; {
		parent := reflect.ValueOf(current).MapKeys()[0].String()
		if _, ex := visitedTable[parent]; ex {
			return parent, nil
		}
		current, ok = rc.backwardRelationMap[parent]
	}
	return "", fmt.Errorf("RelationChain: %s and %s had no common parent", t1, t2)
}

func (rc *RelationChain) FindRelationPath(from string, to string) *[]RelationChainNode {
	if rc.IsDescendantOf(from, to) {
		result := []RelationChainNode{}
		for current, ok := rc.backwardRelationMap[from]; ok; {
			parent := reflect.ValueOf(current).MapKeys()[0].String()
			result = append(result, current[parent])
			if parent == to {
				return &result
			}
			current, ok = rc.backwardRelationMap[parent]
		}
		return nil
	} else if rc.IsParentOf(from, to) {
		path := []string{to}
		for current, ok := rc.backwardRelationMap[to]; ok; {
			parent := reflect.ValueOf(current).MapKeys()[0].String()
			if parent == from {
				last := from
				result := []RelationChainNode{}
				for _, p := range path {
					result = append(result, rc.forwardRelationMap[last][p])
					last = p
				}
				return &result
			}
			path = append([]string{parent}, path...)
			current, ok = rc.backwardRelationMap[parent]
		}
		return nil
	}
	return nil
}

func (rc *RelationChain) IsParentOf(t1 string, t2 string) bool {
	for current, ok := rc.backwardRelationMap[t2]; ok; {
		parent := reflect.ValueOf(current).MapKeys()[0].String()
		if parent == t1 {
			return true
		}
		current, ok = rc.backwardRelationMap[parent]
	}
	return false
}

func (rc *RelationChain) IsDescendantOf(t1 string, t2 string) bool {
	return rc.IsParentOf(t2, t1)
}

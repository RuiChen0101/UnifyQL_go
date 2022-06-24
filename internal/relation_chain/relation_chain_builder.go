package relation_chain

import (
	"fmt"
	"reflect"
	"regexp"

	"github.com/RuiChen0101/UnifyQL_go/pkg/element"
)

func BuildRelationChain(el *element.UnifyQLElement) (*RelationChain, error) {
	completeRelationMap := relationChainMap{}

	definedTables := map[string]bool{
		el.QueryTarget: true,
	}

	for _, w := range el.With {
		definedTables[w] = true
	}

	for _, l := range el.Link {
		relation, err := extractRelation(l)
		if err != nil {
			return nil, err
		}
		if !definedTables[relation.FromTable] || !definedTables[relation.ToTable] {
			return nil, fmt.Errorf("RelationChain: %s using undefined table", l)
		}
		safeMapAssign(&completeRelationMap, relation.FromTable, relation.ToTable, *relation)
		safeMapAssign(&completeRelationMap, relation.ToTable, relation.FromTable, RelationChainNode{
			FromTable: relation.ToTable,
			FromField: relation.ToField,
			ToTable:   relation.FromTable,
			ToField:   relation.FromField,
		})
	}

	forwardMap := relationChainMap{}
	backwardMap := relationChainMap{}

	trackingTable := []string{el.QueryTarget}
	visitedTable := map[string]bool{}

	for len(trackingTable) != 0 {
		t := trackingTable[0]
		trackingTable = trackingTable[1:]
		r, ok := completeRelationMap[t]
		if !ok {
			visitedTable[t] = true
			continue
		}
		desc := reflect.ValueOf(r).MapKeys()
		for _, d := range desc {
			dStr := d.String()
			if visitedTable[dStr] {
				continue
			}
			safeMapAssign(&forwardMap, t, dStr, r[dStr])
			safeMapAssign(&backwardMap, dStr, t, completeRelationMap[dStr][t])
			trackingTable = append(trackingTable, dStr)
		}
		visitedTable[t] = true
	}

	result := NewRelationChain(forwardMap, backwardMap)
	return &result, nil
}

func safeMapAssign(m *relationChainMap, key1 string, key2 string, value RelationChainNode) {
	if _, ok := (*m)[key1]; !ok {
		(*m)[key1] = map[string]RelationChainNode{}
	}
	(*m)[key1][key2] = value
}

func extractRelation(relation string) (*RelationChainNode, error) {
	reg := regexp.MustCompile(`([^\s]+)\.([^\s]+)\s*=\s*([^\s]+)\.([^\s]+)`)
	captureGroup := reg.FindStringSubmatch(relation)
	if len(captureGroup) != 5 {
		return nil, fmt.Errorf("RelationChain: %s invalid format", relation)
	}
	return &RelationChainNode{
		FromTable: captureGroup[1],
		FromField: captureGroup[2],
		ToTable:   captureGroup[3],
		ToField:   captureGroup[4],
	}, nil
}

package relation_chain

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"

	"github.com/RuiChen0101/unfiyql/internal/element"
)

func BuildRelationChain(el element.UnifyQLElement) (*RelationChain, error) {
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
			return nil, errors.New(fmt.Sprintf("RelationChain: %s using undefined table", l))
		}
		completeRelationMap[relation.FromTable][relation.ToTable] = *relation
		completeRelationMap[relation.ToTable][relation.FromTable] = *relation
	}

	forwardMap := relationChainMap{}
	backwardMap := relationChainMap{}

	trackingTable := []string{el.QueryTarget}
	visitedTable := map[string]bool{}

	for len(trackingTable) != 0 {
		t, trackingTable := trackingTable[0], trackingTable[1:]
		r, ok := completeRelationMap[t]
		if !ok {
			visitedTable[t] = true
			continue
		}
		desc := reflect.ValueOf(r).MapKeys()
		newDesc := false
		for _, d := range desc {
			dStr := d.String()
			if visitedTable[dStr] {
				continue
			}
			newDesc = true
			forwardMap[t][dStr] = r[dStr]
		}
	}

	result := NewRelationChain(forwardMap, backwardMap)
	return &result, nil
}

func extractRelation(relation string) (*RelationChainNode, error) {
	reg := regexp.MustCompile(`([^\s]+)\.([^\s]+)\s*=\s*([^\s]+)\.([^\s]+)`)
	captureGroup := reg.FindStringSubmatch(relation)
	if len(captureGroup) != 4 {
		return nil, errors.New("RelationChain: invalid format")
	}
	return &RelationChainNode{
		FromTable: captureGroup[0],
		FromField: captureGroup[1],
		ToTable:   captureGroup[2],
		ToField:   captureGroup[3],
	}, nil
}

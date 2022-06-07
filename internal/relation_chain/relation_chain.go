package relation_chain

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

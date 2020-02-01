package impls

import (
	"fmt"

	ipld "github.com/ipld/go-ipld-prime/_rsrch/nodesolution"
)

// plainMap is a concrete type that provides a map-kind ipld.Node.
// It can contain any kind of value.
// plainMap is also embedded in the 'any' struct and usable from there.
type plainMap struct {
	m map[string]ipld.Node // string key -- even if a runtime schema wrapper is using us for storage, we must have a comparable type here, and string is all we know.
	t []plainMap__Entry    // table for fast iteration, order keeping, and yielding pointers to enable alloc/conv amortization.
}

type plainMap__Entry struct {
	k plainString // address of this used when we return keys as nodes, such as in iterators.  Need in one place to amortize shifts to heap when ptr'ing for iface.
	v ipld.Node   // identical to map values.  keeping them here simplifies iteration.  (in codegen'd maps, this position is also part of amortization, but in this implementation, that's less useful.)
	// note on alternate implementations: 'v' could also use the 'any' type, and thus amortize value allocations.  the memory size trade would be large however, so we don't, here.
}

// -- Node interface methods -->

func (plainMap) ReprKind() ipld.ReprKind {
	return ipld.ReprKind_Map
}
func (n *plainMap) LookupString(key string) (ipld.Node, error) {
	v, exists := n.m[key]
	if !exists {
		return nil, ipld.ErrNotExists{ipld.PathSegmentOfString(key)}
	}
	return v, nil
}
func (n *plainMap) Lookup(key ipld.Node) (ipld.Node, error) {
	ks, err := key.AsString()
	if err != nil {
		return nil, err
	}
	return n.LookupString(ks)
}
func (plainMap) LookupIndex(idx int) (ipld.Node, error) {
	return nil, ipld.ErrWrongKind{TypeName: "map", MethodName: "LookupIndex", AppropriateKind: ipld.ReprKindSet_JustList, ActualKind: ipld.ReprKind_Map}
}
func (n *plainMap) LookupSegment(seg ipld.PathSegment) (ipld.Node, error) {
	return n.LookupString(seg.String())
}
func (n *plainMap) MapIterator() ipld.MapIterator {
	return &plainMap_MapIterator{n, 0}
}
func (plainMap) ListIterator() ipld.ListIterator {
	panic("no")
}
func (n *plainMap) Length() int {
	return len(n.t)
}
func (plainMap) IsUndefined() bool {
	return false
}
func (plainMap) IsNull() bool {
	return false
}
func (plainMap) AsBool() (bool, error) {
	return false, ipld.ErrWrongKind{TypeName: "map", MethodName: "AsBool", AppropriateKind: ipld.ReprKindSet_JustBool, ActualKind: ipld.ReprKind_Map}
}
func (plainMap) AsInt() (int, error) {
	return 0, ipld.ErrWrongKind{TypeName: "map", MethodName: "AsInt", AppropriateKind: ipld.ReprKindSet_JustFloat, ActualKind: ipld.ReprKind_Map}
}
func (plainMap) AsFloat() (float64, error) {
	return 0, ipld.ErrWrongKind{TypeName: "map", MethodName: "AsFloat", AppropriateKind: ipld.ReprKindSet_JustFloat, ActualKind: ipld.ReprKind_Map}
}
func (plainMap) AsString() (string, error) {
	return "", ipld.ErrWrongKind{TypeName: "map", MethodName: "AsString", AppropriateKind: ipld.ReprKindSet_JustString, ActualKind: ipld.ReprKind_Map}
}
func (plainMap) AsBytes() ([]byte, error) {
	return nil, ipld.ErrWrongKind{TypeName: "map", MethodName: "AsBytes", AppropriateKind: ipld.ReprKindSet_JustBytes, ActualKind: ipld.ReprKind_Map}
}
func (plainMap) AsLink() (ipld.Link, error) {
	return nil, ipld.ErrWrongKind{TypeName: "map", MethodName: "AsLink", AppropriateKind: ipld.ReprKindSet_JustLink, ActualKind: ipld.ReprKind_Map}
}
func (plainMap) Style() ipld.NodeStyle {
	return Style__Map{}
}

type plainMap_MapIterator struct {
	n   *plainMap
	idx int
}

func (itr *plainMap_MapIterator) Next() (k ipld.Node, v ipld.Node, _ error) {
	if itr.Done() {
		return nil, nil, ipld.ErrIteratorOverread{}
	}
	k = &itr.n.t[itr.idx].k
	v = itr.n.t[itr.idx].v
	itr.idx++
	return
}
func (itr *plainMap_MapIterator) Done() bool {
	return itr.idx >= len(itr.n.t)
}

// -- NodeStyle -->

type Style__Map struct{}

func (Style__Map) NewBuilder() ipld.NodeBuilder {
	return &plainMap__Builder{plainMap__Assembler{w: &plainMap{}}}
}

// -- NodeBuilder -->

type plainMap__Builder struct {
	plainMap__Assembler
}

func (nb *plainMap__Builder) Build() (ipld.Node, error) {
	if nb.state != maState_done {
		panic("invalid state: assembler must be 'done' before Build can be called!")
	}
	return nb.w, nil
}
func (nb *plainMap__Builder) Reset() {
	*nb = plainMap__Builder{}
	nb.w = &plainMap{}
}

// -- NodeAssembler -->

type plainMap__Assembler struct {
	w *plainMap

	ka plainMap__KeyAssembler
	va plainMap__ValueAssembler

	state maState
}
type plainMap__KeyAssembler struct {
	ma *plainMap__Assembler
}
type plainMap__ValueAssembler struct {
	ma *plainMap__Assembler
}

// maState is an enum of the state machine for a map assembler.
// (this might be something to export reusably, but it's also very much an impl detail that need not be seen, so, dubious.)
type maState uint8

const (
	maState_initial     maState = iota // also the 'expect key or done' state
	maState_midKey                     // waiting for a 'done' state in the KeyAssembler.
	maState_expectValue                // 'AssembleValue' is the only valid next step
	maState_midValue                   // waiting for a 'done' state in the ValueAssembler.
	maState_done                       // 'w' will also be nil, but this is a politer statement
)

func (na *plainMap__Assembler) BeginMap(sizeHint int) (ipld.MapNodeAssembler, error) {
	// Allocate storage space.
	na.w.t = make([]plainMap__Entry, 0, sizeHint)
	na.w.m = make(map[string]ipld.Node, sizeHint)
	// Initialize the key and value assemblers with pointers back to the whole.
	na.ka.ma = na
	na.va.ma = na
	// That's it; return self as the MapNodeAssembler.  We already have all the right methods on this structure.
	return na, nil
}
func (plainMap__Assembler) BeginList(sizeHint int) (ipld.ListNodeAssembler, error) { panic("no") }
func (plainMap__Assembler) AssignNull() error                                      { panic("no") }
func (plainMap__Assembler) AssignBool(bool) error                                  { panic("no") }
func (plainMap__Assembler) AssignInt(int) error                                    { panic("no") }
func (plainMap__Assembler) AssignFloat(float64) error                              { panic("no") }
func (plainMap__Assembler) AssignString(v string) error                            { panic("no") }
func (plainMap__Assembler) AssignBytes([]byte) error                               { panic("no") }
func (na *plainMap__Assembler) Assign(v ipld.Node) error {
	// todo: apply a generic 'copy' function.
	// todo: probably can also shortcut to copying na.t and na.m if it's our same concrete type?
	//  (can't quite just `na.w = v`, because we don't have 'freeze' features, and we don't wanna open door to mutation of 'v'.)
	//   (wait... actually, probably we can?  'Assign' is a "done" method.  we can&should invalidate the wip pointer here.)
	panic("later")
}
func (plainMap__Assembler) Style() ipld.NodeStyle { panic("later") }

// -- MapNodeAssembler -->

// AssembleDirectly is part of conforming to MapAssembler, which we do on
// plainMap__Assembler so that BeginMap can just return a retyped pointer rather than new object.
func (ma *plainMap__Assembler) AssembleDirectly(k string) (ipld.NodeAssembler, error) {
	// Sanity check, then update, assembler state.
	if ma.state != maState_initial {
		panic("misuse")
	}
	ma.state = maState_midValue
	// Check for dup keys; error if so.
	_, exists := ma.w.m[k]
	if exists {
		return nil, ipld.ErrRepeatedMapKey{String(k)}
	}
	//l := len(ma.w.t)
	ma.w.t = append(ma.w.t, plainMap__Entry{k: plainString(k)})
	// configure and return an anyAssembler, similar to below in prepareAssigner
	panic("todo")
}

// AssembleKey is part of conforming to MapAssembler, which we do on
// plainMap__Assembler so that BeginMap can just return a retyped pointer rather than new object.
func (ma *plainMap__Assembler) AssembleKey() ipld.NodeAssembler {
	// Sanity check, then update, assembler state.
	if ma.state != maState_initial {
		panic("misuse")
	}
	ma.state = maState_midKey
	// Extend entry table.
	ma.w.t = append(ma.w.t, plainMap__Entry{})
	// No work to be done to init key assembler; it already points back to whole 'ma'; just yield it.
	return &ma.ka
}

// AssembleValue is part of conforming to MapAssembler, which we do on
// plainMap__Assembler so that BeginMap can just return a retyped pointer rather than new object.
func (ma *plainMap__Assembler) AssembleValue() ipld.NodeAssembler {
	// Sanity check, then update, assembler state.
	if ma.state != maState_expectValue {
		panic("misuse")
	}
	ma.state = maState_midValue
	// No work to be done to init value assembler; it already points back to whole 'ma'; just yield it.
	return &ma.va
}

// Done is part of conforming to MapAssembler, which we do on
// plainMap__Assembler so that BeginMap can just return a retyped pointer rather than new object.
func (ma *plainMap__Assembler) Done() error {
	// Sanity check, then update, assembler state.
	if ma.state != maState_initial {
		panic("misuse")
	}
	ma.state = maState_done
	// validators could run and report errors promptly, if this type had any.
	return nil
}
func (plainMap__Assembler) KeyStyle() ipld.NodeStyle   { panic("later") }
func (plainMap__Assembler) ValueStyle() ipld.NodeStyle { panic("later") }

// -- MapNodeAssembler.KeyAssembler -->

func (plainMap__KeyAssembler) BeginMap(sizeHint int) (ipld.MapNodeAssembler, error)   { panic("no") }
func (plainMap__KeyAssembler) BeginList(sizeHint int) (ipld.ListNodeAssembler, error) { panic("no") }
func (plainMap__KeyAssembler) AssignNull() error                                      { panic("no") }
func (plainMap__KeyAssembler) AssignBool(bool) error                                  { panic("no") }
func (plainMap__KeyAssembler) AssignInt(int) error                                    { panic("no") }
func (plainMap__KeyAssembler) AssignFloat(float64) error                              { panic("no") }
func (mka *plainMap__KeyAssembler) AssignString(v string) error {
	// Check for dup keys; error if so.
	_, exists := mka.ma.w.m[v]
	if exists {
		return ipld.ErrRepeatedMapKey{String(v)}
	}
	// Assign the key into the end of the entry table;
	//  we'll be doing map insertions after we get the value in hand.
	//  (There's no need to delegate to another assembler for the key type,
	//   because we're just at Data Model level here, which only regards plain strings.)
	mka.ma.w.t[len(mka.ma.w.t)-1].k = plainString(v)
	// Update parent assembler state: clear to proceed.
	mka.ma.state = maState_expectValue
	return nil
}
func (plainMap__KeyAssembler) AssignBytes([]byte) error { panic("no") }
func (mka *plainMap__KeyAssembler) Assign(v ipld.Node) error {
	vs, err := v.AsString()
	if err != nil {
		return fmt.Errorf("cannot assign non-string node into map key assembler") // FIXME:errors: this doesn't quite fit in ErrWrongKind cleanly; new error type?
	}
	return mka.AssignString(vs)
}
func (plainMap__KeyAssembler) Style() ipld.NodeStyle { panic("later") } // probably should give the style of plainString, which could say "only stores string kind" (though we haven't made such a feature part of the interface yet).

// -- MapNodeAssembler.ValueAssembler -->

func (mva *plainMap__ValueAssembler) BeginMap(sizeHint int) (ipld.MapNodeAssembler, error) {
	panic("todo") // now please
}
func (mva *plainMap__ValueAssembler) BeginList(sizeHint int) (ipld.ListNodeAssembler, error) {
	panic("todo") // now please
}
func (mva *plainMap__ValueAssembler) AssignNull() error     { panic("todo") }
func (mva *plainMap__ValueAssembler) AssignBool(bool) error { panic("todo") }
func (mva *plainMap__ValueAssembler) AssignInt(v int) error {
	l := len(mva.ma.w.t) - 1
	vb := plainInt(v)
	mva.ma.w.t[l].v = &vb
	mva.ma.w.m[string(mva.ma.w.t[l].k)] = &vb
	mva.ma.state = maState_initial
	return nil
}
func (mva *plainMap__ValueAssembler) AssignFloat(float64) error   { panic("todo") }
func (mva *plainMap__ValueAssembler) AssignString(v string) error { panic("todo") }
func (mva *plainMap__ValueAssembler) AssignBytes([]byte) error    { panic("todo") }
func (mva *plainMap__ValueAssembler) Assign(v ipld.Node) error {
	l := len(mva.ma.w.t) - 1
	mva.ma.w.t[l].v = v
	mva.ma.w.m[string(mva.ma.w.t[l].k)] = v
	mva.ma.state = maState_initial
	return nil
}
func (plainMap__ValueAssembler) Style() ipld.NodeStyle { panic("later") }

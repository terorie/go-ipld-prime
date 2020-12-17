package mixins

import (
	"io"

	ipld "github.com/ipld/go-ipld-prime"
)

type ListTraits struct {
	PkgName    string
	TypeName   string // see doc in kindTraitsGenerator
	TypeSymbol string // see doc in kindTraitsGenerator
}

func (ListTraits) Kind() ipld.Kind {
	return ipld.Kind_List
}
func (g ListTraits) EmitNodeMethodKind(w io.Writer) {
	doTemplate(`
		func ({{ .TypeSymbol }}) Kind() ipld.Kind {
			return ipld.Kind_List
		}
	`, w, g)
}
func (g ListTraits) EmitNodeMethodLookupByString(w io.Writer) {
	kindTraitsGenerator{g.PkgName, g.TypeName, g.TypeSymbol, ipld.Kind_List}.emitNodeMethodLookupByString(w)
}
func (g ListTraits) EmitNodeMethodLookupBySegment(w io.Writer) {
	doTemplate(`
		func (n {{ .TypeSymbol }}) LookupBySegment(seg ipld.PathSegment) (ipld.Node, error) {
			i, err := seg.Index()
			if err != nil {
				return nil, ipld.ErrInvalidSegmentForList{TypeName: "{{ .PkgName }}.{{ .TypeName }}", TroubleSegment: seg, Reason: err}
			}
			return n.LookupByIndex(i)
		}
	`, w, g)
}
func (g ListTraits) EmitNodeMethodMapIterator(w io.Writer) {
	kindTraitsGenerator{g.PkgName, g.TypeName, g.TypeSymbol, ipld.Kind_List}.emitNodeMethodMapIterator(w)
}
func (g ListTraits) EmitNodeMethodIsAbsent(w io.Writer) {
	kindTraitsGenerator{g.PkgName, g.TypeName, g.TypeSymbol, ipld.Kind_List}.emitNodeMethodIsAbsent(w)
}
func (g ListTraits) EmitNodeMethodIsNull(w io.Writer) {
	kindTraitsGenerator{g.PkgName, g.TypeName, g.TypeSymbol, ipld.Kind_List}.emitNodeMethodIsNull(w)
}
func (g ListTraits) EmitNodeMethodAsBool(w io.Writer) {
	kindTraitsGenerator{g.PkgName, g.TypeName, g.TypeSymbol, ipld.Kind_List}.emitNodeMethodAsBool(w)
}
func (g ListTraits) EmitNodeMethodAsInt(w io.Writer) {
	kindTraitsGenerator{g.PkgName, g.TypeName, g.TypeSymbol, ipld.Kind_List}.emitNodeMethodAsInt(w)
}
func (g ListTraits) EmitNodeMethodAsFloat(w io.Writer) {
	kindTraitsGenerator{g.PkgName, g.TypeName, g.TypeSymbol, ipld.Kind_List}.emitNodeMethodAsFloat(w)
}
func (g ListTraits) EmitNodeMethodAsString(w io.Writer) {
	kindTraitsGenerator{g.PkgName, g.TypeName, g.TypeSymbol, ipld.Kind_List}.emitNodeMethodAsString(w)
}
func (g ListTraits) EmitNodeMethodAsBytes(w io.Writer) {
	kindTraitsGenerator{g.PkgName, g.TypeName, g.TypeSymbol, ipld.Kind_List}.emitNodeMethodAsBytes(w)
}
func (g ListTraits) EmitNodeMethodAsLink(w io.Writer) {
	kindTraitsGenerator{g.PkgName, g.TypeName, g.TypeSymbol, ipld.Kind_List}.emitNodeMethodAsLink(w)
}

type ListAssemblerTraits struct {
	PkgName       string
	TypeName      string // see doc in kindAssemblerTraitsGenerator
	AppliedPrefix string // see doc in kindAssemblerTraitsGenerator
}

func (ListAssemblerTraits) Kind() ipld.Kind {
	return ipld.Kind_List
}
func (g ListAssemblerTraits) EmitNodeAssemblerMethodBeginMap(w io.Writer) {
	kindAssemblerTraitsGenerator{g.PkgName, g.TypeName, g.AppliedPrefix, ipld.Kind_List}.emitNodeAssemblerMethodBeginMap(w)
}
func (g ListAssemblerTraits) EmitNodeAssemblerMethodAssignNull(w io.Writer) {
	kindAssemblerTraitsGenerator{g.PkgName, g.TypeName, g.AppliedPrefix, ipld.Kind_List}.emitNodeAssemblerMethodAssignNull(w)
}
func (g ListAssemblerTraits) EmitNodeAssemblerMethodAssignBool(w io.Writer) {
	kindAssemblerTraitsGenerator{g.PkgName, g.TypeName, g.AppliedPrefix, ipld.Kind_List}.emitNodeAssemblerMethodAssignBool(w)
}
func (g ListAssemblerTraits) EmitNodeAssemblerMethodAssignInt(w io.Writer) {
	kindAssemblerTraitsGenerator{g.PkgName, g.TypeName, g.AppliedPrefix, ipld.Kind_List}.emitNodeAssemblerMethodAssignInt(w)
}
func (g ListAssemblerTraits) EmitNodeAssemblerMethodAssignFloat(w io.Writer) {
	kindAssemblerTraitsGenerator{g.PkgName, g.TypeName, g.AppliedPrefix, ipld.Kind_List}.emitNodeAssemblerMethodAssignFloat(w)
}
func (g ListAssemblerTraits) EmitNodeAssemblerMethodAssignString(w io.Writer) {
	kindAssemblerTraitsGenerator{g.PkgName, g.TypeName, g.AppliedPrefix, ipld.Kind_List}.emitNodeAssemblerMethodAssignString(w)
}
func (g ListAssemblerTraits) EmitNodeAssemblerMethodAssignBytes(w io.Writer) {
	kindAssemblerTraitsGenerator{g.PkgName, g.TypeName, g.AppliedPrefix, ipld.Kind_List}.emitNodeAssemblerMethodAssignBytes(w)
}
func (g ListAssemblerTraits) EmitNodeAssemblerMethodAssignLink(w io.Writer) {
	kindAssemblerTraitsGenerator{g.PkgName, g.TypeName, g.AppliedPrefix, ipld.Kind_List}.emitNodeAssemblerMethodAssignLink(w)
}
func (g ListAssemblerTraits) EmitNodeAssemblerMethodPrototype(w io.Writer) {
	kindAssemblerTraitsGenerator{g.PkgName, g.TypeName, g.AppliedPrefix, ipld.Kind_List}.emitNodeAssemblerMethodPrototype(w)
}

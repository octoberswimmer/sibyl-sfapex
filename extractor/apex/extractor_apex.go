package apex

import (
	"github.com/octoberswimmer/go-tree-sitter-sfapex/apex"
	"github.com/opensibyl/sibyl2/pkg/core"
)

const LangApex core.LangType = "APEX"

func init() {
	core.RegisterLang(LangApex, apex.GetLanguage(), ".cls")
}

const (
	KindApexScopeIdentifier      core.KindRepr = "scoped_identifier"
	KindApexIdentifier           core.KindRepr = "identifier"
	KindApexClassDeclaration     core.KindRepr = "class_declaration"
	KindApexClassBody            core.KindRepr = "class_body"
	KindApexFieldDeclaration     core.KindRepr = "field_declaration"
	KindApexEnumDeclaration      core.KindRepr = "enum_declaration"
	KindApexInterfaceDeclaration core.KindRepr = "interface_declaration"
	KindApexMethodDeclaration    core.KindRepr = "method_declaration"
	KindApexFormalParameters     core.KindRepr = "formal_parameters"
	KindApexFormalParameter      core.KindRepr = "formal_parameter"
	KindApexMethodInvocation     core.KindRepr = "method_invocation"
	KindApexModifiers            core.KindRepr = "modifiers"
	KindApexAnnotation           core.KindRepr = "annotation"
	KindApexMarkerAnnotation     core.KindRepr = "marker_annotation"
	KindApexBlock                core.KindRepr = "block"
	KindApexSuperClass           core.KindRepr = "superclass"
	KindApexSuperInterface       core.KindRepr = "super_interfaces"
	KindApexTypeList             core.KindRepr = "type_list"
	KindApexTypeIdentifier       core.KindRepr = "type_identifier"
	KindApexGenericType          core.KindRepr = "generic_type"
	FieldApexType                core.KindRepr = "type"
	FieldApexDimensions          core.KindRepr = "dimensions"
	FieldApexObject              core.KindRepr = "object"
	FieldApexName                core.KindRepr = "name"
	FieldApexArguments           core.KindRepr = "arguments"
	FieldApexDeclarator          core.KindRepr = "declarator"
)

type Extractor struct {
}

func (extractor *Extractor) GetLang() core.LangType {
	return LangApex
}

package apex

import (
	"errors"

	"github.com/opensibyl/sibyl2/pkg/core"
	"github.com/opensibyl/sibyl2/pkg/extractor/object"
)

// FunctionExtras ApexFunctionExtras
type FunctionExtras struct {
	Annotations []string   `json:"annotations"`
	ClassInfo   *ClassInfo `json:"classInfo"`
}

type ClassInfo struct {
	ClassName   string   `json:"className"`
	Annotations []string `json:"annotations"`
}

func (extractor *Extractor) IsFunction(unit *core.Unit) bool {
	// no function in apex
	if unit.Kind == KindApexMethodDeclaration {
		return true
	}
	return false
}

func (extractor *Extractor) ExtractFunctions(units []*core.Unit) ([]*object.Function, error) {
	var ret []*object.Function
	for _, eachUnit := range units {
		if !extractor.IsFunction(eachUnit) {
			continue
		}
		eachFunc, err := extractor.ExtractFunction(eachUnit)
		if err != nil {
			return nil, err
		}
		ret = append(ret, eachFunc)
	}
	return ret, nil
}

func (extractor *Extractor) ExtractFunction(unit *core.Unit) (*object.Function, error) {
	funcUnit := object.NewFunction()
	funcUnit.Span = unit.Span
	funcUnit.Unit = unit
	funcUnit.Lang = extractor.GetLang()

	// body scope
	funcBody := core.FindFirstByKindInSubsWithBfs(unit, KindApexBlock)
	if funcBody != nil {
		funcUnit.BodySpan = funcBody.Span
	}

	clazzName := ""

	// trace its class (the closest one
	clazzDecl := core.FindFirstByOneOfKindInParent(unit, KindApexClassDeclaration, KindApexEnumDeclaration, KindApexInterfaceDeclaration)
	clazzIdentifier := core.FindFirstByKindInSubsWithBfs(clazzDecl, KindApexIdentifier)
	if clazzIdentifier == nil {
		return nil, errors.New("no class found in " + unit.Content)
	}
	clazzName = clazzIdentifier.Content
	funcUnit.Receiver = clazzName

	funcIdentifier := core.FindFirstByKindInSubsWithBfs(unit, KindApexIdentifier)
	if funcIdentifier == nil {
		return nil, errors.New("no func id found in identifier" + unit.Content)
	}
	funcUnit.Name = funcIdentifier.Content

	// returns
	retUnit := core.FindFirstByFieldInSubsWithDfs(unit, FieldApexDimensions)
	valueUnit := &object.ValueUnit{
		Type: retUnit.Content,
		// apex has no named return value
		Name: "",
	}
	funcUnit.Returns = append(funcUnit.Returns, valueUnit)

	// params
	parameters := core.FindFirstByKindInSubsWithDfs(unit, KindApexFormalParameters)
	if parameters != nil {
		for _, each := range core.FindAllByKindInSubsWithDfs(parameters, KindApexFormalParameter) {
			typeName := core.FindFirstByFieldInSubsWithBfs(each, FieldApexType)
			paramName := core.FindFirstByFieldInSubsWithBfs(each, FieldApexDimensions)
			valueUnit = &object.ValueUnit{
				Type: typeName.Content,
				Name: paramName.Content,
			}
			funcUnit.Parameters = append(funcUnit.Parameters, valueUnit)
		}
	}

	// extras
	extras := &FunctionExtras{}
	classInfo := &ClassInfo{
		ClassName:   clazzName,
		Annotations: nil,
	}
	extras.ClassInfo = classInfo

	// class annotations
	classModifiers := core.FindFirstByKindInSubsWithBfs(clazzDecl, KindApexModifiers)
	if classModifiers != nil {
		classAnnotations := core.FindAllByKindsInSubs(classModifiers, KindApexMarkerAnnotation, KindApexAnnotation)
		if len(classAnnotations) != 0 {
			for _, each := range classAnnotations {
				classInfo.Annotations = append(classInfo.Annotations, each.Content)
			}
		}
	}
	// todo: inherit

	modifiers := core.FindFirstByKindInSubsWithBfs(unit, KindApexModifiers)
	if modifiers != nil {
		annotations := core.FindAllByKindsInSubs(modifiers, KindApexMarkerAnnotation, KindApexAnnotation)
		if len(annotations) != 0 {
			for _, each := range annotations {
				extras.Annotations = append(extras.Annotations, each.Content)
			}
		}
	}
	funcUnit.Extras = extras

	return funcUnit, nil
}

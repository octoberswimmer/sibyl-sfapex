package apex

import (
	"errors"
	"strings"

	"github.com/opensibyl/sibyl2/pkg/core"
	"github.com/opensibyl/sibyl2/pkg/extractor/object"
)

type ClassField struct {
	Name        string   `json:"name"`
	Type        string   `json:"type"`
	Annotations []string `json:"annotations"`
	Modifiers   []string `json:"modifiers"`
}

type ClassExtras struct {
	Annotations []string      `json:"annotations"`
	Fields      []*ClassField `json:"fields"`
	Modifiers   []string      `json:"modifiers"`
	Extends     string        `json:"extends"`
	Implements  []string      `json:"implements"`
}

func (extractor *Extractor) IsClass(unit *core.Unit) bool {
	if unit.Kind == KindApexClassDeclaration || unit.Kind == KindApexEnumDeclaration || unit.Kind == KindApexInterfaceDeclaration {
		return true
	}
	return false
}

func (extractor *Extractor) ExtractClasses(units []*core.Unit) ([]*object.Clazz, error) {
	ret := make([]*object.Clazz, 0)
	for _, eachUnit := range units {
		if !extractor.IsClass(eachUnit) {
			continue
		}
		eachClazz, err := extractor.ExtractClass(eachUnit)
		if err != nil {
			return nil, err
		}
		ret = append(ret, eachClazz)
	}
	return ret, nil
}

func (extractor *Extractor) ExtractClass(unit *core.Unit) (*object.Clazz, error) {
	clazz := object.NewClazz()
	clazz.Span = unit.Span
	clazz.Lang = extractor.GetLang()
	clazz.Unit = unit

	// trace its class (the closest one
	clazzDecl := core.FindFirstByOneOfKindInParent(unit, KindApexClassDeclaration, KindApexEnumDeclaration, KindApexInterfaceDeclaration)
	clazzIdentifier := core.FindFirstByKindInSubsWithBfs(clazzDecl, KindApexIdentifier)
	if clazzIdentifier == nil {
		return nil, errors.New("no class found in " + unit.Content)
	}
	clazz.Name = clazzIdentifier.Content

	extras := &ClassExtras{}
	// class annotations
	classModifiers := core.FindFirstByKindInSubsWithBfs(clazzDecl, KindApexModifiers)
	if classModifiers != nil {
		classAnnotations := core.FindAllByKindsInSubs(classModifiers, KindApexMarkerAnnotation, KindApexAnnotation)
		if len(classAnnotations) != 0 {
			for _, each := range classAnnotations {
				extras.Annotations = append(extras.Annotations, each.Content)
			}
		}
	}
	// fields
	body := core.FindFirstByKindInSubsWithBfs(clazzDecl, KindApexClassBody)
	if body != nil {
		fields := core.FindAllByKindsInSubs(body, KindApexFieldDeclaration)
		for _, eachField := range fields {
			typeDecl := core.FindFirstByFieldInSubs(eachField, FieldApexType)
			variableDecl := core.FindFirstByFieldInSubs(eachField, FieldApexDeclarator)
			nameDecl := core.FindFirstByKindInSubsWithBfs(variableDecl, KindApexIdentifier)
			if nameDecl == nil || typeDecl == nil {
				return nil, errors.New("not finished field decl")
			}
			field := &ClassField{
				Name:        nameDecl.Content,
				Type:        typeDecl.Content,
				Annotations: nil,
				Modifiers:   nil,
			}
			extras.Fields = append(extras.Fields, field)

			modifiers := core.FindFirstByKindInSubsWithBfs(eachField, KindApexModifiers)
			if modifiers == nil {
				// no modifiers and annotations
				continue
			}
			modifiersStr := modifiers.Content

			// annotation?
			annotations := core.FindAllByKindsInSubs(modifiers, KindApexMarkerAnnotation, KindApexAnnotation)
			if len(annotations) != 0 {
				for _, each := range annotations {
					field.Annotations = append(extras.Annotations, each.Content)
					// remove it from modifiers
					// currently tree-sitter did not split these nodes
					modifiersStr = strings.Replace(modifiersStr, each.Content, "", 1)
				}
			}
			field.Modifiers = strings.Split(strings.TrimSpace(modifiersStr), " ")
		}
	}
	// extends and implements
	extends := core.FindFirstByKindInSubsWithBfs(unit, KindApexSuperClass)
	if extends != nil {
		typeIdentifier := core.FindFirstByOneOfKindInParent(unit, KindApexTypeIdentifier, KindApexGenericType)
		if typeIdentifier != nil {
			extras.Extends = typeIdentifier.Content
		}
	}
	implements := core.FindFirstByKindInSubsWithBfs(unit, KindApexSuperInterface)
	if implements != nil {
		typeList := core.FindFirstByKindInSubsWithBfs(implements, KindApexTypeList)
		// should not nil
		if typeList == nil {
			return nil, errors.New("implements but not decl")
		}
		for _, each := range core.FindAllByKindInSubs(typeList, KindApexTypeIdentifier) {
			extras.Implements = append(extras.Implements, each.Content)
		}
	}

	clazz.Extras = extras

	return clazz, nil
}

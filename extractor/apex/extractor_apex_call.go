package apex

import (
	"errors"

	"github.com/opensibyl/sibyl2/pkg/core"
	"github.com/opensibyl/sibyl2/pkg/extractor/object"
)

func (extractor *Extractor) IsCall(unit *core.Unit) bool {
	if unit.Kind == KindApexMethodInvocation {
		return true
	}
	return false
}

func (extractor *Extractor) ExtractCalls(units []*core.Unit) ([]*object.Call, error) {
	var ret []*object.Call
	for _, eachUnit := range units {
		if !extractor.IsCall(eachUnit) {
			continue
		}

		eachCall, err := extractor.unit2Call(eachUnit)
		if err != nil {
			core.Log.Warnf("err: %v", err)
			continue
		}
		ret = append(ret, eachCall)
	}
	return ret, nil
}

func (extractor *Extractor) unit2Call(unit *core.Unit) (*object.Call, error) {
	funcUnit := core.FindFirstByOneOfKindInParent(unit, KindApexMethodDeclaration)
	var srcFunc *object.Function
	var err error
	if funcUnit != nil {
		srcFunc, err = extractor.ExtractFunction(funcUnit)
		if err != nil {
			return nil, errors.New("convert func failed: " + funcUnit.Content)
		}
	}

	// headless, give up (temp
	if srcFunc == nil {
		return nil, errors.New("headless call")
	}

	var argumentPart *core.Unit
	var arguments []string
	var caller string

	callerPart := core.FindFirstByFieldInSubs(unit, FieldApexObject)
	if callerPart == nil {
		// b()
		callerPart = core.FindFirstByFieldInSubs(unit, FieldApexName)
		argumentPart = core.FindFirstByFieldInSubs(unit, FieldApexArguments)
		caller = callerPart.Content
	} else {
		// a.b()
		identifiers := core.FindAllByKindInSubs(unit, KindApexIdentifier)
		argumentPart = core.FindFirstByFieldInSubs(unit, FieldApexName)

		if len(identifiers) == 0 {
			return nil, errors.New("no id: " + unit.Content)
		}

		caller = callerPart.Content + "." + identifiers[len(identifiers)-1].Content
	}

	// not perfect, eg: anonymous function call?
	if argumentPart != nil {
		for _, each := range core.FindAllByKindInSubs(argumentPart, KindApexIdentifier) {
			arguments = append(arguments, each.Content)
		}
	}

	ret := &object.Call{
		Src:       srcFunc.GetSignature(),
		Caller:    caller,
		Arguments: arguments,
		Span:      unit.Span,
	}
	return ret, nil
}

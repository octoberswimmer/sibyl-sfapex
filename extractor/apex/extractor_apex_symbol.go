package apex

import (
	"strings"

	"github.com/opensibyl/sibyl2/pkg/core"
	"github.com/opensibyl/sibyl2/pkg/extractor/object"
)

func (extractor *Extractor) IsSymbol(unit *core.Unit) bool {
	// todo: use grammar.js instead
	if strings.HasSuffix(unit.Kind, "identifier") {
		return true
	}
	return false
}

func (extractor *Extractor) ExtractSymbols(units []*core.Unit) ([]*object.Symbol, error) {
	ret := make([]*object.Symbol, 0)
	for _, eachUnit := range units {
		if !extractor.IsSymbol(eachUnit) {
			continue
		}
		symbol := &object.Symbol{
			Symbol:    eachUnit.Content,
			Kind:      eachUnit.Kind,
			Span:      eachUnit.Span,
			FieldName: eachUnit.FieldName,
			Unit:      eachUnit,
		}
		ret = append(ret, symbol)
	}
	return ret, nil
}

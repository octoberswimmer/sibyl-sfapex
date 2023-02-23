package apex_test

import (
	"testing"

	"github.com/octoberswimmer/sibyl-sfapex/extractor/apex"
	"github.com/opensibyl/sibyl2/pkg/core"
	"github.com/stretchr/testify/assert"
)

var apexCode = `
@ClassAnnotationA(argA="yes")
@ClassAnnotationB
public class Apex8SnapshotListener extends Apex8MethodLayerListener<Method> {
	private static final int ABC = 1;

	@InjectMocks
	private static final int ABCD = 1;

	private final String DBCA;

    @Override
	@abcde
	@adeflkjbg(abc = "dfff")
    public void enterMethodDeclarationWithoutMethodBody(
            Apex8Parser.MethodDeclarationWithoutMethodBodyContext ctx) {
        super.enterMethodDeclarationWithoutMethodBody(ctx);
        this.storage.save(curMethodStack.peekLast());
    }

    @Override
    public void enterInterfaceMethodDeclaration(Apex8Parser.InterfaceMethodDeclarationContext ctx) {
        super.enterInterfaceMethodDeclaration(ctx);
        this.storage.save(curMethodStack.peekLast());
    }
}

class D extends B {}

class B implements A, C {
	void abcd() {
	}
}

interface A {
	void abcd();
}

interface C {}
`

func TestApexExtractor_ExtractSymbols(t *testing.T) {
	t.Parallel()
	extractor := &apex.Extractor{}
	parser := core.NewParser(extractor.GetLang())
	units, err := parser.Parse([]byte(apexCode))
	if err != nil {
		panic(err)
	}

	symbols, err := extractor.ExtractSymbols(units)
	assert.Nil(t, err)
	assert.NotEmpty(t, symbols)
}

func TestApexExtractor_ExtractFunctions(t *testing.T) {
	t.Parallel()
	extractor := &apex.Extractor{}
	parser := core.NewParser(extractor.GetLang())
	units, err := parser.Parse([]byte(apexCode))
	if err != nil {
		panic(err)
	}

	data, err := extractor.ExtractFunctions(units)
	assert.Nil(t, err)
	assert.NotEmpty(t, data)
	for _, each := range data {
		core.Log.Debugf("each: %s %s %s", each.Name, each.Extras, each.BodySpan.String())
		// check base info
		if each.Name == "enterMethodDeclarationWithoutMethodBody" {
			assert.Equal(t, "15:71,18:5", each.BodySpan.String())
			assert.NotNil(t, each.Extras.(*apex.FunctionExtras).ClassInfo.Annotations)
		}
	}
}

func TestExtractor_ExtractClasses(t *testing.T) {
	t.Parallel()
	extractor := &apex.Extractor{}
	parser := core.NewParser(extractor.GetLang())
	units, err := parser.Parse([]byte(apexCode))
	if err != nil {
		panic(err)
	}

	data, err := extractor.ExtractClasses(units)
	assert.Nil(t, err)
	for _, each := range data {
		core.Log.Debugf("find class: %v", each.GetSignature())
		for _, field := range each.Extras.(*apex.ClassExtras).Fields {
			core.Log.Infof("field: %v", field)
		}
		core.Log.Debugf("class extends: %v", each.Extras.(*apex.ClassExtras).Extends)
		core.Log.Debugf("class implements: %v", each.Extras.(*apex.ClassExtras).Implements)
	}
}

package graphqlgenerator

import (
	"fmt"
	"strconv"
	"strings"
	"testing"
)

func Test_Parse(t *testing.T) {
	testString := `type Query {
  timeseries: int
  transactions: Transactions! 
}
type Mutation {
  performance(word: int = "100"!, fish: Animal): [PerformanceSummary]! 
}`
	reader := strings.NewReader(testString)
	p := NewParser(reader)
	obj, err := p.Parse()
	if err != nil {
		t.Errorf(err.Error())
	}
	if obj.Name != "Query" {
		t.Errorf("Parse failed, found %s as name instead of Query", obj.Name)
	}
	firstVar := obj.Variables[0]
	if err = testModelVarNoArg(firstVar, "timeseries", INT, "int", false, false); err != nil {
		t.Errorf(err.Error())
	}
	secondVar := obj.Variables[1]
	if err = testModelVarNoArg(secondVar, "transactions", IDENT, "Transactions", true, false); err != nil {
		t.Errorf(err.Error())
	}
	obj, err = p.Parse() // to test that successive calls of Parse works
	if err != nil {
		t.Errorf(err.Error())
	}
	if obj.Name != "Mutation" {
		t.Errorf("Parse failed, found %s as name instead of Mutation", obj.Name)
	}
	firstVar = obj.Variables[0]
	if err = testModelVarNoArg(firstVar, "performance", IDENT, "PerformanceSummary", true, true); err != nil {
		t.Errorf(err.Error())
	}

	if err = testGqlArg(firstVar.Arg[0], "word", INT, "int", `"100"`, true); err != nil {
		t.Errorf(err.Error())
	}

	if err = testGqlArg(firstVar.Arg[1], "fish", IDENT, "Animal", "", false); err != nil {
		t.Errorf(err.Error())
	}

}

func testModelVarNoArg(obj ModelVar, name string, tok Token, lit string, required bool, list bool) error {
	if obj.Name != name {
		return fmt.Errorf("Parse failed, member variable Name %s found instead of %s", obj.Name, name)
	}
	if obj.Tok != tok {
		tokString1 := TokenToString(obj.Tok)
		tokString2 := TokenToString(tok)
		return fmt.Errorf("Parse failed, member variable Token %s found instead of %s", tokString1, tokString2)
	}
	if obj.Lit != lit {
		return fmt.Errorf("Parse failed, member variable literal %s found instead of %s", obj.Lit, lit)
	}
	if obj.Required != required {
		return fmt.Errorf("Parse failed, member variable Required boolean is %s instead of %s", strconv.FormatBool(obj.Required), strconv.FormatBool(required))
	}
	if obj.List != list {
		return fmt.Errorf("Parse failed, member variable List boolean is %s instead of %s", strconv.FormatBool(obj.List), strconv.FormatBool(list))
	}

	return nil
}

func testGqlArg(obj GqlArg, name string, tok Token, lit string, defaultstring string, required bool) error {
	if obj.Name != name {
		return fmt.Errorf("Parse failed, argument Name %s found instead of %s", obj.Name, name)
	}
	if obj.Tok != tok {
		tokString1 := TokenToString(obj.Tok)
		tokString2 := TokenToString(tok)
		return fmt.Errorf("Parse failed, argument Token %s found instead of %s", tokString1, tokString2)
	}
	if obj.Lit != lit {
		return fmt.Errorf("Parse failed, argument literal %s found instead of %s", obj.Lit, lit)
	}
	if obj.Default != defaultstring {
		return fmt.Errorf("Parse failed, argument default string %s found instead of %s", obj.Default, defaultstring)
	}
	if obj.Required != required {
		return fmt.Errorf("Parse failed, argument Required boolean %s found instead of %s", strconv.FormatBool(obj.Required), strconv.FormatBool(required))
	}
	return nil
}

package graphqlgenerator

import (
	"fmt"
	"gateway/gqlmanual/hand"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func tokenType(token Token) string {
	switch token {
	case STRING:
		return "graphql.String"
	case FLOAT:
		return "graphql.Float"
	case MAP:
		return "graphql.String"
	case INT:
		return "graphql.Int"
	case BOOLEAN:
		return "graphql.Boolean"
	case ID:
		return "graphql.String"
	case IDENT:
		return "lit"
	}
	return ""
}

func gqlObjString(elements string, name string) string {
	return fmt.Sprintf(`var %s = graphql.NewObject(graphql.ObjectConfig{
    Name: "%s",
    Fields: graphql.Fields{ 
        %s
    },
})
`, name, name, elements)

}

func gqlListString(name string, typename string, arg string, argCheck bool, required bool) string {
	if required == true {
		typename = fmt.Sprintf(`graphql.NewNonNull(%s)`,typename)
	}
	if argCheck == true {
		return fmt.Sprintf(`"%s": &graphql.Field{
            Type: graphql.NewList(%s),
            Args: graphql.FieldConfigArgument{
            	%s
            },
        },
        `, name, typename, arg)
	}
	return fmt.Sprintf(`"%s": &graphql.Field{
            Type: graphql.NewList(%s),
        },
        `, name, typename)
}

func gqlElementString(name string, typename string, arg string, argCheck bool, required bool) string {
	if required == true {
		typename = fmt.Sprintf(`graphql.NewNonNull(%s)`, typename)
	}
	if argCheck == true {
		return fmt.Sprintf(`"%s": &graphql.Field{
            Type: %s,
            Args: graphql.FieldConfigArgument{
            	%s
            },
        },
        `, name, typename, arg)
	}
	return fmt.Sprintf(`"%s": &graphql.Field{
            Type: %s,
        },
        `, name, typename)
}

func gqlArgString(name string, argType string, required bool, argDefault string) string {
	if required == true {
		argType = fmt.Sprintf(`graphql.NewNonNull(%s)`, argType)
	}
	if argDefault != "" {
		argDefault = fmt.Sprintf(`DefaultValue: %s,`, argDefault)
	}
	return fmt.Sprintf(`"%s": &graphql.ArgumentConfig{
                Type: %s,
                %s
        		},`, name, argType, argDefault)

}

func Generate(schemafile string, outputfile string) {
	data, err := ioutil.ReadFile(schemafile)
	if err != nil {
		log.Fatal(err)
	}
	g := strings.NewReader(string(data))
	p := NewParser(g)
	packageName, err := p.ParsePackage()
	if err != nil{
		fmt.Println(err)
		return
	}

	obj, err := p.Parse()
	if err != nil {
		fmt.Println(err)
		return
	}
	toadd := fmt.Sprintf("package %s \n \n", packageName)
	for obj != nil {
		curadd := ""
		for _, element := range obj.Variables {

			tType := tokenType(element.Var[0])
			if tType == "lit" {
				tType = element.Lit
			}
			argString := ""
			argCheck := false
			for _, arg := range element.Arg {
				argType := tokenType(arg.Type)
				if argType == "lit" {
					argType = arg.Lit
				}
				argString += gqlArgString(arg.Name, argType, arg.Required, arg.Default)
				argCheck = true
			}

			if element.List == true {
				curadd += gqlListString(element.Name, tType, argString, argCheck,element.Required )
			} else {
				curadd += gqlElementString(element.Name, tType, argString, argCheck, element.Required)
			}
		}
		toadd += gqlObjString(curadd, obj.Name)
		obj, err = p.Parse()
		if err != nil {
			fmt.Println(err)
		}
	}
	f, err := os.Create(outputfile)
	if err != nil {
		fmt.Println(err)
		return
	}
	l, err := f.WriteString(toadd)
	if err != nil {
		fmt.Println(err)
		f.Close()
		return
	}
	fmt.Println(l, "bytes written successfully")
	err = f.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
	return
}

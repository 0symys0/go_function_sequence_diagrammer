package main

import (
  "bytes"
  "fmt"
  "log"
  "go/parser"
  "go/ast"
  "go/token"
  "go/printer"
  "strings"
  "flag"
  //"os"


  "github.com/goccy/go-graphviz"
)

func main() {
  var filename_in string
  flag.StringVar(&filename_in, "file", "analyze_me.go", "name of source file to analyze")
  flag.Parse()
  fmt.Println("filename_in: ", filename_in)
  g := graphviz.New()
  graph, err := g.Graph()
  if err != nil {
    log.Fatal(err)
  }
  defer func() {
    if err := graph.Close(); err != nil {
      log.Fatal(err)
    }
    g.Close()
  }()
    fset := token.NewFileSet()
    node, err := parser.ParseFile(fset, filename_in, nil, parser.ParseComments)
    if err != nil {
        log.Fatal(err)
    }
    //fmt.Println("Imports:")
	imp_node, err := graph.CreateNode("Imports:")
	if err != nil {
	  log.Fatal(err)
	}
    for _, i := range node.Imports {
		label := fmt.Sprintf("%v",i.Path.Value)
		imp_node, err = graph.CreateNode(label)
		if err != nil {
		  log.Fatal(err)
		}
        //fmt.Println(i.Path.Value)
    }
	//fmt.Println("Comments:")
	comm_node, err := graph.CreateNode("Comments:")
	for _, c := range node.Comments {
		label := fmt.Sprintf("%v",c.Text())
		comm_node, err = graph.CreateNode(label)
		if err != nil {
		  log.Fatal(err)
		}
		//fmt.Print(c.Text())
	}

	//fmt.Println("Functions:")
	func_node, err := graph.CreateNode("Functions:")
	for _, f := range node.Decls {
		fn, ok := f.(*ast.FuncDecl)
		if !ok {
			continue
		}
		var buf bytes.Buffer
		printer.Fprint(&buf, fset, fn.Body)

		// Remove braces {} enclosing the function body, unindent,
		// and trim leading and trailing white space.
		funcbody := buf.String()
		funcbody = funcbody[1 : len(funcbody)-1]
		funcbody = strings.TrimSpace(strings.ReplaceAll(funcbody, "\n\t", "\n"))
		label := fmt.Sprintf("%v : %v",fn.Name.Name, funcbody)
		func_node, err = graph.CreateNode(label)
		if err != nil {
		  log.Fatal(err)
		}
		_, err = graph.SubNode(func_node,1)
		if err != nil {
		  log.Fatal(err)
		}
		for _, thing := range(fn.Body.List){
			var thingbuf bytes.Buffer
			printer.Fprint(&thingbuf, fset, thing)
			thingstring := thingbuf.String()
			thingstring = strings.TrimSpace(strings.ReplaceAll(thingstring, "\n\t", "\n"))
			outstring := fmt.Sprintf("type(thing): %T\nthingstring:\n%v",thing,thingstring)
			fmt.Println(outstring)
			// this is how we check the type in go, with this type assertion statement:
			thingIf, isIf := thing.(*ast.IfStmt)
            if isIf {
                fmt.Println("^^^ IF STATEMENT!! vvv")
				outstring = fmt.Sprintf("thingIf.Cond: %v\nthingIf.Body.List: %v\n",get_node_string(fset, thingIf.Cond),get_node_string(fset, thingIf.Body.List))
				fmt.Println(outstring)
                
            }
			thingAss, isAss := thing.(*ast.AssignStmt)
            if isAss {
                fmt.Println("^^^ ASSIGNMENT STATEMENT!! vvv")
				outstring = ""
				for lhs_ind, lhs := range(thingAss.Lhs){
					for rhs_ind, rhs := range(thingAss.Rhs){
						outstring = fmt.Sprintf("%vthingAss.Lhs[%v]: %v\nthingAss.Rhs[%v]: %v\n",outstring, lhs_ind,get_node_string(fset, lhs),rhs_ind,get_node_string(fset, rhs))
					}
				}
				fmt.Println(outstring)
            }
			thingExpr, isExpr := thing.(*ast.ExprStmt)
            if isExpr {
                fmt.Println("^^^ EXPRESSION STATEMENT!! vvv")
				outstring = fmt.Sprintf("thingExpr.X: %v\n",get_node_string(fset, thingExpr.X))
				fmt.Println(outstring)
            }
			thingFor, isFor := thing.(*ast.ForStmt)
            if isFor {
                fmt.Println("^^^ FOR STATEMENT!! vvv")
				outstring = fmt.Sprintf("thingFor.Cond: %v\nthingFor.Body.List: %v\n",get_node_string(fset, thingFor.Cond),get_node_string(fset, thingFor.Body.List))
				fmt.Println(outstring)
            }
			thingReturn, isReturn := thing.(*ast.ReturnStmt)
            if isReturn {
                fmt.Println("^^^ RETURN STATEMENT!! vvv")
				str_builder := ""
				for n, result := range(thingReturn.Results){
					str_builder = fmt.Sprintf("%vthingReturn.Results[%v]: %v\n",str_builder,n,get_node_string(fset, result))
				}
				fmt.Println(str_builder)
            }
			thingSwitch, isSwitch := thing.(*ast.SwitchStmt)
            if isSwitch {
                fmt.Println("^^^ SWITCH STATEMENT!! vvv")
				outstring = fmt.Sprintf("thingSwitch.Tag: %v\nthingSwitch.Body.List: %v\n",get_node_string(fset, thingSwitch.Tag),get_node_string(fset, thingSwitch.Body.List))
				fmt.Println(outstring)
            }
		}
	}

  if err != nil {
    log.Fatal(err)
  }
  e, err := graph.CreateEdge("e", imp_node, func_node)
  if err != nil {
    log.Fatal(err)
  }
  e.SetLabel("Last imp to last func")
  f, err := graph.CreateEdge("e", imp_node, comm_node)
  if err != nil {
    log.Fatal(err)
  }
  f.SetLabel("Last imp to last comm")
  var buf bytes.Buffer
  if err := g.Render(graph, "dot", &buf); err != nil {
    log.Fatal(err)
  }
  //fmt.Println(buf.String())
}

func get_node_string(fset *token.FileSet, node interface{}) (outstring string){
	var thingbuf bytes.Buffer
	printer.Fprint(&thingbuf, fset, node)
	outstring = thingbuf.String()
	outstring = strings.TrimSpace(strings.ReplaceAll(outstring, "\n\t", "\n"))
	return
}

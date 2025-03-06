package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strings"
	"text/tabwriter"
)

type Node struct {
	Name      string
	Type      string
	Adjacents []string
}

var (
	grammarRules = make(map[string][]string)
	initSymbol   string
	precGraph    = make(map[string]Node)
	terminals    = []string{}
)

var built = false

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Bienvenido al analizador sintáctico")

	for {
		fmt.Print("$> ")
		scanner.Scan()
		input := scanner.Text()
		parts := strings.Fields(input)

		if len(parts) == 0 {
			continue
		}

		command := strings.ToUpper(parts[0])

		switch command {
		case "EXIT":
			fmt.Println("Saliendo del programa...")
			return
		case "RULE":
			if len(parts) >= 2 {
				if len(parts) >= 2 {
					noTerminal := parts[1]
					simbolos := strings.Join(parts[2:], " ")
					handleRule(noTerminal, simbolos)
				}
			} else {
				fmt.Println("ERROR: Falta el argumento <no-terminal> para la instrucción RULE")
			}
		case "INIT":
			if len(parts) == 2 {
				noTerminal := parts[1]
				handleInit(noTerminal)
			} else {
				fmt.Println("ERROR: Falta el argumento <no-terminal> para la instrucción INIT")
			}
		case "PREC":
			if len(parts) == 4 {
				terminal1 := parts[1]
				op := parts[2]
				terminal2 := parts[3]
				handlePrec(terminal1, op, terminal2)
			} else {
				fmt.Println("ERROR: Argumentos incorrectos para la instrucción PREC")
			}
		case "BUILD":
			handleBuild()
		case "PARSE":
			if len(parts) > 1 {
				inputString := strings.Join(parts[1:], " ")
				handleParse(inputString)
			} else {
				fmt.Println("ERROR: Falta el argumento <string> para la instrucción PARSE")
			}
		default:
			fmt.Println("Instrucción no reconocida:", input)
		}
	}
}

func handleRule(noTerminal string, rhs string) {
	if !IsUpper(noTerminal) {
		fmt.Println("ERROR: \"" + noTerminal + "\" no es un no-terminal")
		return
	}

	simbolos := strings.Fields(rhs)
	for i := range len(simbolos) - 1 {
		r := simbolos[i]
		if (IsLower(r) && IsLower(simbolos[i+1])) || (IsUpper(r) && IsUpper(simbolos[i+1]) || (IsSymbol(r) && IsSymbol(simbolos[i+1])) || r == "$") {
			fmt.Println("ERROR: \"" + noTerminal + " -> " + strings.Join(simbolos, " ") + "\" no corresponde a una gramática de operadores")
			return
		}
		if !IsUpper(r) {
			terminals = append(terminals, r)
		}
	}

	if !IsUpper(simbolos[len(simbolos)-1]) {
		terminals = append(terminals, simbolos[len(simbolos)-1])
	}

	grammarRules[noTerminal] = append(grammarRules[noTerminal], rhs)
	fmt.Println("Regla \"" + noTerminal + " -> " + rhs + "\" agregada a la gramática")
}

func handleInit(noTerminal string) {
	_, ok := grammarRules[noTerminal]
	if !ok {
		fmt.Println("ERROR: \"" + noTerminal + "\" no es un no-terminal")
		return
	}

	initSymbol = noTerminal
	fmt.Println("\"" + noTerminal + "\" es ahora el símbolo inicial de la gramática")
}

func handlePrec(terminal1, op, terminal2 string) {
	if op != "<" && op != ">" && op != "=" {
		fmt.Println("ERROR: \"" + op + " no es un operador válido")
		return
	}

	if IsUpper(terminal1) || IsUpper(terminal2) {
		fmt.Println("ERROR: Los símbolos deben ser terminales")
		return
	}

	strOp := ""
	f1 := "f_" + terminal1
	g1 := "g_" + terminal1
	f2 := "f_" + terminal2
	g2 := "g_" + terminal2

	// Inicializar los nodos con arreglos vacíos si no existen, dejar igual si existen
	if _, ok := precGraph[f1]; !ok {
		precGraph[f1] = Node{
			Name:      terminal1,
			Type:      "f",
			Adjacents: []string{},
		}
	}
	if _, ok := precGraph[g1]; !ok {
		precGraph[g1] = Node{
			Name:      terminal1,
			Type:      "g",
			Adjacents: []string{},
		}
	}
	if _, ok := precGraph[f2]; !ok {
		precGraph[f2] = Node{
			Name:      terminal2,
			Type:      "f",
			Adjacents: []string{},
		}
	}
	if _, ok := precGraph[g2]; !ok {
		precGraph[g2] = Node{
			Name:      terminal2,
			Type:      "g",
			Adjacents: []string{},
		}
	}

	switch op {
	case "<":
		if entry, ok := precGraph[g2]; ok {
			entry.Adjacents = append(entry.Adjacents, f1)
			precGraph[g2] = entry
		}
		strOp = "menor"
	case ">":
		if entry, ok := precGraph[f1]; ok {
			entry.Adjacents = append(entry.Adjacents, g2)
			precGraph[f1] = entry
		}
		strOp = "mayor"
	case "=":
		strOp = "igual"
	}
	fmt.Println("\"" + terminal1 + "\"" + " tiene " + strOp + " precedencia que " + "\"" + terminal2 + "\"")
}

func handleBuild() {
	if hasCycle(precGraph) {
		fmt.Println("ERROR: La gramática de precedencia tiene ciclos")
		return
	}

	fmt.Println("Analizador sintáctico construido")

	fmt.Println("Valores para f:")
	for name, node := range precGraph {
		if node.Type == "f" {
			f := longestPath(precGraph, name)
			fmt.Println(node.Name+":", f)
		}
	}

	fmt.Println("Valores para g:")
	for name, node := range precGraph {
		if node.Type == "g" {
			g := longestPath(precGraph, name)
			fmt.Println(node.Name+":", g)
		}
	}

	built = true
}

func handleParse(inputString string) {
	if !built {
		fmt.Println("ERROR: Aún no se ha construido el analizador sintáctico")
		return
	}

	inputNonTerminals := findNonTerminals(inputString)
	if len(inputNonTerminals) > 0 {
		fmt.Printf("ERROR: Los siguientes símbolos son no-terminales: \"%s\"\n", strings.Join(inputNonTerminals, ",\" "))
		return
	}

	nonTerminalsNotInGramamr := findTerminalsNotInGrammar(inputString, terminals)
	if len(nonTerminalsNotInGramamr) > 0 {
		fmt.Printf("ERROR: Los siguientes símbolos no son terminales de la gramática: \"%s\"\n", strings.Join(nonTerminalsNotInGramamr, ",\" "))
		return
	}

	nonComparables := findNonComparables(inputString)
	if len(nonComparables) > 0 {
		fmt.Printf("ERROR: \"%s\" no es comparable con \"%s\"\n", nonComparables[0], nonComparables[0])
		return
	}

	if len(inputString) == 0 {
		fmt.Println("ERROR: \"$\" no es comparable con \"$\"")
		return
	}

	symbols := calculatePrecedence(inputString)
	parseSymbols(inputString, symbols)
}

func parseSymbols(initialString, symbols string) {
	stack := []string{}
	symbolArray := strings.Fields(symbols)
	displayString := initialString
	startIndex := 0
	endIndex := 0

	// Tabwriter
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 5, ' ', 0)

	fmt.Fprintln(w, "Pila\tEntrada\tAcción")

	for i := range symbolArray {
		// Leer
		if symbolArray[i] != ">" && symbolArray[i] != "<" && symbolArray[i] != "=" && symbolArray[i] != "$" {
			fmt.Fprintf(w, "%s\t%s\tleer\n", strings.Join(stack, " "), calculatePrecedence(displayString))
			stack = append(stack, symbolArray[i])
			endIndex += 1
			startIndex = endIndex

			// Reducir
		} else if symbolArray[i] == ">" {
			subset := ""
			for j := i - 1; j >= 0; j-- {
				if symbolArray[j] != ">" && symbolArray[j] != "<" && symbolArray[j] != "=" && symbolArray[j] != "$" {
					if subset == "" {
						subset = symbolArray[j]
					} else {
						subset = fmt.Sprintf("%s %s", symbolArray[j], subset)
					}
					startIndex -= 1
				} else if symbolArray[j] == "<" {
					rule, found := findRule(grammarRules, subset)
					if found {
						fmt.Fprintf(w, "%s\t%s\treducir %s -> %s\n", strings.Join(stack, " "), calculatePrecedence(displayString), rule[0], rule[1])
						symbolArray = append(append(symbolArray[:j+1], rule[0]), symbolArray[i:]...)
						displayStringParts := slices.Delete(strings.Fields(displayString), startIndex, startIndex+1)
						displayString = strings.Join(displayStringParts, " ")
						endIndex = startIndex
						stack = stack[:len(stack)-len(strings.Fields(subset))]
						stack = append(stack, rule[0])
						break
					} else {
						fmt.Fprintf(w, "%s\t%s\trechazar, no se puede reducir por %s -> %s\n", strings.Join(stack, " "), calculatePrecedence(displayString), rule[0], rule[1])
						w.Flush()
						return
					}
				}
			}
		}
	}

	// Reducir y aceptar
	for i := len(symbolArray) - 1; i >= 0; i-- {
		if symbolArray[i] == ">" {
			subset := ""
			for j := i - 1; j >= 0; j-- {
				if symbolArray[j] != ">" && symbolArray[j] != "<" && symbolArray[j] != "=" && symbolArray[j] != "$" {
					if subset == "" {
						subset = symbolArray[j]
					} else {
						subset = fmt.Sprintf("%s %s", symbolArray[j], subset)
					}
				} else if symbolArray[j] == "<" {
					rule, found := findRule(grammarRules, subset)
					if found {
						fmt.Fprintf(w, "%s\t%s\treducir %s -> %s\n", strings.Join(stack, " "), calculatePrecedence(displayString), rule[0], rule[1])
						symbolArray = append(append(symbolArray[:j+1], rule[0]), symbolArray[i:]...)
						startIndex -= 1
						displayStringParts := slices.Delete(strings.Fields(displayString), startIndex, startIndex+1)
						displayString = strings.Join(displayStringParts, " ")
						stack = stack[:len(stack)-len(strings.Fields(subset))]
						stack = append(stack, rule[0])
						i = len(symbolArray) - 1
						break
					} else {
						partsSubset := strings.Fields(subset)
						if len(partsSubset) == 1 {
							continue
						} else {
							for i := range len(partsSubset) - 1 {
								isTerminal1 := IsLower(partsSubset[i]) || IsSymbol(partsSubset[i])
								isTerminal2 := IsLower(partsSubset[i+1]) || IsSymbol(partsSubset[i+1])
								nearLhs, nearRhs := findClosestRule(subset[1:], grammarRules)
								if isTerminal1 && isTerminal2 {
									fmt.Fprintf(w, "%s\t%s\trechazar, no se puede reducir por %s -> %s\n", strings.Join(stack, " "), calculatePrecedence(displayString), nearLhs, nearRhs)
									w.Flush()
									return
								}
							}
						}
					}
				}
			}
		} else if len(strings.Fields(displayString)) == 0 {
			fmt.Fprintf(w, "%s\t$ $\taceptar\n", strings.Join(stack, " "))
			break
		}
	}

	w.Flush()
}

func calculatePrecedence(inputString string) string {
	input := strings.Fields(inputString)
	symbols := ""
	if len(input) == 0 {
		return symbols
	}

	if inputString[0] != '$' {
		input = append([]string{"$"}, input...)
	}

	if inputString[len(inputString)-1] != '$' {
		input = append(input, "$")
	}

	for i := range input {
		symbols += input[i] + " "
		if i < len(input)-1 {

			f_current := longestPath(precGraph, "f_"+input[i])
			g_next := longestPath(precGraph, "g_"+input[i+1])

			if f_current > g_next {
				symbols += "> "
			} else if f_current < g_next {
				symbols += "< "
			} else {
				symbols += "= "
			}
		}
	}

	return symbols
}

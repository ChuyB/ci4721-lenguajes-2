package main

import (
	"slices"
	"strings"
	"unicode"
)

func IsUpper(s string) bool {
	for _, r := range s {
		if !unicode.IsUpper(r) && unicode.IsLetter(r) {
			return false
		} else if !unicode.IsLetter(r) {
			return false
		}
	}
	return true
}

func IsLower(s string) bool {
	for _, r := range s {
		if !unicode.IsLower(r) && unicode.IsLetter(r) {
			return false
		} else if !unicode.IsLetter(r) {
			return false
		}
	}
	return true
}

func IsSymbol(s string) bool {
	for _, r := range s {
		if unicode.IsLetter(r) {
			return false
		}
	}
	return true
}

func hasCycle(graph map[string]Node) bool {
	visited := make(map[string]bool)
	recStack := make(map[string]bool)

	for node := range graph {
		if !visited[node] && isCyclic(graph, node, visited, recStack) {
			return true
		}
	}
	return false
}

func isCyclic(graph map[string]Node, node string, visited, recStack map[string]bool) bool {
	visited[node] = true
	recStack[node] = true

	for _, neighbour := range graph[node].Adjacents {
		if !visited[neighbour] && isCyclic(graph, neighbour, visited, recStack) {
			return true
		} else if recStack[neighbour] {
			return true
		}
	}
	recStack[node] = false
	return false
}

func longestPath(graph map[string]Node, node string) int {
	visited := make(map[string]bool)
	return dfs(graph, node, visited)
}

func dfs(graph map[string]Node, node string, visited map[string]bool) int {
	visited[node] = true
	maxLength := 0

	for _, neighbour := range graph[node].Adjacents {
		if !visited[neighbour] {
			pathLength := dfs(graph, neighbour, visited)
			if pathLength+1 > maxLength {
				maxLength = pathLength + 1
			}
		}
	}

	visited[node] = false
	return maxLength
}

func findRule(grammarRules map[string][]string, terminal string) ([]string, bool) {
	for lhs, rhs := range grammarRules {
		for _, rule := range rhs {
			if rule == terminal {
				return []string{lhs, rule}, true
			}
		}
	}

	return nil, false
}

func findClosestRule(input string, rules map[string][]string) (string, string) {
	var closestRule string
	var closestNonTerminal string
	shortestDistance := -1

	for nonTerminal, rightSides := range rules {
		for _, rule := range rightSides {
			distance := levenshteinDistance(input, rule)
			if shortestDistance == -1 || distance < shortestDistance {
				shortestDistance = distance
				closestRule = rule
				closestNonTerminal = nonTerminal
			}
		}
	}
	return closestNonTerminal, closestRule
}

func levenshteinDistance(a, b string) int {
	al := len(a)
	bl := len(b)
	d := make([][]int, al+1)
	for i := range d {
		d[i] = make([]int, bl+1)
	}
	for i := 0; i <= al; i++ {
		d[i][0] = i
	}
	for j := 0; j <= bl; j++ {
		d[0][j] = j
	}
	for i := 1; i <= al; i++ {
		for j := 1; j <= bl; j++ {
			if a[i-1] == b[j-1] {
				d[i][j] = d[i-1][j-1]
			} else {
				d[i][j] = min(d[i-1][j-1]+1, min(d[i-1][j]+1, d[i][j-1]+1))
			}
		}
	}
	return d[al][bl]
}

func findNonTerminals(input string) []string {
	nonTerminals := []string{}
	symbols := strings.Fields(input)
	for _, symbol := range symbols {
		if IsUpper(symbol) {
      nonTerminals = append(nonTerminals, symbol)
		}
	}
	return nonTerminals
}

func findTerminalsNotInGrammar(input string, nonTerminals []string) []string {
	notFound := []string{}
	symbols := strings.Fields(input)
	for _, symbol := range symbols {
		found := slices.Contains(nonTerminals, symbol)
		if !found {
			notFound = append(notFound, symbol)
		}
	}
	return notFound
}

func findNonComparables(input string) []string {
  nonComparables := []string{}
  symbols := strings.Fields(input)
  for  i := range len(symbols) - 1  {
    if (symbols[i] == symbols[i+1]) {
      nonComparables = append(nonComparables, symbols[i])
    }
  }

  return nonComparables
}


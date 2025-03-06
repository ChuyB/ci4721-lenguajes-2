package main

import ("testing")

func TestMain(t *testing.T) {
  handleRule("E", "E + E")
  handleRule("E", "E * E")
  handleRule("E", "E E")
  handleRule("E", "n")
  handleRule("E", "n $ n")
  handleInit("e")
  handleInit("E")
  handlePrec("n", ">", "+")
  handlePrec("n", ">", "*")
  handlePrec("n", ">", "$")
  handlePrec("+", "<", "n")
  handlePrec("+", ">", "+")
  handlePrec("+", "<", "*")
  handlePrec("+", ">", "$")
  handlePrec("*", "<", "n")
  handlePrec("*", ">", "+")
  handlePrec("*", ">", "*")
  handlePrec("*", ">", "$")
  handlePrec("$", "<", "n")
  handlePrec("$", "<", "+")
  handlePrec("$", "<", "*")

  handleBuild()

  handleParse("n + n * n")
  handleParse("n + * n")
  handleParse("n")
  handleParse("a + b * c")
  handleParse("n n")
  handleParse("")
}


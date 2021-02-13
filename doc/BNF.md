# BNF

``` ebnf
<SQL> ::= <SELECT Statement>
<SELECT Statement> ::= "SELECT" <Columns> "FROM" <Table Expr> ";"
<Columns> ::= <Column> ["," <Column>]...
<Column> ::= <expr>
<expr> ::=  <Literal> | <Unary Ope Expr> | <Binary Ope Expr> | <Function Call Expr> | <Is Expr> | <Column Expr> | "(" <expr> ")" | "CAST" "(" <expr> "AS" <type> ")"
<Literal> ::= <Bool> | <*Numeric*> | <*String*> | "NULL"
<Bool> ::= "TRUE" | "FALSE"
<Unary Ope Expr> ::= <*Unary Ope> <expr>
<Binary Ope Expr> ::= <expr> <*Binary Ope> <expr>
<Function Call Expr> ::= <*Function Name> "(" [<expr> ["," <expr>]...]... ")"
<Is Expr> ::= <expr> <"IS" ["NOT"] "NULL" | "ISNULL" | "NOTNULL">
<Column Expr> ::= [[[<*schema*> "."] <*database*> "."] <*table*> "."] <*column*>
<Table Expr> ::= [[<*schema*> "."] <*database*> "."] <*table*>
```

## Unary Ope

``` ebnf
<unary ope> ::= "+" | "-" | "NOT"
```
## Binary Ope

``` ebnf
<binary ope> ::= "+" | "-" | "*" | "/" | "%" | "AND" | "OR" | "=" | "<>" | "<" | ">" | ">=" | "<="
```

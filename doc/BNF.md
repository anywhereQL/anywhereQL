# BNF

``` ebnf
<SQL> ::= <SELECT Statement>
<SELECT Statement> ::= "SELECT" <Columns> [<FROM Clause>] [<WHERE Clause>] ";"
<Columns> ::= <Column> ["," <Column>]...
<Column> ::= <expr>
<expr> ::=   <Literal>
           | <Unary Ope Expr>
           | <Binary Ope Expr>
           | <Function Call Expr>
           | <Is Expr>
           | <In Expr>
           | <Column Expr>
           | "(" <expr> ")"
           | "CAST" "(" <expr> "AS" <type> ")"
           | "CASE" [<expr>] ["WHEN" <expr> "THEN" <expr>]+ ["ELSE" <expr>] "END"
           | <expr> ["NOT"] "BWTWEEN" <expr> "AND" <expr>
<Literal> ::= <Bool> | <*Numeric*> | <*String*> | "NULL"
<Bool> ::= "TRUE" | "FALSE"
<Unary Ope Expr> ::= <*Unary Ope> <expr>
<Binary Ope Expr> ::= <expr> <*Binary Ope> <expr>
<Function Call Expr> ::= <*Function Name> "(" [<expr> ["," <expr>]...]... ")"
<Is Expr> ::= <expr> (("IS" ["NOT"] |"NULL" | "ISNULL" | "NOTNULL") | ("IS" ["NOT"] <expr>)
<In Expr> ::= <expr> ["NOT"] "IN" (<Table Expr> | "(" (<expr> ["," <expr>]... ")")
<Column Expr> ::= [[[<*schema*> "."] <*database*> "."] <*table*> "."] <*column*>
<FROM Clause> ::= "FROM" <Table Expr> [["AS"] <*alias*>] [([INNER | ([REFT | RIGHT | FULL] [OUTER]) | CROSS] JOIN <Table Expr> [["AS"] <*alias*>] "ON" <expr>) | ("," <Table Expr>)]...
<Table Expr> ::= [[<*schema*> "."] <*database*> "."] <*table*>
<WHERE Clause> ::= "WHERE" <expr>
```

## Unary Ope

``` ebnf
<unary ope> ::= "+" | "-" | "NOT"
```
## Binary Ope

``` ebnf
<binary ope> ::= "+" | "-" | "*" | "/" | "%" | "AND" | "OR" | "=" | "<>" | "<" | ">" | ">=" | "<="
```

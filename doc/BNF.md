# BNF

``` ebnf
<SQL> ::= <SELECT Statement>
<SELECT Statement> ::= "SELECT" <Columns> [<FROM Clause>] ";"
<Columns> ::= <Column> ["," <Column>]...
<Column> ::= <expr>
<expr> ::=   <Literal>
           | <Unary Ope Expr>
           | <Binary Ope Expr>
           | <Function Call Expr>
           | <Is Expr>
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
<Column Expr> ::= [[[<*schema*> "."] <*database*> "."] <*table*> "."] <*column*>
<FROM Clause> ::= "FROM" <Table Expr> [["AS"] <*alias*>] [INNER JOIN <Table Expr> [["AS"] <*alias*>] "ON" <expr>]...
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

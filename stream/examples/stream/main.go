package main

import (
	"fmt"
	"strings"

	"github.com/tiagods/go-extras/stream"
)

// Definir uma estrutura de pessoa para usar nos exemplos
type Pessoa struct {
	Nome   string
	Idade  int
	Cidade string
}

// Dados de exemplo que serão usados por várias funções
var (
	pessoas = []Pessoa{
		{Nome: "João", Idade: 25, Cidade: "São Paulo"},
		{Nome: "Maria", Idade: 30, Cidade: "Rio de Janeiro"},
		{Nome: "Pedro", Idade: 28, Cidade: "São Paulo"},
		{Nome: "Ana", Idade: 35, Cidade: "Belo Horizonte"},
		{Nome: "Carlos", Idade: 22, Cidade: "Rio de Janeiro"},
	}

	palavras = []string{"zebra", "abacaxi", "maçã", "banana"}
)

func main() {
	fmt.Println("==== Exemplos de Stream ====")

	// Cria uma lista de números de 1 a 99
	lista := criarListaNumeros(1, 99)

	// Exemplos organizados por funcionalidade
	exemploFiltro(lista)
	exemploForEach()
	exemploReduce()
	exemploSort()
	exemploMap()
	exemploFind()
	exemploFlatMap()
	exemploDistinct()
	exemploCombinacaoOperacoes(lista)
	exemploCount(lista)
	exemploJoin()
	exemploGroupBy()
	exemploGroupByWithValueMapper()
}

// Função auxiliar para criar uma lista de números
func criarListaNumeros(inicio, fim int) []int {
	resultado := make([]int, fim-inicio+1)
	for i := inicio; i <= fim; i++ {
		resultado[i-inicio] = i
	}
	return resultado
}

// Exemplo 1: Filter - filtrar elementos
func exemploFiltro(lista []int) {
	fmt.Println("\n1. Filter - números pares:")

	pares := stream.NewStream(lista).
		Filter(func(i int) bool {
			return i%2 == 0
		})

	// Limitando a saída para os primeiros 5 elementos
	primeiros5Pares := stream.Limit(pares, 5)
	fmt.Println("Primeiros 5 números pares:", primeiros5Pares.ToSlice())
}

// Exemplo 2: ForEach - aplicar operação a cada elemento
func exemploForEach() {
	fmt.Println("\n2. ForEach - operação para cada elemento:")

	stream.NewStream([]string{"maçã", "banana", "laranja"}).
		ForEach(func(fruta string) {
			fmt.Printf("Fruta: %s\n", fruta)
		})
}

// Exemplo 3: Reduce - reduzir elementos a um valor
func exemploReduce() {
	fmt.Println("\n3. Reduce - soma de números:")

	soma := stream.NewStream([]int{1, 2, 3, 4, 5}).
		Reduce(func(acc, atual int) int {
			return acc + atual
		}, 0)

	fmt.Println("Soma dos números de 1 a 5:", soma)
}

// Exemplo 4: Sort - ordenar elementos
func exemploSort() {
	fmt.Println("\n4. Sort - ordenação:")

	ordenado := stream.NewStream(palavras).
		Sort(func(a, b string) bool {
			return a < b
		})

	fmt.Println("Palavras ordenadas:", ordenado.ToSlice())
}

// Exemplo 5: Map - transformar elementos
func exemploMap() {
	fmt.Println("\n5. Map - transformação:")

	maiusculas := stream.Map(stream.NewStream(palavras), func(s string) string {
		return strings.ToUpper(s)
	})

	fmt.Println("Palavras em maiúsculas:", maiusculas.ToSlice())
}

// Exemplo 6: Find - encontrar elementos
func exemploFind() {
	fmt.Println("\n6. FindFirst e FindAny:")

	numeros := stream.NewStream([]int{10, 20, 30, 40, 50})

	// FindFirst - obter o primeiro elemento
	primeiro := numeros.FindFirst()
	if val, ok := primeiro.GetIfPresent(); ok {
		fmt.Println("Primeiro número:", val)
	}

	// FindAny - obter um elemento qualquer
	qualquer := numeros.FindAny()
	if val, ok := qualquer.GetIfPresent(); ok {
		fmt.Println("Um número qualquer:", val)
	}
}

// Exemplo 7: FlatMap - achatar estruturas aninhadas
func exemploFlatMap() {
	fmt.Println("\n7. FlatMap - achatamento:")

	matrizNumeros := [][]int{
		{1, 2, 3},
		{4, 5, 6},
		{7, 8, 9},
	}

	achatado := stream.FlatMap(stream.NewStream(matrizNumeros), func(linha []int) []int {
		return linha
	})

	fmt.Println("Matriz achatada:", achatado.ToSlice())
}

// Exemplo 8: Distinct - remover duplicatas
func exemploDistinct() {
	fmt.Println("\n8. Distinct - remover duplicatas:")

	comDuplicatas := stream.NewStream([]int{1, 2, 2, 3, 3, 3, 4, 5, 5})
	semDuplicatas := comDuplicatas.Distinct()

	fmt.Println("Sem duplicatas:", semDuplicatas.ToSlice())
}

// Exemplo 9: Combinação de operações
func exemploCombinacaoOperacoes(lista []int) {
	fmt.Println("\n9. Combinação de operações:")

	resultado := stream.NewStream(lista).
		Filter(func(i int) bool {
			return i%2 == 1 // números ímpares
		}).
		Filter(func(i int) bool {
			return i > 50 // maiores que 50
		}).
		Sort(func(a, b int) bool {
			return a > b // ordem decrescente
		})

	fmt.Println("Ímpares maiores que 50 em ordem decrescente:", resultado.ToSlice())
}

// Exemplo 10: Count - contar elementos
func exemploCount(lista []int) {
	fmt.Println("\n10. Contagem de elementos:")

	pares := stream.NewStream(lista).
		Filter(func(i int) bool {
			return i%2 == 0
		})

	fmt.Println("Quantidade de números pares:", pares.Count())
	fmt.Println("Quantidade de palavras:", stream.NewStream(palavras).Count())
}

// Exemplo 11: Join - concatenar elementos
func exemploJoin() {
	fmt.Println("\n11. Join - concatenar elementos:")

	frutas := stream.NewStream([]string{"maçã", "banana", "uva", "laranja"})
	frutasCSV := frutas.Join(", ")
	fmt.Println("Frutas (CSV):", frutasCSV)

	nums := stream.NewStream([]int{1, 2, 3, 4, 5})
	numerosTraco := nums.Join(" - ")
	fmt.Println("Números separados por traço:", numerosTraco)

	letrasSemSeparador := stream.NewStream([]string{"a", "b", "c", "d"}).Join("")
	fmt.Println("Letras concatenadas:", letrasSemSeparador)
}

// Exemplo 12: GroupBy - agrupar elementos por uma chave
func exemploGroupBy() {
	fmt.Println("\n12. GroupBy - agrupar elementos:")

	// Exemplo com números - agrupar por paridade (par/ímpar)
	numsList := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	numStream := stream.NewStream(numsList)

	// Usando o método genérico GroupBy
	agrupadosPorParidade := numStream.GroupBy(func(n int) interface{} {
		if n%2 == 0 {
			return "Par"
		}
		return "Ímpar"
	})

	fmt.Println("Números agrupados por paridade:")
	for chave, valores := range agrupadosPorParidade {
		fmt.Printf("  %v: %v\n", chave, valores)
	}

	// Exemplo com estruturas - agrupar pessoas por cidade
	pessoasStream := stream.NewStream(pessoas)

	// Também podemos usar o método específico GroupByString para mais segurança de tipo
	pessoasPorCidade := pessoasStream.GroupByString(func(p Pessoa) string {
		return p.Cidade
	})

	fmt.Println("\nPessoas agrupadas por cidade:")
	for cidade, pessoasNaCidade := range pessoasPorCidade {
		fmt.Printf("  %s:\n", cidade)
		for _, p := range pessoasNaCidade {
			fmt.Printf("    - %s (%d anos)\n", p.Nome, p.Idade)
		}
	}

	// Exemplo agrupando por idade (número)
	pessoasPorIdade := pessoasStream.GroupBy(func(p Pessoa) interface{} {
		return p.Idade
	})

	fmt.Println("\nPessoas agrupadas por idade:")
	for idade, pessoasMesmaIdade := range pessoasPorIdade {
		fmt.Printf("  Idade %v:\n", idade)
		for _, p := range pessoasMesmaIdade {
			fmt.Printf("    - %s (%s)\n", p.Nome, p.Cidade)
		}
	}
}

// Exemplo 13: GroupByWithValueMapper - agrupar e transformar elementos
func exemploGroupByWithValueMapper() {
	fmt.Println("\n13. GroupBy com transformação de valores:")

	pessoasStream := stream.NewStream(pessoas)

	// Usando o método específico para string->string para mais segurança de tipo
	nomesPorCidade := pessoasStream.GroupByStringToString(
		func(p Pessoa) string { return p.Cidade }, // keyMapper - extrai a cidade
		func(p Pessoa) string { return p.Nome },   // valueMapper - extrai o nome
	)

	fmt.Println("Nomes agrupados por cidade:")
	for cidade, nomes := range nomesPorCidade {
		fmt.Printf("  %s: %v\n", cidade, nomes)
	}

	// Usando o método genérico GroupByAndTransform
	pessoasPorFaixaEtaria := pessoasStream.GroupByAndTransform(
		func(p Pessoa) interface{} { // keyMapper - determina a faixa etária
			if p.Idade < 30 {
				return "Jovem"
			}
			return "Adulto"
		},
		func(p Pessoa) interface{} { // valueMapper - formata a informação
			return fmt.Sprintf("%s (%d) - %s", p.Nome, p.Idade, p.Cidade)
		},
	)

	fmt.Println("\nPessoas agrupadas por faixa etária:")
	for faixa, info := range pessoasPorFaixaEtaria {
		fmt.Printf("  %v:\n", faixa)
		for _, p := range info {
			fmt.Printf("    - %v\n", p)
		}
	}

	// Agrupando por múltiplos critérios usando uma struct como chave
	type ChaveAgrupamento struct {
		Cidade string
		Jovem  bool
	}

	infoAgrupada := pessoasStream.GroupByAndTransform(
		func(p Pessoa) interface{} {
			return ChaveAgrupamento{
				Cidade: p.Cidade,
				Jovem:  p.Idade < 30,
			}
		},
		func(p Pessoa) interface{} {
			return p.Nome
		},
	)

	fmt.Println("\nPessoas agrupadas por cidade e faixa etária:")
	for chave, nomes := range infoAgrupada {
		info := chave.(ChaveAgrupamento)
		faixaEtaria := "Adulto"
		if info.Jovem {
			faixaEtaria = "Jovem"
		}
		fmt.Printf("  %s (%s): %v\n", info.Cidade, faixaEtaria, nomes)
	}
}

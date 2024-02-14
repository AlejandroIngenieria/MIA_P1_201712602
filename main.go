package main

import (
	"P1/Analyzer"
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("Ingrese un comando: ")
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error al leer la entrada:", err)
			continue
		}

		input = strings.TrimSpace(input)     //quitamos el salto de linea
		lowerInput := strings.ToLower(input) //quitamos el problema de mayusculas/minisculas
		//Para llamar una funcion desde otro archivo este debe ir en mayuscula al inicio
		analyzer.Command(lowerInput)
	}
}

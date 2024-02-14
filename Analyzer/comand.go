package analyzer

import (
	"P1/Functions"
	"flag"
	"fmt"
	"strings"
)

func Command(input string) {

	switch {
	case strings.HasPrefix(input, "mkdisk"):
		handleMkDiskCommand(input)
	case strings.HasPrefix(input, "rmdisk"):
		handleRmDiskCommand(input)
	case strings.HasPrefix(input, "fdisk"):
		handleFDiskCommand(input)
	default:
		fmt.Println("Comando no reconocido:", input)
	}
}

var (
	size        = flag.Int("size", 0, "Tamaño")
	fit         = flag.String("fit", "f", "Ajuste")
	unit        = flag.String("unit", "m", "Unidad")
	tipe        = flag.String("type", "p", "Tipo")
	driveletter = flag.String("driveletter", "", "Busqueda")
	name        = flag.String("name", "", "Nombre")
	delete      = flag.String("delete", "", "Eliminar")
	add         = flag.String("add", "", "Añadir")
)

func handleMkDiskCommand(input string) {

	flag.Parse()
	functions.ProcessMKDISK(input, size, fit, unit)

	// validate fit equals to b/w/f
	if *fit != "b" && *fit != "w" && *fit != "f" {
		fmt.Println("Error: Fit must be b, w or f")
		return
	}

	// validate size > 0
	if *size <= 0 {
		fmt.Println("Error: Size must be greater than 0")
		return
	}

	// validate unit equals to k/m
	if *unit != "k" && *unit != "m" {
		fmt.Println("Error: Unit must be k or m")
		return
	}

	// Print the values of the flags
	fmt.Println("Size:", *size)
	fmt.Println("Fit:", *fit)
	fmt.Println("Unit:", *unit)

}

func handleRmDiskCommand(input string) {
	flag.Parse()
	functions.ProcessMRDISK(input, driveletter)

	if !functions.ValidDriveLetter(*driveletter) {
		fmt.Println("Error: DriveLetter must be a letter")
		return
	}

	// Print the values of the flags
	fmt.Println("Driveletter:", *driveletter)
}

func handleFDiskCommand(input string) {
	flag.Parse()
	functions.ProcessFDISK(input, size, driveletter, name, unit, tipe, fit, delete)
}

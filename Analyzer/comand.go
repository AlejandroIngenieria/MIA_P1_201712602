package analyzer_test

import (
	"P1/Functions"
	"flag"
	"fmt"
	"strings"
)

func Command(input string) {
	comando := input
	input = strings.ToLower(input)
	switch {
	case strings.HasPrefix(input, "mkdisk"):
		handleMKDISKCommand(comando)
	case strings.HasPrefix(input, "rmdisk"):
		handleRMDISKCommand(comando)
	case strings.HasPrefix(input, "fdisk"):
		handleFDISKCommand(comando)
	case strings.HasPrefix(input, "mount"):
		handleMOUNTCommand(comando)
	case strings.HasPrefix(input, "unmount"):
		handleUNMOUNTCommand(comando)
	case strings.HasPrefix(input, "mkfs"):
		handleMKFSCommand(comando)
	case strings.HasPrefix(input, "login"):
		handleLOGINCommand(comando)
	case strings.HasPrefix(input, "logout"):
		handleLOGOUTCommand(comando)
	case strings.HasPrefix(input, "mkgrp"):
		handleMKGRPCommand(comando)
	case strings.HasPrefix(input, "rmgrp"):
		handleRMGRPCommand(comando)
	case strings.HasPrefix(input, "mkusr"):
		handleMKUSRCommand(comando)
	case strings.HasPrefix(input, "rmusr"):
		handleRMUSRCommand(comando)
	case strings.HasPrefix(input, "mkfile"):
		handleMKFILECommand(comando)
	case strings.HasPrefix(input, "cat"):
		handleCATCommand(comando)
	case strings.HasPrefix(input, "remove"):
		handleREMOVECommand(comando)
	case strings.HasPrefix(input, "edit"):
		handleEDITCommand(comando)
	case strings.HasPrefix(input, "rename"):
		handleRENAMECommand(comando)
	case strings.HasPrefix(input, "mkdir"):
		handleMKDIRCommand(comando)
	case strings.HasPrefix(input, "copy"):
		handleCOPYCommand(comando)
	case strings.HasPrefix(input, "move"):
		handleMOVECommand(comando)
	case strings.HasPrefix(input, "find"):
		handleFINDCommand(comando)
	case strings.HasPrefix(input, "chown"):
		handleCHOWNCommand(comando)
	case strings.HasPrefix(input, "chgrp"):
		handleCHGRPCommand(comando)
	case strings.HasPrefix(input, "chmod"):
		handleCHMODCommand(comando)
	case strings.HasPrefix(input, "pause"):
		handlePAUSECommand(comando)
	case strings.HasPrefix(input, "execute"):
		handleEXECUTECommand(comando)
	case strings.HasPrefix(input, "#"):
	default:
		fmt.Println("Comando no reconocido:", input)
	}
}

var (
	size        = flag.Int("size", 0, "Tamaño")
	fit         = flag.String("fit", "ff", "Ajuste")
	unit        = flag.String("unit", "m", "Unidad")
	type_       = flag.String("type", "p", "Tipo")
	driveletter = flag.String("driveletter", "", "Busqueda")
	name        = flag.String("name", "", "Nombre")
	delete      = flag.String("delete", "", "Eliminar")
	add         = flag.String("add", "", "Añadir/Quitar")
	path        = flag.String("path", "", "Directorio")
)

func handleMKDISKCommand(input string) {

	flag.Parse()
	functions_test.ProcessMKDISK(input, size, fit, unit)

	// validate size > 0
	if *size <= 0 {
		fmt.Println("Error: Size must be greater than 0")
		return
	}

	// validate fit equals to b/w/f
	if *fit != "bf" && *fit != "ff" && *fit != "wf" {
		fmt.Println("Error: Fit must be (bf/ff/wf)")
		return
	}

	// validate unit equals to k/m
	if *unit != "k" && *unit != "m" {
		fmt.Println("Error: Unit must be (k/m)")
		return
	}

	// Create the file
	functions_test.CreateBinFile(size, fit, unit)
	*size = 0
	*fit = "f"
	*unit = "m"
}

func handleRMDISKCommand(input string) {
	flag.Parse()
	functions_test.ProcessRMDISK(input, driveletter)
	// validate driveletter be a letter and not empty
	if !functions_test.ValidDriveLetter(*driveletter) {
		fmt.Println("Error: DriveLetter must be a letter")
		return
	} else if len(*driveletter) == 0 {
		fmt.Println("Error: DriveLetter cannot be empty")
		return
	}

	functions_test.DeleteBinFile(driveletter)
	*driveletter = ""
}

func handleFDISKCommand(input string) {
	flag.Parse()
	functions_test.ProcessFDISK(input, size, driveletter, name, unit, type_, fit, delete, add, path)

	// validate size > 0
	if *size <= 0 {
		fmt.Println("Error: Size must be greater than 0")
		return
	}

	// validate driveletter be a letter and not empty
	if !functions_test.ValidDriveLetter(*driveletter) {
		fmt.Println("Error: DriveLetter must be a letter")
		return
	} else if len(*driveletter) == 0 {
		fmt.Println("Error: DriveLetter cannot be empty")
		return
	}

	// validate fit equals to b/w/f
	if *fit != "bf" && *fit != "ff" && *fit != "wf" {
		fmt.Println("Error: Fit must be (BF/FF/WF)")
		return
	}

	// validate unit equals to b/k/m
	if *unit != "b" && *unit != "k" && *unit != "m" {
		fmt.Println("Error: Unit must be (B/K/M)")
		return
	}

	if *type_ != "p" && *type_ != "e" && *type_ != "l" {
		fmt.Println("Error: Type must be (P/E/L)")
		return
	}

	if *delete == "" || *name == "" || *path == "" {
		fmt.Println("Error: Delete -> remember that needs name and path")
		return
	}
}

func handleMOUNTCommand(input string) {
	panic("unimplemented")
}

func handleUNMOUNTCommand(input string) {
	panic("unimplemented")
}

func handleMKFSCommand(input string) {
	panic("unimplemented")
}

func handleLOGINCommand(input string) {
	panic("unimplemented")
}

func handleLOGOUTCommand(input string) {
	panic("unimplemented")
}

func handleMKGRPCommand(input string) {
	panic("unimplemented")
}

func handleRMGRPCommand(input string) {
	panic("unimplemented")
}

func handleMKUSRCommand(input string) {
	panic("unimplemented")
}

func handleRMUSRCommand(input string) {
	panic("unimplemented")
}

func handleMKFILECommand(input string) {
	panic("unimplemented")
}

func handleCATCommand(input string) {
	panic("unimplemented")
}

func handleREMOVECommand(input string) {
	panic("unimplemented")
}

func handleEDITCommand(input string) {
	panic("unimplemented")
}

func handleRENAMECommand(input string) {
	panic("unimplemented")
}

func handleMKDIRCommand(input string) {
	panic("unimplemented")
}

func handleCOPYCommand(input string) {
	panic("unimplemented")
}

func handleMOVECommand(input string) {
	panic("unimplemented")
}

func handleFINDCommand(input string) {
	panic("unimplemented")
}

func handleCHOWNCommand(input string) {
	panic("unimplemented")
}

func handleCHGRPCommand(input string) {
	panic("unimplemented")
}

func handleCHMODCommand(input string) {
	panic("unimplemented")
}

func handlePAUSECommand(comando string) {
	panic("unimplemented")
}

func handleEXECUTECommand(input string) {
	panic("unimplemented")
}

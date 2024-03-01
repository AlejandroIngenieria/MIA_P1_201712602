package analyzer_test

import (
	"P1/Functions"
	"P1/Utilities"
	"bufio"
	"flag"
	"fmt"
	"strings"
)

func Command(input string) {

	// Verificar si el input está vacío
	if input == "" {
		return // No hacer nada si el input está vacío
	}

	comando := input
	input = strings.ToLower(input)
	switch {
	case strings.HasPrefix(input, "mkdisk"):
		fmt.Println(input)
		handleMKDISKCommand(comando)
	case strings.HasPrefix(input, "rmdisk"):
		fmt.Println(input)
		handleRMDISKCommand(comando)
	case strings.HasPrefix(input, "fdisk"):
		fmt.Println(input)
		handleFDISKCommand(comando)
	case strings.HasPrefix(input, "mount"):
		fmt.Println(input)
		handleMOUNTCommand(comando)
	case strings.HasPrefix(input, "unmount"):
		fmt.Println(input)
		handleUNMOUNTCommand(comando)
	case strings.HasPrefix(input, "mkfs"):
		fmt.Println(input)
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
	case strings.HasPrefix(input, "rep"):
		handleREPCommand(comando)
	case strings.HasPrefix(input, "#"):
	default:
		fmt.Println("Comando no reconocido:", input)
	}
}

var (
	size        = flag.Int("size", 0, "Tamaño")
	fit         = flag.String("fit", "", "Ajuste")
	unit        = flag.String("unit", "", "Unidad")
	type_       = flag.String("type", "", "Tipo")
	driveletter = flag.String("driveletter", "", "Busqueda")
	name        = flag.String("name", "", "Nombre")
	delete      = flag.String("delete", "", "Eliminar")
	add         = flag.Int("add", 0, "Añadir/Quitar")
	path        = flag.String("path", "", "Directorio")
	id          = flag.String("id", "", "ID")
	fs          = flag.String("", "", "")
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
	if *fit != "b" && *fit != "f" && *fit != "w" {
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
	*fit = ""
	*unit = ""
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

	//Obligatorio cuando no existe la particion
	// validate size > 0
	if *size <= 0 && *delete != "full" && *add == 0 {
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
	if *fit != "b" && *fit != "f" && *fit != "w" {
		fmt.Println("Error: Fit must be (BF/FF/WF)")
		return
	}

	// validate unit equals to b/k/m
	if *unit != "b" && *unit != "k" && *unit != "m" {
		fmt.Println("Error: Unit must be (B/K/M)")
		return
	}

	//println("ADD", *add)
	// validate type equals to P/E/L
	if *type_ != "p" && *type_ != "e" && *type_ != "l" && *delete != "full" && *add == 0 {
		fmt.Println("Error: Type must be (P/E/L)")
		return
	}

	if *delete != "" {
		if *delete != "full" {
			fmt.Println("Error: Delete must be full")
			return
		}
		if *name == "" && *path == "" {
			println("Error: you need path and name to delete")
			return
		}
	}

	functions_test.CRUD_Partitions(size, driveletter, name, unit, type_, fit, delete, add, path)
	*size = 0
	*driveletter = ""
	*name = ""
	*unit = ""
	*type_ = ""
	*fit = ""
	*delete = ""
	*add = 0
	*path = ""
}

func handleMOUNTCommand(input string) {
	flag.Parse()
	functions_test.ProcessMOUNT(input, driveletter, name)

	// validate driveletter be a letter and not empty
	if !functions_test.ValidDriveLetter(*driveletter) {
		fmt.Println("Error: DriveLetter must be a letter")
		return
	} else if len(*driveletter) == 0 {
		fmt.Println("Error: DriveLetter cannot be empty")
		return
	}

	functions_test.MountPartition(driveletter, name)
	*driveletter = ""
	*name = ""

}

func handleUNMOUNTCommand(input string) {
	flag.Parse()
	functions_test.ProcessUNMOUNT(input, id)

	functions_test.UNMOUNT_Partition(id)
}

func handleMKFSCommand(input string) {
	flag.Parse()
	functions_test.ProcessMKFS(input, id, type_, fs)
	
	if *id == "" {
		println("Error: id cannot be empty")
	}
	
	if *fs != "2fs" && *fs != "3fs"{
		println("Error: fs must be 2fs or 3fs")
	}

	
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

func handlePAUSECommand(input string) {
	panic("unimplemented")
}

func handleEXECUTECommand(input string) {
	flag.Parse()
	functions_test.ProcessExecute(input, path)
	if *path == "" {
		fmt.Println("Error: Path cannot be empty")
		return
	}
	// Open bin file
	file, err := utilities_test.OpenFile(*path)
	if err != nil {
		return
	}

	// Close bin file
	defer file.Close()

	// Crea un nuevo scanner para leer el archivo
	scanner := bufio.NewScanner(file)

	// Itera sobre cada línea del archivo
	for scanner.Scan() {
		linea := scanner.Text() // Lee la línea actual
		//fmt.Println(linea)
		Command(linea)
	}

	// Verifica si hubo algún error durante la lectura
	if err := scanner.Err(); err != nil {
		fmt.Println("Error al leer el archivo:", err)
	}
	*path = ""
}

func handleREPCommand(comando string) {
	panic("unimplemented")
}

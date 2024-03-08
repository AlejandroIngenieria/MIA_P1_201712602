package functions_test

import (
	"P1/Global"
	"P1/Structs"
	"P1/Utilities"
	"encoding/binary"
	"fmt"
	"regexp"
	"strings"
)

var (
	session      = false
	usuario      = Global.UserInfo{}
	groupCounter = 1
	blockIndex   = int32(0)
	// userCounter = 1
	letra = ""
	ID    = ""
)

//?                    ADMINISTRACION DE USUARIOS Y GRUPOS
/* -------------------------------------------------------------------------- */
/*                                COMANDO LOGIN                               */
/* -------------------------------------------------------------------------- */
func ProcessLOGIN(input string, user *string, pass *string, id *string) {
	re := regexp.MustCompile(`-(\w+)=("[^"]+"|\S+)`)

	matches := re.FindAllStringSubmatch(input, -1)

	for _, match := range matches {
		flagName := match[1]
		flagValue := match[2]

		// Delete quotes if they are present in the value
		flagValue = strings.Trim(flagValue, "\"")

		switch flagName {
		case "user":
			*user = flagValue
		case "pass":
			*pass = flagValue
		case "id":
			*id = flagValue
		default:
			fmt.Println("Error: Flag not found: " + flagName)
		}
	}
}

func LOGIN(user *string, pass *string, id *string) {
	fmt.Println("User:", *user)
	fmt.Println("Pass:", *pass)
	fmt.Println("Id:", *id)

	letra = string((*id)[0])
	ID = *id

	filepath := "./Disks/" + letra + ".dsk"
	file, err := utilities_test.OpenFile(filepath)
	if err != nil {
		return
	}

	var TempMBR structs_test.MBR
	// Read object from bin file
	if err := utilities_test.ReadObject(file, &TempMBR, 0); err != nil {
		return
	}

	var index int = -1
	// Iterate over the partitions
	for i := 0; i < 4; i++ {
		if TempMBR.Mbr_particion[i].Part_size != 0 {
			if strings.Contains(string(TempMBR.Mbr_particion[i].Part_id[:]), *id) {
				fmt.Println("Partition found")
				if strings.Contains(string(TempMBR.Mbr_particion[i].Part_status[:]), "1") {
					fmt.Println("Partition is mounted")
					index = i
				} else {
					fmt.Println("Partition is not mounted")
					return
				}
				break
			}
		}
	}

	if index != -1 {
		structs_test.PrintPartition(TempMBR.Mbr_particion[index])
	} else {
		fmt.Println("Partition not found")
		return
	}

	// INICIA LA LLAMADA AL SUPERBLOQUE DE LA PARTICION

	var tempSuperblock structs_test.Superblock
	// Read object from bin file
	if err := utilities_test.ReadObject(file, &tempSuperblock, int64(TempMBR.Mbr_particion[index].Part_start)); err != nil {
		return
	}

	// initSearch /users.txt -> regresa no Inodo
	// initSearch -> 1

	indexInode := int32(1)

	var crrInode structs_test.Inode
	// Read object from bin file
	if err := utilities_test.ReadObject(file, &crrInode, int64(tempSuperblock.S_inode_start+indexInode*int32(binary.Size(structs_test.Inode{})))); err != nil {
		return
	}

	// getInodeFileData -> Iterate the I_Block n concat the data

	var Fileblock structs_test.Fileblock
	// Read object from bin file
	if err := utilities_test.ReadObject(file, &Fileblock, int64(tempSuperblock.S_block_start+crrInode.I_block[0]*int32(binary.Size(structs_test.Fileblock{})))); err != nil {
		return
	}

	fmt.Println("Fileblock------------")
	data := string(Fileblock.B_content[:])
	// Dividir la cadena en líneas
	lines := strings.Split(data, "\n")

	// Iterar a través de las líneas
	for _, line := range lines {
		// Imprimir cada línea
		fmt.Println(line)
		items := strings.Split(line, ",")
		if len(items) > 3 {
			usuario.Nombre = items[len(items)-2]
			usuario.ID = items[0]
			//Global.PrintUser(usuario)
			if usuario.Nombre == *user {
				session = true
				break
			}
		}
	}

	Global.PrintUser(usuario)

	if session {
		fmt.Println("SESION INICIADA " + usuario.Nombre)
	} else {
		println("Error: no se logro iniciar sesion")
		usuario.ID = ""
		usuario.Nombre = ""
	}
	// Print object
	fmt.Println("Inode", crrInode.I_block)

	// Close bin file
	defer file.Close()
}

/* -------------------------------------------------------------------------- */
/*                               COMANDO LOGOUT                               */
/* -------------------------------------------------------------------------- */
func ProcessLOGOUT() {
	if session {
		println("Se ha cerrado la sesion")
		session = false
		return
	}
	println("Error: no hay una sesion activa")
}

/* -------------------------------------------------------------------------- */
/*                                COMANDO MKGRP                               */
/* -------------------------------------------------------------------------- */
func ProcessMKGRP(input string, name *string) {
	if usuario.Nombre == "root" {
		re := regexp.MustCompile(`-(\w+)=("[^"]+"|\S+)`)

		matches := re.FindAllStringSubmatch(input, -1)

		for _, match := range matches {
			flagName := match[1]
			flagValue := match[2]

			// Delete quotes if they are present in the value
			flagValue = strings.Trim(flagValue, "\"")

			switch flagName {
			case "name":
				*name = flagValue
			default:
				fmt.Println("Error: Flag not found: " + flagName)
			}
		}
	} else {
		println("Error: solo el usuario root puede acceder a este comando")
		return
	}
}

func MKGRP(name *string) {
	// Abrir el archivo del disco
	filepath := "./Disks/" + letra + ".dsk"
	file, err := utilities_test.OpenFile(filepath)
	if err != nil {
		fmt.Println("Error opening disk file:", err)
		return
	}
	defer file.Close()

	// Leer el MBR del disco
	var TempMBR structs_test.MBR
	if err := utilities_test.ReadObject(file, &TempMBR, 0); err != nil {
		fmt.Println("Error reading MBR:", err)
		return
	}

	// Encontrar la partición adecuada en el MBR
	index := -1
	for i := 0; i < 4; i++ {
		if TempMBR.Mbr_particion[i].Part_size != 0 && strings.Contains(string(TempMBR.Mbr_particion[i].Part_id[:]), ID) {
			index = i
			break
		}
	}
	if index == -1 {
		fmt.Println("Partition not found")
		return
	}

	// Leer el superbloque de la partición
	var tempSuperblock structs_test.Superblock
	if err := utilities_test.ReadObject(file, &tempSuperblock, int64(TempMBR.Mbr_particion[index].Part_start)); err != nil {
		fmt.Println("Error reading superblock:", err)
		return
	}

	// Leer el inode correspondiente al bloque de archivos
	indexInode := int32(1)
	var crrInode structs_test.Inode
	if err := utilities_test.ReadObject(file, &crrInode, int64(tempSuperblock.S_inode_start+indexInode*int32(binary.Size(structs_test.Inode{})))); err != nil {
		fmt.Println("Error reading inode:", err)
		return
	}

	// Leer el contenido actual del Fileblock
	var Fileblock structs_test.Fileblock
	if err := utilities_test.ReadObject(file, &Fileblock, int64(tempSuperblock.S_block_start+crrInode.I_block[0]*int32(binary.Size(structs_test.Fileblock{})))); err != nil {
		fmt.Println("Error reading Fileblock:", err)
		return
	}

	data := string(Fileblock.B_content[:])
	// Dividir la cadena en líneas
	lines := strings.Split(data, "\n")

	//------------------------------------------------------------------ Iterar a través de las líneas
	for _, line := range lines {
		// Imprimir cada línea
		fmt.Println(line)
		items := strings.Split(line, ",")
		if len(items) == 3 {
			if *name == items[2] {
				println("Error: nombre repetido")
				return
			}
		}
	}

	// Convertir B_content a string y añadir el nuevo grupo
	currentContent := strings.TrimRight(string(Fileblock.B_content[:]), "\x00")
	groupCounter++
	nuevoGrupo := fmt.Sprintf("%d,G,%s\n", groupCounter, *name)
	newContent := currentContent + nuevoGrupo

	// Asegurarse de que el nuevo contenido quepa dentro del Fileblock
	if len(newContent) > len(Fileblock.B_content) {
		for i := int32(0); i < 15; i++ {
			if crrInode.I_block[i] == -1 {
				crrInode.I_block[i] = 1
				blockIndex = i
				break
			}
		}
		fmt.Println("New content is too large to fit in Fileblock")
		return
	}

	// Copiar el nuevo contenido al Fileblock
	copy(Fileblock.B_content[:], newContent)

	// Escribir el Fileblock actualizado de nuevo en el disco
	if err := utilities_test.WriteObject(file, &Fileblock, int64(tempSuperblock.S_block_start+crrInode.I_block[0]*int32(binary.Size(structs_test.Fileblock{})))); err != nil {
		fmt.Println("Error writing Fileblock to disk:", err)
		return
	}

	// Leer y mostrar el contenido actualizado del Fileblock desde el disco
	var updatedFileblock structs_test.Fileblock
	if err := utilities_test.ReadObject(file, &updatedFileblock, int64(tempSuperblock.S_block_start+crrInode.I_block[0]*int32(binary.Size(structs_test.Fileblock{})))); err != nil {
		fmt.Println("Error reading updated Fileblock:", err)
		return
	}

	// Mostrar el contenido actualizado del Fileblock
	updatedContent := strings.TrimRight(string(updatedFileblock.B_content[:]), "\x00")
	fmt.Println("Updated Fileblock content:")
	fmt.Println(updatedContent)
}

/* -------------------------------------------------------------------------- */
/*                                COMANDO RMGRP                               */
/* -------------------------------------------------------------------------- */
func ProcessRMGRP(input string, name *string) {
	if usuario.Nombre == "root" {
		re := regexp.MustCompile(`-(\w+)=("[^"]+"|\S+)`)

		matches := re.FindAllStringSubmatch(input, -1)

		for _, match := range matches {
			flagName := match[1]
			flagValue := match[2]

			// Delete quotes if they are present in the value
			flagValue = strings.Trim(flagValue, "\"")

			switch flagName {
			case "name":
				*name = flagValue
			default:
				fmt.Println("Error: Flag not found: " + flagName)
			}
		}
	} else {
		println("Error: solo el usuario root puede acceder a este comando")
		return
	}
}

func RMGRP(name *string) {
	// Abrir el archivo del disco
	filepath := "./Disks/" + letra + ".dsk"
	file, err := utilities_test.OpenFile(filepath)
	if err != nil {
		fmt.Println("Error opening disk file:", err)
		return
	}
	defer file.Close()

	// Leer el MBR del disco
	var TempMBR structs_test.MBR
	if err := utilities_test.ReadObject(file, &TempMBR, 0); err != nil {
		fmt.Println("Error reading MBR:", err)
		return
	}

	// Encontrar la partición adecuada en el MBR
	index := -1
	for i := 0; i < 4; i++ {
		if TempMBR.Mbr_particion[i].Part_size != 0 && strings.Contains(string(TempMBR.Mbr_particion[i].Part_id[:]), ID) {
			index = i
			break
		}
	}
	if index == -1 {
		fmt.Println("Partition not found")
		return
	}

	// Leer el superbloque de la partición
	var tempSuperblock structs_test.Superblock
	if err := utilities_test.ReadObject(file, &tempSuperblock, int64(TempMBR.Mbr_particion[index].Part_start)); err != nil {
		fmt.Println("Error reading superblock:", err)
		return
	}

	// Leer el inode correspondiente al bloque de archivos
	indexInode := int32(1)
	var crrInode structs_test.Inode
	if err := utilities_test.ReadObject(file, &crrInode, int64(tempSuperblock.S_inode_start+indexInode*int32(binary.Size(structs_test.Inode{})))); err != nil {
		fmt.Println("Error reading inode:", err)
		return
	}

	// Leer el contenido actual del Fileblock
	var Fileblock structs_test.Fileblock
	if err := utilities_test.ReadObject(file, &Fileblock, int64(tempSuperblock.S_block_start+crrInode.I_block[0]*int32(binary.Size(structs_test.Fileblock{})))); err != nil {
		fmt.Println("Error reading Fileblock:", err)
		return
	}

	// Convertir B_content a string y buscar y eliminar el grupo especificado
	currentContent := strings.TrimRight(string(Fileblock.B_content[:]), "\x00")
	lines := strings.Split(currentContent, "\n")
	deleted := false
	for i, line := range lines {
		if strings.Contains(line, *name) {
			lines[i] = "0,G," + *name
			deleted = true
			break
		}
	}

	// Si el grupo no fue encontrado, mostrar un mensaje y salir
	if !deleted {
		fmt.Println("Group not found")
		return
	}

	// Reunir las líneas actualizadas en un nuevo contenido con saltos de línea
	newContent := strings.Join(lines, "\n")

	// Copiar el nuevo contenido al Fileblock
	copy(Fileblock.B_content[:], newContent)

	// Escribir el Fileblock actualizado de nuevo en el disco
	if err := utilities_test.WriteObject(file, &Fileblock, int64(tempSuperblock.S_block_start+crrInode.I_block[0]*int32(binary.Size(structs_test.Fileblock{})))); err != nil {
		fmt.Println("Error writing Fileblock to disk:", err)
		return
	}

	// Leer y mostrar el contenido actualizado del Fileblock desde el disco
	var updatedFileblock structs_test.Fileblock
	if err := utilities_test.ReadObject(file, &updatedFileblock, int64(tempSuperblock.S_block_start+crrInode.I_block[0]*int32(binary.Size(structs_test.Fileblock{})))); err != nil {
		fmt.Println("Error reading updated Fileblock:", err)
		return
	}

	// Mostrar el contenido actualizado del Fileblock
	updatedContent := strings.TrimRight(string(updatedFileblock.B_content[:]), "\x00")
	fmt.Println("Updated Fileblock content:")
	fmt.Println(updatedContent)
}

/* -------------------------------------------------------------------------- */
/*                                COMANDO MKUSR                               */
/* -------------------------------------------------------------------------- */
func ProcessMKUSR(input string, user *string, pass *string, grp *string) {
	if usuario.Nombre == "root" {
		re := regexp.MustCompile(`-(\w+)=("[^"]+"|\S+)`)

		matches := re.FindAllStringSubmatch(input, -1)

		for _, match := range matches {
			flagName := match[1]
			flagValue := match[2]

			// Delete quotes if they are present in the value
			flagValue = strings.Trim(flagValue, "\"")

			switch flagName {
			case "user":
				*user = flagValue
			case "pass":
				*pass = flagValue
			case "grp":
				*grp = flagValue
			default:
				fmt.Println("Error: Flag not found: " + flagName)
			}
		}
	} else {
		println("Error: solo el usuario root puede acceder a este comando")
		return
	}
}

func MKUSR(user *string, pass *string, grp *string) {
	// Abrir el archivo del disco
	filepath := "./Disks/" + letra + ".dsk"
	file, err := utilities_test.OpenFile(filepath)
	if err != nil {
		fmt.Println("Error opening disk file:", err)
		return
	}
	defer file.Close()

	// Leer el MBR del disco
	var TempMBR structs_test.MBR
	if err := utilities_test.ReadObject(file, &TempMBR, 0); err != nil {
		fmt.Println("Error reading MBR:", err)
		return
	}

	// Encontrar la partición adecuada en el MBR
	index := -1
	for i := 0; i < 4; i++ {
		if TempMBR.Mbr_particion[i].Part_size != 0 && strings.Contains(string(TempMBR.Mbr_particion[i].Part_id[:]), ID) {
			index = i
			break
		}
	}
	if index == -1 {
		fmt.Println("Partition not found")
		return
	}

	// Leer el superbloque de la partición
	var tempSuperblock structs_test.Superblock
	if err := utilities_test.ReadObject(file, &tempSuperblock, int64(TempMBR.Mbr_particion[index].Part_start)); err != nil {
		fmt.Println("Error reading superblock:", err)
		return
	}

	// Leer el inode correspondiente al bloque de archivos
	indexInode := int32(1)
	var crrInode structs_test.Inode
	if err := utilities_test.ReadObject(file, &crrInode, int64(tempSuperblock.S_inode_start+indexInode*int32(binary.Size(structs_test.Inode{})))); err != nil {
		fmt.Println("Error reading inode:", err)
		return
	}

	// Leer el contenido actual del Fileblock
	// Leer el contenido actual del Fileblock
	var Fileblock structs_test.Fileblock
	if err := utilities_test.ReadObject(file, &Fileblock, int64(tempSuperblock.S_block_start+crrInode.I_block[0]*int32(binary.Size(structs_test.Fileblock{})))); err != nil {
		fmt.Println("Error reading Fileblock:", err)
		return
	}

	data := string(Fileblock.B_content[:])
	// Dividir la cadena en líneas
	lines := strings.Split(data, "\n")

	// Buscar el grupo especificado en el Fileblock
	groupFound := false
	var newFileblockContent string

	//------------------------------------------------------------------ Iterar a través de las líneas
	newUserLine := ""
	for _, line := range lines {
		// Imprimir cada línea
		fmt.Println(line)
		items := strings.Split(line, ",")
		if len(items) == 3 {
			fmt.Println("items[2]->"+items[2])
			if *grp == items[2] {
				groupFound = true
				newUserLine = fmt.Sprintf("%s,G,%s,%s,%s\n", items[0], *grp, *user, *pass)
				newFileblockContent += newUserLine
				break
			}
		}
	}

	// Si el grupo no fue encontrado, mostrar un mensaje de error y salir
	if !groupFound {
		fmt.Println("Group", *grp, "not found")
		return
	}

	// Convertir B_content a string y añadir el nuevo grupo
	currentContent := strings.TrimRight(string(Fileblock.B_content[:]), "\x00")
	newContent := currentContent + newUserLine

	// Copiar el nuevo contenido al Fileblock
	copy(Fileblock.B_content[:], newContent)

	// Escribir el Fileblock actualizado de nuevo en el disco
	if err := utilities_test.WriteObject(file, &Fileblock, int64(tempSuperblock.S_block_start+crrInode.I_block[0]*int32(binary.Size(structs_test.Fileblock{})))); err != nil {
		fmt.Println("Error writing Fileblock to disk:", err)
		return
	}

	// Leer y mostrar el contenido actualizado del Fileblock desde el disco
	var updatedFileblock structs_test.Fileblock
	if err := utilities_test.ReadObject(file, &updatedFileblock, int64(tempSuperblock.S_block_start+crrInode.I_block[0]*int32(binary.Size(structs_test.Fileblock{})))); err != nil {
		fmt.Println("Error reading updated Fileblock:", err)
		return
	}

	// Mostrar el contenido actualizado del Fileblock
	updatedContent := strings.TrimRight(string(updatedFileblock.B_content[:]), "\x00")
	fmt.Println("Updated Fileblock content:")
	fmt.Println(updatedContent)
}

/* -------------------------------------------------------------------------- */
/*                                COMANDO RMUSR                               */
/* -------------------------------------------------------------------------- */
func ProcessRMUSR(input string, user *string) {
	if usuario.Nombre == "root" {
		re := regexp.MustCompile(`-(\w+)=("[^"]+"|\S+)`)

		matches := re.FindAllStringSubmatch(input, -1)

		for _, match := range matches {
			flagName := match[1]
			flagValue := match[2]

			// Delete quotes if they are present in the value
			flagValue = strings.Trim(flagValue, "\"")

			switch flagName {
			case "user":
				*user = flagValue
			default:
				fmt.Println("Error: Flag not found: " + flagName)
			}
		}
	} else {
		println("Error: solo el usuario root puede acceder a este comando")
		return
	}
}

func RMUSR(user *string) {
	// Abrir el archivo del disco
	filepath := "./Disks/" + letra + ".dsk"
	file, err := utilities_test.OpenFile(filepath)
	if err != nil {
		fmt.Println("Error opening disk file:", err)
		return
	}
	defer file.Close()

	// Leer el MBR del disco
	var TempMBR structs_test.MBR
	if err := utilities_test.ReadObject(file, &TempMBR, 0); err != nil {
		fmt.Println("Error reading MBR:", err)
		return
	}

	// Encontrar la partición adecuada en el MBR
	index := -1
	for i := 0; i < 4; i++ {
		if TempMBR.Mbr_particion[i].Part_size != 0 && strings.Contains(string(TempMBR.Mbr_particion[i].Part_id[:]), ID) {
			index = i
			break
		}
	}
	if index == -1 {
		fmt.Println("Partition not found")
		return
	}

	// Leer el superbloque de la partición
	var tempSuperblock structs_test.Superblock
	if err := utilities_test.ReadObject(file, &tempSuperblock, int64(TempMBR.Mbr_particion[index].Part_start)); err != nil {
		fmt.Println("Error reading superblock:", err)
		return
	}

	// Leer el inode correspondiente al bloque de archivos
	indexInode := int32(1)
	var crrInode structs_test.Inode
	if err := utilities_test.ReadObject(file, &crrInode, int64(tempSuperblock.S_inode_start+indexInode*int32(binary.Size(structs_test.Inode{})))); err != nil {
		fmt.Println("Error reading inode:", err)
		return
	}

	// Leer el contenido actual del Fileblock
	var Fileblock structs_test.Fileblock
	if err := utilities_test.ReadObject(file, &Fileblock, int64(tempSuperblock.S_block_start+crrInode.I_block[0]*int32(binary.Size(structs_test.Fileblock{})))); err != nil {
		fmt.Println("Error reading Fileblock:", err)
		return
	}

	// Convertir B_content a string
	currentContent := strings.TrimRight(string(Fileblock.B_content[:]), "\x00")

	// Buscar el usuario especificado en el Fileblock
	userFound := false
	var newFileblockContent string
	for _, line := range strings.Split(currentContent, "\n") {
		if strings.Contains(line, *user) {
			userFound = true
			// Cambiar el número de inicio de la línea por 0
			lineItems := strings.Split(line, ",")
			if len(lineItems) > 0 {
				lineItems[0] = "0"
			}
			newLine := strings.Join(lineItems, ",") + "\n"
			newFileblockContent += newLine
		} else {
			newFileblockContent += line + "\n"
		}
	}

	// Si el usuario no fue encontrado, mostrar un mensaje de error y salir
	if !userFound {
		fmt.Println("User", *user, "not found")
		return
	}

	// Asegurarse de que el nuevo contenido quepa dentro del Fileblock
	if len(newFileblockContent) > len(Fileblock.B_content) {
		fmt.Println("New content is too large to fit in Fileblock")
		return
	}

	// Copiar el nuevo contenido al Fileblock
	copy(Fileblock.B_content[:], newFileblockContent)

	// Escribir el Fileblock actualizado en el disco
	if err := utilities_test.WriteObject(file, &Fileblock, int64(tempSuperblock.S_block_start+crrInode.I_block[0]*int32(binary.Size(structs_test.Fileblock{})))); err != nil {
		fmt.Println("Error writing Fileblock to disk:", err)
		return
	}

	fmt.Println("User", *user, "removed")
}
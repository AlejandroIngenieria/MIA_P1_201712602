package functions_test

import (
	"P1/Structs"
	"P1/Utilities"
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
	"math/rand"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var fileCounter int = 0

//?                          APLICACION DE COMANDOS
/* -------------------------------------------------------------------------- */
/*                               COMANDO MKDISK                               */
/* -------------------------------------------------------------------------- */
func ProcessMKDISK(input string, size *int, fit *string, unit *string) {
	input = strings.ToLower(input) //quitamos el problema de mayusculas/minisculas
	re := regexp.MustCompile(`-(\w+)=("[^"]+"|\S+)`)

	matches := re.FindAllStringSubmatch(input, -1)

	for _, match := range matches {
		flagName := match[1]
		flagValue := match[2]

		// Delete quotes if they are present in the value
		flagValue = strings.Trim(flagValue, "\"")

		switch flagName {
		case "size":
			sizeValue := 0
			fmt.Sscanf(flagValue, "%d", &sizeValue)
			*size = sizeValue
		case "fit":
			flagValue = flagValue[:1]
			*fit = flagValue
		case "unit":
			*unit = flagValue
		default:
			fmt.Println("Error: Flag not found")
		}
	}

	if *fit == "" {
		*fit = "f"
	}
	if *unit == "" {
		*unit = "m"
	}
}

func CreateBinFile(size *int, fit *string, unit *string) {
	// Letras del alfabeto
	letters := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	// Contador para archivos
	if *unit == "k" {
		*size = *size * 1024
	} else {
		*size = *size * 1024 * 1024
	}

	if err := createFile(fmt.Sprintf("./Disks/%c.dsk", letters[fileCounter]), *size, *fit); err != nil {
		fmt.Printf("Error al crear archivo de %d %s: %e", *size, *unit, err)
		return
	}

	// Incrementar el contador
	fileCounter++
}

func createFile(filename string, size int, fit string) error {
	// Crear el archivo con el nombre proporcionado
	err := utilities_test.CreateFile(filename)
	if err != nil {
		return err
	}

	// Open bin file
	file, err := utilities_test.OpenFile(filename)
	if err != nil {
		return nil
	}

	// Write 0 binary data to the file

	// create array of byte(0)
	for i := 0; i < size; i++ {
		err := utilities_test.WriteObject(file, byte(0), int64(i))
		if err != nil {
			fmt.Println("Error: ", err)
		}
	}

	// Obtener la hora actual
	currentTime := time.Now()
	// Formatear la hora actual como una cadena
	timeString := currentTime.Format("2006-01-02 15:04:05")
	//Asignacion de datos al MBR
	var TempMBR structs_test.MBR
	TempMBR.Mbr_tamano = int32(size)
	copy(TempMBR.Mbr_fecha_creacion[:], []byte(timeString))
	TempMBR.Mbr_dsk_signature = int32(GenerateUniqueID())
	copy(TempMBR.Dsk_fit[:], fit)

	// Write object in bin file
	if err := utilities_test.WriteObject(file, TempMBR, 0); err != nil {
		return nil
	}

	var mbr structs_test.MBR
	// Read object from bin file
	if err := utilities_test.ReadObject(file, &mbr, 0); err != nil {
		return nil
	}

	// Print object
	// structs_test.PrintMBR(TempMBR)

	// Close bin file
	defer file.Close()

	defer file.Close()

	return nil
}

/* -------------------------------------------------------------------------- */
/*                               COMANDO RMDISK                               */
/* -------------------------------------------------------------------------- */
func ProcessRMDISK(input string, driveletter *string) {
	input = strings.ToLower(input) //quitamos el problema de mayusculas/minisculas
	re := regexp.MustCompile(`-(\w+)=("[^"]+"|\S+)`)
	match := re.FindStringSubmatch(input)
	if len(match) != 3 {
		fmt.Println("Comando rmdisk incorrecto")
		return
	}
	flagValue := match[2]
	*driveletter = flagValue
}

func DeleteBinFile(driveletter *string) {
	// Archivo a buscar y eliminar
	*driveletter = strings.ToUpper(*driveletter)
	filename := "./Disks/" + *driveletter + ".dsk"
	// Buscar el archivo
	if _, err := os.Stat(filename); err == nil {
		// El archivo existe, intenta eliminarlo

		fmt.Print("Desea eliminar el archivo " + *driveletter + ".dsk(y/n)?")
		var input string
		fmt.Print("Ingrese 'y' para continuar o 'n' para cancelar: ")
		_, err := fmt.Scanln(&input)
		if err != nil {
			fmt.Println("Error al leer la entrada:", err)
			return
		}

		if input == "y" {
			if err := os.Remove(filename); err != nil {
				fmt.Println("Error al eliminar el archivo:", err)
				return
			}
		} else {
			fmt.Println("Entrada no válida.")
		}

	} else if os.IsNotExist(err) {
		// El archivo no existe
		fmt.Printf("El archivo %s no existe.\n", filename)
	} else {
		// Otro error ocurrió
		fmt.Println("Error al verificar la existencia del archivo:", err)
	}
}

/* -------------------------------------------------------------------------- */
/*                                COMANDO FDISK                               */
/* -------------------------------------------------------------------------- */
func ProcessFDISK(input string, size *int, driveletter *string, name *string, unit *string, type_ *string, fit *string, delete *string, add *int, path *string) {
	input = strings.ToLower(input) //quitamos el problema de mayusculas/minisculas
	re := regexp.MustCompile(`-(\w+)=("[^"]+"|\S+)`)
	matches := re.FindAllStringSubmatch(input, -1)

	for _, match := range matches {
		flagName := match[1]
		flagValue := match[2]

		// Delete quotes if they are present in the value
		flagValue = strings.Trim(flagValue, "\"")
		switch flagName {
		case "size":
			sizeValue := 0
			fmt.Sscanf(flagValue, "%d", &sizeValue)
			*size = sizeValue
		case "driveletter":
			*driveletter = flagValue
		case "name":
			*name = flagValue
		case "unit":
			*unit = flagValue
		case "type":
			*type_ = flagValue
		case "fit":
			flagValue = flagValue[:1]
			*fit = flagValue
		case "delete":
			*delete = flagValue
		case "add":
			addValue := 0
			fmt.Sscanf(flagValue, "%d", &addValue)
			*add = addValue
		case "path":
			*path = flagValue
		default:
			fmt.Println("Error: Flag not found")
		}
		if *unit == "" {
			*unit = "k"
		}
		if *fit == "" {
			*fit = "w"
		}
		if *type_ == "" {
			*type_ = "p"
		}
	}
}

func CRUD_Partitions(size *int, driveletter *string, name *string, unit *string, type_ *string, fit *string, delete *string, add *int, path *string) {
	//println(*unit)

	if *unit == "k" {
		*size = *size * 1024
	} else if *unit == "m" {
		*size = *size * 1024 * 1024
	}
	if *unit == "k" {
		*add = *add * 1024
	} else if *unit == "m" {
		*add = *add * 1024 * 1024
	}

	//println("Size partition: ", *size)

	// Open bin file
	*driveletter = strings.ToUpper(*driveletter)
	filepath := "./Disks/" + *driveletter + ".dsk"
	file, err := utilities_test.OpenFile(filepath)
	if err != nil {
		return
	}

	var compareMBR structs_test.MBR
	copy(compareMBR.Mbr_particion[0].Part_name[:], *name)
	copy(compareMBR.Mbr_particion[0].Part_type[:], "p")
	copy(compareMBR.Mbr_particion[1].Part_type[:], "e")
	copy(compareMBR.Mbr_particion[2].Part_type[:], "l")
	var TempMBR structs_test.MBR
	// Read object from bin file
	if err := utilities_test.ReadObject(file, &TempMBR, 0); err != nil {
		return
	}

	// Verificar si el nombre de la partición ya está en uso
	for _, partition := range TempMBR.Mbr_particion {
		if bytes.Equal(partition.Part_name[:], compareMBR.Mbr_particion[0].Part_name[:]) && *delete == "" && *add == 0 {
			fmt.Println("Error: El nombre de la partición ya está en uso!!!!!!!!!!!!!!!!!!!!!!!!!!")
			return
		}
	}

	//Validar si existe una particion extendida
	var EPartition = false
	var EPartitionStart int
	for _, partition := range TempMBR.Mbr_particion {
		if bytes.Equal(partition.Part_type[:], compareMBR.Mbr_particion[1].Part_type[:]) {
			EPartition = true
			EPartitionStart = int(partition.Part_start)
			//fmt.Println("¡Existe una particion extendida!")
		}
	}

	//? Print object
	// fmt.Println(">>>>>ANTES")
	// structs_test.PrintMBR(TempMBR)

	var Partition structs_test.Partition
	// Si la operación es de eliminación y se especifica eliminar completamente
	if *delete == "full" {
		// Buscar la partición por nombre y eliminarla
		for i := range TempMBR.Mbr_particion {
			if bytes.Equal(TempMBR.Mbr_particion[i].Part_name[:], compareMBR.Mbr_particion[0].Part_name[:]) {
				TempMBR.Mbr_particion[i] = Partition // Vaciar la partición
				break
			}
		}
	} else if *add != 0 {
		//println("ADD", *add)
		// Añadir o quitar espacio en las particiones
		for i := range TempMBR.Mbr_particion {
			if bytes.Equal(TempMBR.Mbr_particion[i].Part_name[:], compareMBR.Mbr_particion[0].Part_name[:]) {
				// Validar que no queden números negativos en el espacio de las particiones
				if TempMBR.Mbr_particion[i].Part_size+int32(*add) < 0 {
					fmt.Println("Error: El espacio de la partición no puede ser negativo")
					return
				}
				// Validar que al añadir no se sobrepase el start de la siguiente partición
				if i < len(TempMBR.Mbr_particion)-1 && TempMBR.Mbr_particion[i+1].Part_start < TempMBR.Mbr_particion[i].Part_start+TempMBR.Mbr_particion[i].Part_size+int32(*add) {
					if TempMBR.Mbr_particion[i+1].Part_start != 0 {
						fmt.Println("Error: Al añadir espacio, se sobrepasa el start de la siguiente partición")
						return
					}
				}
				TempMBR.Mbr_particion[i].Part_size += int32(*add)
				if TempMBR.Mbr_particion[i].Part_size > TempMBR.Mbr_tamano {
					println("Supera el tamaño del disco")
					return
				}
				break
			}
		}
	} else {
		var count = 0
		var gap = int32(0)
		// Iterate over the partitions
		for i := 0; i < 4; i++ {

			if TempMBR.Mbr_particion[i].Part_size != 0 {
				count++
				gap = TempMBR.Mbr_particion[i].Part_start + TempMBR.Mbr_particion[i].Part_size
			}
		}

		for i := 0; i < 4; i++ {

			if TempMBR.Mbr_particion[i].Part_size == 0 {
				TempMBR.Mbr_particion[i].Part_size = int32(*size)

				if count == 0 {
					TempMBR.Mbr_particion[i].Part_start = int32(binary.Size(TempMBR))
				} else {
					TempMBR.Mbr_particion[i].Part_start = gap
				}

				suma := int32(*size) + int32(binary.Size(TempMBR))
				//println("Tamaño del disco:", TempMBR.Mbr_tamano)
				//println("Suma:", suma)
				if suma > TempMBR.Mbr_tamano {
					println("Error: La particion exede el tamaño del disco!!!!!!!!!!!!!!!!!!")
					return
				}

				copy(TempMBR.Mbr_particion[i].Part_name[:], *name)
				copy(TempMBR.Mbr_particion[i].Part_fit[:], *fit)
				copy(TempMBR.Mbr_particion[i].Part_status[:], "0")
				copy(TempMBR.Mbr_particion[i].Part_type[:], *type_)
				TempMBR.Mbr_particion[i].Part_correlative = int32(count + 1)
				break
			}
		}

		// Validar que si no existe una particion extendida no se puede crear una logica
		for _, partition := range TempMBR.Mbr_particion {
			if bytes.Equal(partition.Part_type[:], compareMBR.Mbr_particion[2].Part_type[:]) && *type_ == "l" {
				fmt.Printf("Original: %s Comparacion: %s \n", partition.Part_type[:], compareMBR.Mbr_particion[2].Part_type[:])
				if !EPartition {
					println("Error: No se puede crear una parcicion logica si no existe una extendida!!!!!!!!!!!!!!!!!!!!!!!!!")
					return
				}
				//?EBR verificacion
				var x = 0
				for x < 1 {
					var TempEBR structs_test.EBR
					if err := utilities_test.ReadObject(file, &TempEBR, int64(EPartitionStart)); err != nil {
						return
					}

					if TempEBR.Part_s != 0 {
						// Escribir un nuevo EBR en el archivo binario
						var newEBR structs_test.EBR
						copy(newEBR.Part_mount[:], "0")                                   // Indica si la partición está montada o no
						copy(newEBR.Part_fit[:], "l")                                     // Tipo de ajuste de la partición
						newEBR.Part_start = int32(EPartitionStart) + 1                    // Indica en qué byte del disco inicia la partición
						newEBR.Part_s = TempEBR.Part_s                                    // Contiene el tamaño total de la partición en bytes
						newEBR.Part_next = int32(EPartitionStart) + int32(TempEBR.Part_s) // Byte en el que está el próximo EBR (-1 si no hay siguiente)
						copy(newEBR.Part_name[:], TempEBR.Part_name[:])                   // Nombre de la partición

						// Escribir el nuevo EBR en el archivo binario
						if err := utilities_test.WriteObject(file, newEBR, int64(EPartitionStart)); err != nil {
							return
						}
						EPartitionStart = EPartitionStart + int(TempEBR.Part_s)
						structs_test.PrintEBR(newEBR)
					} else {
						// Escribir un nuevo EBR en el archivo binario
						var newEBR structs_test.EBR
						copy(newEBR.Part_mount[:], "0")                // Indica si la partición está montada o no
						copy(newEBR.Part_fit[:], "l")                  // Tipo de ajuste de la partición
						newEBR.Part_start = int32(EPartitionStart) + 1 // Indica en qué byte del disco inicia la partición
						newEBR.Part_s = int32(*size)                   // Contiene el tamaño total de la partición en bytes
						newEBR.Part_next = -1                          // Byte en el que está el próximo EBR (-1 si no hay siguiente)
						copy(newEBR.Part_name[:], *name)               // Nombre de la partición

						// Escribir el nuevo EBR en el archivo binario
						if err := utilities_test.WriteObject(file, newEBR, int64(EPartitionStart)); err != nil {
							return
						}
						x = 1
						structs_test.PrintEBR(newEBR)
					}

				}
				break
			}
		}

		// Validar que no exista mas de 1 particion extendida por disco
		var Ecount = 0
		for _, partition := range TempMBR.Mbr_particion {
			if bytes.Equal(partition.Part_type[:], compareMBR.Mbr_particion[1].Part_type[:]) {
				if EPartition {
					Ecount += 1
				}
				if Ecount > 1 {
					println("Error: No se puede tener mas de 1 particion extendida por disco!!!!!!!!!!!!!!!!!!!!")
					return
				}
			}
		}

	}

	// Overwrite the MBR
	if err := utilities_test.WriteObject(file, TempMBR, 0); err != nil {
		return
	}

	var TempMBR2 structs_test.MBR
	// Read object from bin file
	if err := utilities_test.ReadObject(file, &TempMBR2, 0); err != nil {
		return
	}

	// Print object
	// fmt.Println(">>>>>DESPUES")
	// structs_test.PrintMBR(TempMBR2)

	// Close bin file
	defer file.Close()
}

/* -------------------------------------------------------------------------- */
/*                                COMANDO MOUNT                               */
/* -------------------------------------------------------------------------- */

func ProcessMOUNT(input string, driveletter *string, name *string) {
	input = strings.ToLower(input) //quitamos el problema de mayusculas/minisculas
	re := regexp.MustCompile(`-(\w+)=("[^"]+"|\S+)`)
	matches := re.FindAllStringSubmatch(input, -1)

	for _, match := range matches {
		flagName := match[1]
		flagValue := match[2]

		// Delete quotes if they are present in the value
		flagValue = strings.Trim(flagValue, "\"")
		switch flagName {
		case "driveletter":
			*driveletter = flagValue
		case "name":
			*name = flagValue
		default:
			fmt.Println("Error: Flag not found")
		}
	}
}

func MountPartition(driveletter *string, name *string) {
	// Open bin file
	*driveletter = strings.ToUpper(*driveletter)
	filepath := "./Disks/" + *driveletter + ".dsk"
	file, err := utilities_test.OpenFile(filepath)
	if err != nil {
		return
	}

	var TempMBR structs_test.MBR
	// Read object from bin file
	if err := utilities_test.ReadObject(file, &TempMBR, 0); err != nil {
		return
	}

	var compareMBR structs_test.MBR
	copy(compareMBR.Mbr_particion[0].Part_name[:], *name)

	for i := 0; i < 4; i++ {

		if bytes.Equal(TempMBR.Mbr_particion[i].Part_name[:], compareMBR.Mbr_particion[0].Part_name[:]) {
			println("entro a la igualacion")
			copy(TempMBR.Mbr_particion[i].Part_status[:], "1")
			ID := fmt.Sprintf("%s%d%s", *driveletter, TempMBR.Mbr_particion[i].Part_correlative, "02")
			println(ID)
			copy(TempMBR.Mbr_particion[i].Part_id[:], ID)
			break
		}
	}

	// Overwrite the MBR
	if err := utilities_test.WriteObject(file, TempMBR, 0); err != nil {
		return
	}

	structs_test.PrintMBR(TempMBR)
}

/* -------------------------------------------------------------------------- */
/*                               COMANDO UNMOUNT                              */
/* -------------------------------------------------------------------------- */

func ProcessUNMOUNT(input string, id *string) {
	input = strings.ToLower(input) //quitamos el problema de mayusculas/minisculas
	re := regexp.MustCompile(`-(\w+)=("[^"]+"|\S+)`)
	match := re.FindStringSubmatch(input)
	if len(match) != 3 {
		fmt.Println("Comando unmount incorrecto")
		return
	}
	flagValue := match[2]
	*id = flagValue
}

func UNMOUNT_Partition(id *string) {
	letra := string((*id)[0])
	letra = strings.ToUpper(letra)

	correlativo, err := strconv.ParseInt(string((*id)[len(*id)-3]), 10, 32)
	if err != nil {
		fmt.Println("Error al convertir la cadena a int32:", err)
		return
	}
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

	var compareMBR structs_test.MBR
	compareMBR.Mbr_particion[0].Part_correlative = int32(correlativo)

	for i := 0; i < 4; i++ {

		if TempMBR.Mbr_particion[i].Part_correlative == compareMBR.Mbr_particion[0].Part_correlative {
			println("entro a la igualacion")
			copy(TempMBR.Mbr_particion[i].Part_status[:], "0")
			break
		}
	}

	// Overwrite the MBR
	if err := utilities_test.WriteObject(file, TempMBR, 0); err != nil {
		return
	}

	structs_test.PrintMBR(TempMBR)
}

/* -------------------------------------------------------------------------- */
/*                                COMANDO MKFS                                */
/* -------------------------------------------------------------------------- */
func ProcessMKFS(input string, id *string, type_ *string, fs *string) {
	input = strings.ToLower(input) //quitamos el problema de mayusculas/minisculas
	re := regexp.MustCompile(`-(\w+)=("[^"]+"|\S+)`)
	matches := re.FindAllStringSubmatch(input, -1)

	for _, match := range matches {
		flagName := match[1]
		flagValue := match[2]

		// Delete quotes if they are present in the value
		flagValue = strings.Trim(flagValue, "\"")
		switch flagName {
		case "id":
			*id = flagValue
		case "type":
			*type_ = flagValue
		case "fs":
			*fs = flagValue
		default:
			fmt.Println("Error: Flag not found")
		}
	}
}

func MKFS(id *string, type_ *string, fs *string) {

	fmt.Println("Id:", id)
	fmt.Println("Type:", type_)
	fmt.Println("Fs:", fs)

	driveletter := string((*id)[0])

	// Open bin file
	filepath := "./Disks/" + strings.ToUpper(driveletter) + ".dsk"
	file, err := utilities_test.OpenFile(filepath)
	if err != nil {
		return
	}

	var TempMBR structs_test.MBR
	// Read object from bin file
	if err := utilities_test.ReadObject(file, &TempMBR, 0); err != nil {
		return
	}

	// Print object
	structs_test.PrintMBR(TempMBR)

	fmt.Println("-------------")

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

	// numerador = (partition_montada.size - sizeof(Structs::Superblock)
	// denrominador base = (4 + sizeof(Structs::Inodes) + 3 * sizeof(Structs::Fileblock))
	// temp = "2" ? 0 : sizeof(Structs::Journaling)
	// denrominador = base + temp
	// n = floor(numerador / denrominador)

	numerador := int32(TempMBR.Mbr_particion[index].Part_size - int32(binary.Size(structs_test.S_block{})))
	denrominador_base := int32(4 + int32(binary.Size(structs_test.Inode{})) + 3*int32(binary.Size(structs_test.B_files{})))
	var temp int32 = 0
	if *fs == "2fs" {
		temp = 0
	} else {
		temp = int32(binary.Size(structs_test.Journaling{}))
	}
	denrominador := denrominador_base + temp
	n := int32(numerador / denrominador)

	fmt.Println("N:", n)

	// var newMRB Structs.MRB
	var newSuperblock structs_test.S_block
	newSuperblock.S_inodes_count = 0
	newSuperblock.S_blocks_count = 0

	newSuperblock.S_free_blocks_count = 3 * n
	newSuperblock.S_free_inodes_count = n

	copy(newSuperblock.S_mtime[:], "28/02/2024")
	copy(newSuperblock.S_umtime[:], "28/02/2024")
	newSuperblock.S_mnt_count = 0

	if *fs == "2fs" {
		create_ext2(n, TempMBR.Mbr_particion[index], newSuperblock, "28/02/2024", file)
	} else {
		create_ext3()
	}

	// Close bin file
	defer file.Close()

}

func create_ext2(n int32, partition structs_test.Partition, newSuperblock structs_test.S_block, date string, file *os.File) {
	fmt.Println("N:", n)
	fmt.Println("Superblock:", newSuperblock)
	fmt.Println("Date:", date)

	newSuperblock.S_filesystem_type = 2
	newSuperblock.S_bm_inode_start = partition.Part_start + int32(binary.Size(structs_test.S_block{}))
	newSuperblock.S_bm_block_start = newSuperblock.S_bm_inode_start + n
	newSuperblock.S_inode_start = newSuperblock.S_bm_block_start + 3*n
	newSuperblock.S_block_start = newSuperblock.S_inode_start + n*int32(binary.Size(structs_test.Inode{}))

	newSuperblock.S_free_inodes_count -= 1
	newSuperblock.S_free_blocks_count -= 1
	newSuperblock.S_free_inodes_count -= 1
	newSuperblock.S_free_blocks_count -= 1

	for i := int32(0); i < n; i++ {
		err := utilities_test.WriteObject(file, byte(0), int64(newSuperblock.S_bm_inode_start+i))
		if err != nil {
			fmt.Println("Error: ", err)
		}
	}

	for i := int32(0); i < 3*n; i++ {
		err := utilities_test.WriteObject(file, byte(0), int64(newSuperblock.S_bm_block_start+i))
		if err != nil {
			fmt.Println("Error: ", err)
		}
	}

	var newInode structs_test.Inode
	for i := int32(0); i < 15; i++ {
		newInode.I_block[i] = -1
	}

	for i := int32(0); i < n; i++ {
		err := utilities_test.WriteObject(file, newInode, int64(newSuperblock.S_inode_start+i*int32(binary.Size(structs_test.Inode{}))))
		if err != nil {
			fmt.Println("Error: ", err)
		}
	}

	var newFileblock structs_test.B_files
	for i := int32(0); i < 3*n; i++ {
		err := utilities_test.WriteObject(file, newFileblock, int64(newSuperblock.S_block_start+i*int32(binary.Size(structs_test.B_files{}))))
		if err != nil {
			fmt.Println("Error: ", err)
		}
	}

	var Inode0 structs_test.Inode //Inode 0
	Inode0.I_uid = 1
	Inode0.I_gid = 1
	Inode0.I_s = 0
	copy(Inode0.I_atime[:], date)
	copy(Inode0.I_ctime[:], date)
	copy(Inode0.I_mtime[:], date)
	copy(Inode0.I_perm[:], "0")
	copy(Inode0.I_perm[:], "664")

	for i := int32(0); i < 15; i++ {
		Inode0.I_block[i] = -1
	}

	Inode0.I_block[0] = 0

	// . | 0
	// .. | 0
	// users.txt | 1
	//

	var Folderblock0 structs_test.B_files //Bloque 0 -> carpetas
	Folderblock0.B_content[0].B_inodo = 0
	copy(Folderblock0.B_content[0].B_name[:], ".")
	Folderblock0.B_content[1].B_inodo = 0
	copy(Folderblock0.B_content[1].B_name[:], "..")
	Folderblock0.B_content[1].B_inodo = 1
	copy(Folderblock0.B_content[1].B_name[:], "users.txt")

	var Inode1 structs_test.Inode //Inode 1
	Inode1.I_uid = 1
	Inode1.I_gid = 1
	Inode1.I_s = int32(binary.Size(structs_test.B_files{}))
	copy(Inode1.I_atime[:], date)
	copy(Inode1.I_ctime[:], date)
	copy(Inode1.I_mtime[:], date)
	copy(Inode1.I_perm[:], "0")
	copy(Inode1.I_perm[:], "664")

	for i := int32(0); i < 15; i++ {
		Inode1.I_block[i] = -1
	}

	Inode0.I_block[0] = 1

	data := "1,G,root\n1,U,root,root,123\n"
	var Fileblock1 structs_test.B_docs //Bloque 1 -> archivo
	copy(Fileblock1.B_content[:], data)

	// Inodo 0 -> Bloque 0 -> Inodo 1 -> Bloque 1
	// Crear la carpeta raiz /
	// Crear el archivo users.txt "1,G,root\n1,U,root,root,123\n"

	// write superblock
	err := utilities_test.WriteObject(file, newSuperblock, int64(partition.Part_start))
	if err != nil {
		fmt.Println("Error: ", err)
	}

	// write bitmap inodes
	err = utilities_test.WriteObject(file, byte(1), int64(newSuperblock.S_bm_inode_start))
	if err != nil {
		fmt.Println("Error: ", err)
	}

	err = utilities_test.WriteObject(file, byte(1), int64(newSuperblock.S_bm_inode_start+1))
	if err != nil {
		fmt.Println("Error: ", err)
	}

	// write bitmap blocks
	err = utilities_test.WriteObject(file, byte(1), int64(newSuperblock.S_bm_block_start))
	if err != nil {
		fmt.Println("Error: ", err)
	}

	err = utilities_test.WriteObject(file, byte(1), int64(newSuperblock.S_bm_block_start+1))
	if err != nil {
		fmt.Println("Error: ", err)
	}

	// write inodes
	err = utilities_test.WriteObject(file, Inode0, int64(newSuperblock.S_inode_start)) //Inode 0

	if err != nil {
		fmt.Println("Error: ", err)
	}

	err = utilities_test.WriteObject(file, Inode1, int64(newSuperblock.S_inode_start+int32(binary.Size(structs_test.Inode{})))) //Inode 1

	if err != nil {
		fmt.Println("Error: ", err)
	}
	// write blocks
	err = utilities_test.WriteObject(file, Folderblock0, int64(newSuperblock.S_block_start)) //Bloque 0

	if err != nil {
		fmt.Println("Error: ", err)
	}

	err = utilities_test.WriteObject(file, Fileblock1, int64(newSuperblock.S_block_start+int32(binary.Size(structs_test.B_files{})))) //Bloque 1

	if err != nil {
		fmt.Println("Error: ", err)
	}

	//mkfs -type=full -id=A119
}

func create_ext3()  {
	
}

//?                    ADMINISTRACION DE USUARIOS Y GRUPOS
/* -------------------------------------------------------------------------- */
/*                                COMANDO LOGIN                               */
/* -------------------------------------------------------------------------- */

/* -------------------------------------------------------------------------- */
/*                               COMANDO LOGOUT                               */
/* -------------------------------------------------------------------------- */

/* -------------------------------------------------------------------------- */
/*                                COMANDO MKGRP                               */
/* -------------------------------------------------------------------------- */

/* -------------------------------------------------------------------------- */
/*                                COMANDO RMGRP                               */
/* -------------------------------------------------------------------------- */

/* -------------------------------------------------------------------------- */
/*                                COMANDO MKUSR                               */
/* -------------------------------------------------------------------------- */

/* -------------------------------------------------------------------------- */
/*                                COMANDO RMUSR                               */
/* -------------------------------------------------------------------------- */

/* -------------------------------------------------------------------------- */
/*                               COMANDO EXECUTE                              */
/* -------------------------------------------------------------------------- */

//?               ADMINISTRACION DE CARPETAS, ARCHIVOS Y PERMISOS
/* -------------------------------------------------------------------------- */
/*                               COMANDO MKFILE                               */
/* -------------------------------------------------------------------------- */

/* -------------------------------------------------------------------------- */
/*                                 COMANDO CAT                                */
/* -------------------------------------------------------------------------- */

/* -------------------------------------------------------------------------- */
/*                               COMANDO REMOVE                               */
/* -------------------------------------------------------------------------- */

/* -------------------------------------------------------------------------- */
/*                                COMANDO EDIT                                */
/* -------------------------------------------------------------------------- */

/* -------------------------------------------------------------------------- */
/*                               COMANDO RENAME                               */
/* -------------------------------------------------------------------------- */

/* -------------------------------------------------------------------------- */
/*                                COMANDO MKDIR                               */
/* -------------------------------------------------------------------------- */

/* -------------------------------------------------------------------------- */
/*                                COMANDO COPY                                */
/* -------------------------------------------------------------------------- */

/* -------------------------------------------------------------------------- */
/*                                COMANDO MOVE                                */
/* -------------------------------------------------------------------------- */

/* -------------------------------------------------------------------------- */
/*                                COMANDO FIND                                */
/* -------------------------------------------------------------------------- */

/* -------------------------------------------------------------------------- */
/*                                COMANDO CHOWN                               */
/* -------------------------------------------------------------------------- */

/* -------------------------------------------------------------------------- */
/*                                COMANDO CHGRP                               */
/* -------------------------------------------------------------------------- */

/* -------------------------------------------------------------------------- */
/*                                COMANDO CHMOD                               */
/* -------------------------------------------------------------------------- */

/* -------------------------------------------------------------------------- */
/*                                COMANDO PAUSE                               */
/* -------------------------------------------------------------------------- */

//?                     			REPORTES
/* -------------------------------------------------------------------------- */
/*                                 REPORTE MBR                                */
/* -------------------------------------------------------------------------- */
/* -------------------------------------------------------------------------- */
/*                                REPORTE DISK                                */
/* -------------------------------------------------------------------------- */
/* -------------------------------------------------------------------------- */
/*                                REPORTE INODE                               */
/* -------------------------------------------------------------------------- */
/* -------------------------------------------------------------------------- */
/*                                REPORTE BLOCK                               */
/* -------------------------------------------------------------------------- */
/* -------------------------------------------------------------------------- */
/*                              REPORTE BM_INODE                              */
/* -------------------------------------------------------------------------- */
/* -------------------------------------------------------------------------- */
/*                               REPORTE BM_BLOC                              */
/* -------------------------------------------------------------------------- */
/* -------------------------------------------------------------------------- */
/*                                REPORTE TREE                                */
/* -------------------------------------------------------------------------- */
/* -------------------------------------------------------------------------- */
/*                                 REPORTE SB                                 */
/* -------------------------------------------------------------------------- */
/* -------------------------------------------------------------------------- */
/*                                REPORTE FILE                                */
/* -------------------------------------------------------------------------- */
/* -------------------------------------------------------------------------- */
/*                                 REPORTE LS                                 */
/* -------------------------------------------------------------------------- */
/* -------------------------------------------------------------------------- */
/*                             REPORTE JOURNALING                             */
/* -------------------------------------------------------------------------- */

func ProcessExecute(input string, path *string) {
	re := regexp.MustCompile(`-(\w+)=("[^"]+"|\S+)`)
	match := re.FindStringSubmatch(input)
	if len(match) > 1 {
		*path = match[2]
	}
}

func GenerateUniqueID() int {
	// Obtener la marca de tiempo actual
	currentTime := time.Now()
	// Generar un número aleatorio entre 0 y 9999
	randomNumber := rand.Intn(10000)
	// Combinar la marca de tiempo y el número aleatorio para crear un identificador único
	uniqueID := currentTime.UnixNano() * int64(randomNumber) % (1 << 31)
	// Tomar el valor absoluto para asegurarse de que sea positivo
	uniqueID = int64(math.Abs(float64(uniqueID)))
	return int(uniqueID)
}

func ValidDriveLetter(str string) bool {
	return regexp.MustCompile(`^[a-zA-Z]$`).MatchString(str)
}

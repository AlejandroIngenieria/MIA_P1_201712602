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
	"strings"
	"time"
)

var fileCounter int = 0

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
	}
}

func CRUD_Partitions(size *int, driveletter *string, name *string, unit *string, type_ *string, fit *string, delete *string, add *int, path *string) {
	println(*unit)

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

	println("Size partition: ", *size)

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
	for _, partition := range TempMBR.Mbr_particion {
		if bytes.Equal(partition.Part_type[:], compareMBR.Mbr_particion[1].Part_type[:]) {
			EPartition = true
			fmt.Println("¡Existe una particion extendida!")
		}
	}
	println("EPartition: ", EPartition)

	// Print object
	fmt.Println(">>>>>ANTES")
	structs_test.PrintMBR(TempMBR)

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
		println("ADD", *add)
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
				println("Tamaño del disco:", TempMBR.Mbr_tamano)
				println("Suma:", suma)
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
			if bytes.Equal(partition.Part_type[:], compareMBR.Mbr_particion[2].Part_type[:]) {
				if !EPartition {
					println("Error: No se puede crear una parcicion logica si no existe una extendida!!!!!!!!!!!!!!!!!!!!!!!!!")
					return
				} else {
					//De momento vamos a ignorar las particiones logicas
					return
				}
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
	fmt.Println(">>>>>DESPUES")
	structs_test.PrintMBR(TempMBR2)

	// Close bin file
	defer file.Close()
}

/* -------------------------------------------------------------------------- */
/*                               COMANDO EXECUTE                              */
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

package functions_test

import (
	//"bufio"
	"P1/Structs"
	"encoding/binary"
	"fmt"
	"math"
	"math/rand"
	"os"
	"path/filepath"
	"regexp"

	//"strconv"
	"strings"
	"time"
)

var fileCounter int = 0

/* -------------------------------------------------------------------------- */
/*                                  COMANDOS                                  */
/* -------------------------------------------------------------------------- */
// comando mkdisk
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
			*fit = flagValue
		case "unit":
			*unit = flagValue
		default:
			fmt.Println("Error: Flag not found")
		}
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
		fmt.Println("Error al crear archivo de 1"+*unit+"b: ", err)
		return
	}

	// Incrementar el contador
	fileCounter++
}

func createFile(filename string, size int, fit string) error {
	// Crear el archivo con el nombre proporcionado
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	// Obtener la hora actual
	currentTime := time.Now()
	// Formatear la hora actual como una cadena
	timeString := currentTime.Format("2006-01-02 15:04:05")
	//Asignacion de datos al MBR
	TempMBR := structs_test.MBR{}
	TempMBR.Mbr_tamano = int32(size)
	copy(TempMBR.Mbr_fecha_creacion[:], []byte(timeString))
	TempMBR.Mbr_dsk_signature = int32(GenerateUniqueID())
	bytes := []byte(fit)
	TempMBR.Dsk_fit = [1]byte{bytes[0]}
	fmt.Println("DiskFit: ",string(TempMBR.Dsk_fit[:]))
	// Escribir el objeto TempMBR al inicio del archivo
	if err := writeMBR(file, TempMBR); err != nil {
		return err
	}

	// Escribir bytes nulos en el resto del archivo
	data := make([]byte, size-binary.Size(TempMBR))
	if _, err := file.Write(data); err != nil {
		return err
	}

	return nil
}

func writeMBR(file *os.File, TempMBR structs_test.MBR) error {
	// Convertir TempMBR a una secuencia de bytes
	mbrBytes := make([]byte, binary.Size(TempMBR))
	offset := 0

	// Escribir campos del MBR en el slice de bytes
	binary.LittleEndian.PutUint32(mbrBytes[offset:offset+4], uint32(TempMBR.Mbr_tamano))
	offset += 4
	copy(mbrBytes[offset:offset+10], TempMBR.Mbr_fecha_creacion[:])
	offset += 10
	binary.LittleEndian.PutUint32(mbrBytes[offset:offset+4], uint32(TempMBR.Mbr_dsk_signature))
	offset += 4
	copy(mbrBytes[offset:offset+1], TempMBR.Dsk_fit[:])

	for i, partition := range TempMBR.Mbr_particion {
		offset = 18 + i*binary.Size(partition)
		copy(mbrBytes[offset:offset+1], partition.Part_status[:])
		copy(mbrBytes[offset+1:offset+2], partition.Part_type[:])
		copy(mbrBytes[offset+2:offset+3], partition.Part_fit[:])
		binary.LittleEndian.PutUint32(mbrBytes[offset+3:offset+7], uint32(partition.Part_start))
		binary.LittleEndian.PutUint32(mbrBytes[offset+7:offset+11], uint32(partition.Part_size))
		copy(mbrBytes[offset+11:offset+27], partition.Part_name[:])
		binary.LittleEndian.PutUint32(mbrBytes[offset+27:offset+31], uint32(partition.Part_correlative))
		copy(mbrBytes[offset+31:offset+35], partition.Part_id[:])
	}

	// Escribir la secuencia de bytes en el archivo
	_, err := file.Write(mbrBytes)
	return err
}



// comando rmdisk
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

// comando fdisk
func ProcessFDISK(input string, size *int, driveletter *string,
	name *string, unit *string, tipe *string, fit *string, delete *string,
	add *string, path *string) {
	input = strings.ToLower(input) //quitamos el problema de mayusculas/minisculas
	re := regexp.MustCompile(`-(\w+)=("[^"]+"|\S+)`)
	matches := re.FindAllStringSubmatch(input, -1)

	for _, match := range matches {
		flagName := match[1]
		flagValue := match[2]

		// Delete quotes if they are present in the value
		flagValue = strings.Trim(flagValue, "\"")
		var unitB = false
		var fitB = false
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
			unitB = true
		case "type":
			*tipe = flagValue
		case "fit":
			flagValue = flagValue[:1]
			*fit = flagValue
			fitB = true
		case "delete":
			*delete = flagValue
		case "add":
			*add = flagValue
		case "path":
			*path = flagValue
		default:
			fmt.Println("Error: Flag not found")
		}
		if !unitB {
			*unit = "k"
		}
		if !fitB {
			*fit = "wf"
		}
	}
}

/* -------------------------------------------------------------------------- */
/*                                 AUXILIARES                                 */
/* -------------------------------------------------------------------------- */
// Funcion para crear el archivo binario
func CreateFile(name string) error {
	//Se asegura que el directorio no exista
	dir := filepath.Dir(name)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		fmt.Println("Err CreateFile dir==", err)
		return err
	}

	// Crea el archivo
	if _, err := os.Stat(name); os.IsNotExist(err) {
		file, err := os.Create(name)
		if err != nil {
			fmt.Println("Err CreateFile create==", err)
			return err
		}
		defer file.Close()
	}
	return nil
}

// Funcion para abrir el archivo binario en el modo (lectura/escritura)
func OpenFile(name string) (*os.File, error) {
	file, err := os.OpenFile(name, os.O_RDWR, 0644)
	if err != nil {
		fmt.Println("Err OpenFile==", err)
		return nil, err
	}
	return file, nil
}

// Funcion para escribir un objeto en el archivo binario
func WriteObject(file *os.File, data interface{}, position int64) error {
	file.Seek(position, 0)
	err := binary.Write(file, binary.LittleEndian, data)
	if err != nil {
		fmt.Println("Err WriteObject==", err)
		return err
	}
	return nil
}

// Funcion para leer un objeto desde un archivo binario
func ReadObject(file *os.File, data interface{}, position int64) error {
	file.Seek(position, 0)
	err := binary.Read(file, binary.LittleEndian, data)
	if err != nil {
		fmt.Println("Err ReadObject==", err)
		return err
	}
	return nil
}

func GenerateUniqueID() int {
	// Obtener la marca de tiempo actual
	currentTime := time.Now()
	// Generar un número aleatorio entre 0 y 9999
	randomNumber := rand.Intn(10000)
	// Combinar la marca de tiempo y el número aleatorio para crear un identificador único
	uniqueID := currentTime.UnixNano()*int64(randomNumber) % (1 << 31)
	// Tomar el valor absoluto para asegurarse de que sea positivo
	uniqueID = int64(math.Abs(float64(uniqueID)))
	return int(uniqueID)
}

func ValidDriveLetter(str string) bool {
	return regexp.MustCompile(`^[a-zA-Z]$`).MatchString(str)
}

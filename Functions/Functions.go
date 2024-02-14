package functions

import (
	"encoding/binary"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

/* -------------------------------------------------------------------------- */
/*                                  COMANDOS                                  */
/* -------------------------------------------------------------------------- */
// comando mkdisk
func ProcessMKDISK(input string, size *int, fit *string, unit *string) {
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
			flagValue = strings.ToLower(flagValue)
			*fit = flagValue
		case "unit":
			flagValue = strings.ToLower(flagValue)
			*unit = flagValue
		default:
			fmt.Println("Error: Flag not found")
		}
	}
}

// comando mrdisk
func ProcessMRDISK(input string, driveletter *string) {
	re := regexp.MustCompile(`-(\w+)=("[^"]+"|\S+)`)
	match := re.FindStringSubmatch(input)
	if len(match) != 3 {
		fmt.Println("Comando rmdisk incorrecto")
		return
	}
	flagValue := match[2]
	flagValue = strings.ToLower(flagValue)
	*driveletter = flagValue

}

// comando fdisk
func ProcessFDISK(input string, size *int, driveletter *string, name *string, unit *string, tipe *string, fit *string, delete *string) {
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
			flagValue = strings.ToLower(flagValue)
			*driveletter = flagValue
		case "name":
			flagValue = strings.ToLower(flagValue)
			*name = flagValue
		case "unit":
			flagValue = strings.ToLower(flagValue)
			*unit = flagValue
		case "type":
			flagValue = strings.ToLower(flagValue)
			*tipe = flagValue
		case "fit":
			flagValue = flagValue[:1]
			flagValue = strings.ToLower(flagValue)
			*fit = flagValue
		case "delete":
			flagValue = strings.ToLower(flagValue)
			*delete = flagValue
		default:
			fmt.Println("Error: Flag not found")
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

func CreateDisk() {
	fmt.Print("Ingrese el tamaño en bytes: ")
	var input string
	fmt.Scanln(&input)
	// Convertir la entrada a un tipo int64
	fileSize, err := strconv.ParseInt(input, 10, 64)
	if err != nil {
		fmt.Println("Error: entrada no válida.")
	}
	randomNumber := rand.Intn(10000)
	nuevoArchivo := "./disks/" + strconv.Itoa(randomNumber) + ".dsk"

	// Crear un slice de bytes lleno de ceros binarios
	block := make([]byte, 1024)

	// Create bin file
	if err := CreateFile(nuevoArchivo); err != nil {
		return
	}

	// Open bin file
	file, err := OpenFile(nuevoArchivo)
	if err != nil {
		return
	}

	// Close bin file
	defer file.Close()

	// Escribir bloques de bytes llenos de ceros binarios hasta alcanzar el tamaño deseado
	for written := int64(0); written < fileSize; written += 1024 {
		_, err := file.Write(block)
		if err != nil {
			fmt.Println("Error al escribir en el archivo:", err)
			return
		}
	}
	return
}

func generateUniqueID() int {
	// Obtener la marca de tiempo actual
	currentTime := time.Now()
	// Generar un número aleatorio entre 0 y 9999
	randomNumber := rand.Intn(10000)
	// Combinar la marca de tiempo y el número aleatorio para crear un identificador único
	uniqueID := currentTime.UnixNano() * int64(randomNumber) % (1 << 31)
	return int(uniqueID)
}

func ValidDriveLetter(str string) bool {
	return regexp.MustCompile(`^[a-zA-Z]$`).MatchString(str)
}

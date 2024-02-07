package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"
)

func main() {

	// Lector para leer la entrada del usuario
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println("---Proyecto1----------------------------------------------------------------------------------------")
		fmt.Println("1.	Crear disco")
		fmt.Println("2.	Montar disco")
		fmt.Println("3.	Desmontar disco")
		fmt.Println("4.	Crear particion")
		fmt.Println("5.	Eliminar particion")
		fmt.Println("6.	Mostrar informacion del disco")
		fmt.Println("7.	Listar archivos y carpetas")
		fmt.Println("8.	Crear archivo")
		fmt.Println("9.	Eliminar archivo")
		fmt.Println("10.Salir")

		// Solicitar al usuario que seleccione una opción
		fmt.Print("Seleccione una opción: ")

		// Leer la entrada del usuario
		input, _ := reader.ReadString('\n')

		// Convertir la entrada del usuario a un número entero
		choice, err := strconv.Atoi(input[:len(input)-1]) // Eliminar el carácter de nueva línea

		// Verificar si hubo un error al convertir la entrada del usuario a un número entero
		if err != nil {
			fmt.Println("---ERROR! INGRESE UN NUMERO VALIDO------------------------------------------------------------------")
			continue
		}

		// Procesar la opción seleccionada por el usuario
		switch choice {
		case 1:
			fmt.Println("---Crear disco--------------------------------------------------------------------------------------")
		case 2:
			fmt.Println("---Montar disco-------------------------------------------------------------------------------------")
		case 3:
			fmt.Println("---Desmontar disco----------------------------------------------------------------------------------")
		case 4:
			fmt.Println("---Crear particion----------------------------------------------------------------------------------")
		case 5:
			fmt.Println("---Eliminar particion-------------------------------------------------------------------------------")
		case 6:
			fmt.Println("---Mostrar informacion del disco--------------------------------------------------------------------")
		case 7:
			fmt.Println("---Listar archivos y carpetas-----------------------------------------------------------------------")
		case 8:
			fmt.Println("---Crear archivo------------------------------------------------------------------------------------")
		case 9:
			fmt.Println("---Eliminar archivo---------------------------------------------------------------------------------")
		case 10:
			fmt.Println("---Programa finalizado------------------------------------------------------------------------------")
			return
		default:
			fmt.Println("---ERROR! OPCION NO VALIDA--------------------------------------------------------------------------")
		}
	}

}

// ? DISCOS
// Master Boot Record (MBR)
type mbr struct {
	mbr_tamano         int
	mbr_fecha_creacion time.Time
	mbr_dsk_signature  int
	dsk_fit            rune
	mbr_particion      [4]partition
}

// Partition
type partition struct {
	part_status      rune
	part_type        rune
	part_fit         rune
	part_start       int
	part_s           int
	part_name        [16]byte
	part_correlative int
	part_id          [4]byte
}

// Extended Boot Record (EBR)
type ebr struct {
	part_mount rune
	part_fit   rune
	part_start int
	part_s     int
	part_next  int
	part_name  [16]byte
}

// ? CARPETAS Y ARCHIVOS (EXT3|EXT2)
// Super bloque
type s_block struct {
	s_filesystem_type   int
	s_inodes_count      int
	s_blocks_count      int
	s_free_blocks_count int
	s_free_inodes_count int
	s_mtime             time.Time
	s_umtime            time.Time
	s_mnt_count         int
	s_magic             int
	s_inode_s           int
	s_block_s           int
	s_firts_ino         int
	s_firts_blo         int
	s_bm_inode_start    int
	s_bm_block_start    int
	s_inode_start       int
	s_block_start       int
}

// Inodos
type inode struct {
	i_uid   int
	i_gid   int
	i_s     int
	i_atime time.Time
	i_ctime time.Time
	i_mtime time.Time
	i_block int
	i_type  rune
	i_perm  [3]byte
}

// ? BLOQUES
// Bloque de carpetas
type b_files struct {
	b_content [4]content
}

type content struct {
	b_name  [12]byte
	b_inodo int
}

// Bloque de archivos
type b_docs struct {
	b_content [64]byte
}

// Bloque de apuntadores
type b_pointer struct {
	b_pointers [16]int
}

// ? FUNCIONES DE APOYO
func generateUniqueID() int {
	// Obtener la marca de tiempo actual
	currentTime := time.Now()
	// Generar un número aleatorio entre 0 y 9999
	randomNumber := rand.Intn(10000)
	// Combinar la marca de tiempo y el número aleatorio para crear un identificador único
	uniqueID := currentTime.UnixNano() * int64(randomNumber) % (1 << 31)
	return int(uniqueID)
}

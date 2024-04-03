package functions_test

import (
	"P1/Structs"
	"P1/Utilities"
	"bytes"
	//"encoding/binary"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

var (
	inodos     string = ""
	bloques    string = ""
	tree       string = ""
	relaciones string = ""
)

// ?                     			REPORTES
func ProcessREP(input string, name *string, path *string, id *string, ruta *string, flagN *bool) {
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
		case "path":
			*path = flagValue
		case "id":
			*id = flagValue
		case "ruta":
			*ruta = flagValue
		default:
			fmt.Println("Error: Flag not found: " + flagName)
			*flagN = true
		}
	}
}

func GenerateReports(name *string, path *string, id *string, ruta *string) {

	switch *name {
	//1
	case "mbr":
		REPORT_MBR(id, path)
	//2
	case "disk":
		REPORT_DISK(id, path)
	//3
	case "inode":
		REPORT_INODE(id, path)
	//4
	case "Journaling":
		REPORT_JOURNALING(id, path)
	//5
	case "block":
		REPORT_BLOCK(id, path)
	//6
	case "bm_inode":
		REPORT_BM_INODE(id, path)
	//7
	case "bm_block":
		REPORT_BM_BLOCK(id, path)
	//8
	case "tree":
		REPORT_TREE(id, path)
	//9
	case "sb":
		REPORT_SB(id, path)
	//10
	case "file":
		REPORT_FILE(id, path, ruta)
	//11
	case "ls":
		REPORT_LS(id, path, ruta)
	default:
		println("Reporte no reconocido:", *name)
	}
}

/* -------------------------------------------------------------------------- */
/*                               1 REPORTE MBR                                */
/* -------------------------------------------------------------------------- */
func REPORT_MBR(id *string, path *string) {
	letra := string((*id)[0])
	letra = strings.ToUpper(letra)
	filepath := "./Disks/" + letra + ".dsk"

	/* -------------------------------------------------------------------------- */
	/*                              BUSCAMOS EL DISCO                             */
	/* -------------------------------------------------------------------------- */

	file, err := os.Open(filepath)
	if err != nil {
		return
	}
	defer file.Close()

	/* -------------------------------------------------------------------------- */
	/*                               CARGAMOS EL MBR                              */
	/* -------------------------------------------------------------------------- */

	var TempMBR structs_test.MBR
	// Read object from bin file
	if err := utilities_test.ReadObject(file, &TempMBR, 0); err != nil {
		return
	}

	var EPartition = false
	var EPartitionStart int

	var compareMBR structs_test.MBR
	copy(compareMBR.Mbr_particion[0].Part_type[:], "p")
	copy(compareMBR.Mbr_particion[1].Part_type[:], "e")
	copy(compareMBR.Mbr_particion[2].Part_type[:], "l")

	/* -------------------------------------------------------------------------- */
	/*                 BUSCAMOS SI EXISTE UNA PARTICION EXTENDIDA                 */
	/* -------------------------------------------------------------------------- */

	for _, partition := range TempMBR.Mbr_particion {
		if bytes.Equal(partition.Part_type[:], compareMBR.Mbr_particion[1].Part_type[:]) {
			EPartition = true
			EPartitionStart = int(partition.Part_start)
		}
	}

	/* -------------------------------------------------------------------------- */
	/*                           ANALISIS DE CODIGO DOT                           */
	/* -------------------------------------------------------------------------- */

	strP := ""
	strE := ""

	/* -------------------------------------------------------------------------- */
	/*                                 PARTICIONES                                */
	/* -------------------------------------------------------------------------- */

	for _, partition := range TempMBR.Mbr_particion {
		partNameClean := strings.Trim(string(partition.Part_name[:]), "\x00")
		if partition.Part_correlative == 0 {
			continue
		} else {
			strP += fmt.Sprintf(`
		|Particion %d
		|{part_status|%s}
		|{part_type|%s}
		|{part_fit|%s}
		|{part_start|%d}
		|{part_size|%d}
		|{part_name|%s}`,
				partition.Part_correlative,
				string(partition.Part_status[:]),
				string(partition.Part_type[:]),
				string(partition.Part_fit[:]),
				partition.Part_start,
				partition.Part_size,
				partNameClean,
			)
		}
		/* -------------------------------------------------------------------------- */
		/*                             PARTICION EXTENDIDA                            */
		/* -------------------------------------------------------------------------- */
		//?EBR verificacion
		if bytes.Equal(partition.Part_type[:], compareMBR.Mbr_particion[1].Part_type[:]) && EPartition {
			// Validar que si no existe una particion extendida no se puede crear una logica
			//?EBR verificacion
			var x = 0
			for x < 1 {
				var TempEBR structs_test.EBR
				if err := utilities_test.ReadObject(file, &TempEBR, int64(EPartitionStart)); err != nil {
					return
				}

				if EPartitionStart != 0 && TempEBR.Part_next != -1 {
					partNameClean := strings.Trim(string(TempEBR.Part_name[:]), "\x00")
					strE += fmt.Sprintf(`
		|Particion Logica
		|{part_status|%s}
		|{part_next|%d}
		|{part_fit|%s}
		|{part_start|%d}
		|{part_size|%d}
		|{part_name|%s}`,
						string(TempEBR.Part_mount[:]),
						TempEBR.Part_next,
						string(TempEBR.Part_fit[:]),
						TempEBR.Part_start,
						TempEBR.Part_s,
						partNameClean,
					)
					//print("fit logica")
					//println(string(TempEBR.Part_fit[:]))
					EPartitionStart = int(TempEBR.Part_next)
				} else {
					//print("fit logica")
					//println(string(TempEBR.Part_fit[:]))
					partNameClean := strings.Trim(string(TempEBR.Part_name[:]), "\x00")
					strE += fmt.Sprintf(`
		|Particion Logica
		|{part_status|%s}
		|{part_next|%d}
		|{part_fit|%s}
		|{part_start|%d}
		|{part_size|%d}
		|{part_name|%s}`,
						string(TempEBR.Part_mount[:]),
						TempEBR.Part_next,
						string(TempEBR.Part_fit[:]),
						TempEBR.Part_start,
						TempEBR.Part_s,
						partNameClean,
					)
					strP += strE
					x = 1
				}
			}

		}

	}

	/* -------------------------------------------------------------------------- */
	/*                               CODIGO DOT BASE                              */
	/* -------------------------------------------------------------------------- */

	dotCode := fmt.Sprintf(`
		digraph G {
 			fontname="Helvetica,Arial,sans-serif"
			node [fontname="Helvetica,Arial,sans-serif"]
			edge [fontname="Helvetica,Arial,sans-serif"]
			concentrate=True;
			rankdir=TB;
			node [shape=record];

			title [label="Reporte MBR" shape=plaintext fontname="Helvetica,Arial,sans-serif"];

  			mbr[label="
				{MBR: %s.dsk|
					{mbr_tamaño|%d}
					|{mbr_fecha_creacion|%s}
					|{mbr_disk_signature|%d}
								%s
				}
			"];
			title2 [label="Reporte EBR" shape=plaintext fontname="Helvetica,Arial,sans-serif"];
			
			ebr[label="
				{EBR%s}
			"];

			title -> mbr [style=invis];
    		mbr -> title2[style=invis];
			title2 -> ebr[style=invis];
		}`,
		letra,
		TempMBR.Mbr_tamano,
		TempMBR.Mbr_fecha_creacion,
		TempMBR.Mbr_dsk_signature,
		strP,
		strE,
	)

	/* -------------------------------------------------------------------------- */
	/*                          ESCRIBIMOS EL CODIGO DOT                          */
	/* -------------------------------------------------------------------------- */

	dotFilePath := "./Reports/Rep1/mbr_rep.dot" // Ruta donde deseas guardar el archivo DOT
	dotFile, err := os.Create(dotFilePath)
	if err != nil {
		fmt.Println("Error al crear el archivo DOT:", err)
		return
	}
	defer dotFile.Close()

	_, err = dotFile.WriteString(dotCode)
	if err != nil {
		fmt.Println("Error al escribir en el archivo DOT:", err)
		return
	}

	/* -------------------------------------------------------------------------- */
	/*                         GENERAMOS EL GRAFO COMO PNG                        */
	/* -------------------------------------------------------------------------- */

	pngFilePath := *path // Ruta donde deseas guardar el archivo PNG
	cmd := exec.Command("dot", "-Tpng", "-o", pngFilePath, dotFilePath)
	err = cmd.Run()
	if err != nil {
		fmt.Println("Error al generar el gráfico:", err)
		return
	}

	fmt.Println("--------------------------------------------------------------------------")
	fmt.Println("               REPORTE-MBR/EBR: GENERADO CORRECTAMENTE                    ")
	fmt.Printf("                          %s                          ", pngFilePath)
	fmt.Println("--------------------------------------------------------------------------")
}

/* -------------------------------------------------------------------------- */
/*                              2 REPORTE DISK                                */
/* -------------------------------------------------------------------------- */

func REPORT_DISK(id *string, path *string) {
	letra := string((*id)[0])
	letra = strings.ToUpper(letra)

	/* -------------------------------------------------------------------------- */
	/*                              BUSCAMOS EL DISCO                             */
	/* -------------------------------------------------------------------------- */

	filepath := "./Disks/" + letra + ".dsk"
	file, err := os.Open(filepath)
	if err != nil {
		return
	}
	defer file.Close()

	/* -------------------------------------------------------------------------- */
	/*                                LEEMOS EL MBR                               */
	/* -------------------------------------------------------------------------- */

	var TempMBR structs_test.MBR
	// Read object from bin file
	if err := utilities_test.ReadObject(file, &TempMBR, 0); err != nil {
		return
	}

	var EPartition = false
	var EPartitionStart int

	/* -------------------------------------------------------------------------- */
	/*                          ESTRUCTURA PARA COMPARAR                          */
	/* -------------------------------------------------------------------------- */

	var compareMBR structs_test.MBR
	copy(compareMBR.Mbr_particion[0].Part_type[:], "p")
	copy(compareMBR.Mbr_particion[1].Part_type[:], "e")
	copy(compareMBR.Mbr_particion[2].Part_type[:], "l")

	/* -------------------------------------------------------------------------- */
	/*                 BUSCAMOS SI EXISTE UNA PARTICION EXTENDIDA                 */
	/* -------------------------------------------------------------------------- */

	for _, partition := range TempMBR.Mbr_particion {
		if bytes.Equal(partition.Part_type[:], compareMBR.Mbr_particion[1].Part_type[:]) {
			EPartition = true
			EPartitionStart = int(partition.Part_start)
		}
	}

	/* -------------------------------------------------------------------------- */
	/*              ITERAMOS LAS PARTICIONES PARA VER ORDEN Y ESPACIO             */
	/* -------------------------------------------------------------------------- */

	strP := ""
	lastSize := int(TempMBR.Mbr_tamano)
	counter := -1
	for _, partition := range TempMBR.Mbr_particion {
		counter++
		if partition.Part_correlative == 0 {
			porcentaje := utilities_test.CalcularPorcentaje(int64(partition.Part_size), int64(TempMBR.Mbr_tamano))
			lastSize -= int(partition.Part_size)
			if porcentaje > 0 {
				strP += fmt.Sprintf(`|Libre %d%%`, porcentaje)
			}
		}

		if bytes.Equal(partition.Part_type[:], compareMBR.Mbr_particion[0].Part_type[:]) {
			//println("primaria: " + string(partition.Part_name[:]))
			porcentaje := utilities_test.CalcularPorcentaje(int64(partition.Part_size), int64(TempMBR.Mbr_tamano))
			lastSize -= int(partition.Part_size)
			strP += fmt.Sprintf(`|Primaria\n%d%%`, porcentaje)
		}

		//?EBR verificacion
		if bytes.Equal(partition.Part_type[:], compareMBR.Mbr_particion[1].Part_type[:]) && EPartition {
			porcentaje := utilities_test.CalcularPorcentaje(int64(partition.Part_size), int64(TempMBR.Mbr_tamano))
			lastSize -= int(partition.Part_size)
			//println("extendida size")
			//println(partition.Part_size)
			strP += fmt.Sprintf(`|{Extendida %d%%|{`, porcentaje)
			// Validar que si no existe una particion extendida no se puede crear una logica
			//?EBR verificacion
			var x = 0
			for x < 1 {
				var TempEBR structs_test.EBR
				if err := utilities_test.ReadObject(file, &TempEBR, int64(EPartitionStart)); err != nil {
					return
				}

				if TempEBR.Part_next != -1 {
					if !bytes.Equal(TempEBR.Part_name[:], compareMBR.Mbr_particion[0].Part_name[:]) {
						porcentaje := utilities_test.CalcularPorcentaje(int64(TempEBR.Part_s), int64(TempMBR.Mbr_tamano))
						strP += fmt.Sprintf(`|EBR|Particion logica %d%%`, porcentaje)
					} else {
						porcentaje := utilities_test.CalcularPorcentaje(int64(TempEBR.Part_s), int64(TempMBR.Mbr_tamano))
						if porcentaje > 0 {
							strP += fmt.Sprintf(`|Libre %d%%`, porcentaje)
						}
					}
					EPartitionStart = int(TempEBR.Part_next)
				} else {
					if !bytes.Equal(TempEBR.Part_name[:], compareMBR.Mbr_particion[0].Part_name[:]) {
						porcentaje := utilities_test.CalcularPorcentaje(int64(TempEBR.Part_s), int64(TempMBR.Mbr_tamano))
						strP += fmt.Sprintf(`|EBR|Particion logica %d%%`, porcentaje)
					} else {
						porcentaje := utilities_test.CalcularPorcentaje(int64(TempEBR.Part_s), int64(TempMBR.Mbr_tamano))
						if porcentaje > 0 {
							strP += fmt.Sprintf(`|Libre %d%%`, porcentaje)
						}
					}
					strP += "}}"
					x = 1
				}
			}
		}
	}

	/* -------------------------------------------------------------------------- */
	/*                        EL ESPACIO RESTANTE DEL DISCO                       */
	/* -------------------------------------------------------------------------- */

	porcentaje := utilities_test.CalcularPorcentaje(int64(lastSize), int64(TempMBR.Mbr_tamano))
	fmt.Print("PORCENTAJE RESTANTE: ")
	println(porcentaje)
	if porcentaje > 0 {
		strP += fmt.Sprintf(`|Libre %d%%`, porcentaje)
	}
	strP += "}"

	/* -------------------------------------------------------------------------- */
	/*                             CODIGO DE GRAPHVIZ                             */
	/* -------------------------------------------------------------------------- */

	dotCode := fmt.Sprintf(`
		digraph G {
 			fontname="Helvetica,Arial,sans-serif"
			node [fontname="Helvetica,Arial,sans-serif"]
			edge [fontname="Helvetica,Arial,sans-serif"]
			concentrate=True;
			rankdir=TB;
			node [shape=record];

			title [label="Reporte DISK %s" shape=plaintext fontname="Helvetica,Arial,sans-serif"];

  			dsk[label="
				{MBR}%s
				}
			"];
			
			title -> dsk [style=invis];
		}`,
		letra,
		strP,
	)

	/* -------------------------------------------------------------------------- */
	/*                          ESCRITURA DEL ARCHIVO DOT                         */
	/* -------------------------------------------------------------------------- */

	dotFilePath := "./Reports/Rep2/disk_rep.dot" // Ruta donde deseas guardar el archivo DOT
	dotFile, err := os.Create(dotFilePath)
	if err != nil {
		fmt.Println("Error al crear el archivo DOT:", err)
		return
	}
	defer dotFile.Close()

	_, err = dotFile.WriteString(dotCode)
	if err != nil {
		fmt.Println("Error al escribir en el archivo DOT:", err)
		return
	}

	/* -------------------------------------------------------------------------- */
	/*                        GENERAMOS EL GRAFICO COMO PNG                       */
	/* -------------------------------------------------------------------------- */

	pngFilePath := *path // Ruta donde deseas guardar el archivo PNG
	cmd := exec.Command("dot", "-Tpng", "-o", pngFilePath, dotFilePath)
	err = cmd.Run()
	if err != nil {
		fmt.Println("Error al generar el gráfico:", err)
		return
	}

	fmt.Println("--------------------------------------------------------------------------")
	fmt.Println("                 REPORTE-DISK: GENERADO CORRECTAMENTE                     ")
	fmt.Printf("                          %s                          ", pngFilePath)
	fmt.Println("--------------------------------------------------------------------------")
}

/* -------------------------------------------------------------------------- */
/*                              3 REPORTE INODE                               */
/* -------------------------------------------------------------------------- */

func REPORT_INODE(id *string, path *string) {
	inodos = ""
	letra := string((*id)[0])
	letra = strings.ToUpper(letra)
	filepath := "./Disks/" + letra + ".dsk"

	/* -------------------------------------------------------------------------- */
	/*                              BUSCAMOS EL DISCO                             */
	/* -------------------------------------------------------------------------- */

	file, err := os.Open(filepath)
	if err != nil {
		return
	}
	defer file.Close()
	/* -------------------------------------------------------------------------- */
	/*                              CARGAMOS EL DISCO                             */
	/* -------------------------------------------------------------------------- */
	var TempMBR structs_test.MBR
	if err := utilities_test.ReadObject(file, &TempMBR, 0); err != nil {
		fmt.Println("Error reading MBR:", err)
		return
	}

	/* -------------------------------------------------------------------------- */
	/*                       BUSCAMOS LA PARTICION CON EL ID                      */
	/* -------------------------------------------------------------------------- */
	index := -1
	for i := 0; i < 4; i++ {
		if TempMBR.Mbr_particion[i].Part_size != 0 && strings.Contains(string(TempMBR.Mbr_particion[i].Part_id[:]), string((*id)[1])) {
			index = i
			break
		}
	}
	if index == -1 {
		fmt.Println("Partition not found")
		return
	}

	/* -------------------------------------------------------------------------- */
	/*                           CARGAMOS EL SUPERBLOQUE                          */
	/* -------------------------------------------------------------------------- */
	var tempSuperblock structs_test.Superblock
	if err := utilities_test.ReadObject(file, &tempSuperblock, int64(TempMBR.Mbr_particion[index].Part_start)); err != nil {
		fmt.Println("Error reading superblock:", err)
		return
	}

	CargarArbol(tempSuperblock, file, 0)

	/* -------------------------------------------------------------------------- */
	/*                               CODIGO DOT BASE                              */
	/* -------------------------------------------------------------------------- */

	dotCode := fmt.Sprintf(`
		digraph G {
 			fontname="Helvetica,Arial,sans-serif"
			node [fontname="Helvetica,Arial,sans-serif"]
			edge [fontname="Helvetica,Arial,sans-serif"]
			concentrate=True;
			rankdir=LR;
			node [shape=record];
			
			%s
		}`,
		inodos,
	)

	/* -------------------------------------------------------------------------- */
	/*                          ESCRIBIMOS EL CODIGO DOT                          */
	/* -------------------------------------------------------------------------- */

	dotFilePath := "./Reports/Rep3/inodes_rep.dot" // Ruta donde deseas guardar el archivo DOT
	dotFile, err := os.Create(dotFilePath)
	if err != nil {
		fmt.Println("Error al crear el archivo DOT:", err)
		return
	}
	defer dotFile.Close()

	_, err = dotFile.WriteString(dotCode)
	if err != nil {
		fmt.Println("Error al escribir en el archivo DOT:", err)
		return
	}

	/* -------------------------------------------------------------------------- */
	/*                         GENERAMOS EL GRAFO COMO PNG                        */
	/* -------------------------------------------------------------------------- */

	pngFilePath := *path // Ruta donde deseas guardar el archivo PNG
	cmd := exec.Command("dot", "-Tsvg", "-o", pngFilePath, dotFilePath)
	err = cmd.Run()
	if err != nil {
		fmt.Println("Error al generar el gráfico:", err)
		return
	}

	fmt.Println("--------------------------------------------------------------------------")
	fmt.Println("               REPORTE-INODOS: GENERADO CORRECTAMENTE                    ")
	fmt.Printf("                          %s                          ", pngFilePath)
	fmt.Println("--------------------------------------------------------------------------")
}

/* -------------------------------------------------------------------------- */
/*                          4 REPORTE JOURNALING                             */
/* -------------------------------------------------------------------------- */

func REPORT_JOURNALING(id *string, path *string) {
}

/* -------------------------------------------------------------------------- */
/*                              5 REPORTE BLOCK                               */
/* -------------------------------------------------------------------------- */

func REPORT_BLOCK(id *string, path *string) {
	bloques = ""
	letra := string((*id)[0])
	letra = strings.ToUpper(letra)
	filepath := "./Disks/" + letra + ".dsk"

	/* -------------------------------------------------------------------------- */
	/*                              BUSCAMOS EL DISCO                             */
	/* -------------------------------------------------------------------------- */

	file, err := os.Open(filepath)
	if err != nil {
		return
	}
	defer file.Close()
	/* -------------------------------------------------------------------------- */
	/*                              CARGAMOS EL DISCO                             */
	/* -------------------------------------------------------------------------- */
	var TempMBR structs_test.MBR
	if err := utilities_test.ReadObject(file, &TempMBR, 0); err != nil {
		fmt.Println("Error reading MBR:", err)
		return
	}

	/* -------------------------------------------------------------------------- */
	/*                       BUSCAMOS LA PARTICION CON EL ID                      */
	/* -------------------------------------------------------------------------- */
	index := -1
	for i := 0; i < 4; i++ {
		if TempMBR.Mbr_particion[i].Part_size != 0 && strings.Contains(string(TempMBR.Mbr_particion[i].Part_id[:]), string((*id)[1])) {
			index = i
			break
		}
	}
	if index == -1 {
		fmt.Println("Partition not found")
		return
	}

	/* -------------------------------------------------------------------------- */
	/*                           CARGAMOS EL SUPERBLOQUE                          */
	/* -------------------------------------------------------------------------- */
	var tempSuperblock structs_test.Superblock
	if err := utilities_test.ReadObject(file, &tempSuperblock, int64(TempMBR.Mbr_particion[index].Part_start)); err != nil {
		fmt.Println("Error reading superblock:", err)
		return
	}

	CargarArbol(tempSuperblock, file, 0)
	/* -------------------------------------------------------------------------- */
	/*                               CODIGO DOT BASE                              */
	/* -------------------------------------------------------------------------- */

	dotCode := fmt.Sprintf(`
		digraph G {
 			fontname="Helvetica,Arial,sans-serif"
			node [fontname="Helvetica,Arial,sans-serif"]
			edge [fontname="Helvetica,Arial,sans-serif"]
			concentrate=True;
			rankdir=LR;
			node [shape=record];
			title [label="Reporte INODOS" shape=plaintext fontname="Helvetica,Arial,sans-serif"];
			%s
		}`,
		bloques,
	)

	/* -------------------------------------------------------------------------- */
	/*                          ESCRIBIMOS EL CODIGO DOT                          */
	/* -------------------------------------------------------------------------- */

	dotFilePath := "./Reports/Rep5/bloques_rep.dot" // Ruta donde deseas guardar el archivo DOT
	dotFile, err := os.Create(dotFilePath)
	if err != nil {
		fmt.Println("Error al crear el archivo DOT:", err)
		return
	}
	defer dotFile.Close()

	_, err = dotFile.WriteString(dotCode)
	if err != nil {
		fmt.Println("Error al escribir en el archivo DOT:", err)
		return
	}

	/* -------------------------------------------------------------------------- */
	/*                         GENERAMOS EL GRAFO COMO PNG                        */
	/* -------------------------------------------------------------------------- */

	pngFilePath := *path // Ruta donde deseas guardar el archivo PNG
	cmd := exec.Command("dot", "-Tsvg", "-o", pngFilePath, dotFilePath)
	err = cmd.Run()
	if err != nil {
		fmt.Println("Error al generar el gráfico:", err)
		return
	}

	fmt.Println("--------------------------------------------------------------------------")
	fmt.Println("               REPORTE-BLOQUES: GENERADO CORRECTAMENTE                    ")
	fmt.Printf("                          %s                          ", pngFilePath)
	fmt.Println("--------------------------------------------------------------------------")
}

/* -------------------------------------------------------------------------- */
/*                            6 REPORTE BM_INODE                              */
/* -------------------------------------------------------------------------- */

func REPORT_BM_INODE(id *string, path *string) {
	letra := string((*id)[0])
	letra = strings.ToUpper(letra)
	filepath := "./Disks/" + letra + ".dsk"

	/* -------------------------------------------------------------------------- */
	/*                              BUSCAMOS EL DISCO                             */
	/* -------------------------------------------------------------------------- */

	file, err := os.Open(filepath)
	if err != nil {
		return
	}
	defer file.Close()
	/* -------------------------------------------------------------------------- */
	/*                              CARGAMOS EL DISCO                             */
	/* -------------------------------------------------------------------------- */
	var TempMBR structs_test.MBR
	if err := utilities_test.ReadObject(file, &TempMBR, 0); err != nil {
		fmt.Println("Error reading MBR:", err)
		return
	}

	/* -------------------------------------------------------------------------- */
	/*                       BUSCAMOS LA PARTICION CON EL ID                      */
	/* -------------------------------------------------------------------------- */
	index := -1
	for i := 0; i < 4; i++ {
		if TempMBR.Mbr_particion[i].Part_size != 0 && strings.Contains(string(TempMBR.Mbr_particion[i].Part_id[:]), string((*id)[1])) {
			index = i
			break
		}
	}
	if index == -1 {
		fmt.Println("Partition not found")
		return
	}

	/* -------------------------------------------------------------------------- */
	/*                           CARGAMOS EL SUPERBLOQUE                          */
	/* -------------------------------------------------------------------------- */
	var tempSuperblock structs_test.Superblock
	if err := utilities_test.ReadObject(file, &tempSuperblock, int64(TempMBR.Mbr_particion[index].Part_start)); err != nil {
		fmt.Println("Error reading superblock:", err)
		return
	}

	var bitmap string
	print("CREADOS: ")
	println(tempSuperblock.S_inodes_count)
	print("LIBRES: ")
	println(tempSuperblock.S_free_inodes_count)
	cont := 0
	for i := 0; i <= int(tempSuperblock.S_inodes_count); i++ {
		bitmap += " 1 "
		cont++
		if cont == 20 {
			bitmap += "\n"
			cont = 0
		}
	}
	for i := 0; i < int(tempSuperblock.S_free_inodes_count); i++ {
		if cont == 20 {
			bitmap += "\n"
			cont = 0
		}
		bitmap += " 0 "
		cont++
	}

	file1, err := os.Create(*path)
	if err != nil {
		fmt.Println("Error al crear el archivo:", err)
		return
	}
	defer file1.Close()
	//println(bitmap)
	_, err = file1.WriteString(bitmap)
	if err != nil {
		fmt.Println("Error al escribir en el archivo:", err)
		return
	}
}

/* -------------------------------------------------------------------------- */
/*                             7 REPORTE BM_BLOC                              */
/* -------------------------------------------------------------------------- */

func REPORT_BM_BLOCK(id *string, path *string) {
	letra := string((*id)[0])
	letra = strings.ToUpper(letra)
	filepath := "./Disks/" + letra + ".dsk"

	/* -------------------------------------------------------------------------- */
	/*                              BUSCAMOS EL DISCO                             */
	/* -------------------------------------------------------------------------- */

	file, err := os.Open(filepath)
	if err != nil {
		return
	}
	defer file.Close()
	/* -------------------------------------------------------------------------- */
	/*                              CARGAMOS EL DISCO                             */
	/* -------------------------------------------------------------------------- */
	var TempMBR structs_test.MBR
	if err := utilities_test.ReadObject(file, &TempMBR, 0); err != nil {
		fmt.Println("Error reading MBR:", err)
		return
	}

	/* -------------------------------------------------------------------------- */
	/*                       BUSCAMOS LA PARTICION CON EL ID                      */
	/* -------------------------------------------------------------------------- */
	index := -1
	for i := 0; i < 4; i++ {
		if TempMBR.Mbr_particion[i].Part_size != 0 && strings.Contains(string(TempMBR.Mbr_particion[i].Part_id[:]), string((*id)[1])) {
			index = i
			break
		}
	}
	if index == -1 {
		fmt.Println("Partition not found")
		return
	}

	/* -------------------------------------------------------------------------- */
	/*                           CARGAMOS EL SUPERBLOQUE                          */
	/* -------------------------------------------------------------------------- */
	var tempSuperblock structs_test.Superblock
	if err := utilities_test.ReadObject(file, &tempSuperblock, int64(TempMBR.Mbr_particion[index].Part_start)); err != nil {
		fmt.Println("Error reading superblock:", err)
		return
	}

	var bitmap string
	print("CREADOS: ")
	println(tempSuperblock.S_blocks_count)
	print("LIBRES: ")
	println(tempSuperblock.S_free_blocks_count)
	cont := 0
	for i := 0; i <= int(tempSuperblock.S_blocks_count); i++ {
		bitmap += " 1 "
		cont++
		if cont == 20 {
			bitmap += "\n"
			cont = 0
		}
	}
	for i := 0; i < int(tempSuperblock.S_free_blocks_count); i++ {
		if cont == 20 {
			bitmap += "\n"
			cont = 0
		}
		bitmap += " 0 "
		cont++
	}

	file1, err := os.Create(*path)
	if err != nil {
		fmt.Println("Error al crear el archivo:", err)
		return
	}
	defer file1.Close()
	//println(bitmap)
	_, err = file1.WriteString(bitmap)
	if err != nil {
		fmt.Println("Error al escribir en el archivo:", err)
		return
	}
}

/* -------------------------------------------------------------------------- */
/*                              8 REPORTE TREE                                */
/* -------------------------------------------------------------------------- */
func REPORT_TREE(id *string, path *string) {
	/* -------------------------------------------------------------------------- */
	/*                              BUSCAMOS EL DISCO                             */
	/* -------------------------------------------------------------------------- */
	filepath := "./Disks/" + letra + ".dsk"
	file, err := utilities_test.OpenFile(filepath)
	if err != nil {
		fmt.Println("Error opening disk file:", err)
		return
	}
	defer file.Close()

	/* -------------------------------------------------------------------------- */
	/*                              CARGAMOS EL DISCO                             */
	/* -------------------------------------------------------------------------- */
	var TempMBR structs_test.MBR
	if err := utilities_test.ReadObject(file, &TempMBR, 0); err != nil {
		fmt.Println("Error reading MBR:", err)
		return
	}

	/* -------------------------------------------------------------------------- */
	/*                       BUSCAMOS LA PARTICION CON EL ID                      */
	/* -------------------------------------------------------------------------- */
	index := -1
	for i := 0; i < 4; i++ {
		if TempMBR.Mbr_particion[i].Part_size != 0 && strings.Contains(string(TempMBR.Mbr_particion[i].Part_id[:]), string((*id)[1])) {
			index = i
			break
		}
	}
	if index == -1 {
		fmt.Println("Partition not found")
		return
	}

	/* -------------------------------------------------------------------------- */
	/*                           CARGAMOS EL SUPERBLOQUE                          */
	/* -------------------------------------------------------------------------- */
	var tempSuperblock structs_test.Superblock
	if err := utilities_test.ReadObject(file, &tempSuperblock, int64(TempMBR.Mbr_particion[index].Part_start)); err != nil {
		fmt.Println("Error reading superblock:", err)
		return
	}

	CargarArbol(tempSuperblock, file, 0)

	/* -------------------------------------------------------------------------- */
	/*                               CODIGO DOT BASE                              */
	/* -------------------------------------------------------------------------- */
	tree += inodos
	tree += bloques
	tree += relaciones
	dotCode := fmt.Sprintf(`
		digraph G {
 			fontname="Helvetica,Arial,sans-serif"
			node [fontname="Helvetica,Arial,sans-serif"]
			edge [fontname="Helvetica,Arial,sans-serif"]
			concentrate=True;
			rankdir=LR;
			node [shape=record];
			%s
		}`,
		tree,
	)

	/* -------------------------------------------------------------------------- */
	/*                          ESCRIBIMOS EL CODIGO DOT                          */
	/* -------------------------------------------------------------------------- */

	dotFilePath := "./Reports/Rep8/tree_rep.dot" // Ruta donde deseas guardar el archivo DOT
	dotFile, err := os.Create(dotFilePath)
	if err != nil {
		fmt.Println("Error al crear el archivo DOT:", err)
		return
	}
	defer dotFile.Close()

	_, err = dotFile.WriteString(dotCode)
	if err != nil {
		fmt.Println("Error al escribir en el archivo DOT:", err)
		return
	}

	/* -------------------------------------------------------------------------- */
	/*                         GENERAMOS EL GRAFO COMO PNG                        */
	/* -------------------------------------------------------------------------- */

	pngFilePath := *path // Ruta donde deseas guardar el archivo PNG
	cmd := exec.Command("dot", "-Tsvg", "-o", pngFilePath, dotFilePath)
	err = cmd.Run()
	if err != nil {
		fmt.Println("Error al generar el gráfico:", err)
		return
	}

	fmt.Println("--------------------------------------------------------------------------")
	fmt.Println("               REPORTE-TREE: GENERADO CORRECTAMENTE                    ")
	fmt.Printf("                          %s                          ", pngFilePath)
	fmt.Println("--------------------------------------------------------------------------")

}

/* -------------------------------------------------------------------------- */
/*                               9 REPORTE SB                                 */
/* -------------------------------------------------------------------------- */

func REPORT_SB(id *string, path *string) {
	letra := string((*id)[0])
	letra = strings.ToUpper(letra)

	/* -------------------------------------------------------------------------- */
	/*                              BUSCAMOS EL DISCO                             */
	/* -------------------------------------------------------------------------- */
	filepath := "./Disks/" + letra + ".dsk"
	file, err := os.Open(filepath)
	if err != nil {
		return
	}
	defer file.Close()
	/* -------------------------------------------------------------------------- */
	/*                              CARGAMOS EL DISCO                             */
	/* -------------------------------------------------------------------------- */
	var TempMBR structs_test.MBR
	if err := utilities_test.ReadObject(file, &TempMBR, 0); err != nil {
		fmt.Println("Error reading MBR:", err)
		return
	}

	/* -------------------------------------------------------------------------- */
	/*                       BUSCAMOS LA PARTICION CON EL ID                      */
	/* -------------------------------------------------------------------------- */
	index := -1
	for i := 0; i < 4; i++ {
		if TempMBR.Mbr_particion[i].Part_size != 0 && strings.Contains(string(TempMBR.Mbr_particion[i].Part_id[:]), string((*id)[1])) {
			index = i
			break
		}
	}
	if index == -1 {
		fmt.Println("Partition not found")
		return
	}

	/* -------------------------------------------------------------------------- */
	/*                           CARGAMOS EL SUPERBLOQUE                          */
	/* -------------------------------------------------------------------------- */
	var tempSuperblock structs_test.Superblock
	if err := utilities_test.ReadObject(file, &tempSuperblock, int64(TempMBR.Mbr_particion[index].Part_start)); err != nil {
		fmt.Println("Error reading superblock:", err)
		return
	}

	/* -------------------------------------------------------------------------- */
	/*                      GENERAMOS EL REPORTE EN GRAPHVIZ                      */
	/* -------------------------------------------------------------------------- */

	dotCode := fmt.Sprintf(`
		digraph G {
 			fontname="Helvetica,Arial,sans-serif"
			node [fontname="Helvetica,Arial,sans-serif"]
			edge [fontname="Helvetica,Arial,sans-serif"]
			concentrate=True;
			rankdir=TB;
			node [shape=record];

			title [label="Reporte SUPERBLOCK" shape=plaintext fontname="Helvetica,Arial,sans-serif"];

  			sb[label="
				{Superblock|
					{S_filesystem_type|%d}
					|{S_inodes_count|%d}
					|{S_blocks_count|%d}
					|{S_free_blocks_count|%d}
					|{S_free_inodes_count|%d}
					|{S_mtime|%s}
					|{S_umtime|%s}
					|{S_mnt_count|%d}
					|{S_magic|%d}
					|{S_inode_size|%d}
					|{S_block_size|%d}
					|{S_fist_ino|%d}
					|{S_first_blo|%d}
					|{S_bm_inode_start|%d}
					|{S_bm_block_start|%d}
					|{S_inode_start|%d}
					|{S_block_start|%d}
				}
			"];
			

			title -> sb [style=invis];
		}`,
		int(tempSuperblock.S_filesystem_type),
		int(tempSuperblock.S_inodes_count),
		int(tempSuperblock.S_blocks_count),
		int(tempSuperblock.S_free_blocks_count),
		int(tempSuperblock.S_free_inodes_count),
		tempSuperblock.S_mtime[:],
		tempSuperblock.S_umtime[:],
		int(tempSuperblock.S_mnt_count),
		int(tempSuperblock.S_magic),
		int(tempSuperblock.S_inode_size),
		int(tempSuperblock.S_block_size),
		int(tempSuperblock.S_fist_ino),
		int(tempSuperblock.S_first_blo),
		int(tempSuperblock.S_bm_inode_start),
		int(tempSuperblock.S_bm_block_start),
		int(tempSuperblock.S_inode_start),
		int(tempSuperblock.S_block_start),
	)

	// Escribir el contenido en el archivo DOT
	dotFilePath := "./Reports/Rep9/sb_rep.dot" // Ruta donde deseas guardar el archivo DOT
	dotFile, err := os.Create(dotFilePath)
	if err != nil {
		fmt.Println("Error al crear el archivo DOT:", err)
		return
	}
	defer dotFile.Close()

	_, err = dotFile.WriteString(dotCode)
	if err != nil {
		fmt.Println("Error al escribir en el archivo DOT:", err)
		return
	}

	// Llamar a Graphviz para generar el gráfico
	pngFilePath := *path // Ruta donde deseas guardar el archivo PNG
	cmd := exec.Command("dot", "-Tpng", "-o", pngFilePath, dotFilePath)
	err = cmd.Run()
	if err != nil {
		fmt.Println("Error al generar el gráfico:", err)
		return
	}

	fmt.Println("Reporte SUPERBLOCK generado en", pngFilePath)
}

/* -------------------------------------------------------------------------- */
/*                              10 REPORTE FILE                                */
/* -------------------------------------------------------------------------- */

func REPORT_FILE(id *string, path *string, ruta *string) {
}

/* -------------------------------------------------------------------------- */
/*                              11 REPORTE LS                                 */
/* -------------------------------------------------------------------------- */

func REPORT_LS(id *string, path *string, ruta *string) {
}

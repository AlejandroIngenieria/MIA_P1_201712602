package functions_test

import (
	structs_test "P1/Structs"
	utilities_test "P1/Utilities"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

// ?                     			REPORTES
func ProcessREP(input string, name *string, path *string, id *string, ruta *string) {
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
			fmt.Println("Error: Flag not found")
		}
	}
}

func GenerateReports(name *string, path *string, id *string, ruta *string) {

	switch *name {
	case "mbr":
		REPORT_MBR(id, path)
	case "disk":
		REPORT_DISK(id, path)
	case "inode":
		REPORT_INODE(id, path)
	case "Journaling":
		REPORT_JOURNALING()
	case "block":
		REPORT_BLOCK(id, path)
	case "bm_inode":
		REPORT_BM_INODE(id, path)
	case "bm_block":
		REPORT_BM_BLOCK(id, path)
	case "tree":
		REPORT_TREE()
	case "sb":
		REPORT_SB(id, path)
	case "file":
		REPORT_FILE(id, path, ruta)
	case "ls":
		REPORT_LS(id, path, ruta)
	default:
		println("Reporte no reconocido:", *name)
	}
}

/* -------------------------------------------------------------------------- */
/*                                 REPORTE MBR                                */
/* -------------------------------------------------------------------------- */
func REPORT_MBR(id *string, path *string) {
	letra := string((*id)[0])
	letra = strings.ToUpper(letra)

	filepath := "./Disks/" + letra + ".dsk"
	file, err := os.Open(filepath)
	if err != nil {
		return
	}
	defer file.Close()

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

	for _, partition := range TempMBR.Mbr_particion {
		if bytes.Equal(partition.Part_type[:], compareMBR.Mbr_particion[1].Part_type[:]) {
			EPartition = true
			EPartitionStart = int(partition.Part_start)
		}
	}

	strP := ""
	strE := ""

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
					print("fit logica")
					println(string(TempEBR.Part_fit[:]))
					EPartitionStart = int(TempEBR.Part_next)
				} else {
					print("fit logica")
					println(string(TempEBR.Part_fit[:]))
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

	//structs_test.PrintMBR(TempMBR)

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

	// Escribir el contenido en el archivo DOT
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

	// Llamar a Graphviz para generar el gráfico
	pngFilePath := *path // Ruta donde deseas guardar el archivo PNG
	cmd := exec.Command("dot", "-Tpng", "-o", pngFilePath, dotFilePath)
	err = cmd.Run()
	if err != nil {
		fmt.Println("Error al generar el gráfico:", err)
		return
	}

	fmt.Println("Reporte MBR, EBR generado en", pngFilePath)
}

/* -------------------------------------------------------------------------- */
/*                                REPORTE DISK                                */
/* -------------------------------------------------------------------------- */

func REPORT_DISK(id *string, path *string) {
	letra := string((*id)[0])
	letra = strings.ToUpper(letra)

	filepath := "./Disks/" + letra + ".dsk"
	file, err := os.Open(filepath)
	if err != nil {
		return
	}
	defer file.Close()

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

	for _, partition := range TempMBR.Mbr_particion {
		if bytes.Equal(partition.Part_type[:], compareMBR.Mbr_particion[1].Part_type[:]) {
			EPartition = true
			EPartitionStart = int(partition.Part_start)
		}
	}

	strP := ""

	for _, partition := range TempMBR.Mbr_particion {
		if partition.Part_correlative == 0 {
			porcentaje := utilities_test.CalcularPorcentaje(int64(partition.Part_size), int64(TempMBR.Mbr_tamano))
			strP += fmt.Sprintf(`|Libre\n%d%%`, porcentaje)
		}

		if bytes.Equal(partition.Part_type[:], compareMBR.Mbr_particion[0].Part_type[:]) {
			//println("primaria: " + string(partition.Part_name[:]))
			porcentaje := utilities_test.CalcularPorcentaje(int64(partition.Part_size), int64(TempMBR.Mbr_tamano))
			strP += fmt.Sprintf(`|Primaria\n%d%%`, porcentaje)
		}

		//?EBR verificacion
		if bytes.Equal(partition.Part_type[:], compareMBR.Mbr_particion[1].Part_type[:]) && EPartition {
			porcentaje := utilities_test.CalcularPorcentaje(int64(partition.Part_size), int64(TempMBR.Mbr_tamano))
			println("extendida size")
			println(partition.Part_size)
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
						strP += fmt.Sprintf(`|Libre %d%%`, porcentaje)
					}
					EPartitionStart = int(TempEBR.Part_next)
				} else {
					if !bytes.Equal(TempEBR.Part_name[:], compareMBR.Mbr_particion[0].Part_name[:]) {
						porcentaje := utilities_test.CalcularPorcentaje(int64(TempEBR.Part_s), int64(TempMBR.Mbr_tamano))
						strP += fmt.Sprintf(`|EBR|Particion logica %d%%`, porcentaje)
					} else {
						porcentaje := utilities_test.CalcularPorcentaje(int64(TempEBR.Part_s), int64(TempMBR.Mbr_tamano))
						strP += fmt.Sprintf(`|Libre %d%%`, porcentaje)
					}
					x = 1
				}
			}
			strP += "}}"
		}

	}

	//structs_test.PrintMBR(TempMBR)

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

	// Escribir el contenido en el archivo DOT
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

	// Llamar a Graphviz para generar el gráfico
	pngFilePath := *path // Ruta donde deseas guardar el archivo PNG
	cmd := exec.Command("dot", "-Tpng", "-o", pngFilePath, dotFilePath)
	err = cmd.Run()
	if err != nil {
		fmt.Println("Error al generar el gráfico:", err)
		return
	}

	fmt.Println("Reporte DISK generado en", pngFilePath)
	println("MBR")
	structs_test.PrintMBR(TempMBR)
}

/* -------------------------------------------------------------------------- */
/*                                REPORTE INODE                               */
/* -------------------------------------------------------------------------- */

func REPORT_INODE(id *string, path *string) {

}

/* -------------------------------------------------------------------------- */
/*                                REPORTE BLOCK                               */
/* -------------------------------------------------------------------------- */

func REPORT_BLOCK(id *string, path *string) {

}

/* -------------------------------------------------------------------------- */
/*                              REPORTE BM_INODE                              */
/* -------------------------------------------------------------------------- */

func REPORT_BM_INODE(id *string, path *string) {

}

/* -------------------------------------------------------------------------- */
/*                               REPORTE BM_BLOC                              */
/* -------------------------------------------------------------------------- */

func REPORT_BM_BLOCK(id *string, path *string) {

}

/* -------------------------------------------------------------------------- */
/*                                REPORTE TREE                                */
/* -------------------------------------------------------------------------- */

func REPORT_TREE() {

}

/* -------------------------------------------------------------------------- */
/*                                 REPORTE SB                                 */
/* -------------------------------------------------------------------------- */

func REPORT_SB(id *string, path *string) {

}

/* -------------------------------------------------------------------------- */
/*                                REPORTE FILE                                */
/* -------------------------------------------------------------------------- */

func REPORT_FILE(id *string, path *string, ruta *string) {

}

/* -------------------------------------------------------------------------- */
/*                                 REPORTE LS                                 */
/* -------------------------------------------------------------------------- */

func REPORT_LS(id *string, path *string, ruta *string) {

}

/* -------------------------------------------------------------------------- */
/*                             REPORTE JOURNALING                             */
/* -------------------------------------------------------------------------- */

func REPORT_JOURNALING() {

}

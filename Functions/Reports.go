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
		REPORT_MBR(path, id)
	case "disk":
	case "inode":
	case "Journaling":
	case "block":
	case "bm_inode":
	case "bm_block":
	case "tree":
	case "sb":
	case "file":
	case "ls":
	default:
		println("Reporte no reconocido:", *name)
	}
}

/* -------------------------------------------------------------------------- */
/*                                 REPORTE MBR                                */
/* -------------------------------------------------------------------------- */
func REPORT_MBR(path *string, id *string) {
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

	// Verificar si el nombre de la partición ya está en uso
	strP := ""

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
			for i := 0; i < 4; i++ {
				//?EBR verificacion
				var x = 0
				for x < 1 {
					var TempEBR structs_test.EBR
					if err := utilities_test.ReadObject(file, &TempEBR, int64(EPartitionStart)); err != nil {
						return
					}

					if EPartitionStart != 0 && TempEBR.Part_next != -1 {
						partNameClean := strings.Trim(string(TempEBR.Part_name[:]), "\x00")
						strP += fmt.Sprintf(`
					|Particion Logica
					|{part_status|%s}
					|{part_next|%d}
					|{part_fit|%s}
					|{part_start|%d}
					|{part_size|%d}
					|{part_name|%s}`,
							TempEBR.Part_mount[:],
							TempEBR.Part_next,
							TempEBR.Part_fit[:],
							TempEBR.Part_start,
							TempEBR.Part_s,
							partNameClean,
						)
						EPartitionStart = int(TempEBR.Part_next)
					} else {
						x = 1
					}
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
  			mbr[label="
				{MBR: %s.dsk|
					{mbr_tamaño|%d}
					|{mbr_fecha_creacion|%s}
					|{mbr_disk_signature|%d}
					%s
				}
			"];
		}`,
		letra,
		TempMBR.Mbr_tamano,
		TempMBR.Mbr_fecha_creacion,
		TempMBR.Mbr_dsk_signature,
		strP,
	)

	// Guardar el código DOT en un archivo temporal
	// tmpFile, err := os.CreateTemp("", "graph-*.dot")
	// if err != nil {
	// 	return
	// }
	// defer os.Remove(tmpFile.Name())

	// _, err = tmpFile.WriteString(dotCode)
	// if err != nil {
	// 	return
	// }

	// Llamar a Graphviz desde la línea de comandos para renderizar el gráfico
	// println("path: ", *path)
	// cmd := exec.Command("dot", "-Tpng", "-o", *path, tmpFile.Name())
	// err = cmd.Run()
	// if err != nil {
	// 	println("Error: no se genero el reporte")
	// 	return
	// }
	// fmt.Println("Reporte MBR generado")

	// Escribir el contenido en el archivo DOT
	dotFilePath := "./Reports/Rep1/graph.dot" // Ruta donde deseas guardar el archivo DOT
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
	pngFilePath := *path // Ruta donde deseas guardar el archivo SVG
	cmd := exec.Command("dot", "-Tpng", "-o", pngFilePath, dotFilePath)
	err = cmd.Run()
	if err != nil {
		fmt.Println("Error al generar el gráfico:", err)
		return
	}

	fmt.Println("Reporte MBR generado en", pngFilePath)
}

/* -------------------------------------------------------------------------- */
/*                                REPORTE DISK                                */
/* -------------------------------------------------------------------------- */

func REPORT_DISK() {

}

/* -------------------------------------------------------------------------- */
/*                                REPORTE INODE                               */
/* -------------------------------------------------------------------------- */

func REPORT_INODE() {

}

/* -------------------------------------------------------------------------- */
/*                                REPORTE BLOCK                               */
/* -------------------------------------------------------------------------- */

func REPORT_BLOCK() {

}

/* -------------------------------------------------------------------------- */
/*                              REPORTE BM_INODE                              */
/* -------------------------------------------------------------------------- */

func REPORT_BM_INODE() {

}

/* -------------------------------------------------------------------------- */
/*                               REPORTE BM_BLOC                              */
/* -------------------------------------------------------------------------- */

func REPORT_BM_BLOCK() {

}

/* -------------------------------------------------------------------------- */
/*                                REPORTE TREE                                */
/* -------------------------------------------------------------------------- */

func REPORT_TREE() {

}

/* -------------------------------------------------------------------------- */
/*                                 REPORTE SB                                 */
/* -------------------------------------------------------------------------- */

func REPORT_SB() {

}

/* -------------------------------------------------------------------------- */
/*                                REPORTE FILE                                */
/* -------------------------------------------------------------------------- */

func REPORT_FILE() {

}

/* -------------------------------------------------------------------------- */
/*                                 REPORTE LS                                 */
/* -------------------------------------------------------------------------- */

func REPORT_LS() {

}

/* -------------------------------------------------------------------------- */
/*                             REPORTE JOURNALING                             */
/* -------------------------------------------------------------------------- */

func REPORT_JOURNALING() {

}

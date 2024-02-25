package structs_test
import "fmt"

// ? DISCOS extension .dsk

// Master Boot Record (MBR)
type MBR struct {
	Mbr_tamano         int32
	Mbr_fecha_creacion [10]byte
	Mbr_dsk_signature  int32
	Dsk_fit            [1]byte
	Mbr_particion      [4]Partition
}

func PrintMBR(data MBR){
	fmt.Printf("CreationDate: %s, fit: %s, size: %d \n", string(data.Mbr_fecha_creacion[:]), string(data.Dsk_fit[:]), data.Mbr_tamano)
	for i := 0; i < 4; i++ {
		fmt.Printf("Partition %d, Name: %s, Tipo: %s, Start: %d, Size: %d \n",
		 i, 
		string(data.Mbr_particion[i].Part_name[:]), 
		string(data.Mbr_particion[i].Part_type[:]), 
		data.Mbr_particion[i].Part_start, 
		data.Mbr_particion[i].Part_size,
		)
	}
}

// Partition
type Partition struct {
	Part_status      [1]byte
	Part_type        [1]byte
	Part_fit         [1]byte
	Part_start       int32
	Part_size        int32
	Part_name        [16]byte
	Part_correlative int32
	Part_id          [4]byte
}

// Extended Boot Record (EBR)
type EBR struct {
	Part_mount [1]byte
	Part_fit   [1]byte
	Part_start int32
	Part_s     int32
	Part_next  int32
	Part_name  [16]byte
}

// ? CARPETAS Y ARCHIVOS (EXT3|EXT2)
// Super bloque
type S_block struct {
	S_filesystem_type   int32
	S_inodes_count      int32
	S_blocks_count      int32
	S_free_blocks_count int32
	S_free_inodes_count int32
	S_mtime             [10]byte
	S_umtime            [10]byte
	S_mnt_count         int32
	S_magic             int32
	S_inode_s           int32
	S_block_s           int32
	S_firts_ino         int32
	S_firts_blo         int32
	S_bm_inode_start    int32
	S_bm_block_start    int32
	S_inode_start       int32
	S_block_start       int32
}

// Inodos
type Inode struct {
	I_uid   int32
	I_gid   int32
	I_s     int32
	I_atime [10]byte
	I_ctime [10]byte
	I_mtime [10]byte
	I_block int32
	I_type  [1]byte
	I_perm  [3]byte
}

// ? BLOQUES
// Bloque de carpetas
type B_files struct {
	B_content [4]Content
}

type Content struct {
	B_name  [12]byte
	B_inodo int32
}

// Bloque de archivos
type B_docs struct {
	B_content [64]byte
}

// Bloque de apuntadores
type B_pointer struct {
	B_pointers [16]int
}

package structs_test

import (
	"time"
)

// ? DISCOS extension .dsk
type Disk struct {
	Capacidad [1024]byte
	Mbr       MBR
}

// Master Boot Record (MBR)
type MBR struct {
	Mbr_tamano         int32
	Mbr_fecha_creacion [10]byte
	Mbr_dsk_signature  int32
	Dsk_fit            [1]byte
	Mbr_particion      [4]Partition
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
	Part_mount rune
	Part_fit   rune
	Part_start int
	Part_s     int
	Part_next  int
	Part_name  [16]byte
}

// ? CARPETAS Y ARCHIVOS (EXT3|EXT2)
// Super bloque
type S_block struct {
	S_filesystem_type   int
	S_inodes_count      int
	S_blocks_count      int
	S_free_blocks_count int
	S_free_inodes_count int
	S_mtime             time.Time
	S_umtime            time.Time
	S_mnt_count         int
	S_magic             int
	S_inode_s           int
	S_block_s           int
	S_firts_ino         int
	S_firts_blo         int
	S_bm_inode_start    int
	S_bm_block_start    int
	S_inode_start       int
	S_block_start       int
}

// Inodos
type Inode struct {
	I_uid   int
	I_gid   int
	I_s     int
	I_atime time.Time
	I_ctime time.Time
	I_mtime time.Time
	I_block int
	I_type  rune
	I_perm  [3]byte
}

// ? BLOQUES
// Bloque de carpetas
type B_files struct {
	B_content [4]Content
}

type Content struct {
	B_name  [12]byte
	B_inodo int
}

// Bloque de archivos
type B_docs struct {
	B_content [64]byte
}

// Bloque de apuntadores
type B_pointer struct {
	B_pointers [16]int
}

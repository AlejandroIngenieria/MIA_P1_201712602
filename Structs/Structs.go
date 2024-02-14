package main

import (
	"time"
)

// ? DISCOS
type disk struct {
	capacidad [1024]byte
	MBR       mbr
}

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

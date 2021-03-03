package tid

import (
	"crypto/md5"
	"encoding/binary"
	"errors"
	"fmt"
	"os"
	"sync"
	"unsafe"
)

const (
	byteLen = 12

	maxSeq uint32 = 1 << 24 - 1

	// encoding stores a custom version of the base32 encoding with lower case
	// letters.
	encoding = "0123456789abcdefghijklmnopqrstuv"
)

var (
	encodedLen = 20

	machineID = getMachineID()

	pid = os.Getpid()

	seq uint32 = 0

	lastSecond uint32

	dgMutex = &sync.Mutex{}

)

func getMachineID() []byte {
	id := make([]byte, 3)
	hostname, err := os.Hostname()
	if err != nil {
		panic(fmt.Errorf("cannot get hostname: %v", err))
	}
	hw := md5.New()
	hw.Write([]byte(hostname))
	copy(id, hw.Sum(nil))
	return id
}


type TIDGenerator interface {

	get() (string, error)

}

func NewDefaultGenerator(address string) TIDGenerator {
	return &defaultGenerator{NewDefaultTimer(address)}
}

func NewDefaultGeneratorWithTimer(timer TIDTimer) TIDGenerator {
	return &defaultGenerator{timer}
}

type defaultGenerator struct {
	timer TIDTimer
}

func (dg *defaultGenerator) get() (string, error) {
	var (
		nowSecond uint32
		nowSeq uint32
		err error
	)
	if nowSecond, err = dg.timer.Second(); err != nil  {
		return "", err
	}
	if nowSeq, nowSecond, err = incrSeqAndUpdateLastTime(nowSecond); err != nil {
		return "", err
	}
	var id [byteLen]byte
	binary.BigEndian.PutUint32(id[:], nowSecond)
	id[4] = machineID[0]
	id[5] = machineID[1]
	id[6] = machineID[2]
	id[7] = byte(pid >> 8)
	id[8] = byte(pid)
	id[9] = byte(nowSeq >> 16)
	id[10] = byte(nowSeq >> 8)
	id[11] = byte(nowSeq)

	text := make([]byte, encodedLen)
	encode(text, id[:])
	return *(*string)(unsafe.Pointer(&text)), nil
}

func incrSeqAndUpdateLastTime(nowSecond uint32) (uint32, uint32, error) {
	dgMutex.Lock()
	defer dgMutex.Unlock()
	// 判断上次时间是否大于现在时间，如果大于，直接用上次时间
	if lastSecond > nowSecond {
		nowSecond = lastSecond
	} else if nowSecond > lastSecond {
		lastSecond = nowSecond
		// 时间更新，seq清零
		seq = 0
	}
	tSeq := seq + 1
	if tSeq == maxSeq {
		return 0, 0, errors.New("seq big")
	}
	seq = tSeq
	return seq, lastSecond, nil
}


// encode base32编码
func encode(dst, id []byte) {
	_ = dst[19]
	_ = id[11]
	dst[19] = encoding[(id[11]<<4)&0x1F]
	dst[18] = encoding[(id[11]>>1)&0x1F]
	dst[17] = encoding[(id[11]>>6)&0x1F|(id[10]<<2)&0x1F]
	dst[16] = encoding[id[10]>>3]
	dst[15] = encoding[id[9]&0x1F]
	dst[14] = encoding[(id[9]>>5)|(id[8]<<3)&0x1F]
	dst[13] = encoding[(id[8]>>2)&0x1F]
	dst[12] = encoding[id[8]>>7|(id[7]<<1)&0x1F]
	dst[11] = encoding[(id[7]>>4)&0x1F|(id[6]<<4)&0x1F]
	dst[10] = encoding[(id[6]>>1)&0x1F]
	dst[9] = encoding[(id[6]>>6)&0x1F|(id[5]<<2)&0x1F]
	dst[8] = encoding[id[5]>>3]
	dst[7] = encoding[id[4]&0x1F]
	dst[6] = encoding[id[4]>>5|(id[3]<<3)&0x1F]
	dst[5] = encoding[(id[3]>>2)&0x1F]
	dst[4] = encoding[id[3]>>7|(id[2]<<1)&0x1F]
	dst[3] = encoding[(id[2]>>4)&0x1F|(id[1]<<4)&0x1F]
	dst[2] = encoding[(id[1]>>1)&0x1F]
	dst[1] = encoding[(id[1]>>6)&0x1F|(id[0]<<2)&0x1F]
	dst[0] = encoding[id[0]>>3]
}
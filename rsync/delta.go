package rsync

import (
	"bytes"
	"os"

	"bitbucket.org/kardianos/rsync"
	"bitbucket.org/kardianos/rsync/proto"
)

const maxDataOp = 1024 * 16

func createRsync() *rsync.RSync {
	return &rsync.RSync{
		MaxDataOp: maxDataOp,
	}
}

func signature(basis string, blockSizeKiB int) (sigBuffer bytes.Buffer, err error) {
	rs := createRsync()
	rs.BlockSize = 1024 * blockSizeKiB

	basisFile, err := os.Open(basis)
	if err != nil {
		return
	}
	defer basisFile.Close()

	sigEncode := &proto.Writer{Writer: sigBuffer}

	err = sigEncode.Header(proto.TypeSignature, proto.CompNone, rs.BlockSize)
	if err != nil {
		return
	}
	defer sigEncode.Close()

	err = rs.CreateSignature(basisFile, sigEncode.SignatureWriter())
	return
}

func delta(sigBuffer bytes.Buffer, newfile string, checkFile bool, comp proto.Comp) (deltaBuffer bytes.Buffer, err error) {
	rs := createRsync()

	nfFile, err := os.Open(newfile)
	if err != nil {
		return
	}
	defer nfFile.Close()

	// Load signature hash list.
	sigDecode := &proto.Reader{Reader: sigBuffer}
	rs.BlockSize, err = sigDecode.Header(proto.TypeSignature)
	if err != nil {
		if err == io.EOF {
			err = io.ErrUnexpectedEOF
		}
		return
	}
	defer sigDecode.Close()

	hl, err := sigDecode.ReadAllSignatures()
	if err != nil {
		return
	}
	if *verbose {
		fmt.Printf("Signature Count: %d\n", len(hl))
	}

	// Save operations to buffer.
	opsEncode := &proto.Writer{Writer: deltaBuffer}
	err = opsEncode.Header(proto.TypeDelta, comp, rs.BlockSize)
	if err != nil {
		return
	}
	defer opsEncode.Close()

	var hasher hash.Hash
	if checkFile {
		hasher = md5.New()
	}
	opF := opsEncode.OperationWriter()
	err = rs.CreateDelta(nfFile, hl, opF, hasher)
	if err != nil {
		return
	}
	if checkFile {
		err = opF(rsync.Operation{
			Type: rsync.OpHash,
			Data: hasher.Sum(nil),
		})
	}
	return
}

func patch(deltaBuffer bytes.Buffer, basis string, newfile string, checkFile bool) (err error) {
	rs := createRsync()
	basisFile, err := os.Open(basis)
	if err != nil {
		return
	}
	defer basisFile.Close()

	fsFile, err := os.Create(newfile)
	if err != nil {
		return
	}
	defer fsFile.Close()

	deltaDecode := proto.Reader{Reader: deltaBuffer}
	rs.BlockSize, err = deltaDecode.Header(proto.TypeDelta)
	if err != nil {
		if err == io.EOF {
			err = io.ErrUnexpectedEOF
		}
		return
	}
	defer deltaDecode.Close()

	hashOps := make(chan rsync.Operation, 2)
	ops := make(chan rsync.Operation)
	// Load operations from file.
	var decodeError error
	go func() {
		defer close(ops)
		decodeError = deltaDecode.ReadOperations(ops, hashOps)
	}()

	var hasher hash.Hash
	if checkFile {
		hasher = md5.New()
	}
	err = rs.ApplyDelta(fsFile, basisFile, ops, hasher)
	if err != nil {
		return
	}
	if decodeError != nil {
		err = decodeError
		return
	}
	if checkFile == false {
		err = nil
		return
	}
	hashOp := <-hashOps
	if hashOp.Data == nil {
		err = NoTargetSumError
		return
	}
	if bytes.Equal(hashOp.Data, hasher.Sum(nil)) == false {
		err = HashNoMatchError
		return
	}

	return
}

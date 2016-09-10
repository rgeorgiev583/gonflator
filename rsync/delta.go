package rsync

import (
	"bitbucket.org/kardianos/rsync"
	"bitbucket.org/kardianos/rsync/proto"
)

func getRsync() *rsync.RSync {
	return &rsync.RSync{
		MaxDataOp: 1024 * 16,
	}
}

func signature(basis, signature string, blockSizeKiB int) error {
	rs := getRsync()
	rs.BlockSize = 1024 * blockSizeKiB

	basisFile, err := os.Open(basis)
	if err != nil {
		return err
	}
	defer basisFile.Close()

	sigFile, err := os.Create(signature)
	if err != nil {
		return err
	}
	defer sigFile.Close()

	sigEncode := &proto.Writer{Writer: sigFile}

	err = sigEncode.Header(proto.TypeSignature, proto.CompNone, rs.BlockSize)
	if err != nil {
		return err
	}
	defer sigEncode.Close()

	return rs.CreateSignature(basisFile, sigEncode.SignatureWriter())
}

func delta(signature, newfile, delta string, checkFile bool, comp proto.Comp) error {
	rs := getRsync()
	sigFile, err := os.Open(signature)
	if err != nil {
		return err
	}
	defer sigFile.Close()

	nfFile, err := os.Open(newfile)
	if err != nil {
		return err
	}
	defer nfFile.Close()

	deltaFile, err := os.Create(delta)
	if err != nil {
		return err
	}
	defer deltaFile.Close()

	// Load signature hash list.
	sigDecode := &proto.Reader{Reader: sigFile}
	rs.BlockSize, err = sigDecode.Header(proto.TypeSignature)
	if err != nil {
		if err == io.EOF {
			return io.ErrUnexpectedEOF
		}
		return err
	}
	defer sigDecode.Close()

	hl, err := sigDecode.ReadAllSignatures()
	if err != nil {
		return err
	}
	if *verbose {
		fmt.Printf("Signature Count: %d\n", len(hl))
	}

	// Save operations to file.
	opsEncode := &proto.Writer{Writer: deltaFile}
	err = opsEncode.Header(proto.TypeDelta, comp, rs.BlockSize)
	if err != nil {
		return err
	}
	defer opsEncode.Close()

	var hasher hash.Hash
	if checkFile {
		hasher = md5.New()
	}
	opF := opsEncode.OperationWriter()
	err = rs.CreateDelta(nfFile, hl, opF, hasher)
	if err != nil {
		return err
	}
	if checkFile {
		return opF(rsync.Operation{
			Type: rsync.OpHash,
			Data: hasher.Sum(nil),
		})
	}
	return nil
}

func patch(basis, delta, newfile string, checkFile bool) error {
	rs := getRsync()
	basisFile, err := os.Open(basis)
	if err != nil {
		return err
	}
	defer basisFile.Close()

	deltaFile, err := os.Open(delta)
	if err != nil {
		return err
	}
	defer deltaFile.Close()

	fsFile, err := os.Create(newfile)
	if err != nil {
		return err
	}
	defer fsFile.Close()

	deltaDecode := proto.Reader{Reader: deltaFile}
	rs.BlockSize, err = deltaDecode.Header(proto.TypeDelta)
	if err != nil {
		if err == io.EOF {
			return io.ErrUnexpectedEOF
		}
		return err
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
		return err
	}
	if decodeError != nil {
		return decodeError
	}
	if checkFile == false {
		return nil
	}
	hashOp := <-hashOps
	if hashOp.Data == nil {
		return NoTargetSumError
	}
	if bytes.Equal(hashOp.Data, hasher.Sum(nil)) == false {
		return HashNoMatchError
	}

	return nil
}

package core

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/icon-project/rewardcalculator/common"
	"github.com/icon-project/rewardcalculator/common/codec"
	"github.com/icon-project/rewardcalculator/common/db"
	"github.com/syndtr/goleveldb/leveldb/util"
)

const (
	IISSDataVersion         uint64 = 2

	IISSDataRevisionDefault uint64 = 0

	IISSDataDBPrefix =  "iiss_rc_db_"
	IISSDataDBFormat = IISSDataDBPrefix + "%d"	// $BH
)

type IISSHeader struct {
	Version     uint64		// version of RC data
	BlockHeight uint64
	Revision    uint64		// revision of ICON Service
}

func (ih *IISSHeader) ID() []byte {
	return []byte("")
}

func (ih *IISSHeader) Bytes() ([]byte, error) {
	var bytes []byte
	if bs, err := codec.MarshalToBytes(ih); err != nil {
		return nil, err
	} else {
		bytes = bs
	}
	return bytes, nil
}

func (ih *IISSHeader) String() string {
	b, err := json.Marshal(ih)
	if err != nil {
		return "Can't covert Message to json"
	}
	return string(b)
}

func (ih *IISSHeader) SetBytes(bs []byte) error {
	_, err := codec.UnmarshalFromBytes(bs, ih)
	if err != nil {
		return err
	}

	// for backward compatibility
	if ih.Version == 1 {
		ih.Revision = IISSDataRevisionDefault
	}

	return nil
}

func loadIISSHeader(iissDB db.Database) (*IISSHeader, error) {
	bucket, _ := iissDB.GetBucket(db.PrefixIISSHeader)
	data, err := bucket.Get([]byte(""))
	if err != nil {
		return nil, err
	}
	if data == nil {
		err = fmt.Errorf("There is no header data in IISS data\n")
		return nil, err
	}
	header := new(IISSHeader)
	err = header.SetBytes(data)
	if err != nil {
		return nil, err
	}

	return header, nil
}

type IISSGVData struct {
	IncentiveRep  uint64
	RewardRep     uint64
	MainPRepCount uint64
	SubPRepCount  uint64
}

type IISSGovernanceVariable struct {
	BlockHeight   uint64
	IISSGVData
}

func (gv *IISSGovernanceVariable) ID() []byte {
	bs := make([]byte, 8)
	id := common.Uint64ToBytes(gv.BlockHeight)
	copy(bs[len(bs)-len(id):], id)
	return bs
}

func (gv *IISSGovernanceVariable) Bytes() ([]byte, error) {
	var bytes []byte
	if bs, err := codec.MarshalToBytes(&gv.IISSGVData); err != nil {
		return nil, err
	} else {
		bytes = bs
	}
	return bytes, nil
}

func (gv *IISSGovernanceVariable) String() string {
	b, err := json.Marshal(gv)
	if err != nil {
		return "Can't covert Message to json"
	}
	return string(b)
}

func (gv *IISSGovernanceVariable) SetBytes(bs []byte, version uint64) error {
	_, err := codec.UnmarshalFromBytes(bs, &gv.IISSGVData)
	if err != nil {
		return err
	}

	// for backward compatibility
	if version == 1 {
		gv.MainPRepCount = 0
		gv.SubPRepCount = 0
	}

	return nil
}

func loadIISSGovernanceVariable(iissDB db.Database, version uint64) ([]*IISSGovernanceVariable, error) {
	gvList := make([]*IISSGovernanceVariable, 0)
	iter, err := iissDB.GetIterator()
	if err != nil {
		return nil, err
	}
	prefix := util.BytesPrefix([]byte(db.PrefixIISSGV))
	iter.New(prefix.Start, prefix.Limit)
	for entries := 0; iter.Next(); entries++ {
		gv := new(IISSGovernanceVariable)
		err = gv.SetBytes(iter.Value(), version)
		if err != nil {
			return nil, err
		}
		gv.BlockHeight = common.BytesToUint64(iter.Key()[len(db.PrefixIISSGV):])
		gvList = append(gvList, gv)
	}
	iter.Release()
	err = iter.Error()
	if err != nil {
		log.Printf("There is error while read IISS GV iteration. %+v", err)
	}

	return gvList, nil
}

type IISSBlockProduceInfoData struct {
	Generator common.Address
	Validator []common.Address
}

type IISSBlockProduceInfo struct {
	BlockHeight uint64
	IISSBlockProduceInfoData
}

func (bp *IISSBlockProduceInfo) ID() []byte {
	bs := make([]byte, 8)
	id := common.Uint64ToBytes(bp.BlockHeight)
	copy(bs[len(bs)-len(id):], id)
	return bs
}

func (bp *IISSBlockProduceInfo) Bytes() ([]byte, error) {
	var bytes []byte
	if bs, err := codec.MarshalToBytes(&bp.IISSBlockProduceInfoData); err != nil {
		return nil, err
	} else {
		bytes = bs
	}
	return bytes, nil
}

func (bp *IISSBlockProduceInfo) String() string {
	b, err := json.Marshal(bp)
	if err != nil {
		return "Can't covert Message to json"
	}
	return string(b)
}

func (bp *IISSBlockProduceInfo) SetBytes(bs []byte) error {
	_, err := codec.UnmarshalFromBytes(bs, &bp.IISSBlockProduceInfoData)
	if err != nil {
		return err
	}
	return nil
}

const (
	TXDataTypeDelegate  = 0
	TXDataTypePrepReg   = 1
	TXDataTypePrepUnReg = 2
)

type IISSTXData struct {
	Address     common.Address
	BlockHeight uint64
	DataType    uint64
	Data        *codec.TypedObj
}

type IISSTX struct {
	Index uint64
	IISSTXData
}

func (tx *IISSTX) ID() []byte {
	bs := make([]byte, 8)
	id := common.Uint64ToBytes(tx.Index)
	copy(bs[len(bs)-len(id):], id)
	return bs
}

func (tx *IISSTX) Bytes() ([]byte, error) {
	var bytes []byte
	if bs, err := codec.MarshalToBytes(&tx.IISSTXData); err != nil {
		return nil, err
	} else {
		bytes = bs
	}
	return bytes, nil
}

func (tx *IISSTX) String() string {
	b, err := json.Marshal(tx)
	if err != nil {
		return "Can't covert Message to json"
	}

	return fmt.Sprintf("%s\n\t Data: %+v", string(b), common.MustDecodeAny(tx.Data))
}

func (tx *IISSTX) SetBytes(bs []byte) error {
	_, err := codec.UnmarshalFromBytes(bs, &tx.IISSTXData)
	if err != nil {
		return err
	}
	return nil
}

func OpenIISSData(path string) db.Database {
	dbPath := filepath.Clean(path)
	dbDir, dbName := filepath.Split(dbPath)
	return db.Open(dbDir, string(db.GoLevelDBBackend), dbName)
}

func LoadIISSData(iissDB db.Database) (*IISSHeader, []*IISSGovernanceVariable, []*PRep) {
	// Load IISS Data
	header, err := loadIISSHeader(iissDB)
	if err != nil {
		log.Printf("Failed to read header from IISS Data. err=%+v", err)
		return nil, nil, nil
	}
	log.Printf("Header: %s\n", header.String())

	// Governance Variable
	gvList, err := loadIISSGovernanceVariable(iissDB, header.Version)
	if err != nil {
		log.Printf("Failed to read governance variable from IISS Data. err=%+v\n", err)
		return nil, nil, nil
	}
	log.Printf("Governance variable:\n")
	for i, gv := range gvList {
		log.Printf("\t%d: %s", i, gv.String())
	}

	// Main/Sub P-Rep
	pRepList, err := LoadPRep(iissDB)
	if err != nil {
		log.Printf("Failed to read P-Rep list from IISS Data. err=%+v\n", err)
		return nil, nil, nil
	}
	log.Printf("Main/Sub P-Rep list:\n")
	for i, preps:= range pRepList {
		log.Printf("\t%d: %s\n", i, preps.String())
	}

	return header, gvList, pRepList
}

func findIISSData(dir string, prefix string) []os.FileInfo {
	iissData := make([]os.FileInfo, 0)

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil
	}

	for _, f := range files {
		if f.IsDir() == true {
			if prefix != "" && strings.HasPrefix(f.Name(), prefix) != true {
				continue
			}
			iissData = append(iissData, f)
		}
	}

	return iissData
}

func cleanupIISSData(path string) {
	var blockHeight, backupBH int
	dir, name := filepath.Split(path)

	fmt.Sscanf(name, IISSDataDBFormat, &blockHeight)

	// delete old backup data
	for _, backup := range findIISSData(dir, IISSDataDBPrefix) {
		fmt.Sscanf(backup.Name(), IISSDataDBFormat, &backupBH)
		if backupBH < blockHeight {
			newPath := filepath.Join(dir, backup.Name())
			log.Printf("remove backup %s", newPath)
			os.RemoveAll(newPath)
		}
	}
}

func WriteIISSHeader(iiss db.Database, version uint64, blockHeight uint64, revision uint64) error {
	bucket, _ := iiss.GetBucket(db.PrefixIISSHeader)

	header := new(IISSHeader)
	header.Version = version
	header.BlockHeight = blockHeight
	header.Revision = revision

	key := []byte("")
	value, _ := header.Bytes()
	return bucket.Set(key, value)
}

func WriteIISSGV(iiss db.Database,
	blockHeight uint64, incentive uint64, reward uint64, mainPRepCount uint64, subPRepCount uint64) error {
	bucket, _ := iiss.GetBucket(db.PrefixIISSGV)

	gv := new(IISSGovernanceVariable)
	gv.BlockHeight = blockHeight
	gv.IncentiveRep = incentive
	gv.RewardRep = reward
	gv.MainPRepCount = mainPRepCount
	gv.SubPRepCount = subPRepCount

	value, _ := gv.Bytes()
	return bucket.Set(gv.ID(), value)
}

func WriteIISSBP(iiss db.Database, blockHeight uint64, generator string, validator []string) error {
	bucket, _ := iiss.GetBucket(db.PrefixIISSBPInfo)

	prep := new(IISSBlockProduceInfo)
	prep.BlockHeight = blockHeight

	key := prep.ID()
	if generator == "" || len(validator) == 0 {
		return fmt.Errorf("need generator and validator parameter\n")
	}

	prep.Generator = *common.NewAddressFromString(generator)
	prep.Validator = make([]common.Address, 0)
	for _, v := range validator {
		prep.Validator = append(prep.Validator, *common.NewAddressFromString(v))
	}

	value, _ := prep.Bytes()
	return bucket.Set(key, value)
}

func WriteIISSPRep(iiss db.Database, blockHeight uint64, totalDelegation uint64, preps []*PRepDelegationInfo) error {
	if len(preps) == 0 {
		return fmt.Errorf("need preps parameter")
	}

	bucket, _ := iiss.GetBucket(db.PrefixIISSPRep)

	prep := new(PRep)
	prep.BlockHeight = blockHeight
	prep.TotalDelegation.SetUint64(totalDelegation)
	prep.List = make([]PRepDelegationInfo, 0)
	for _, value := range preps {
		prep.List = append(prep.List, *value)
	}

	key := prep.ID()
	value, _ := prep.Bytes()
	return bucket.Set(key, value)
}

func WriteIISSTX(iiss db.Database, index uint64, address string, blockHeight uint64, dataType uint64,
	delegations []*PRepDelegationInfo) error {

	bucket, _ := iiss.GetBucket(db.PrefixIISSTX)

	tx := new(IISSTX)
	tx.Address = *common.NewAddressFromString(address)
	tx.BlockHeight = blockHeight
	tx.DataType = dataType
	tx.Data = new(codec.TypedObj)
	tx.Index = index

	switch tx.DataType {
	case TXDataTypeDelegate:
		var delegation []interface{}
		var dgData []interface{}

		for _, dg := range delegations {
			dgData = append(dgData, &dg.Address)
			dgData = append(dgData, dg.DelegatedAmount.Uint64())
			delegation = append(delegation, dgData)
		}

		var err error
		tx.Data, err = common.EncodeAny(delegation)
		if err != nil {
			return fmt.Errorf("failed to encode stake %+v", err)
		}
	case TXDataTypePrepReg:
		tx.Data.Type = codec.TypeNil
		tx.Data.Object = []byte("")
	case TXDataTypePrepUnReg:
		tx.Data.Type = codec.TypeNil
		tx.Data.Object = []byte("")
	}

	key := tx.ID()
	value, _ := tx.Bytes()
	return bucket.Set(key, value)
}

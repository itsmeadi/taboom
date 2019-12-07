package repositries

import (
	"github.com/omise/go-tamboon/cipher"
	"github.com/omise/go-tamboon/src/custom/customError"
	"github.com/omise/go-tamboon/src/entities"
	"os"
)

type FileSysInterface interface {
	ReadFile(filePath string, userChannel chan entities.UserInfo) error
}

type FileSys struct {
	BufferSize    int                    //Buffer length for fileReading
	UserChannel   chan entities.UserInfo //Channel to push User data
	ChannelLength int
}

func InitFile(buffer int, channelSize int) FileSys {
	return FileSys{
		BufferSize:    buffer,
		UserChannel:   make(chan entities.UserInfo, channelSize),
		ChannelLength: channelSize,
	}
}

//ReadFile will read the filePath and start pushing the data to fileSys.UserChannel channel
func (fileSys *FileSys) ReadFile(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	reader, err := cipher.NewRot128Reader(file)
	if err != nil {
		return err
	}

	buf := make([]byte, fileSys.BufferSize)
	var i int
	length, err := reader.Read(buf)

	if err != nil {
		return err
	}
	go func() {
		for {
			var user entities.UserInfo
			var err error
			iTemp := i

			user, i, err = fileSys.getUserStruct(i, length, buf)

			if err == customError.ErrorInvalidFile {
				break
			}
			if err == customError.ErrorBufferExhausted {
				buf, length, err=fileSys.fillBuffer(reader, iTemp, buf)
				if err != nil {
					break
				}

				i = 0
				continue
			}
			//for len(fileSys.UserChannel) >= (fileSys.ChannelLength-fileSys.ChannelLengthRetry) {
			//	time.Sleep(10 * time.Millisecond) 	//Not occupying complete channel, leaving the remaining for Retry user
			//}
			//log.Print(user)
			fileSys.UserChannel <- user

		}
		//fileSys.UserChannel <- entities.UserInfo{IsLastUser: true} //cannot close the channel now, bcoz retry user still remaining
		close(fileSys.UserChannel)
	}()
	return nil
}

func (fileSys *FileSys) fillBuffer(reader *cipher.Rot128Reader, after int, buf []byte) ([]byte, int, error) {

	newbuf := make([]byte, fileSys.BufferSize)

	residue := buf[after:]
	for j, r := range residue {
		newbuf[j] = r
	}
	length, err := reader.Read(newbuf[fileSys.BufferSize-after:])
	length += fileSys.BufferSize - after

	return newbuf, length, err
}

func (fileSys *FileSys) getUserStruct(i, length int, buf []byte) (entities.UserInfo, int, error) {

	var err error
	var user entities.UserInfo
	user.Name, i, err = readKey(i, length, ",", buf)
	if err != nil {
		return user, i, err
	}
	user.Amount, i, err = readKeyInt(i, length, ",", buf)
	if err != nil {
		return user, i, err
	}
	user.CCNumber, i, err = readKeyByteArr(i, length, ",", buf)
	if err != nil {
		return user, i, err
	}
	user.CVV, i, err = readKeyByteArr(i, length, ",", buf)
	if err != nil {
		return user, i, err
	}
	user.ExpMonth, i, err = readKeyInt(i, length, ",", buf)
	if err != nil {
		return user, i, err
	}
	user.ExpYear, i, err = readKeyInt(i, length, "\n", buf)
	if err != nil {
		return user, i, err
	}
	return user, i, err
}

func readKey(i, length int, breakChar string, buf []byte) (string, int, error) {
	str := ""
	iTemp := i
	gracefully := false
	for ; i < length; i++ {
		w := buf[i]
		if string(w) == breakChar {
			gracefully = true
			break
		}
		str += string(w)
	}
	if len(str) > 100 {
		return "", iTemp, customError.ErrorInvalidFile
	}
	if !gracefully {
		return "", iTemp, customError.ErrorBufferExhausted
	}
	return str, i + 1, nil
}
func readKeyInt(i, length int, breakChar string, buf []byte) (int64, int, error) {
	var val int64
	iTemp := i
	gracefully := false
	for ; i < length; i++ {
		w := buf[i]
		if string(w) == breakChar {
			gracefully = true
			break
		}
		val = val*10 + int64(w) - 48
	}
	if val > 1<<50 {
		return val, iTemp, customError.ErrorInvalidFile
	}
	if !gracefully {
		return -1, iTemp, customError.ErrorBufferExhausted
	}
	return val, i + 1, nil
}

func readKeyByteArr(i, length int, breakChar string, buf []byte) ([]byte, int, error) {
	var val []byte
	iTemp := i
	gracefully := false
	for ; i < length; i++ {
		w := buf[i]
		if string(w) == breakChar {
			gracefully = true
			break
		}
		val = append(val, w)
	}
	if len(val) > 100 {
		return val, iTemp, customError.ErrorInvalidFile
	}
	if !gracefully {
		return val, iTemp, customError.ErrorBufferExhausted
	}
	return val, i + 1, nil
}

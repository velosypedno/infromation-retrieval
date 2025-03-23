package internal

type ExtToReaderMapper struct{}

func (m ExtToReaderMapper) Map(ext string) FileReader {
	fileReaderMap := map[string]FileReader{
		".txt": DefaultFileReader{},
		".xml": XMLFileReader{},
		".fb2": XMLFileReader{},
	}
	reader := fileReaderMap[ext]
	return reader

}

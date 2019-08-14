package fake

type FakeRequestReader struct {
	data string
}

func (self FakeRequestReader) Read() string {
	return self.data
}

func CreateFakeRequestReader(data string) FakeRequestReader {
	return FakeRequestReader{
		data: data,
	}
}

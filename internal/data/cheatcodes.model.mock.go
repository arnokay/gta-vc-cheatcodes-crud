package data

type MockCheatcodeModel struct{}

func (m MockCheatcodeModel) Insert(cheatcode *Cheatcode) error {
	return nil
}

func (m MockCheatcodeModel) Get(id int64) (*Cheatcode, error) {
	return nil, nil
}

func (m MockCheatcodeModel) Update(cheatcode *Cheatcode) error {
	return nil
}

func (m MockCheatcodeModel) Delete(id int64) error {
	return nil
}

func (m MockCheatcodeModel) GetAll(code string, description string, tags []string, filters Filters) ([]*Cheatcode, Metadata, error) {
	return nil, Metadata{}, nil
}

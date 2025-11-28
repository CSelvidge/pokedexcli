package actors


func NewUser() (*User, error) {
	return &User{
		CaughtPokemon: make(map[string]Pokemon),
	}, nil
}

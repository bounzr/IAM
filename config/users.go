package config

type Users struct{
	Implementation	string	`yaml:"implementation"`
	Admin			Admin	`yaml:"admin"`
}

type Admin struct{
	Username		string	`yaml:"username"`
	Password		string	`yaml:"password"`
	Repository		string	`yaml:"repository"`
	HideAfterInit	bool	`yaml:"hideAfterInit"`
}

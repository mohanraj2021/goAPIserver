package model

type LoUser struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	// Isadmin  string `json:"isadmin"`
}

type Users struct {
	Id       int
	Username string
	Email    string
	Password string
	Isadmin  string
}

type Product struct {
	Id          int
	Productname string
	Description string
	Price       int
	Owner       string
	Sqrft       int
	Bedroom     int
	Bathroom    int
	Location    string
	City        string
}

type Addcart struct {
	Productid  int
	Customerid int
}

type Cartdetails struct {
	Productid      int
	Customerid     int
	Eventtimestamp string
}

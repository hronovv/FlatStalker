package models

type Payment struct {
	ID        int64
	UserID    int64
	ChatID    int64
	Payload   string
	Plan      string
	Days      int
	AmountKop int
	Status    string
}

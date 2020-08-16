package models

type SmtpServer struct {
	SmtpAddr    string `json:"m1"`
	Port        string `json:"m2"`
	From     	string `json:"m3"`
	Pass 		string `json:"m4"`
}
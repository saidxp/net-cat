package helpers

import "database/sql"
import "unicode"

func Check(str string) bool{
	if len(str) <= 1 {
		return false
	} 
	i := 0
	for i < len(str)-1 {
		if (str[i] < 65 || str[i] > 90) && (str[i] < 97 || str[i] > 122){
			return false  
		}
		i++
	}
	return true
}
func CheckMessage(s string) bool {
	/*for _, r := range s { 
		if r < 32 || r == 127 {
			return false
		}
	}
	*/
	i := 0
	for i < len(s) {
		if s[i] == 27 {  
			if i+1 < len(s) && s[i+1] == '[' {
				i += 2			 
				for i < len(s) && (unicode.IsDigit(rune(s[i])) || s[i] == ';') {
					i++
				}
				if i < len(s) && unicode.IsLetter(rune(s[i])) {
					return false
				}
			}
		}
		i++
	}
	return true
}

	
func Checkmap(name string, auth *Authentication) bool {
	for key , _ := range auth.Con {
		if key == name {
			return false
		}
	}
	return true
}
func Getmessages(db *sql.DB, groube string) ([]string, error) {
	tablerow, err := db.Query("SELECT content FROM messages;")
	if err != nil {
		return nil, err
	}
	defer tablerow.Close()
	var messages []string
	for tablerow.Next() {
		var msg string
		if err := tablerow.Scan(&msg); err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}
	return messages, nil
}
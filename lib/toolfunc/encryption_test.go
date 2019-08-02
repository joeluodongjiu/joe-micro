package toolfunc

import (
	"fmt"
	"testing"
)

func TestGenerateSalt(t *testing.T) {
	fmt.Println(GenerateSalt())
}



func TestEncUserPwdc(t *testing.T) {
	fmt.Println(EncUserPwd("super","6d16bc7"))
}